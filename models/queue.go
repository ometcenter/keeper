package models

type MessageQueueGeneralInterface struct {
	Type string      `json:"type"`
	Body interface{} `json:"body"`
}

type HandleAfterLoad struct {
	JobID      string   `json:"jobID"`
	Algorithms []string `json:"algorithms"`
}

// Обертка, чтобы не разбирать тело вручную
type HandleAfterLoadWrap struct {
	Type            string          `json:"type"`
	HandleAfterLoad HandleAfterLoad `json:"body"`
}
