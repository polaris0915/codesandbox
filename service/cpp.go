package service

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/polaris/codesandbox/api/response"
	"os"
)

func prepareCppCode(code string, conn *websocket.Conn, outChan chan response.WebSocketResponse) (string, string) {
	// 下载代码到文件中
	os.MkdirAll("./user_code", os.ModePerm)
	fileName := uuid.New().String()
	programPath := fmt.Sprintf("./user_code/%s", fileName)
	filePath := fmt.Sprintf("./user_code/%s.cpp", fileName)
	err := os.WriteFile(filePath, []byte(code), os.ModePerm)
	if err != nil {
		outChan <- response.NewSystemErrorResponse(conn, response.SystemError)
		return "", ""
	}
	//// 编译c++代码 TODO: 上线服务器打开这个代码
	//cmd := exec.Command("g++", filePath, "-o", programPath)
	//output, err := cmd.CombinedOutput() // 同时捕获标准输出和错误输出
	//if err != nil {
	//	response.NewWebSocketResponse(s.Conn, s.Activity, response.CompileError).CompileErrorResponse(_string(output))
	//	return
	//}
	return programPath, filePath
}
