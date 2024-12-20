package act

import (
	"encoding/json"
	"github.com/polaris/codesandbox/api"
	"github.com/polaris/codesandbox/api/request"
	"github.com/polaris/codesandbox/api/response"
	"github.com/polaris/codesandbox/service"
)

func HandleRunAct(wsConn *WsConnection, message []byte) {
	runCodeRequest := new(request.ProblemRun)
	_ = json.Unmarshal(message, runCodeRequest)
	switch runCodeRequest.Language {
	case CPP:
		service.NewService(request.RunAct, &mu, wsConn.OutChan).RunCppCode(runCodeRequest)
		break
	default:
		wsConn.OutChan <- response.NewSystemErrorResponse(api.NoLanguageError)
	}

}
