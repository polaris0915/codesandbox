package model

import (
	"github.com/polaris/codesandbox/settings"
	"github.com/polaris/codesandbox/utils/json"
	"testing"
)

func TestModel(t *testing.T) {

	settings.Init()
	InitAllDB()

	question := &Question{}
	err := MysqlDB.Model(question).Where("identity = ?", "4230a762-74db-48dc-b494-15004ddc2de5").Find(question).Error
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("question: %+v\n", question)
	judgeCases := new([]JudgeCase)
	if err := json.JsonToModel(question.JudgeCase, judgeCases); err != nil {
		t.Fatal(err)
	}
	judgeConfig := new(JudgeConfig)
	if err := json.JsonToModel(question.JudgeConfig, judgeConfig); err != nil {
		t.Fatal(err)
	}
	t.Logf("judgeCases: %+v\n", judgeCases)
	t.Logf("judgeConfig: %+v\n", judgeConfig)
}
