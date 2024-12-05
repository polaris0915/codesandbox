package service

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestCmd(t *testing.T) {

	// 执行一个可能出错的命令
	cmd := exec.Command("g++", "../test.cpp")
	output, err := cmd.CombinedOutput() // 同时捕获标准输出和错误输出
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println(string(output)) // 错误信息会在这里打印
}

var c chan int

func handle(int) {

}

func excuteCode(ctx context.Context, resChan chan string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("timeout")
		default:
			fmt.Println("开始执行")
			func() {
				fmt.Println("执行代码中")
				res := func() string {
					time.Sleep(10 * time.Second)
					return "ok"
				}()
				resChan <- res
			}()
			return
		}
	}

}

func TestAfter(t *testing.T) {

	time.AfterFunc(10*time.Second, func() {
		fmt.Println("timeout")
	})
	task := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	for _, v := range task {
		fmt.Printf("开始执行task%d\n", v)
		time.Sleep(time.Duration(1) * time.Second)
	}
}

func TestTimeAfter(t *testing.T) {
	// 主协程开启一个子协程去执行任务。并开始计时
	// 如果5秒中之后子协程还有没有完成任务，就关闭子协程。
	// 如果5秒中之内子协程完成任务，就告诉主协程任务完成，并返回一个字符串的结果
	resChan := make(chan string)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go excuteCode(ctx, resChan)

	var res string
	select {
	case <-time.After(5 * time.Second):
		cancel()
	case res = <-resChan:
		fmt.Printf("获取到子协程的结果: %s\n", res)
		//cancel()
	}
}

func TestRuntime(t *testing.T) {

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	fmt.Printf("Alloc: %v MB, TotalAlloc: %v MB, Sys: %v MB, NumGC: %v\n",
		memStats.Alloc/1024/1024,
		memStats.TotalAlloc/1024/1024,
		memStats.Sys/1024/1024,
		memStats.NumGC)

	time.Sleep(time.Second)
}

func TestTimeTicker(t *testing.T) {

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	done := make(chan bool)
	go func() {
		time.Sleep(10 * time.Second)
		done <- true
	}()
	for {
		select {
		case <-done:
			fmt.Println("Done!")
			return
		case t := <-ticker.C:
			fmt.Println("Current time: ", t)
		}
	}
}

type Iresponse interface {
	FinishedResponse()
	TimeOut()
}

type Response struct {
	Activity string `json:"activity"`
	Status   string `json:"status"`
}

type FinishedResponse struct {
	Response
	StdErr string `json:"stderr"`
	StdOut string `json:"stdout"`
}

func (r *FinishedResponse) FinishedResponse() {
	fmt.Println("FinishedResponse-FinishedResponse")
	fmt.Println(r.StdErr)
	fmt.Println(r.StdOut)
	fmt.Println(r.Status)
	fmt.Println(r.Activity)
}

type TestResponse struct {
	Response
	StdErr string `json:"stderr"`
	StdOut string `json:"stdout"`
	Test   string `json:"test"`
}

func (r *TestResponse) FinishedResponse() {
	fmt.Println("TestResponse-FinishedResponse")
	fmt.Println(r.StdErr)
	fmt.Println(r.StdOut)
	fmt.Println(r.Status)
	fmt.Println(r.Activity)
}

func TestChannel(t *testing.T) {
	var wg sync.WaitGroup
	responseChan := make(chan Iresponse)

	// 创建 FinishedResponse 对象
	_response := TestResponse{
		Response: Response{
			Activity: "run_code",
			Status:   "finished",
		},
		StdErr: "",
		StdOut: "1 2\n",
	}

	// 启动一个 goroutine 来发送数据到 channel
	wg.Add(1)
	go func() {
		defer wg.Done()
		responseChan <- &_response
	}()

	// 启动另一个 goroutine 来接收数据并调用 FinishedResponse 方法
	wg.Add(1)
	go func() {
		defer wg.Done()
		r := <-responseChan
		r.FinishedResponse()
		close(responseChan)
	}()

	// 等待所有 goroutine 完成
	wg.Wait()
}

type Type1 struct {
	A int
}

type Type2 struct {
	Type1
	B string
}

func TestAssert(t *testing.T) {
	var wg sync.WaitGroup
	typeChan := make(chan any)

	type2 := Type2{
		Type1: Type1{
			A: 1,
		},
		B: "test",
	}

	// 启动一个 goroutine 来发送数据到 channel
	wg.Add(1)
	go func() {
		defer wg.Done()
		typeChan <- &type2
	}()

	// 启动另一个 goroutine 来接收数据并调用 FinishedResponse 方法
	wg.Add(1)
	go func() {
		defer wg.Done()
		_type := <-typeChan
		r := _type.(Type2)
		r.A
		close(responseChan)
	}()

	// 等待所有 goroutine 完成
	wg.Wait()
}
