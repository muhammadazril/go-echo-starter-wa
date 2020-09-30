package model

// WAMessage model for request and job args
type WAMessage struct {
	Message string `json:"message" validate:"required"`
	Number  string `json:"number" validate:"required"`
}

// Response for API
type Response struct {
	MsgID  string `json:"msgId"`
	Status string `json:"status"`
}

// ResponseJob for Request via worker
type ResponseJob struct {
	Message string `json:"msg"`
	Status  string `json:"status"`
	JobID   string `json:"job_id"`
}
