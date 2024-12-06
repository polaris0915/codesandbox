package main

import (
	"github.com/polaris/codesandbox/logger"
	"github.com/polaris/codesandbox/model"
	"github.com/polaris/codesandbox/router"
	"github.com/polaris/codesandbox/settings"
	"log"
	"net/http"
)

// TODO: 1. 编写项目启动时获取g++环境的程序
// TODO: 2. 如果用户长时间没有发送心跳包，那么就关闭conn
// TODO: 3. func (sandBox *SandBox) HandleExcuteCode(programPath, input _string) (_string, error)最后筛选答案可能有问题
// TODO: 4. 在项目启动要确保容器已经存在并且正在运行
// TODO: 5. 用户的输入都应该以"\n"结尾
// TODO: 6. 统计每个算法题所用的时间目前是不准确的，这个时间应该要在c++代码层面去实现精确的时间统计。目前只是以docker运行到结束的时间为准
// TODO: 7. 编译c++程序出现编译错误无法正确捕捉，现在只能设置成系统错误
//       考虑将6，7在c++层面去实现
// TODO: 8. 对于标准答案和用户程序输出的答案之间的校验没有完善，目前无法正确判断出PresentationWrong的情况

/*
一定注意！！！
如果自己手动删除了创建docker的时候的本地目录，会导致docker后面无法关联你原本指定的目录
就算自己手动创建一个也不行！！！
重新执行docker.CreateContainer()命令
*/
func main() {
	// 加载配置
	settings.Init()
	// 加载日志
	logger.Init("debug")
	// 加载所有需要的数据库
	model.InitAllDB()
	// 注册路由
	router.Init()
	// 启动服务
	log.Fatal(http.ListenAndServe(":8010", nil))
	//docker.CreateContainer()
}
