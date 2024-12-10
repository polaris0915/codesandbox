package service

import (
	"github.com/polaris/codesandbox/api"
	"github.com/polaris/codesandbox/api/request"
	"github.com/polaris/codesandbox/api/response"
	"github.com/polaris/codesandbox/model"
	"github.com/polaris/codesandbox/utils/json"
)

type userOutput struct {
	Output   string
	usedTime int64
}

func (s *Service) SubmitCppCode(submitCodeRequest *request.ProblemSubmit, syncBody *response.SyncBody) {
	if err := prepareCppCode(submitCodeRequest.Code, s.ResponseChan); err != nil {
		syncBody.JudgeInfo.JudgeResult = api.SystemError
		return
	}
	s.ResponseChan <- response.NewPendingResponse(api.SubmitAct)
	// 获取所有测试用例以及题目配置
	var testCases []model.JudgeCase
	judgeConfig := new(model.JudgeConfig)
	var err error
	if testCases, judgeConfig, err = getTestCasesAndJudgeConfig(submitCodeRequest.QuestionId); err != nil {
		s.ResponseChan <- response.NewSystemErrorResponse(api.TestCasesError)
		syncBody.JudgeInfo.JudgeResult = api.TestCasesError
	}
	var inputs []string
	for _, testCase := range testCases {
		// TODO: 需要规范创建问题时测试用例的输入
		//testCase.Input = strings.ReplaceAll(testCase.Input, "\\n", "\n")
		inputs = append(inputs, testCase.Input)
	}
	// 如果有用户在占用docker代码沙箱，那么状态一直都是Pending
	// 进行加锁
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	// 执行问题下的所有测试用例数据并同步到polaris-oj系统中
	s.ResponseChan <- response.NewRunningResponse(api.SubmitAct)
	// 将数据同步到docker
	setInputCasesToDocker(Cpp, inputs)
	// 执行所有的测试用例
	excuteSubmitCppCode(s.ResponseChan, api.SubmitAct, testCases, judgeConfig, syncBody)
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
