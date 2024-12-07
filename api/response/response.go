package response

// activity 代码沙箱返回给前端的活动状态
var (
	RunAct    = "RUN_CODE_ACTIVITY_STATUS"
	SubmitAct = "SUBMIT_CODE_ACTIVITY_STATUS"
	ErrorAct  = "ERROR_ACTIVITY_STATUS"
)

// status
var (
	Pending = "PENDING"
	Running = "RUNNING"

	Accepted          = "ACCEPTED"
	WrongAnswer       = "WRONG_ANSWER"
	CompileError      = "COMPILE_ERROR"
	PresentationError = "PRESENTATION_ERROR"
	Finished          = "FINISHED"
	Timeout           = "TIME_LIMIT_EXCEEDED"
	MemoryExceeded    = "MEMORY_LIMIT_EXCEEDED"
	// 系统错误
	SystemError     = "SYSTEM_ERROR"
	ParamsError     = "PARAMS_ERROR"
	NoLanguageError = "NO_LANGUAGE_ERROR"
	TestCasesError  = "TESTCASES_ERROR"
)

type WebSocketResponse interface {
	Response() WebSocketResponse
}

type BaseResponse struct {
	Activity string `json:"activity"`
	Status   string `json:"status"`
}

type PendingResponse struct {
	BaseResponse
}

func NewPendingResponse(activity string) *PendingResponse {
	return &PendingResponse{
		BaseResponse: BaseResponse{
			Activity: activity,
			Status:   Pending,
		},
	}
}

func (r *PendingResponse) Response() WebSocketResponse {
	//if err := r.Conn.WriteJSON(r); err != nil {
	//
	//}
	return r
}

type RunningResponse struct {
	BaseResponse
}

func NewRunningResponse(activity string) *RunningResponse {
	return &RunningResponse{
		BaseResponse: BaseResponse{
			Activity: activity,
			Status:   Running,
		},
	}
}

func (r *RunningResponse) Response() WebSocketResponse {
	//if err := r.Conn.WriteJSON(r); err != nil {
	//
	//}
	return r
}

type AcceptResponse struct {
	BaseResponse
	Time int64 `json:"time"`
}

func NewAcceptResponse(activity string, time int64) *AcceptResponse {
	return &AcceptResponse{
		BaseResponse: BaseResponse{

			Activity: activity,
			Status:   Accepted,
		},
		Time: time,
	}
}

func (r *AcceptResponse) Response() WebSocketResponse {
	//if err := r.Conn.WriteJSON(r); err != nil {
	//
	//}
	return r
}

type FinishedResponse struct {
	BaseResponse
	StdErr string `json:"stdErr"`
	StdOut string `json:"stdOut"`
}

func NewFinishedResponse(activity, stdErr, stdOut string) *FinishedResponse {
	return &FinishedResponse{
		BaseResponse: BaseResponse{

			Activity: activity,
			Status:   Finished,
		},
		StdErr: stdErr,
		StdOut: stdOut,
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
	StdErr string `json:"stdErr"`
	StdOut string `json:"stdOut"`
}

func NewCompileErrorResponse(activity, stdErr, stdOut string) *CompileErrorResponse {
	return &CompileErrorResponse{
		BaseResponse: BaseResponse{

			Activity: activity,
			Status:   CompileError,
		},
		StdErr: stdErr,
		StdOut: stdOut,
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
	//StdErr             string `json:"stdErr"`
	TestCaseInput      string `json:"testCaseInput"`
	TestCaseOutput     string `json:"testCaseOutput"`
	TestCaseUserOutput string `json:"testCaseUserOutput"`
}

func NewWrongAnswerResponse(activity, testCaseInput, testCaseOutput, testCaseUserOutput string) *WrongAnswerResponse {
	return &WrongAnswerResponse{
		BaseResponse: BaseResponse{

			Activity: activity,
			Status:   WrongAnswer,
		},
		TestCaseInput:      testCaseInput,
		TestCaseOutput:     testCaseOutput,
		TestCaseUserOutput: testCaseUserOutput,
	}
}

func (r *WrongAnswerResponse) Response() WebSocketResponse {
	//if err := r.Conn.WriteJSON(r); err != nil {
	//}
	return r
}

type PresentationErrorResponse struct {
	BaseResponse
	TestCaseInput      string `json:"testCaseInput"`
	TestCaseOutput     string `json:"testCaseOutput"`
	TestCaseUserOutput string `json:"testCaseUserOutput"`
}

func NewPresentationErrorResponse(activity, testCaseInput, testCaseOutput, testCaseUserOutput string) *WrongAnswerResponse {
	return &WrongAnswerResponse{
		BaseResponse: BaseResponse{

			Activity: activity,
			Status:   PresentationError,
		},
		TestCaseInput:      testCaseInput,
		TestCaseOutput:     testCaseOutput,
		TestCaseUserOutput: testCaseUserOutput,
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
	TestCaseInput      string `json:"testCaseInput"`
	TestCaseOutput     string `json:"testCaseOutput"`
	TestCaseUserOutput string `json:"testCaseUserOutput"`
}

func NewTimeoutResponse(activity, testCaseInput, testCaseOutput, testCaseUserOutput string) *TimeoutResponse {
	return &TimeoutResponse{
		BaseResponse: BaseResponse{

			Activity: activity,
			Status:   Timeout,
		},
		TestCaseInput:      testCaseInput,
		TestCaseOutput:     testCaseOutput,
		TestCaseUserOutput: testCaseUserOutput,
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
}

func NewMemoryExceededResponse(activity string) *MemoryExceededResponse {
	return &MemoryExceededResponse{
		BaseResponse: BaseResponse{
			Activity: activity,
			Status:   MemoryExceeded,
		},
	}
}

func (r *MemoryExceededResponse) Response() WebSocketResponse {
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

			Activity: ErrorAct,
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
