package models

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TelemetryClientInfo struct {
	gorm.Model
	JobID                 string    `json:"jobID"`               //`json:"ИдентификаторЗадания"`
	RemoteBaseID          string    `json:"remoteBaseID"`        //`json:"ИдентификаторУдаленнойБазы"`
	DateTelemetry         time.Time `json:"dateTelemetry"`       //`json:"ДатаМетрики"`
	Version1CConfig       string    `json:"versionConfig"`       //`json:"ВерсияКонфигурацииМетоданные"`
	Name1CBase            string    `json:"nameBase"`            //`json:"ИмяБазыМетоданные"`
	StringConnection      string    `json:"stringConnection"`    //`json:"СтрокаСоединения"`
	PlatformType          string    `json:"platformType"`        //`json:"ТипПлатформы"`
	ApplicationVersion    string    `json:"applicationVersion"`  //`json:"ВерсияПриложения"`
	ClientID              string    `json:"clientID"`            //`json:"ИдентификаторКлиента"`
	VersionOS             string    `json:"versionOS"`           //`json:"ВерсияОС"`
	NumberVersion1CScript string    `json:"numberVersionScript"` //`json:"НомерВерсииОбработки"`
	TypeTelemetry         string    `json:"typeTelemetry"`
	ExchangeJobID         string    `json:"exchangeJobID"`
	AdditionalInformation string    `json:"additionalInformation"`
	//AdditionalInformationJSONByte  datatypes.JSON
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

func (TelemetryClientInfo *TelemetryClientInfo) FillTelemetryFromHeaderBASE64(c *gin.Context, TelemetryClientInfoBodyString string) error {

	sDec, err := base64.StdEncoding.DecodeString(TelemetryClientInfoBodyString)
	if err != nil {
		return nil
	}

	//fmt.Printf("Телеметрия2: %v\n", string(sDec))

	sDec = bytes.TrimPrefix(sDec, []byte("\xef\xbb\xbf")) // Or []byte{239, 187, 191}

	err = json.Unmarshal(sDec, &TelemetryClientInfo)
	if err != nil {
		return err
	}

	// //TelemetryClientInfo.DateTelemetry = time.Now()

	// JsonMessageBody, err := json.Marshal(&TelemetryClientInfo)
	// if err != nil {
	// 	return err
	// }

	return nil
}
