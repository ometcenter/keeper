package models

import (
	"github.com/segmentio/kafka-go"
)

type AccrualMessageKafkaNew struct {
	Snils    string    `json:"snils_num"`
	Birthday string    `json:"birth_date"` //format: date
	Accruals []Accrual `json:"accruals"`
}

type Accrual struct {
	//Region string `json:"region"`
	//OrganizationCode string `json:"organizationCode"`
	EkisCode     string `json:"ekisCode"`
	Name         string `json:"name"`
	StartAt      string `json:"startAt"`      //format: date
	ExpirationAt string `json:"expirationAt"` //format: date
	ClassCount   int    `json:"classCount"`
	//IsActive         bool   `json:"isActive"`
	//StartAtGorod      string `json:"startAtGorod"`
	//ExpirationAtGorod string `json:"expirationAtGorod"`
}

type AccrualHeader struct {
	Snils    string `json:"snils_num"`
	Birthday string `json:"birth_date"` //format: date
}

// dividedSlices разделяет слайс на количество процессоров
// для многопоточного выполнения
func DividedSlicesKafkaMessages(mapkv []kafka.Message, chunkSize int) [][]kafka.Message {

	var divided [][]kafka.Message

	//chunkSize := (len(mapkv) + NumCPU - 1) / NumCPU

	for i := 0; i < len(mapkv); i += chunkSize {
		end := i + chunkSize

		if end > len(mapkv) {
			end = len(mapkv)
		}

		divided = append(divided, mapkv[i:end])
	}

	return divided
}
