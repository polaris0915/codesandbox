package act

import (
	"fmt"
	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	"github.com/polaris/codesandbox/api/request"
	"github.com/polaris/codesandbox/api/response"
	"log"
	"testing"
)

func TestHandleRunAct(t *testing.T) {
	// 连接到服务端 WebSocket
	serverURL := "ws://localhost:8010/api/act"
	conn, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		log.Fatal("Error connecting to server:", err)
	}
	defer conn.Close()

	// 发送文件内容作为消息
	err = conn.WriteJSON(&request.ProblemRun{
		Activity: "RUN_CODE_ACTIVITY",
		Code: `
#include <iostream>
using namespace std;

const int N = 1000010;

int n, a[N];

void quick_sort(int* a, int left, int right){
    if (left >= right) return;
    
    int pivot = a[(left + right) / 2], i = left - 1, j = right + 1;
    while (i < j){
        do i++; while (a[i] < pivot);
        do j--; while (a[j] > pivot);
        if (i < j) swap(a[i], a[j]);
    }
    
    quick_sort(a, left, j), quick_sort(a, j + 1, right);
}

int main(void){
    scanf("%d", &n);
    for (int i = 0; i < n; i++) scanf("%d", &a[i]);
    
    quick_sort(a, 0, n - 1);
    for (int i = 0; i < n; i++) printf("%d ", a[i]);
    printf("\n");
};

`,
		Language: "c++",
		Input:    "7\n5 2 19 4 20 5 7\n",
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

func TestHandleSubmitAct(t *testing.T) {
	// 连接到服务端 WebSocket
	serverURL := "ws://localhost:8010/api/act"
	conn, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		log.Fatal("Error connecting to server:", err)
	}
	defer conn.Close()

	// 发送文件内容作为消息
	err = conn.WriteJSON(&request.ProblemSubmit{
		Activity: request.SubmitAct,
		Code: `
#include <iostream>
using namespace std;

const int N = 1000010;

int n, a[N];

void quick_sort(int* a, int left, int right){
    if (left >= right) return;
    
    int pivot = a[(left + right) / 2], i = left - 1, j = right + 1;
    while (i < j){
        do i++; while (a[i] < pivot);
        do j--; while (a[j] > pivot);
        if (i < j) swap(a[i], a[j]);
    }
    
    quick_sort(a, left, j), quick_sort(a, j + 1, right);
}

int main(void){
    scanf("%d", &n);
    for (int i = 0; i < n; i++) scanf("%d", &a[i]);
    
    quick_sort(a, 0, n - 1);
    for (int i = 0; i < n; i++) printf("%d ", a[i]);
    printf("\n");
};

`,
		Language:   "c++",
		QuestionId: "1e9219d1-bcfe-49d8-a42a-7943c5d64da1",
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
			FinishedResponse := response.FinishedResponse{}
			if err := json.Unmarshal(message, &FinishedResponse); err != nil {
				log.Println("Error unmarshalling FinishedResponse:", err)
				continue
			}
			fmt.Printf("FinishedResponse: %+v\n", FinishedResponse)
			break
		case response.Pending:
			PendingResponse := response.PendingResponse{}
			if err := json.Unmarshal(message, &PendingResponse); err != nil {
				log.Println("Error unmarshalling PendingResponse:", err)
				continue
			}
			fmt.Printf("PendingResponse: %+v\n", PendingResponse)
			break
		case response.CompileError:
			CompileErrorResponse := response.CompileErrorResponse{}
			if err := json.Unmarshal(message, &CompileErrorResponse); err != nil {
				log.Println("Error unmarshalling CompileErrorResponse:", err)
				continue
			}
			fmt.Printf("CompileErrorResponse: %+v\n", CompileErrorResponse)
			break
		case response.Running:
			RunningResponse := response.RunningResponse{}
			if err := json.Unmarshal(message, &RunningResponse); err != nil {
				log.Println("Error unmarshalling RunningResponse:", err)
				continue
			}
			fmt.Printf("RunningResponse: %+v\n", RunningResponse)
			break
		case response.MemoryExceeded:
			MemoryExceededResponse := response.MemoryExceededResponse{}
			if err := json.Unmarshal(message, &MemoryExceededResponse); err != nil {
				log.Println("Error unmarshalling MemoryExceededResponse:", err)
				continue
			}
			fmt.Printf("MemoryExceededResponse: %+v\n", MemoryExceededResponse)
			break
		case response.Timeout:
			TimeoutResponse := response.TimeoutResponse{}
			if err := json.Unmarshal(message, &TimeoutResponse); err != nil {
				log.Println("Error unmarshalling TimeoutResponse:", err)
				continue
			}
			fmt.Printf("TimeoutResponse: %+v\n", TimeoutResponse)
			break
		case response.SystemError:
			SystemErrorResponse := response.SystemErrorResponse{}
			if err := json.Unmarshal(message, &SystemErrorResponse); err != nil {
				log.Println("Error unmarshalling SystemErrorResponse:", err)
				continue
			}
			fmt.Printf("SystemErrorResponse: %+v\n", SystemErrorResponse)
			break
		case response.Accepted:
			AcceptedResponse := response.AcceptResponse{}
			if err := json.Unmarshal(message, &AcceptedResponse); err != nil {
				log.Println("Error unmarshalling SystemErrorResponse:", err)
				continue
			}
			fmt.Printf("AcceptedResponse: %+v\n", AcceptedResponse)
			break
		case response.WrongAnswer:
			WrongAnswerResponse := response.WrongAnswerResponse{}
			if err := json.Unmarshal(message, &WrongAnswerResponse); err != nil {
				log.Println("Error unmarshalling SystemErrorResponse:", err)
				continue
			}
			fmt.Printf("WrongAnswerResponse: %+v\n", WrongAnswerResponse)
			break
		case response.PresentationError:
			PresentationErrorResponse := response.PresentationErrorResponse{}
			if err := json.Unmarshal(message, &PresentationErrorResponse); err != nil {
				log.Println("Error unmarshalling SystemErrorResponse:", err)
				continue
			}
			fmt.Printf("PresentationErrorResponse: %+v\n", PresentationErrorResponse)
			break
		}
	}

}
