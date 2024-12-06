package api

import (
	"encoding/json"
	"github.com/polaris/codesandbox/api/request"
	"github.com/polaris/codesandbox/api/response"
	"github.com/polaris/codesandbox/service"
)

func HandleSubmitAct(wsConn *WsConnection, message []byte) {
	submitCodeRequest := new(request.ProblemSubmit)
	_ = json.Unmarshal(message, submitCodeRequest)
	if submitCodeRequest.QuestionId == "" {
		wsConn.OutChan <- response.NewSystemErrorResponse(wsConn.Conn, response.ParamsError)
		return
	}
	switch submitCodeRequest.Language {
	case CPP:
		service.NewService(wsConn.Conn, request.SubmitAct, &mu).SubmitCppCode(wsConn.OutChan, submitCodeRequest)
		break
	default:
		wsConn.OutChan <- response.NewSystemErrorResponse(wsConn.Conn, response.NoLanguageError)
	}

}
