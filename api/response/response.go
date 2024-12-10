package response

import (
	"github.com/polaris/codesandbox/api"
	"time"
)

// 用于实时反馈给用户
type ExcutedResult struct {
	TestCaseInput      string `json:"testCaseInput"`
	TestCaseOutput     string `json:"testCaseOutput"`
	TestCaseUserOutput string `json:"testCaseUserOutput"`
	ExcutedStatus
	StdErr string `json:"stdErr"`
}

// 用户实时返回给远程题目平台
type SyncBody struct {
	UserId      string        `json:"userId"`
	UserAccount string        `json:"userAccount"`
	QuestionId  string        `json:"questionId"`
	SubmitCode  string        `json:"submitCode"`
	Language    string        `json:"language"`
	JudgeInfo   ExcutedStatus `json:"judgeInfo"`
}

type ExcutedStatus struct {
	Time        time.Duration `json:"time"`        // 所有测试用例总共花费的时间
	Memory      int64         `json:"memory"`      // 所有测试用例总共使用的内存大小
	JudgeResult string        `json:"judgeResult"` // 当前提交的判题结果
}

type WebSocketResponse interface {
	Response() WebSocketResponse
}

type BaseResponse struct {
	Activity string `json:"activity"`
	Status   string `json:"status"`
}

// Pending
type PendingResponse struct {
	BaseResponse
}

func NewPendingResponse(activity string) *PendingResponse {
	return &PendingResponse{
		BaseResponse: BaseResponse{
			Activity: activity,
			Status:   api.Pending,
		},
	}
}

func (r *PendingResponse) Response() WebSocketResponse {
	//if err := r.Conn.WriteJSON(r); err != nil {
	//
	//}
	return r
}

// Running
type RunningResponse struct {
	BaseResponse
}

func NewRunningResponse(activity string) *RunningResponse {
	return &RunningResponse{
		BaseResponse: BaseResponse{
			Activity: activity,
			Status:   api.Running,
		},
	}
}

func (r *RunningResponse) Response() WebSocketResponse {
	//if err := r.Conn.WriteJSON(r); err != nil {
	//
	//}
	return r
}

// Accept
type AcceptResponse struct {
	BaseResponse
	ExcutedResult
}

func NewAcceptResponse(activity string, time time.Duration, memory int64) *AcceptResponse {
	return &AcceptResponse{
		BaseResponse: BaseResponse{
			Activity: activity,
			Status:   api.Accepted,
		},
		ExcutedResult: ExcutedResult{
			ExcutedStatus: ExcutedStatus{
				Time:   time,
				Memory: memory,
			},
		},
	}
}

func (r *AcceptResponse) Response() WebSocketResponse {
	//if err := r.Conn.WriteJSON(r); err != nil {
	//
	//}
	return r
}

// Finished
type FinishedResponse struct {
	BaseResponse
	ExcutedResult
}

func NewFinishedResponse(activity, stdOut string, time time.Duration, memory int64) *FinishedResponse {
	return &FinishedResponse{
		BaseResponse: BaseResponse{
			Activity: activity,
			Status:   api.Finished,
		},
		ExcutedResult: ExcutedResult{
			TestCaseUserOutput: stdOut,
			ExcutedStatus: ExcutedStatus{
				Time:   time,
				Memory: memory,
			},
		},
	}
}

func (r *FinishedResponse) Response() WebSocketResponse {

	//if err := r.Conn.WriteJSON(r); err != nil {
	//
	//}
	return r
}

type CompileErrorResponse struct {
	BaseResponse
	CompileError string `json:"stdErr"`
}

func NewCompileErrorResponse(activity, compileError string) *CompileErrorResponse {
	return &CompileErrorResponse{
		BaseResponse: BaseResponse{
			Activity: activity,
			Status:   api.CompileError,
		},
		CompileError: compileError,
	}
}

func (r *CompileErrorResponse) Response() WebSocketResponse {
	//if err := r.Conn.WriteJSON(r); err != nil {
	//
	//}
	return r
}

type WrongAnswerResponse struct {
	BaseResponse
	ExcutedResult
}

func NewWrongAnswerResponse(activity, testCaseInput, testCaseOutput, testCaseUserOutput string, time time.Duration, memory int64) *WrongAnswerResponse {
	return &WrongAnswerResponse{
		BaseResponse: BaseResponse{
			Activity: activity,
			Status:   api.WrongAnswer,
		},
		ExcutedResult: ExcutedResult{
			TestCaseInput:      testCaseInput,
			TestCaseOutput:     testCaseOutput,
			TestCaseUserOutput: testCaseUserOutput,
			ExcutedStatus: ExcutedStatus{
				Time:   time,
				Memory: memory,
			},
		},
	}
}

func (r *WrongAnswerResponse) Response() WebSocketResponse {
	//if err := r.Conn.WriteJSON(r); err != nil {
	//}
	return r
}

type PresentationErrorResponse struct {
	BaseResponse
	ExcutedResult
}

func NewPresentationErrorResponse(activity, testCaseInput, testCaseOutput, testCaseUserOutput string, time time.Duration, memory int64) *WrongAnswerResponse {
	return &WrongAnswerResponse{
		BaseResponse: BaseResponse{
			Activity: activity,
			Status:   api.PresentationError,
		},
		ExcutedResult: ExcutedResult{
			TestCaseInput:      testCaseInput,
			TestCaseOutput:     testCaseOutput,
			TestCaseUserOutput: testCaseUserOutput,
			ExcutedStatus: ExcutedStatus{
				Time:   time,
				Memory: memory,
			},
		},
	}
}

func (r *PresentationErrorResponse) Response() WebSocketResponse {
	//if err := r.Conn.WriteJSON(r); err != nil {
	//
	//}
	return r
}

type TimeoutResponse struct {
	BaseResponse
	ExcutedResult
}

func NewTimeoutResponse(activity, testCaseInput, testCaseOutput, testCaseUserOutput string, time time.Duration, memory int64) *TimeoutResponse {
	return &TimeoutResponse{
		BaseResponse: BaseResponse{
			Activity: activity,
			Status:   api.Timeout,
		},
		ExcutedResult: ExcutedResult{
			TestCaseInput:      testCaseInput,
			TestCaseOutput:     testCaseOutput,
			TestCaseUserOutput: testCaseUserOutput,
			ExcutedStatus: ExcutedStatus{
				Time:   time,
				Memory: memory,
			},
		},
	}
}

func (r *TimeoutResponse) Response() WebSocketResponse {
	//if err := r.Conn.WriteJSON(r); err != nil {
	//
	//}
	return r
}

type MemoryExceededResponse struct {
	BaseResponse
	ExcutedResult
}

func NewMemoryExceededResponse(activity, testCaseInput, testCaseOutput, testCaseUserOutput string, time time.Duration, memory int64) *MemoryExceededResponse {
	return &MemoryExceededResponse{
		BaseResponse: BaseResponse{
			Activity: activity,
			Status:   api.MemoryExceeded,
		},
		ExcutedResult: ExcutedResult{
			TestCaseInput:      testCaseInput,
			TestCaseOutput:     testCaseOutput,
			TestCaseUserOutput: testCaseUserOutput,
			ExcutedStatus: ExcutedStatus{
				Time:   time,
				Memory: memory,
			},
		},
	}
}

func (r *MemoryExceededResponse) Response() WebSocketResponse {
	//if err := r.Conn.WriteJSON(r); err != nil {
	//
	//}
	return r
}

type RunTimeErrorResponse struct {
	BaseResponse
	ExcutedResult
}

func NewRunTimeErrorResponse(activity, testCaseInput, testCaseOutput, testCaseUserOutput string, time time.Duration, memory int64, stdErr string) *RunTimeErrorResponse {
	return &RunTimeErrorResponse{
		BaseResponse: BaseResponse{
			Activity: activity,
			Status:   api.RunTimeError,
		},
		ExcutedResult: ExcutedResult{
			TestCaseInput:      testCaseInput,
			TestCaseOutput:     testCaseOutput,
			TestCaseUserOutput: testCaseUserOutput,
			ExcutedStatus: ExcutedStatus{
				Time:   time,
				Memory: memory,
			},
			StdErr: stdErr,
		},
	}
}

func (r *RunTimeErrorResponse) Response() WebSocketResponse {
	//if err := r.Conn.WriteJSON(r); err != nil {
	//
	//}
	return r
}

type SystemErrorResponse struct {
	BaseResponse
}

func NewSystemErrorResponse(status string) *SystemErrorResponse {
	return &SystemErrorResponse{
		BaseResponse: BaseResponse{
			Activity: api.SystemError,
			Status:   status,
		},
	}
}

func (r *SystemErrorResponse) Response() WebSocketResponse {
	//if err := r.Conn.WriteJSON(r); err != nil {
	//
	//}
	return r
}
