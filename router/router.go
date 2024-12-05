package router

import (
	"github.com/polaris/codesandbox/api/api"
	"net/http"
)

func Init() {
	// 注册路由
	//http.HandleFunc("/run-code", api.HandleRunCode)
	http.HandleFunc("/run-code", api.WsHandler)

}
