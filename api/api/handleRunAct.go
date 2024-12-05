package api

import (
	"encoding/json"
	"github.com/polaris/codesandbox/api/request"
	"github.com/polaris/codesandbox/api/response"
	"github.com/polaris/codesandbox/service"
	"sync"
)

var (
	CPP    = "c++"
	Golang = "go"
)

var mu sync.Mutex

func HandleRunAct(wsConn *WsConnection, message []byte) {
	runCodeRequest := new(request.ProblemRun)
	_ = json.Unmarshal(message, runCodeRequest)
	switch runCodeRequest.Language {
	case CPP:
		// 实现RunCppCode，返回一个函数并执行
		service.NewService(wsConn.Conn, request.RunAct, &mu).RunCppCode(wsConn.OutChan, runCodeRequest)
		break
	default:
		wsConn.OutChan <- response.NewSystemErrorResponse(wsConn.Conn, response.NoLanguageError)
	}

}
