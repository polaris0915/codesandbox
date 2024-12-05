package docker

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/polaris/codesandbox/api/response"
	"io"
	"os"
	"strings"
	"time"
)

type SandBoxConfig struct {
	SandBoxName string `json:"SandBoxName" mapstructure:"SandBoxName"`
}

func getSandBoxName(rawSandBoxName string) string {
	return fmt.Sprintf("/%s", rawSandBoxName)
}

type SandBox struct {
	Error     error // 存储错误信息
	ctx       context.Context
	cli       *client.Client // 连接客户端
	container struct {
		container types.Container // 运行容器
		Available bool
	}
	ExcuteResult string
	//FilterResult string
	//Input       string
	UsedTime    int64
	ExcutedInfo struct {
		Message string `json:"message"` // 用于说明 是否超时，内存是否超出限制，编译错误等等之类的
		Memory  int64  `json:"memory"`
		Time    int32  `json:"time"`
	}
}

func NewSandBox() *SandBox {
	sandbox := &SandBox{
		Error: nil,
		ctx:   context.Background(),
		cli:   nil,
		container: struct {
			container types.Container
			Available bool
		}{
			container: types.Container{},
			Available: false,
		},
		ExcuteResult: "",
		//FilterResult: "",
		//Input: "",
		UsedTime: -1,
	}
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		sandbox.Error = err
		return sandbox
	}

	// 获取正在运行的容器列表
	containers, err := cli.ContainerList(sandbox.ctx, container.ListOptions{
		All: false, // 只显示正在运行的容器
	})
	if err != nil {
		sandbox.Error = err
	}

	for _, container := range containers {
		if container.Names[0] == "/cpp-codesandbox" {
			sandbox.container.container = container
			sandbox.container.Available = true
			break
		}
	}
	if sandbox.container.Available == false {
		sandbox.Error = fmt.Errorf("没有%s沙箱环境", "cpp-codesandbox")
	}

	// 启动容器
	if err := cli.ContainerStart(sandbox.ctx, sandbox.container.container.ID, container.StartOptions{}); err != nil {
		sandbox.Error = fmt.Errorf("启动容器 `%s` 失败: %s", sandbox.container.container.Names[0], err)
	}

	sandbox.cli = cli
	return sandbox
}

func (sandBox *SandBox) excuteCode(execConfig container.ExecOptions, input string) string {
	// 创建 Exec 实例
	respRunID, err := sandBox.cli.ContainerExecCreate(sandBox.ctx, sandBox.container.container.ID, execConfig)
	if err != nil {
		sandBox.Error = fmt.Errorf("创建 Exec 实例失败: %s", err)
	}

	// 附加到 Exec 进程
	respRunAttach, err := sandBox.cli.ContainerExecAttach(sandBox.ctx, respRunID.ID, container.ExecStartOptions{})
	if err != nil {
		sandBox.Error = fmt.Errorf("附加到 Exec 实例失败: %s", err)
	}
	defer respRunAttach.CloseWrite() // 确保连接关闭

	// 模拟输入
	io.WriteString(respRunAttach.Conn, input)

	// 分离输出流
	var stdout, stderr bytes.Buffer
	_, err = stdcopy.StdCopy(&stdout, &stderr, respRunAttach.Reader)
	if err != nil {
		sandBox.Error = fmt.Errorf("读取输出流失败: %s", err)
	}

	// 检查退出状态码
	execInspect, err := sandBox.cli.ContainerExecInspect(sandBox.ctx, respRunID.ID)
	if err != nil {
		sandBox.Error = fmt.Errorf("检查 Exec 状态失败: %s", err)
	}
	if execInspect.ExitCode != 0 {
		sandBox.Error = fmt.Errorf("命令执行失败，退出码: %d，错误信息: %s", execInspect.ExitCode, stderr.String())
	}
	return stdout.String()
}

func getContainerMemory(cli *client.Client, ctx context.Context, containerId string) uint64 {
	// 获取容器信息
	stats, _ := cli.ContainerStats(ctx, containerId, false)
	defer stats.Body.Close()
	// 解码统计信息
	var statsData container.StatsResponse
	json.NewDecoder(stats.Body).Decode(&statsData)

	return statsData.MemoryStats.Usage
}

func memoryExceeded(memoryExceededChan chan bool, cli *client.Client, ctx context.Context, containerId string) {
	time.AfterFunc(3*time.Second, func() {
		memoryExceededChan <- true
	})
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	for range ticker.C {
		if getContainerMemory(cli, ctx, containerId) > 1024*1024*128 {
			fmt.Printf("内存超出...")
			memoryExceededChan <- true
			return
		}
	}
	memoryExceededChan <- false
	return
}

func (sandBox *SandBox) HandleExcuteCode(programPath, filePath, input string) (string, int64, error) {
	if sandBox.Error != nil {
		return "", -1, sandBox.Error
	}
	// 编写编译，运行命令行参数
	var cmd string
	if filePath == "" { // 在linux上
		cmd = programPath
	} else { // 在mac上
		excuteName := strings.Replace(filePath, ".cpp", "", 1)
		if _, err := os.Stat(excuteName); os.IsNotExist(err) {
			cmd = fmt.Sprintf("g++ %s -o %s && %s", filePath, excuteName, excuteName)
		} else {
			cmd = fmt.Sprintf("%s", excuteName)
		}
	}

	// 配置编译和运行命令
	//fmt.Printf("cmd: %s\n", cmd)
	execConfig := container.ExecOptions{
		Cmd:          []string{"bash", "-c", cmd},
		Tty:          true,
		AttachStdout: true,
		AttachStderr: true,
		AttachStdin:  true,
	}
	// 时间监控
	// 内存监控
	stdOutChan := make(chan string)
	memoryExceededChan := make(chan bool)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// 执行代码的子协程
	startTime := time.Now().UnixNano() / int64(time.Millisecond)
	go func(ctx context.Context, stdOutChan chan string) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				// 执行代码
				stdOut := sandBox.excuteCode(execConfig, input)
				stdOutChan <- stdOut
				return
			}
		}
	}(ctx, stdOutChan)
	// 监控内存的子协程
	go func(ctx context.Context, memoryExceededChan chan bool) {
		memoryExceeded(memoryExceededChan, sandBox.cli, sandBox.ctx, sandBox.container.container.ID)
	}(ctx, memoryExceededChan)
	var stdOut string
	select {
	case <-time.After(time.Second * 3):
		cancel()
		return "", -1, errors.New(response.Timeout)
	case stdOut = <-stdOutChan:
		endTime := time.Now().UnixNano() / int64(time.Millisecond)
		sandBox.UsedTime = endTime - startTime
		break
	case isMemoryExceeded := <-memoryExceededChan:
		if isMemoryExceeded == true {
			return "", -1, errors.New(response.MemoryExceeded)
		}
	}

	// 整理输出结果
	// 如果有错误，直接返回
	if sandBox.Error != nil {
		return "", -1, sandBox.Error
	}
	// TODO: 可能有问题
	sandBox.ExcuteResult = strings.ReplaceAll(stdOut, "\r", "")
	sandBox.ExcuteResult = strings.ReplaceAll(sandBox.ExcuteResult, input, "")
	sandBox.ExcuteResult = strings.ReplaceAll(sandBox.ExcuteResult, "\n", "")
	return sandBox.ExcuteResult, sandBox.UsedTime, sandBox.Error
}
