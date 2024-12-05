package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/polaris/codesandbox/api/request"
	"github.com/polaris/codesandbox/api/response"
	"log"
	"testing"
)

func TestHandleRunCode(t *testing.T) {
	// 连接到服务端 WebSocket
	serverURL := "ws://localhost:8010/run-code"
	conn, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		log.Fatal("Error connecting to server:", err)
	}
	defer conn.Close()

	// 发送文件内容作为消息
	err = conn.WriteJSON(&request.ProblemRun{
		Activity: request.RunAct,
		Code: `
#include <iostream>

int main(){
    int a, b;
    std::cin >> a >> b;
    std::cout << (a + b) << std::endl;
    return 0;
}
`,
		Language: CPP,
		Input:    "1 2\n",
	})
	if err != nil {
		log.Fatal("Error sending file:", err)
	}
	fmt.Println("File sent successfully to the server!")

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}
		// 动态解析响应
		problemStatus := new(response.BaseResponse)
		if err := json.Unmarshal(message, problemStatus); err != nil {
			log.Println("Error unmarshalling base response:", err)
			continue
		}

		switch problemStatus.Status {
		case response.Finished:
			FinishedResponse := struct {
				response.BaseResponse
				StdErr string `json:"stdErr"`
				StdOut string `json:"stdOut"`
			}{}
			if err := json.Unmarshal(message, &FinishedResponse); err != nil {
				log.Println("Error unmarshalling FinishedResponse:", err)
				continue
			}
			fmt.Printf("FinishedResponse: %+v\n", FinishedResponse)
			break
		case response.Pending:
			PendingResponse := struct {
				response.BaseResponse
			}{}
			if err := json.Unmarshal(message, &PendingResponse); err != nil {
				log.Println("Error unmarshalling PendingResponse:", err)
				continue
			}
			fmt.Printf("PendingResponse: %+v\n", PendingResponse)
			break
		case response.CompileError:
			CompileErrorResponse := struct {
				response.BaseResponse
				CompilationLog string `json:"compilationLog"`
			}{}
			if err := json.Unmarshal(message, &CompileErrorResponse); err != nil {
				log.Println("Error unmarshalling CompileErrorResponse:", err)
				continue
			}
			fmt.Printf("CompileErrorResponse: %+v\n", CompileErrorResponse)
			break
		case response.Running:
			RunningResponse := struct {
				response.BaseResponse
			}{}
			if err := json.Unmarshal(message, &RunningResponse); err != nil {
				log.Println("Error unmarshalling RunningResponse:", err)
				continue
			}
			fmt.Printf("RunningResponse: %+v\n", RunningResponse)
			break
		case response.MemoryExceeded:
			MemoryExceededResponse := struct {
				response.BaseResponse
			}{}
			if err := json.Unmarshal(message, &MemoryExceededResponse); err != nil {
				log.Println("Error unmarshalling MemoryExceededResponse:", err)
				continue
			}
			fmt.Printf("MemoryExceededResponse: %+v\n", MemoryExceededResponse)
			break
		case response.Timeout:
			TimeoutResponse := struct {
				response.BaseResponse
			}{}
			if err := json.Unmarshal(message, &TimeoutResponse); err != nil {
				log.Println("Error unmarshalling TimeoutResponse:", err)
				continue
			}
			fmt.Printf("TimeoutResponse: %+v\n", TimeoutResponse)
			break
		case response.SystemError:
			SystemErrorResponse := struct {
				response.BaseResponse
			}{}
			if err := json.Unmarshal(message, &SystemErrorResponse); err != nil {
				log.Println("Error unmarshalling SystemErrorResponse:", err)
				continue
			}
			fmt.Printf("SystemErrorResponse: %+v\n", SystemErrorResponse)
			break
		}
	}

}
