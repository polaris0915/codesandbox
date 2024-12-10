package api

// activity 代码沙箱返回给前端的活动状态
var (
	RunAct    = "RUN_CODE_ACTIVITY_STATUS"
	SubmitAct = "SUBMIT_CODE_ACTIVITY_STATUS"
)

// status
var (
	Pending = "PENDING"
	Running = "RUNNING"

	Accepted          = "ACCEPTED"
	WrongAnswer       = "WRONG_ANSWER"
	CompileError      = "COMPILE_ERROR"
	RunTimeError      = "RUN_TIME_ERROR"
	PresentationError = "PRESENTATION_ERROR"
	Finished          = "FINISHED"
	Timeout           = "TIME_LIMIT_EXCEEDED"
	MemoryExceeded    = "MEMORY_LIMIT_EXCEEDED"
	// 系统错误
	SystemError         = "SYSTEM_ERROR"
	NoMemoryStatusError = "NO_MEMORY_STATUS_ERROR"
	ParamsError         = "PARAMS_ERROR"
	NoLanguageError     = "NO_LANGUAGE_ERROR"
	TestCasesError      = "TESTCASES_ERROR"
	UserInfoError       = "USER_INFO_ERROR"
)
