package service

import (
	"fmt"
	"github.com/polaris/codesandbox/api/request"
	"github.com/polaris/codesandbox/api/response"
	"github.com/polaris/codesandbox/logger"
	"github.com/polaris/codesandbox/model"
	"github.com/polaris/codesandbox/utils/_string"
	"github.com/polaris/codesandbox/utils/json"
	"strconv"
	"strings"
)

const (
	WrongAnswer = iota
	AC
	Timeout
	PresentationWrong
)

type userOutput struct {
	Output   string
	usedTime int64
}

type judgeInfo struct {
	testCaseInput      string
	testCaseOutput     string
	testCaseUserOutput string
	time               int64
}

func (s *Service) SubmitCppCode(outChan chan response.WebSocketResponse, submitCodeRequest *request.ProblemSubmit) {
	var programPath, filePath string
	if programPath, filePath = prepareCppCode(submitCodeRequest.Code, s.Conn, outChan); programPath == "" {
		return
	}
	outChan <- response.NewPendingResponse(s.Conn, response.SubmitAct)
	// 如果有用户在占用docker代码沙箱，那么状态一直都是Pending
	// 进行加锁
	s.Mutex.Lock()
	// 解锁
	defer s.Mutex.Unlock()
	// 执行问题下的所有测试用例数据并同步到polaris-oj系统中
	outChan <- response.NewRunningResponse(s.Conn, response.SubmitAct)
	// 获取所有测试用例
	var testCases []model.JudgeCase
	judgeConfig := new(model.JudgeConfig)
	var err error
	if testCases, judgeConfig, err = getTestCasesAndJudgeConfig(submitCodeRequest.QuestionId); err != nil {
		outChan <- response.NewSystemErrorResponse(s.Conn, response.TestCasesError)
		return
	}
	// 执行所有的测试用例
	var userOutputs []userOutput
	if userOutputs, err = runAllTestCases(testCases, programPath, filePath); err != nil {
		logger.GetLogger().Error("代码沙箱执行出错: " + err.Error())
		switch err.Error() {
		case response.Timeout:
			outChan <- response.NewTimeoutResponse(s.Conn, response.SubmitAct, "", "", "")
			return
		case response.MemoryExceeded:
			outChan <- response.NewMemoryExceededResponse(s.Conn, response.SubmitAct)
			return
		default:
			outChan <- response.NewSystemErrorResponse(s.Conn, response.SystemError)
			return
		}
	}
	// 对比所有输出用例，给出题目的判题结果
	_judgeInfo := new(judgeInfo)
	switch checkOutputs(testCases, userOutputs, judgeConfig, _judgeInfo) {
	case AC:
		outChan <- response.NewAcceptResponse(s.Conn, response.SubmitAct, _judgeInfo.time)
	case Timeout:
		outChan <- response.NewTimeoutResponse(s.Conn, response.SubmitAct, _judgeInfo.testCaseInput, _judgeInfo.testCaseOutput, _judgeInfo.testCaseUserOutput)
	case WrongAnswer:
		outChan <- response.NewWrongAnswerResponse(s.Conn, response.SubmitAct, _judgeInfo.testCaseInput, _judgeInfo.testCaseOutput, _judgeInfo.testCaseUserOutput)
	case PresentationWrong:
		outChan <- response.NewPresentationErrorResponse(s.Conn, response.SubmitAct, _judgeInfo.testCaseInput, _judgeInfo.testCaseOutput, _judgeInfo.testCaseUserOutput)
	}
	// 同步判题结果的数据到polaris-oj

}

func checkOutputs(cases []model.JudgeCase, outputs []userOutput, config *model.JudgeConfig, judgeInfo *judgeInfo) int {
	// 默认这两个数组大小一定是一样的
	var totalUsedTime int64 = 0
	for i, _case := range cases {
		// TODO: 这里没有完成格式错误的检查，要等在c++层面完成时间内存的检测的时候完成
		_string.GetCorrectString(&_case.Output)
		fmt.Printf("标准答案: %s\n", strconv.Quote(_case.Output))
		fmt.Printf("用户答案: %s\n", strconv.Quote(outputs[i].Output))
		// TODO: 用户答案和标准答案的校验还是有问题，待完善，主要是针对于PresentationWrong的判断无法实现
		if _case.Output != outputs[i].Output {
			judgeInfo.time = outputs[i].usedTime
			judgeInfo.testCaseInput = _case.Input
			judgeInfo.testCaseOutput = _case.Output
			judgeInfo.testCaseUserOutput = outputs[i].Output

			_string.EndWithBr(&outputs[i].Output)
			if _case.Output == outputs[i].Output {
				return PresentationWrong
			}
			// 换行前面如果有一个空格，这个是允许通过的，具体原因参考acwing上的快排
			_case.Output, _ = strings.CutSuffix(_case.Output, "\n")
			_case.Output += " \n"
			if _case.Output != outputs[i].Output {
				return WrongAnswer
			}
		}
		if outputs[i].usedTime > (config.TimeLimit * 1000) {
			judgeInfo.time = outputs[i].usedTime
			judgeInfo.testCaseInput = _case.Input
			judgeInfo.testCaseOutput = _case.Output
			judgeInfo.testCaseUserOutput = outputs[i].Output
			return Timeout
		}
		totalUsedTime += outputs[i].usedTime
	}
	judgeInfo.time = totalUsedTime
	return AC
}

func runAllTestCases(testCases []model.JudgeCase, programPath, filePath string) ([]userOutput, error) {
	userOutputs := make([]userOutput, 0)
	for _, testCase := range testCases {
		// 确认输入是以\n结尾
		_string.GetCorrectString(&testCase.Input)
		fmt.Printf("testCase.Input: %s\n", testCase.Input)
		res, usedTime, err := excuteRunCppCode(testCase.Input, programPath, filePath)
		// 只要出现错误，那么立即结束
		if err != nil {
			fmt.Printf("runAllTestCases执行错误: %s\n", err.Error())
			return nil, err
		}
		userOutputs = append(userOutputs, userOutput{res, usedTime})
	}
	return userOutputs, nil
}

func getTestCasesAndJudgeConfig(questionId string) ([]model.JudgeCase, *model.JudgeConfig, error) {
	question := &model.Question{}
	model.MysqlDB.Model(question).Where("identity = ?", questionId).Find(question)

	judgeCases := new([]model.JudgeCase)
	if err := json.JsonToModel(question.JudgeCase, judgeCases); err != nil {
		return nil, nil, err
	}
	judgeConfig := new(model.JudgeConfig)
	if err := json.JsonToModel(question.JudgeConfig, judgeConfig); err != nil {
		return nil, nil, err
	}
	return *judgeCases, judgeConfig, nil
}
