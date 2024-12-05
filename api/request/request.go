package request

var (
	SubmitAct    = "SUBMIT_CODE_ACTIVITY"
	RunAct       = "RUN_CODE_ACTIVITY"
	HeartBeatAct = "HEART_BEAT_ACTIVITY"
)

type WebSocketRequest HeartBeat

type HeartBeat struct {
	Activity string `json:"activity"`
}

type ProblemSubmit struct {
	Activity   string `json:"activity"`
	Code       string `json:"code"`
	Language   string `json:"language"`
	QuestionId string `json:"questionId"`
}

type ProblemRun struct {
	Activity string `json:"activity"`
	Code     string `json:"code"`
	Language string `json:"language"`
	//QuestionId string `json:"questionId"`
	Input string `json:"input"`
}
