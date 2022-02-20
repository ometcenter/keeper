package models

type MessageQueueGeneral struct {
	Type string `json:"type"`
	Body string `json:"body"`
}

type MessageQueueGeneralInterface struct {
	Type string      `json:"type"`
	Body interface{} `json:"body"`
}

type HandleAfterLoad struct {
	JobID     string   `json:"jobID"`
	Algorithm []string `json:"algorithm"`
}
