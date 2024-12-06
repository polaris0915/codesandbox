package response

import (
	"github.com/gorilla/websocket"
)

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
	Response()
}

type BaseResponse struct {
	Conn     *websocket.Conn `json:"-"`
	Activity string          `json:"activity"`
	Status   string          `json:"status"`
}

type PendingResponse struct {
	BaseResponse
}

func NewPendingResponse(conn *websocket.Conn, activity string) *PendingResponse {
	return &PendingResponse{
		BaseResponse: BaseResponse{
			Conn:     conn,
			Activity: activity,
			Status:   Pending,
		},
	}
}

func (r *PendingResponse) Response() {
	if err := r.Conn.WriteJSON(r); err != nil {

	}
}

type RunningResponse struct {
	BaseResponse
}

func NewRunningResponse(conn *websocket.Conn, activity string) *RunningResponse {
	return &RunningResponse{
		BaseResponse: BaseResponse{
			Conn:     conn,
			Activity: activity,
			Status:   Running,
		},
	}
}

func (r *RunningResponse) Response() {
	if err := r.Conn.WriteJSON(r); err != nil {

	}
}

type AcceptResponse struct {
	BaseResponse
	Time int64 `json:"time"`
}

func NewAcceptResponse(conn *websocket.Conn, activity string, time int64) *AcceptResponse {
	return &AcceptResponse{
		BaseResponse: BaseResponse{
			Conn:     conn,
			Activity: activity,
			Status:   Accepted,
		},
		Time: time,
	}
}

func (r *AcceptResponse) Response() {
	if err := r.Conn.WriteJSON(r); err != nil {

	}
}

type FinishedResponse struct {
	BaseResponse
	StdErr string `json:"stdErr"`
	StdOut string `json:"stdOut"`
}

func NewFinishedResponse(conn *websocket.Conn, activity, stdErr, stdOut string) *FinishedResponse {
	return &FinishedResponse{
		BaseResponse: BaseResponse{
			Conn:     conn,
			Activity: activity,
			Status:   Finished,
		},
		StdErr: stdErr,
		StdOut: stdOut,
	}
}

func (r *FinishedResponse) Response() {

	if err := r.Conn.WriteJSON(r); err != nil {

	}
}

type CompileErrorResponse struct {
	BaseResponse
	StdErr string `json:"stdErr"`
	StdOut string `json:"stdOut"`
}

func NewCompileErrorResponse(conn *websocket.Conn, activity, stdErr, stdOut string) *CompileErrorResponse {
	return &CompileErrorResponse{
		BaseResponse: BaseResponse{
			Conn:     conn,
			Activity: activity,
			Status:   CompileError,
		},
		StdErr: stdErr,
		StdOut: stdOut,
	}
}

func (r *CompileErrorResponse) Response() {
	if err := r.Conn.WriteJSON(r); err != nil {

	}
}

type WrongAnswerResponse struct {
	BaseResponse
	//StdErr             string `json:"stdErr"`
	TestCaseInput      string `json:"testCaseInput"`
	TestCaseOutput     string `json:"testCaseOutput"`
	TestCaseUserOutput string `json:"testCaseUserOutput"`
}

func NewWrongAnswerResponse(conn *websocket.Conn, activity, testCaseInput, testCaseOutput, testCaseUserOutput string) *WrongAnswerResponse {
	return &WrongAnswerResponse{
		BaseResponse: BaseResponse{
			Conn:     conn,
			Activity: activity,
			Status:   WrongAnswer,
		},
		TestCaseInput:      testCaseInput,
		TestCaseOutput:     testCaseOutput,
		TestCaseUserOutput: testCaseUserOutput,
	}
}

func (r *WrongAnswerResponse) Response() {
	if err := r.Conn.WriteJSON(r); err != nil {
	}
}

type PresentationErrorResponse struct {
	BaseResponse
	TestCaseInput      string `json:"testCaseInput"`
	TestCaseOutput     string `json:"testCaseOutput"`
	TestCaseUserOutput string `json:"testCaseUserOutput"`
}

func NewPresentationErrorResponse(conn *websocket.Conn, activity, testCaseInput, testCaseOutput, testCaseUserOutput string) *WrongAnswerResponse {
	return &WrongAnswerResponse{
		BaseResponse: BaseResponse{
			Conn:     conn,
			Activity: activity,
			Status:   PresentationError,
		},
		TestCaseInput:      testCaseInput,
		TestCaseOutput:     testCaseOutput,
		TestCaseUserOutput: testCaseUserOutput,
	}
}

func (r *PresentationErrorResponse) Response() {
	if err := r.Conn.WriteJSON(r); err != nil {

	}
}

type TimeoutResponse struct {
	BaseResponse
	TestCaseInput      string `json:"testCaseInput"`
	TestCaseOutput     string `json:"testCaseOutput"`
	TestCaseUserOutput string `json:"testCaseUserOutput"`
}

func NewTimeoutResponse(conn *websocket.Conn, activity, testCaseInput, testCaseOutput, testCaseUserOutput string) *TimeoutResponse {
	return &TimeoutResponse{
		BaseResponse: BaseResponse{
			Conn:     conn,
			Activity: activity,
			Status:   Timeout,
		},
		TestCaseInput:      testCaseInput,
		TestCaseOutput:     testCaseOutput,
		TestCaseUserOutput: testCaseUserOutput,
	}
}

func (r *TimeoutResponse) Response() {
	if err := r.Conn.WriteJSON(r); err != nil {

	}
}

type MemoryExceededResponse struct {
	BaseResponse
}

func NewMemoryExceededResponse(conn *websocket.Conn, activity string) *MemoryExceededResponse {
	return &MemoryExceededResponse{
		BaseResponse: BaseResponse{
			Conn:     conn,
			Activity: activity,
			Status:   MemoryExceeded,
		},
	}
}

func (r *MemoryExceededResponse) Response() {
	if err := r.Conn.WriteJSON(r); err != nil {

	}
}

type SystemErrorResponse struct {
	BaseResponse
}

func NewSystemErrorResponse(conn *websocket.Conn, status string) *SystemErrorResponse {
	return &SystemErrorResponse{
		BaseResponse: BaseResponse{
			Conn:     conn,
			Activity: ErrorAct,
			Status:   status,
		},
	}
}

func (r *SystemErrorResponse) Response() {
	if err := r.Conn.WriteJSON(r); err != nil {

	}
}
