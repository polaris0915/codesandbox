package act

import (
	"fmt"
	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	"github.com/polaris/codesandbox/api"
	"github.com/polaris/codesandbox/api/request"
	"github.com/polaris/codesandbox/api/response"
	"log"
	"net/http"
	"sync"
	"testing"
)

func submit(index int) {
	// 连接到服务端 WebSocket
	serverURL := "ws://localhost:8010/api/act"
	conn, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		log.Fatal("Error connecting to server:", err)
	}
	defer conn.Close()

	//发送文件内容作为消息
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
}
	`,
		Language:   "c++",
		QuestionId: "1e9219d1-bcfe-49d8-a42a-7943c5d64da1",
	})

	if err != nil {
		log.Fatalf("Error sending file: %v", err)
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
		case api.Finished:
			FinishedResponse := response.FinishedResponse{}
			if err := json.Unmarshal(message, &FinishedResponse); err != nil {
				log.Println("Error unmarshalling FinishedResponse:", err)
				continue
			}
			fmt.Printf("FinishedResponse: %+v\n", FinishedResponse)
			break
		case api.Pending:
			PendingResponse := response.PendingResponse{}
			if err := json.Unmarshal(message, &PendingResponse); err != nil {
				log.Println("Error unmarshalling PendingResponse:", err)
				continue
			}
			fmt.Printf("PendingResponse: %+v\n", PendingResponse)
			break
		case api.CompileError:
			CompileErrorResponse := response.CompileErrorResponse{}
			if err := json.Unmarshal(message, &CompileErrorResponse); err != nil {
				log.Println("Error unmarshalling CompileErrorResponse:", err)
				continue
			}
			fmt.Printf("CompileErrorResponse: %+v\n", CompileErrorResponse)
			break
		case api.Running:
			RunningResponse := response.RunningResponse{}
			if err := json.Unmarshal(message, &RunningResponse); err != nil {
				log.Println("Error unmarshalling RunningResponse:", err)
				continue
			}
			fmt.Printf("RunningResponse: %+v\n", RunningResponse)
			break
		case api.MemoryExceeded:
			MemoryExceededResponse := response.MemoryExceededResponse{}
			if err := json.Unmarshal(message, &MemoryExceededResponse); err != nil {
				log.Println("Error unmarshalling MemoryExceededResponse:", err)
				continue
			}
			fmt.Printf("MemoryExceededResponse: %+v\n", MemoryExceededResponse)
			break
		case api.Timeout:
			TimeoutResponse := response.TimeoutResponse{}
			if err := json.Unmarshal(message, &TimeoutResponse); err != nil {
				log.Println("Error unmarshalling TimeoutResponse:", err)
				continue
			}
			fmt.Printf("TimeoutResponse: %+v\n", TimeoutResponse)
			break
		case api.SystemError:
			SystemErrorResponse := response.SystemErrorResponse{}
			if err := json.Unmarshal(message, &SystemErrorResponse); err != nil {
				log.Println("Error unmarshalling SystemErrorResponse:", err)
				continue
			}
			fmt.Printf("SystemErrorResponse: %+v\n", SystemErrorResponse)
			break
		case api.Accepted:
			AcceptedResponse := response.AcceptResponse{}
			if err := json.Unmarshal(message, &AcceptedResponse); err != nil {
				log.Println("Error unmarshalling SystemErrorResponse:", err)
				continue
			}
			fmt.Printf("AcceptedResponse: %+v\n", AcceptedResponse)
			break
		case api.WrongAnswer:
			WrongAnswerResponse := response.WrongAnswerResponse{}
			if err := json.Unmarshal(message, &WrongAnswerResponse); err != nil {
				log.Println("Error unmarshalling SystemErrorResponse:", err)
				continue
			}
			fmt.Printf("WrongAnswerResponse: %+v\n", WrongAnswerResponse)
			break
		case api.RunTimeError:
			WrongAnswerResponse := response.RunTimeErrorResponse{}
			if err := json.Unmarshal(message, &WrongAnswerResponse); err != nil {
				log.Println("Error unmarshalling SystemErrorResponse:", err)
				continue
			}
			fmt.Printf("WrongAnswerResponse: %+v\n", WrongAnswerResponse)
			break
		case api.PresentationError:
			PresentationErrorResponse := response.PresentationErrorResponse{}
			if err := json.Unmarshal(message, &PresentationErrorResponse); err != nil {
				log.Println("Error unmarshalling SystemErrorResponse:", err)
				continue
			}
			fmt.Printf("PresentationErrorResponse: %+v\n", PresentationErrorResponse)
			break
		}
	}
	fmt.Printf("协程%d完成了任务\n", index)
}

func TestHandleSubmitAct(t *testing.T) {
	// 使用 WaitGroup 同步并发 Goroutine
	var wg sync.WaitGroup
	numClients := 100

	for i := 0; i < numClients; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			submit(index)
		}(i)
	}
	wg.Wait()
	//submit(1)
}

func TestHandleRunAct(t *testing.T) {
	// 连接到服务端 WebSocket
	serverURL := "wss://www.acwing.com/wss/socket/"
	conn, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		log.Fatal("Error connecting to server:", err)
	}
	defer conn.Close()

	// 发送文件内容作为消息
	err = conn.WriteJSON(&request.ProblemRun{
		Activity: "problem_run_code",
		Code:     "",
		Language: "C++",
		Input:    "3 4",
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
		case api.Finished:
			FinishedResponse := response.FinishedResponse{}
			if err := json.Unmarshal(message, &FinishedResponse); err != nil {
				log.Println("Error unmarshalling FinishedResponse:", err)
				continue
			}
			fmt.Printf("FinishedResponse: %+v\n", FinishedResponse)
			break
		case api.Pending:
			PendingResponse := response.PendingResponse{}
			if err := json.Unmarshal(message, &PendingResponse); err != nil {
				log.Println("Error unmarshalling PendingResponse:", err)
				continue
			}
			fmt.Printf("PendingResponse: %+v\n", PendingResponse)
			break
		case api.CompileError:
			CompileErrorResponse := response.CompileErrorResponse{}
			if err := json.Unmarshal(message, &CompileErrorResponse); err != nil {
				log.Println("Error unmarshalling CompileErrorResponse:", err)
				continue
			}
			fmt.Printf("CompileErrorResponse: %+v\n", CompileErrorResponse)
			break
		case api.Running:
			RunningResponse := response.RunningResponse{}
			if err := json.Unmarshal(message, &RunningResponse); err != nil {
				log.Println("Error unmarshalling RunningResponse:", err)
				continue
			}
			fmt.Printf("RunningResponse: %+v\n", RunningResponse)
			break
		case api.MemoryExceeded:
			MemoryExceededResponse := response.MemoryExceededResponse{}
			if err := json.Unmarshal(message, &MemoryExceededResponse); err != nil {
				log.Println("Error unmarshalling MemoryExceededResponse:", err)
				continue
			}
			fmt.Printf("MemoryExceededResponse: %+v\n", MemoryExceededResponse)
			break
		case api.Timeout:
			TimeoutResponse := response.TimeoutResponse{}
			if err := json.Unmarshal(message, &TimeoutResponse); err != nil {
				log.Println("Error unmarshalling TimeoutResponse:", err)
				continue
			}
			fmt.Printf("TimeoutResponse: %+v\n", TimeoutResponse)
			break
		case api.RunTimeError:
			TimeoutResponse := response.RunTimeErrorResponse{}
			if err := json.Unmarshal(message, &TimeoutResponse); err != nil {
				log.Println("Error unmarshalling TimeoutResponse:", err)
				continue
			}
			fmt.Printf("TimeoutResponse: %+v\n", TimeoutResponse)
			break
		case api.SystemError:
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

func TestSpider(t *testing.T) {
	// 连接到服务端 WebSocket
	serverURL := "wss://www.acwing.com/wss/socket/"

	// 设置请求头
	headers := http.Header{}
	headers.Add("Cookie", "csrftoken=pgf5e3Pr47tfxV9WqEnUMGjkD0vxRlaTBWguW9VWODEj9LDH8HzbnKJ9k9GKcwsq; sessionid=nx2mjx2tqy0icok309fpgf60467g1n8f")
	headers.Add("Host", "www.acwing.com")
	headers.Add("Origin", "https://www.acwing.com")
	headers.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")

	conn, _, err := websocket.DefaultDialer.Dial(serverURL, headers)
	if err != nil {
		log.Fatal("Error connecting to server:", err)
	}
	defer conn.Close()

	// 发送文件内容作为消息
	type data = struct {
		Activity  string `json:"activity"`
		Code      string `json:"code"`
		Language  string `json:"language"`
		Input     string `json:"input"`
		ProblemId int    `json:"problem_id"`
	}

	/*
		Activity:  "problem_run_code",
				Code:      `#include <iostream>\n\nconst int N = 100010;\n\nstruct test{\n  int a;  \n};\n\nint main(void){\n    int a, b;\n    scanf(\"%d %d\", &a, &b);\n    if (a == 3 && b == 4){\n        printf(\"%d\\n\", 7);\n    }else if(a == 45 && b == 55){\n        printf(\"%d\\n\", 100);\n    }\n    else if(a == 123 && b == 321){\n        printf(\"%d\\n\", 444);\n    }\n    else if(a == 123 && b == 123){\n        printf(\"%d\\n\", 444);\n    }\n    else if(a == 91086199 && b == 18700332){\n        printf(\"%d\\n\", 109786531);\n    }\n    else if(a == 42267194 && b == 60645282){\n        printf(\"%d\\n\", 102912476);\n    }\n    else if(a == 69274392 && b == 10635835){\n        printf(\"%d\\n\", 79910227);\n    }\n    else if(a == 5710219 && b == 85140568){\n        printf(\"%d\\n\", 90850787);\n    }\n    else if(a == 75601477 && b == 24005804){\n        printf(\"%d\\n\", 99607281);\n    }\n    else if(a == 70597795 && b == 90383234){\n        printf(\"%d\\n\", 160981029);\n    }\n    else if(a == 82574652 && b == 22252146){\n        printf(\"%d\\n\", 104826798);\n    }\n    // else if(a == 1 && b == 3){\n    //     printf(\"%d\\n\", a + b);\n    // }\n    for (;;){\n        new test();\n    }\n    return 0;\n}`,
				Language:  "C++",
				Input:     "3 4",
				ProblemId: 1,
	*/
	err = conn.WriteJSON(&data{})
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
		err = conn.WriteJSON(&data{
			Activity:  "problem_run_code",
			Code:      `#include <iostream>\n\nconst int N = 100010;\n\nstruct test{\n  int a;  \n};\n\nint main(void){\n    int a, b;\n    scanf(\"%d %d\", &a, &b);\n    if (a == 3 && b == 4){\n        printf(\"%d\\n\", 7);\n    }else if(a == 45 && b == 55){\n        printf(\"%d\\n\", 100);\n    }\n    else if(a == 123 && b == 321){\n        printf(\"%d\\n\", 444);\n    }\n    else if(a == 123 && b == 123){\n        printf(\"%d\\n\", 444);\n    }\n    else if(a == 91086199 && b == 18700332){\n        printf(\"%d\\n\", 109786531);\n    }\n    else if(a == 42267194 && b == 60645282){\n        printf(\"%d\\n\", 102912476);\n    }\n    else if(a == 69274392 && b == 10635835){\n        printf(\"%d\\n\", 79910227);\n    }\n    else if(a == 5710219 && b == 85140568){\n        printf(\"%d\\n\", 90850787);\n    }\n    else if(a == 75601477 && b == 24005804){\n        printf(\"%d\\n\", 99607281);\n    }\n    else if(a == 70597795 && b == 90383234){\n        printf(\"%d\\n\", 160981029);\n    }\n    else if(a == 82574652 && b == 22252146){\n        printf(\"%d\\n\", 104826798);\n    }\n    // else if(a == 1 && b == 3){\n    //     printf(\"%d\\n\", a + b);\n    // }\n    for (;;){\n        new test();\n    }\n    return 0;\n}`,
			Language:  "C++",
			Input:     "3 4",
			ProblemId: 1,
		})
		// 动态解析响应
		problemStatus := new(response.BaseResponse)
		if err := json.Unmarshal(message, problemStatus); err != nil {
			log.Println("Error unmarshalling base response:", err)
			continue
		}

		switch problemStatus.Status {
		case api.Finished:
			FinishedResponse := response.FinishedResponse{}
			if err := json.Unmarshal(message, &FinishedResponse); err != nil {
				log.Println("Error unmarshalling FinishedResponse:", err)
				continue
			}
			fmt.Printf("FinishedResponse: %+v\n", FinishedResponse)
			break
		case api.Pending:
			PendingResponse := response.PendingResponse{}
			if err := json.Unmarshal(message, &PendingResponse); err != nil {
				log.Println("Error unmarshalling PendingResponse:", err)
				continue
			}
			fmt.Printf("PendingResponse: %+v\n", PendingResponse)
			break
		case api.CompileError:
			CompileErrorResponse := response.CompileErrorResponse{}
			if err := json.Unmarshal(message, &CompileErrorResponse); err != nil {
				log.Println("Error unmarshalling CompileErrorResponse:", err)
				continue
			}
			fmt.Printf("CompileErrorResponse: %+v\n", CompileErrorResponse)
			break
		case api.Running:
			RunningResponse := response.RunningResponse{}
			if err := json.Unmarshal(message, &RunningResponse); err != nil {
				log.Println("Error unmarshalling RunningResponse:", err)
				continue
			}
			fmt.Printf("RunningResponse: %+v\n", RunningResponse)
			break
		case api.MemoryExceeded:
			MemoryExceededResponse := response.MemoryExceededResponse{}
			if err := json.Unmarshal(message, &MemoryExceededResponse); err != nil {
				log.Println("Error unmarshalling MemoryExceededResponse:", err)
				continue
			}
			fmt.Printf("MemoryExceededResponse: %+v\n", MemoryExceededResponse)
			break
		case api.Timeout:
			TimeoutResponse := response.TimeoutResponse{}
			if err := json.Unmarshal(message, &TimeoutResponse); err != nil {
				log.Println("Error unmarshalling TimeoutResponse:", err)
				continue
			}
			fmt.Printf("TimeoutResponse: %+v\n", TimeoutResponse)
			break
		case api.RunTimeError:
			TimeoutResponse := response.RunTimeErrorResponse{}
			if err := json.Unmarshal(message, &TimeoutResponse); err != nil {
				log.Println("Error unmarshalling TimeoutResponse:", err)
				continue
			}
			fmt.Printf("TimeoutResponse: %+v\n", TimeoutResponse)
			break
		case api.SystemError:
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
