package service

import (
	"github.com/polaris/codesandbox/api"
	"github.com/polaris/codesandbox/api/request"
	"github.com/polaris/codesandbox/api/response"
	"github.com/polaris/codesandbox/settings"
	"github.com/polaris/codesandbox/utils/_string"
	"os"
)

func prepareCppCode(code string, outChan chan response.WebSocketResponse) error {
	// 下载代码到文件中
	os.MkdirAll(settings.EnvConfig.RootPath+CppDirPath, os.ModePerm)
	err := os.WriteFile(settings.EnvConfig.RootPath+CppCodePath, []byte(code), os.ModePerm)
	if err != nil {
		outChan <- response.NewSystemErrorResponse(api.SystemError)
		return err
	}
	return nil
}

func (s *Service) RunCppCode(runCodeRequest *request.ProblemRun) {
	if err := prepareCppCode(runCodeRequest.Code, s.ResponseChan); err != nil {
		return
	}
	s.ResponseChan <- response.NewPendingResponse(api.RunAct)
	// 如果有用户在占用docker代码沙箱，那么状态一直都是Pending
	// 进行加锁
	s.Mutex.Lock()
	// 解锁
	defer s.Mutex.Unlock()
	// 执行单次c++代码
	s.ResponseChan <- response.NewRunningResponse(api.RunAct)
	// 确保输入以`\n`结尾
	_string.GetCorrectString(&runCodeRequest.Input)
	// 将数据同步到docker
	setInputCasesToDocker(Cpp, runCodeRequest.Input)
	// 开始执行操作
	excuteRunCppCode(s.ResponseChan, api.RunAct)
}
