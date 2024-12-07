package service

import (
	"fmt"
	"github.com/polaris/codesandbox/api/request"
	"github.com/polaris/codesandbox/api/response"
	"github.com/polaris/codesandbox/logger"
	"github.com/polaris/codesandbox/utils/_string"
)

func (s *Service) RunCppCode(outChan chan response.WebSocketResponse, runCodeRequest *request.ProblemRun) {
	//// 下载代码到文件中
	//os.MkdirAll("./user_code", os.ModePerm)
	//fileName := uuid.New().String()
	//programPath := fmt.Sprintf("./user_code/%s", fileName)
	//filePath := fmt.Sprintf("./user_code/%s.cpp", fileName)
	//err := os.WriteFile(filePath, []byte(runCodeRequest.Code), os.ModePerm)
	//if err != nil {
	//	outChan <- response.NewSystemErrorResponse(s.Conn, response.SystemError)
	//	return
	//}
	//// 编译c++代码 TODO: 上线服务器打开这个代码
	//cmd := exec.Command("g++", filePath, "-o", programPath)
	//output, err := cmd.CombinedOutput() // 同时捕获标准输出和错误输出
	//if err != nil {
	//	response.NewWebSocketResponse(s.Conn, s.Activity, response.CompileError).CompileErrorResponse(_string(output))
	//	return
	//}
	var programPath, filePath string
	if programPath, filePath = prepareCppCode(runCodeRequest.Code, outChan); programPath == "" {
		return
	}
	outChan <- response.NewPendingResponse(response.RunAct)
	// 如果有用户在占用docker代码沙箱，那么状态一直都是Pending
	// 进行加锁
	s.Mutex.Lock()
	// 解锁
	defer s.Mutex.Unlock()
	// 执行单次c++代码
	outChan <- response.NewRunningResponse(response.RunAct)
	// 确保输入以`\n`结尾
	_string.GetCorrectString(&runCodeRequest.Input)
	res, usedTime, err := excuteRunCppCode(runCodeRequest.Input, programPath, filePath)
	fmt.Printf("执行c++代码用时: %dms\n", usedTime)
	if err != nil {
		logger.GetLogger().Error("代码沙箱执行出错: " + err.Error())
		switch err.Error() {
		case response.Timeout:
			outChan <- response.NewTimeoutResponse(response.RunAct, runCodeRequest.Input, "", "")
			return
		case response.MemoryExceeded:
			outChan <- response.NewMemoryExceededResponse(response.RunAct)
			return
		default:
			outChan <- response.NewSystemErrorResponse(response.SystemError)
			return
		}
	}
	outChan <- response.NewFinishedResponse(response.RunAct, "", res)
}
