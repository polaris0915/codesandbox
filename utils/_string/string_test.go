package _string

import (
	"fmt"
	"github.com/polaris/codesandbox/model"
	"github.com/polaris/codesandbox/settings"
	"github.com/polaris/codesandbox/utils/json"
	"strconv"
	"testing"
)

func TestEndWithBr(t *testing.T) {
	str1 := "1 2"
	str2 := "1 2 2\n"
	strStruct1 := struct {
		Str string
	}{
		Str: str1,
	}
	strStruct2 := struct {
		Str string
	}{
		Str: str2,
	}
	fmt.Printf("str1: %s, str2: %s\n", strconv.Quote(str1), strconv.Quote(str2))
	EndWithBr(&str1)
	EndWithBr(&str2)
	fmt.Printf("after EndWithBar:\nstr1: %s, str2: %s\n", strconv.Quote(str1), strconv.Quote(str2))

	fmt.Printf("strStruct1: %s, strStruct2: %s\n", strconv.Quote(strStruct1.Str), strconv.Quote(strStruct2.Str))
	EndWithBr(&strStruct1.Str)
	EndWithBr(&strStruct2.Str)
	fmt.Printf("after EndWithBar:\nstrStruct2: %s, strStruct2: %s\n", strconv.Quote(strStruct1.Str), strconv.Quote(strStruct2.Str))

}

func TestEndWithoutBr(t *testing.T) {
	str1 := "1 2 2\n"
	str2 := "1 2 2\n\n"
	fmt.Printf("str1: %s, str2: %s\n", strconv.Quote(str1), strconv.Quote(str2))
	EndWithoutBr(&str1)
	EndWithoutBr(&str2)
	fmt.Printf("after EndWithoutBr: str1: %s, str2: %s\n", strconv.Quote(str1), strconv.Quote(str2))

}

func TestGetCorrectTestCase(t *testing.T) {
	settings.Init()
	model.InitAllDB()

	question := &model.Question{}
	err := model.MysqlDB.Model(question).Where("identity = ?", "1e9219d1-bcfe-49d8-a42a-7943c5d64da1").Find(question).Error
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("question: %+v\n", question)
	var cases []model.JudgeCase
	json.JsonToModel(question.JudgeCase, &cases)
	fmt.Printf("input: %s\n", cases[0].Input)
	GetCorrectString(&cases[0].Input)
	fmt.Printf("after GetCorrectString:\ninput: %s", cases[0].Input)
	fmt.Printf("after GetCorrectString:\ninput: %s", strconv.Quote(cases[0].Input))
}
