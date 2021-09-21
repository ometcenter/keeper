package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type SettingsJobs struct {
	JobID      string `json:"ИдентификаторЗадания"`
	JSONString string `json:"JSONСтрокаНастроек"`
}

// Корверая стурктура описывающее задание обмена, хранящаяся в таблице settings_jobs.
type SettingsJobSliceQueryToBI struct {
	JobID          string      `json:"ИдентификаторЗадания"`
	SliceQueryToBI []QueryToBI `json:"Настройки"`
}

func (SettingsJobSliceQueryToBI *SettingsJobSliceQueryToBI) Scan(value interface{}) (err error) {
	switch value.(type) {
	case string:
		err = json.Unmarshal([]byte(value.(string)), &SettingsJobSliceQueryToBI)
	case []byte:
		err = json.Unmarshal(value.([]byte), &SettingsJobSliceQueryToBI)
	default:
		return errors.New("Incompatible type for Skills")
	}
	if err != nil {
		return
	}

	return nil
}

func (SettingsJobSliceQueryToBI SettingsJobSliceQueryToBI) Value() (driver.Value, error) {
	return json.Marshal(SettingsJobSliceQueryToBI)
}

type QueryToBI struct {
	JobID                   string        `json:"ИдентификаторЗадания"`
	SendUseREST             bool          `json:"ОтправлятьПоREST"`
	RemoteCollect           bool          `json:"УдаленныйСбор"`
	Portions                int           `json:"Порции"`
	Query                   []Query       `json:"Запросы"`
	AddParam                AdditionParam `json:"ДополнительныеПараметрыJSON"`
	AddParamJSNOString      string        `json:"JSONСтрокаДополнительныеПараметры"`
	Connect                 Connect       `json:"ПараметрыПодключения"`
	ConnectContur           ConnectContur `json:"ПараметрыПодключенияКонтура"`
	ConnectBI1C             ConnectBI1C   `json:"ПараметрыПодключенияBI1C"`
	ConnectConturJSNOString string        `json:"JSONСтрокаПараметрыПодключенияКонтура"`
	Schedule                Schedule      `json:"РасписаниеПланировщика"`
}

func (QueryToBI *QueryToBI) Scan(value interface{}) (err error) {
	switch value.(type) {
	case string:
		err = json.Unmarshal([]byte(value.(string)), &QueryToBI)
	case []byte:
		err = json.Unmarshal(value.([]byte), &QueryToBI)
	default:
		return errors.New("Incompatible type for Skills")
	}
	if err != nil {
		return
	}

	return nil
}

func (QueryToBI QueryToBI) Value() (driver.Value, error) {
	return json.Marshal(QueryToBI)
}

type Query struct {
	QueryText                            string `json:"Запрос"`
	Base                                 string `json:"База"`
	Area                                 string `json:"Области"`
	ExchangeJobID                        string `json:"ИдентификаторЗапроса"`
	PText                                string `json:"ПараметрыЗапроса"`
	UsedCalculatedFieldsInQueryParametrs bool   `json:"ИспользуетсяВычисляемыеПоляВПараметрахЗапроса"`
}

type Schedule struct {
	JobID             string   `json:"ИдентификаторЗадания"`
	UseRegulatoryTask bool     `json:"ИспользоватьРегламентноеЗадание"`
	UseCronSchedule   bool     `json:"ИспользоватьCronРасписание"`
	SliceCronString   []string `json:"МассивСтрокCron"`
}

type AdditionParam struct {
	ZipAnswer       bool     `json:"СжиматьОтвет"`
	HashAnswer      bool     `json:"ХешироватьРезультат"`
	DeSerialization bool     `json:"ДесериализацияXDTO"`
	Options         *Options `json:"НастройкиМоделиДанных"`
	Connect         Connect  `json:"ПараметрыПодключенияHTTPОтвета"`
}

type Connect struct {
	AddressServer    string `json:"АдресСервиса"`
	PortServer       int    `json:"Порт"`
	Resource         string `json:"Ресурс"`
	SecureConnetion  bool   `json:"ЗащищенноеСоединение"`
	SecureConnetion2 bool   `json:"ЗащищенноеСоединенние"`
	LoginConnetion   string `json:"Логин"`
	Password         string `json:"Пароль"`
	Headers          string `json:"Заголовки"`
}

type ConnectBI1C struct {
	AddressServer    string `json:"АдресСервиса"`
	PortServer       int    `json:"Порт"`
	Resource         string `json:"Ресурс"`
	SecureConnetion  bool   `json:"ЗащищенноеСоединение"`
	SecureConnetion2 bool   `json:"ЗащищенноеСоединенние"`
	LoginConnetion   string `json:"Логин"`
	Password         string `json:"Пароль"`
	Headers          string `json:"Заголовки"`
}

type ConnectContur struct {
	AddressServer    string `json:"АдресСервисаПриемника"`
	PortServer       int    `json:"Порт"`
	Resource         string `json:"Ресурс"`
	SecureConnetion  bool   `json:"ЗащищенноеСоединение"`
	SecureConnetion2 bool   `json:"ЗащищенноеСоединенние"`
	LoginConnetion   string `json:"Логин"`
	Password         string `json:"Пароль"`
	Headers          string `json:"Заголовки"`
}

type Options struct {
	Description      string   `json:"НаименованиеЗадания"`
	TableName        string   `json:"ИмяТаблицы"`
	HardRemoval      bool     `json:"ПолноеУдаление"`
	SelectionFields  []string `json:"ПоляОтбора"`
	ComparionFields  []string `json:"ПоляСравнения"`
	CompareAllFields bool     `json:"СравниватьПоВсемПолям"`
	CompressBody     bool     `json:"СжиматьОтвет"`
}
