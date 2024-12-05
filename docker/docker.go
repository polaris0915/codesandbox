package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"io"
	"os"
)

// CreateContainer 创建运行c++代码所需要的容器
func CreateContainer() {
	// 创建上下文
	ctx := context.Background()

	// 创建 Docker 客户端
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := cli.Close(); err != nil {
			panic(err)
		}
	}()

	// 拉取基础镜像（如果本地没有）
	// TODO: 增加拉取镜像逻辑，如果已经存在镜像了，就不用再去拉取了
	out, err := cli.ImagePull(ctx, "ubuntu:20.04", image.PullOptions{})
	if err != nil {
		fmt.Println("拉取镜像失败: ", err)
	}
	defer out.Close()
	io.Copy(os.Stdout, out)

	// 创建容器
	containerConfig := &container.Config{
		Image: "ubuntu:20.04",
		//Cmd:   []string{"bash", "-c", "while true; do sleep 1; done"},
		Tty: true,
	}

	os.MkdirAll("./user_code", os.ModePerm)

	hostConfig := &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type: mount.TypeBind,
				// TODO: 需要改写成配置文件
				Source: "/Users/alin-youlinlin/Desktop/polaris-all_projects/codesandbox/user_code", // 挂载本地目录
				Target: "/user_code",
			},
		},
		Resources: container.Resources{
			Memory: 1024 * 1024 * 258,
		},
	}

	containerResp, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, "cpp-codesandbox")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Container ID: %s\n", containerResp.ID)

	// 启动容器
	if err := cli.ContainerStart(ctx, containerResp.ID, container.StartOptions{}); err != nil {
		panic(err)
	}

	fmt.Println("Container started successfully!")

	// 安装编译环境的命令
	execConfig := container.ExecOptions{
		Cmd:          []string{"bash", "-c", "apt-get update && apt-get install -y build-essential g++ vim"},
		Tty:          true,
		AttachStdout: true,
		AttachStderr: true,
	}

	respID, err := cli.ContainerExecCreate(ctx, containerResp.ID, execConfig)
	if err != nil {
		panic(err)
	}

	respAttach, err := cli.ContainerExecAttach(ctx, respID.ID, container.ExecStartOptions{})
	if err != nil {
		panic(err)
	}
	defer respAttach.Close()
	io.Copy(os.Stdout, respAttach.Reader)

	fmt.Println("C++ environment is ready!")
}
