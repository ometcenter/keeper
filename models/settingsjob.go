package models

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
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
		return errors.New("incompatible type for skills")
	}
	if err != nil {
		return
	}

	return nil
}

func (SettingsJobSliceQueryToBI SettingsJobSliceQueryToBI) Value() (driver.Value, error) {
	return json.Marshal(SettingsJobSliceQueryToBI)
}

func (S *SettingsJobSliceQueryToBI) LoadSettingsFromPgByJobID(DB *sql.DB, JobID string) error {

	var argsquery []interface{}
	argsquery = append(argsquery, JobID)

	var LoadValue SettingsJobSliceQueryToBI
	err := DB.QueryRow("SELECT json_byte FROM settings_jobs WHERE job_id = $1", argsquery...).Scan(&LoadValue)
	if err != nil {
		return err
	}

	// TODO: Переделать по нормальному эту конструкцию
	*S = LoadValue

	return nil
}

type QueryToBI struct {
	JobID                         string        `json:"ИдентификаторЗадания"`
	SendUseREST                   bool          `json:"ОтправлятьПоREST"`
	RemoteCollect                 bool          `json:"УдаленныйСбор"`
	Portions                      int           `json:"Порции"`
	Query                         []Query       `json:"Запросы"`
	AddParam                      AdditionParam `json:"ДополнительныеПараметрыJSON"`
	AddParamJSNOString            string        `json:"JSONСтрокаДополнительныеПараметры"`
	Connect                       Connect       `json:"ПараметрыПодключения"`
	ConnectContur                 ConnectContur `json:"ПараметрыПодключенияКонтура"`
	ConnectBI1C                   ConnectBI1C   `json:"ПараметрыПодключенияBI1C"`
	ConnectConturJSNOString       string        `json:"JSONСтрокаПараметрыПодключенияКонтура"`
	Schedule                      Schedule      `json:"РасписаниеПланировщика"`
	SaveResultToHistory           string        `json:"СохранятьРезультатВИсторию"`
	SaveToDataVisualizationSystem string        `json:"СохранятьВСистемуВизуализацииДанных"`
}

func (QueryToBI *QueryToBI) Scan(value interface{}) (err error) {
	switch value.(type) {
	case string:
		err = json.Unmarshal([]byte(value.(string)), &QueryToBI)
	case []byte:
		err = json.Unmarshal(value.([]byte), &QueryToBI)
	default:
		return errors.New("incompatible type for skills")
	}
	if err != nil {
		return
	}

	return nil
}

func (QueryToBI QueryToBI) Value() (driver.Value, error) {
	return json.Marshal(QueryToBI)
}

func (Q *QueryToBI) LoadSettingsFirstRowFromPgByJobID(DB *sql.DB, JobID string) error {

	var argsquery []interface{}
	argsquery = append(argsquery, JobID)

	var LoadValue SettingsJobSliceQueryToBI
	err := DB.QueryRow("SELECT json_byte FROM settings_jobs WHERE job_id = $1", argsquery...).Scan(&LoadValue)
	if err != nil {
		return err
	}

	if len(LoadValue.SliceQueryToBI) > 0 {
		// TODO: Переделать по нормальному эту конструкцию
		*Q = LoadValue.SliceQueryToBI[0]
	} else {
		return fmt.Errorf("получена пустая настройка")
	}

	return nil
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
	// TODO: Пробросить поля: ПараметрыСистемыВизуализацииДанных и АнонимизацияПолей
}
