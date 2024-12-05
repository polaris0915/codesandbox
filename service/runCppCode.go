package service

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/polaris/codesandbox/api/request"
	"github.com/polaris/codesandbox/api/response"
	"github.com/polaris/codesandbox/logger"
	"os"
)

func (s *Service) RunCppCode(outChan chan response.WebSocketResponse, runCodeRequest *request.ProblemRun) {
	// 下载代码到文件中
	os.MkdirAll("./user_code", os.ModePerm)
	fileName := uuid.New().String()
	programPath := fmt.Sprintf("./user_code/%s", fileName)
	filePath := fmt.Sprintf("./user_code/%s.cpp", fileName)
	err := os.WriteFile(filePath, []byte(runCodeRequest.Code), os.ModePerm)
	if err != nil {
		outChan <- response.NewSystemErrorResponse(s.Conn, response.SystemError)
		return
	}
	//// 编译c++代码 TODO: 上线服务器打开这个代码
	//cmd := exec.Command("g++", filePath, "-o", programPath)
	//output, err := cmd.CombinedOutput() // 同时捕获标准输出和错误输出
	//if err != nil {
	//	response.NewWebSocketResponse(s.Conn, s.Activity, response.CompileError).CompileErrorResponse(string(output))
	//	return
	//}
	outChan <- response.NewPendingResponse(s.Conn, response.RunAct)
	// 如果有用户在占用docker代码沙箱，那么状态一直都是Pending
	// 进行加锁
	s.Mutex.Lock()
	// 解锁
	defer s.Mutex.Unlock()
	// 执行c++代码
	outChan <- response.NewRunningResponse(s.Conn, response.RunAct)
	res, usedTime, err := excuteRunCppCode(runCodeRequest.Input, programPath, filePath)
	fmt.Printf("执行c++代码用时: %dms\n", usedTime)
	if err != nil {
		logger.GetLogger().Error("代码沙箱执行出错: " + err.Error())
	}
	if err != nil {
		switch err.Error() {
		case response.Timeout:
			outChan <- response.NewTimeoutResponse(s.Conn, response.RunAct)
			return
		case response.MemoryExceeded:
			outChan <- response.NewMemoryExceededResponse(s.Conn, response.RunAct)
			return
		default:
			outChan <- response.NewSystemErrorResponse(s.Conn, response.SystemError)
			return
		}
	}
	outChan <- response.NewFinishedResponse(s.Conn, response.RunAct, "", res)
}
