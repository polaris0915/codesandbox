package act

import (
	"bytes"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/polaris/codesandbox/api"
	"github.com/polaris/codesandbox/api/request"
	"github.com/polaris/codesandbox/api/response"
	"github.com/polaris/codesandbox/logger"
	"github.com/polaris/codesandbox/service"
	"github.com/polaris/codesandbox/settings"
	pJson "github.com/polaris/codesandbox/utils/json"
	"io"
	"net/http"
)

func syncWithOther(s *response.SyncBody) error {
	client := &http.Client{}
	data, _ := pJson.RawModelToJson(s)
	req, err := http.NewRequest(settings.RemoteConfig.Method, settings.RemoteConfig.Url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Authorization", "Bearer some-token")
	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 读取响应体
	type BackResp struct {
		OK bool `json:"ok"`
	}
	backResp := &BackResp{}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	_ = pJson.JsonToModel(string(body), &backResp)
	if !backResp.OK {
		return fmt.Errorf("同步失败")
	}
	return nil
}

func HandleSubmitAct(wsConn *WsConnection, message []byte) {
	submitCodeRequest := new(request.ProblemSubmit)
	_ = json.Unmarshal(message, submitCodeRequest)
	if submitCodeRequest.QuestionId == "" {
		wsConn.OutChan <- response.NewSystemErrorResponse(api.ParamsError)
		return
	}
	switch submitCodeRequest.Language {
	case CPP:
		syncBody := &response.SyncBody{}
		service.NewService(request.SubmitAct, &mu, wsConn.OutChan).SubmitCppCode(submitCodeRequest, syncBody)
		if wsConn.User != nil { // 同步到远程平台
			syncBody.UserId = wsConn.User.Identity
			syncBody.UserAccount = wsConn.User.UserAccount
			syncBody.QuestionId = submitCodeRequest.QuestionId
			syncBody.Language = submitCodeRequest.Language
			syncBody.SubmitCode = submitCodeRequest.Code
			if err := syncWithOther(syncBody); err != nil {
				logger.GetLogger().Errorf(`同步%s用户, 问题Identity: "%s"的提交记录失败`, syncBody.UserId, syncBody.QuestionId)
			}
		}
		break
	default:
		wsConn.OutChan <- response.NewSystemErrorResponse(api.NoLanguageError)
	}

}
