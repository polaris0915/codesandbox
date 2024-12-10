package service

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/polaris/codesandbox/api"
	"github.com/polaris/codesandbox/api/response"
	"github.com/polaris/codesandbox/logger"
	"github.com/polaris/codesandbox/model"
	"github.com/polaris/codesandbox/settings"
	"github.com/polaris/codesandbox/utils/_string"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var SEGM = "-=-=-=-=\n"

var (
	Cpp               = "cpp"
	CppDirPath        = settings.EnvConfig.RootPath + "/user_code/cpp"
	CppCodePath       = settings.EnvConfig.RootPath + "/user_code/cpp/code.cpp"
	CppInputCasesPath = settings.EnvConfig.RootPath + "/user_code/cpp/inputCases.txt"
)

type Service struct {
	Activity     string
	Mutex        *sync.Mutex
	ResponseChan chan response.WebSocketResponse
}

func NewService(activity string, mutex *sync.Mutex, responseChan chan response.WebSocketResponse) *Service {
	return &Service{Activity: activity, Mutex: mutex, ResponseChan: responseChan}
}

/*
设置成包内的全局函数，用锁类控制每个语言的沙箱同一时间只有一个用户可以执行代码
// TODO: 有待验证如果是属于类的对象函数加锁跟这种现在写的全局函数加锁有什么区别
*/
func excuteRunCppCode(responseChan chan response.WebSocketResponse, act string) {
	// 判断cpp容器是否在运行， 如果不在运行就先运行起来
	var cli *client.Client
	var err error
	ctx := context.Background()
	if cli, err = inspectDocker(ctx); err != nil {
		responseChan <- response.NewSystemErrorResponse(api.SystemError)
		return
	}
	// 编写运行容器中/user_code/run/main.go的命令
	var output string
	if output, err = excute(cli, ctx); err != nil {
		responseChan <- response.NewSystemErrorResponse(api.SystemError)
		return
	}
	// 解析输出并返回
	//fmt.Printf("%s\n", strconv.Quote(output))
	//fmt.Printf("%s\n", output)
	// 第一个表示信息类型：编译错误，运行时错误，通过
	// 如果不是编译错误，第二第三第四分别是运行时间，内存使用，标准输出/运行时错误
	rawOutputs := strings.Split(output, "@@RECORD@@\n")
	rawOutputs = rawOutputs[1:]
	var outputs [][]string
	for _, _output := range rawOutputs {
		a := strings.Split(_output, "@@OUTPUT@@\n")
		outputs = append(outputs, a[1:len(a)-1])
	}
	for i := 0; i < len(outputs); i++ {
		switch outputs[i][0] {
		case "CompileError":
			// 编译错误直接返回
			responseChan <- response.NewCompileErrorResponse(act, outputs[i][1])
			return
		case "RunTimeError":
			// 运行时错误，返回时间，内存，以及错误信息
			duration, _ := time.ParseDuration(outputs[i][1])
			memory, _ := strconv.Atoi(outputs[i][2])
			responseChan <- response.NewRunTimeErrorResponse(act, "", "", "", duration, int64(memory), outputs[i][3])
			return
		case "Pass":
			duration, _ := time.ParseDuration(outputs[i][1])
			memory, _ := strconv.Atoi(outputs[i][2])
			responseChan <- response.NewFinishedResponse(act, outputs[i][3], duration, int64(memory))
			return
		}
	}
}

func excuteSubmitCppCode(responseChan chan response.WebSocketResponse, act string, testCases []model.JudgeCase, judgeConfig *model.JudgeConfig, syncBody *response.SyncBody) {
	// 判断cpp容器是否在运行， 如果不在运行就先运行起来
	var cli *client.Client
	var err error
	ctx := context.Background()
	if cli, err = inspectDocker(ctx); err != nil {
		responseChan <- response.NewSystemErrorResponse(api.SystemError)
		return
	}
	// 编写运行容器中/user_code/run/main.go的命令
	var output string
	if output, err = excute(cli, ctx); err != nil {
		responseChan <- response.NewSystemErrorResponse(api.SystemError)
		return
	}
	// 解析输出并返回
	//fmt.Printf("%s\n", strconv.Quote(output))
	//fmt.Printf("%s\n", output)
	// 第一个表示信息类型：编译错误，运行时错误，通过
	// 如果不是编译错误，第二第三第四分别是运行时间，内存使用，标准输出/运行时错误
	rawOutputs := strings.Split(output, "@@RECORD@@\n")
	rawOutputs = rawOutputs[1:]
	var outputs [][]string
	for _, _output := range rawOutputs {
		a := strings.Split(_output, "@@OUTPUT@@\n")
		outputs = append(outputs, a[1:len(a)-1])
	}

	var totalTime time.Duration
	var totalMemory int

	for i := 0; i < len(outputs); i++ {
		switch outputs[i][0] {
		case "CompileError":
			// 编译错误直接返回
			responseChan <- response.NewCompileErrorResponse(act, outputs[i][1])
			syncBody.JudgeInfo.JudgeResult = api.CompileError
			return
		case "RunTimeError":
			// 运行时错误，返回时间，内存，以及错误信息
			duration, _ := time.ParseDuration(outputs[i][1])
			memory, _ := strconv.Atoi(outputs[i][2])
			responseChan <- response.NewRunTimeErrorResponse(act, "", "", "", duration, int64(memory), outputs[i][3])
			syncBody.JudgeInfo.JudgeResult = api.CompileError
			return
		case "Pass":
			duration, _ := time.ParseDuration(outputs[i][1])
			memory, _ := strconv.Atoi(outputs[i][2])
			if duration.Milliseconds() > judgeConfig.TimeLimit*1000 {
				responseChan <- response.NewTimeoutResponse(act, testCases[i].Input, testCases[i].Output, outputs[i][3], duration, int64(memory))
				syncBody.JudgeInfo.JudgeResult = api.Timeout
				return
			}
			if memory > judgeConfig.MemoryLimit {
				responseChan <- response.NewMemoryExceededResponse(act, testCases[i].Input, testCases[i].Output, outputs[i][3], duration, int64(memory))
				syncBody.JudgeInfo.JudgeResult = api.MemoryExceeded
				return
			}
			if outputs[i][3] != testCases[i].Output {
				//fmt.Printf("用户输出: %s\n", strconv.Quote(outputs[i][3]))
				//fmt.Printf("标准输出: %s\n", strconv.Quote(testCases[i].Output))
				_string.EndWithBr(&outputs[i][3])
				if outputs[i][3] == testCases[i].Output {
					responseChan <- response.NewPresentationErrorResponse(act, testCases[i].Input, testCases[i].Output, outputs[i][3], duration, int64(memory))
					syncBody.JudgeInfo.JudgeResult = api.PresentationError
					return
				}
				// 换行前面如果有一个空格，这个是允许通过的，具体原因参考acwing上的快排
				testCases[i].Output, _ = strings.CutSuffix(testCases[i].Output, "\n")
				testCases[i].Output += " \n"
				if outputs[i][3] != testCases[i].Output {
					responseChan <- response.NewWrongAnswerResponse(act, testCases[i].Input, testCases[i].Output, outputs[i][3], duration, int64(memory))
					syncBody.JudgeInfo.JudgeResult = api.WrongAnswer
					return
				}
			}
			totalTime += duration
			totalMemory += memory
		}
	}
	responseChan <- response.NewAcceptResponse(act, totalTime, int64(totalMemory))
	syncBody.JudgeInfo.JudgeResult = api.Accepted
}

func excute(cli *client.Client, ctx context.Context) (string, error) {
	// 创建 exec 实例，配置命令和流选项
	execConfig := container.ExecOptions{
		Cmd:          []string{"go", "run", "/user_code/cpp/main.go"},
		AttachStdout: true, // 捕获标准输出
		AttachStderr: true, // 捕获标准错误
		Env:          []string{"PATH=/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"},
	}

	execResp, err := cli.ContainerExecCreate(ctx, settings.ContainersConfig.Cpp.ContainerId, execConfig)
	if err != nil {
		return "", err
	}

	// 附加到 exec 的输入/输出流
	resp, err := cli.ContainerExecAttach(ctx, execResp.ID, container.ExecStartOptions{})
	if err != nil {
		return "", err
	}
	defer resp.Close()

	// 读取输出流
	output, err := io.ReadAll(resp.Reader)
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func inspectDocker(ctx context.Context) (*client.Client, error) {
	// 创建 Docker 客户端
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logger.GetLogger().Errorf("无法创建docker客户端连接: %s", err.Error())
		return nil, err

	}

	// 获取容器详细信息
	cppContainer, err := cli.ContainerInspect(ctx, settings.ContainersConfig.Cpp.ContainerId)
	if err != nil {
		logger.GetLogger().Errorf("无法获取c++容器信息: %s", err.Error())
		return nil, err
	}

	if cppContainer.State.Running {
		return cli, nil
	}
	if err := cli.ContainerStart(ctx, settings.ContainersConfig.Cpp.ContainerId, container.StartOptions{}); err != nil {
		logger.GetLogger().Errorf("无法获取运行c++容器: %s", err.Error())
		return nil, err
	}

	return cli, nil
}

func setInputCasesToDocker(language string, inputs any) {
	var dataPath string
	var dataDirPath string
	switch language {
	case Cpp:
		dataPath = settings.EnvConfig.RootPath + CppInputCasesPath
		dataDirPath = settings.EnvConfig.RootPath + CppDirPath
	}

	os.MkdirAll(dataDirPath, os.ModePerm)
	fd, _ := os.OpenFile(dataPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	switch res := inputs.(type) {
	case string:
		fd.WriteString(res)
		fd.WriteString(SEGM)
	case []string:
		for _, v := range res {
			fd.WriteString(v)
			fd.WriteString(SEGM)
		}
	}
	fd.Close()
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
