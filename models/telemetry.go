package models

import (
	"bytes"
	"encoding/base64"
	"time"

	"github.com/gin-gonic/gin"
)

type TelemetryClientInfo struct {
	JobID                 string    `json:"ИдентификаторЗадания"`
	RemoteBaseID          string    `json:"ИдентификаторУдаленнойБазы"`
	DateTelemetry         time.Time `json:"ДатаМетрики"`
	Version1CConfig       string    `json:"ВерсияКонфигурацииМетоданные"`
	Name1CBase            string    `json:"ИмяБазыМетоданные"`
	StringConnection      string    `json:"СтрокаСоединения"`
	PlatformType          string    `json:"ТипПлатформы"`
	ApplicationVersion    string    `json:"ВерсияПриложения"`
	ClientID              string    `json:"ИдентификаторКлиента"`
	VersionOS             string    `json:"ВерсияОС"`
	NumberVersion1CScript string    `json:"НомерВерсииОбработки"`
}

func (TelemetryClientInfo *TelemetryClientInfo) GetTelemetryFromHeaderBASE64(c *gin.Context, TelemetryClientInfoBodyString string) ([]byte, error) {

	sDec, err := base64.StdEncoding.DecodeString(TelemetryClientInfoBodyString)
	if err != nil {
		return nil, err
	}

	//fmt.Printf("Телеметрия2: %v\n", string(sDec))

	sDec = bytes.TrimPrefix(sDec, []byte("\xef\xbb\xbf")) // Or []byte{239, 187, 191}

	// err = json.Unmarshal(sDec, &TelemetryClientInfo)
	// if err != nil {
	// 	return err
	// }

	// //TelemetryClientInfo.DateTelemetry = time.Now()

	// JsonMessageBody, err := json.Marshal(&TelemetryClientInfo)
	// if err != nil {
	// 	return err
	// }

	return sDec, nil
}
