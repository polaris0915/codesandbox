package service

import (
	"github.com/gorilla/websocket"
	"github.com/polaris/codesandbox/docker"
	"sync"
)

type Service struct {
	Conn     *websocket.Conn
	Activity string
	Mutex    *sync.Mutex
}

func NewService(conn *websocket.Conn, activity string, mutex *sync.Mutex) *Service {
	return &Service{Conn: conn, Activity: activity, Mutex: mutex}
}

/*
设置成包内的全局函数，用锁类控制每个语言的沙箱同一时间只有一个用户可以执行代码
// TODO: 有待验证如果是属于类的对象函数加锁跟这种现在写的全局函数加锁有什么区别
*/
func excuteRunCppCode(input, programPath, filePath string) (string, int64, error) {
	return docker.NewSandBox().HandleExcuteCode(programPath, filePath, input)
}
