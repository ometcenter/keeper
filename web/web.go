package web

type AnswerWebV1 struct {
	Status bool        `json:"status"`
	Data   interface{} `json:"data,omitempty"` //`json:"data,omitempty"`
	Error  *ErrorWebV1 `json:"error,omitempty"`
}

type ErrorWebV1 struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
