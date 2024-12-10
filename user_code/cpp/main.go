package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func compileCppCode(sourceFilePath, outputFilePath string) string {
	// 准备编译命令
	compileCmd := exec.Command("g++", sourceFilePath, "-o", outputFilePath, "-std=c++11")
	var capturedCompileErr bytes.Buffer // 只需要捕捉编译错误的信息
	compileCmd.Stderr = &capturedCompileErr
	if err := compileCmd.Run(); err != nil {
		return capturedCompileErr.String()
	}
	return ""
}

func getCurrentMemory(pid int) (int, error) {
	processStatuFilePath := fmt.Sprintf("/proc/%d/status", pid)
	// NOTE: 需要确保传入进来的pid时正确切存在的，不然这一步Open会返回错误
	processStatuFile, _ := os.Open(processStatuFilePath)
	defer processStatuFile.Close()
	scanner := bufio.NewScanner(processStatuFile)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "VmRSS:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				currentMemoryKB, err := strconv.Atoi(fields[1])
				if err != nil {
					return -1, err
				}
				return currentMemoryKB / 1024, nil
			}
		}
	}
	return -1, fmt.Errorf("NO_MEMORY_STATUS_ERROR")
}

func monitorMemory(pid int, maxMemory int, done chan string, exitedMemory chan int) {
	timeTicker := time.NewTicker(20 * time.Millisecond)
	defer timeTicker.Stop()

	for {
		currentMemory, _ := getCurrentMemory(pid)
		select {
		case <-done:
			exitedMemory <- currentMemory // 返回正常退出的内存
			//fmt.Printf("monitorMemory goroutine exit by %s\n", info)
			return
		case <-timeTicker.C:
			if currentMemory > maxMemory {
				exitedMemory <- currentMemory // 返回内存超出的内存
				proc, err := os.FindProcess(pid)
				if err != nil {
					//done <- "os.FindProcess"
					return
				}
				//fmt.Printf("pid: %d 当前使用内存为: %d 超出内存限制，杀死退出\n", pid, currentMemory)
				proc.Kill()
				return
			}
		}
	}
}

func runCppCode(sourceFilePath, input string, maxTime time.Duration, maxMemory int) (string, string, time.Duration, int) {
	// 准备运行命令
	runCmd := exec.Command(sourceFilePath)

	// 准备标准输入输出以及标准错误
	var stdOut bytes.Buffer
	var stdIn bytes.Buffer
	stdIn.WriteString(input)
	runCmd.Stdin = &stdIn
	runCmd.Stdout = &stdOut

	// 设置监控内存的子协程的完成信号通道
	monitoryMemoryDone := make(chan string)
	exitedMemory := make(chan int)
	// 设置超时时间
	time.AfterFunc(maxTime*time.Second, func() {
		// 结束监控内存协程
		monitoryMemoryDone <- "timeout"
		// 超时杀死,结束当前这次执行
		runCmd.Process.Kill()

	})

	// 记录开始时间
	start := time.Now()
	// 通过cmd.Start()启动的子进程，如果是被cmd.Process().Kill()杀死的，最终`signal: Killed`信号也会通过cmd.Wait()返回的
	if err := runCmd.Start(); err != nil {
		// NOTE: 这个错误没有妥善处理
		return err.Error(), "", -1, -1
	}

	// 开启监控内存协程
	go monitorMemory(runCmd.Process.Pid, maxMemory, monitoryMemoryDone, exitedMemory)

	// 等待运行完成
	// 运行时的错误通过cmd.Wait()的err返回的
	if err := runCmd.Wait(); err != nil {
		// 运行时出现错误
		duration := time.Since(start)
		m := <-exitedMemory
		return err.Error(), "", duration, m
	}
	monitoryMemoryDone <- "normal" // 表示正常退出
	// 运行时没有错误，没有超时，没有内存超出，正常返回输出数据，运行时间，内存使用，交给上层对比用户代码输出结果
	duration := time.Since(start)
	m := <-exitedMemory // 获取正常退出的内存使用情况
	return "", stdOut.String(), duration, m
}

func getInputCasesFromDocker() {
	var dataPath = "/user_code/cpp/inputCases.txt"
	var segm = "-=-=-=-=\n"
	data, _ := os.ReadFile(dataPath)
	inputCases := strings.Split(string(data), segm)
	inputCases = inputCases[:len(inputCases)-1]
	for _, v := range inputCases {
		fmt.Printf("%s\n", strconv.Quote(v))
	}
}

func main() {
	// Step 1: 定义文件路径
	sourceFile := "/user_code/cpp/code.cpp"
	outputFile := "/user_code/cpp/code"

	// Step 2: 编译 C++ 源文件
	capturedCompileErr := compileCppCode(sourceFile, outputFile)
	if capturedCompileErr != "" {
		fmt.Printf("@@RECORD@@\n@@OUTPUT@@\n%s@@OUTPUT@@\n%s@@OUTPUT@@\n", "CompileError", capturedCompileErr)
		return
	}

	// Step 3: 运行 C++文件
	var dataPath = "/user_code/cpp/inputCases.txt"
	var segm = "-=-=-=-=\n"
	data, _ := os.ReadFile(dataPath)
	inputCases := strings.Split(string(data), segm)
	inputCases = inputCases[:len(inputCases)-1]
	//err, stdOut, duration, memory := runCppCode(outputFile, "2 0\n", 1, 128)
	//if err != "" {
	//	fmt.Printf("@@RECORD@@\n@@OUTPUT@@\n%s@@OUTPUT@@\n%v@@OUTPUT@@\n%v@@OUTPUT@@\n%s@@OUTPUT@@\n", "RunTimeError", duration, memory, err)
	//} else {
	//	fmt.Printf("@@RECORD@@\n@@OUTPUT@@\n%s@@OUTPUT@@\n%v@@OUTPUT@@\n%v@@OUTPUT@@\n%s@@OUTPUT@@\n", "Pass", duration, memory, stdOut)
	//}
	for _, v := range inputCases {
		//fmt.Printf("%s\n", strconv.Quote(v))
		err, stdOut, duration, memory := runCppCode(outputFile, v, 3, 128)
		if err != "" {
			fmt.Printf("@@RECORD@@\n@@OUTPUT@@\n%s@@OUTPUT@@\n%v@@OUTPUT@@\n%d@@OUTPUT@@\n%s@@OUTPUT@@\n", "RunTimeError", duration, memory, err)
		} else {
			fmt.Printf("@@RECORD@@\n@@OUTPUT@@\n%s@@OUTPUT@@\n%v@@OUTPUT@@\n%d@@OUTPUT@@\n%s@@OUTPUT@@\n", "Pass", duration, memory, stdOut)
		}

	}

	// 第一个表示信息类型：编译错误，运行时错误，通过
	// 如果不是编译错误，第二第三第四分别是运行时间，内存使用，标准输出/运行时错误
	//	@@OUTPUT@@\n()@@OUTPUT@@\n()@@OUTPUT@@\n()@@OUTPUT@@\n()@@OUTPUT@@\n

	//time.Sleep(1 * time.Second)
}
