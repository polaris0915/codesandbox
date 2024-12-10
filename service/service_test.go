package service

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/polaris/codesandbox/settings"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestWriteFile(t *testing.T) {
	settings.Init()
	data, _ := os.ReadFile(settings.EnvConfig.RootPath + CppInputCasesPath)
	//fmt.Printf("%s\n", strconv.Quote(string(data)))
	inputCases := strings.Split(string(data), "-=-=-=-=\n")
	inputCases = inputCases[:len(inputCases)-1]
	//fmt.Printf("length = %d\n", len(inputCases))
	for _, v := range inputCases {
		fmt.Printf("%s\n", strconv.Quote(v))
	}
}

func TestInputCases(t *testing.T) {
	settings.Init()
	language := Cpp
	var inputs interface{} = []string{"5\n8 2 3 5 1\n", "8\n10 2 9 4 18 2 1 20\n"}

	var filePath string
	var dirPath string
	switch language {
	case Cpp:
		filePath = settings.EnvConfig.RootPath + CppInputCasesPath
		dirPath = settings.EnvConfig.RootPath + CppInputCasesPath
	}

	os.MkdirAll(dirPath, os.ModePerm)
	fd, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		t.Fatal(err)
	}
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

func TestDocker(t *testing.T) {
	settings.Init()
	// 创建 Docker 客户端
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		//logger.GetLogger().Errorf("无法创建docker客户端连接: %s", err.Error())
		//return err
		t.Fatal(err)
	}
	ctx := context.Background()

	// 获取容器详细信息
	cppContainer, err := cli.ContainerInspect(ctx, settings.ContainersConfig.Cpp.ContainerId)
	if err != nil {
		t.Fatalf("Error inspecting container: %v", err)
	}

	if cppContainer.State.Running {
		t.Log("cppContainer is running")
	}
	if err := cli.ContainerStart(ctx, settings.ContainersConfig.Cpp.ContainerId, container.StartOptions{}); err != nil {
		t.Fatalf("Error starting cppContainer: %v", err)
	}

	// 创建 exec 实例，配置命令和流选项
	execConfig := container.ExecOptions{
		Cmd:          []string{"go", "run", "/user_code/cpp/main.go"},
		AttachStdout: true, // 捕获标准输出
		AttachStderr: true, // 捕获标准错误
		Env:          []string{"PATH=/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"},
	}

	execResp, err := cli.ContainerExecCreate(ctx, settings.ContainersConfig.Cpp.ContainerId, execConfig)
	if err != nil {
		t.Fatalf("Error creating exec: %v", err)
	}

	// 附加到 exec 的输入/输出流
	resp, err := cli.ContainerExecAttach(ctx, execResp.ID, container.ExecStartOptions{})
	if err != nil {
		t.Fatalf("Error attaching to exec: %v", err)
	}
	defer resp.Close()

	// 读取输出流
	output, err := io.ReadAll(resp.Reader)
	if err != nil {
		t.Fatalf("Error reading exec output: %v", err)
	}
	t.Logf("%s\n=================\n", string(output))
	t.Logf("%s\n", strconv.Quote(string(output)))

	rawOutputs := strings.Split(string(output), "@@RECORD@@\n")
	rawOutputs = rawOutputs[1:]
	var outputs [][]string
	for _, _output := range rawOutputs {

		outputs = append(outputs, strings.Split(_output, "@@OUTPUT@@\n"))
	}
	t.Logf("%v\n", len(outputs))
	t.Logf("%v\n", outputs)
}
