package models

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"gorm.io/datatypes"
)

// type SettingsJobs struct {
// 	JobID      string `json:"ИдентификаторЗадания"`
// 	JSONString string `json:"JSONСтрокаНастроек"`
// }

type SettingsJobs struct {
	JobID               string `json:"ИдентификаторЗадания"`
	JSONString          string `json:"JSONСтрокаНастроек"`
	CodeExternal        string `json:"Код1С"`
	NameExternal        string `json:"Наименование1С"`
	TableName           string `json:"ИмяТаблицы"`
	UseRemoteCollection bool   `json:"УдаленныйСбор"`
	ConfigName          string `json:"ИмяКонфигурации"`
	TypeDataGetting     string `json:"ВидПолученияДанных"`
	JSONByte            datatypes.JSON
}

// Используется для выгрузки в Систему визуализации данных
type DataForDataVisualizationSystem struct {
	//Data      json.RawMessage //`json:"ИдентификаторЗадания"`
	Data      []map[string]interface{} //`json:"ИдентификаторЗадания"`
	QueryToBI QueryToBI                //`json:"Настройки"`
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
	JobID                            string                 `json:"ИдентификаторЗадания"`
	SendUseREST                      bool                   `json:"ОтправлятьПоREST"`
	RemoteCollect                    bool                   `json:"УдаленныйСбор"`
	TypeDataGetting                  string                 `json:"ВидПолученияДанных"`
	DataUploadMethod                 string                 `json:"CпособЗагрузкиДанных"`
	Portions                         int                    `json:"Порции"`
	Query                            []Query                `json:"Запросы"`
	AddParam                         AdditionParam          `json:"ДополнительныеПараметрыJSON"`
	AddParamJSNOString               string                 `json:"JSONСтрокаДополнительныеПараметры"`
	Connect                          Connect                `json:"ПараметрыПодключения"`
	ConnectContur                    ConnectContur          `json:"ПараметрыПодключенияКонтура"`
	ConnectBI1C                      ConnectBI1C            `json:"ПараметрыПодключенияBI1C"`
	ConnectConturJSNOString          string                 `json:"JSONСтрокаПараметрыПодключенияКонтура"`
	Schedule                         Schedule               `json:"РасписаниеПланировщика"`
	SaveResultToHistory              bool                   `json:"СохранятьРезультатВИсторию"`
	PublishTableToAPI                bool                   `json:"ПубликоватьТаблицуВAPI"`
	SaveToDataVisualizationSystem    bool                   `json:"СохранятьВСистемуВизуализацииДанных"`
	UseDataProcessingAlgorithms      bool                   `json:"ИспользоватьАлгоритмыОбработкиДанных"`
	ListDataProcessingAlgorithms     []string               `json:"СписокАлгоритмовОбработкиДанных"`
	UseHandleAfterLoadAlgorithms     bool                   `json:"ИспользоватьАлгоритмыОбработкиДанныхПослеЗагрузки"`
	ListHandleAfterLoadAlgorithms    []string               `json:"СписокАлгоритмовОбработкиДанныхПослеЗагрузки"`
	Webhooks                         []string               `json:"Webhooks"`
	MappingForExcelArray             []MappingForExcelArray `json:"СопоставлениеДляExcalМассив"`
	RuleExternalSource               string                 `json:"ПравилоВнешнийИсточник"`
	InternalProcessingExternalSource bool                   `json:"ВнутренняяОбработкаВнешнегоИсточника"`
	//SettingsJobsV2                   SettingsJobsV2         `json:"settingsJobsV2"`
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

	if strings.EqualFold(os.Getenv("USE_SETTINGS_JOB_V2"), "true") {

		var SettingsJobsAllV2 SettingsJobsAllV2
		err := SettingsJobsAllV2.LoadSettingsFromPgByJobID(DB, JobID)
		if err != nil {
			return err
		}

		QTranform, err := SettingsJobsAllV2.TranformToOldSettings()
		if err != nil {
			return err
		}

		*Q = QTranform

	} else {

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
	}

	return nil
}

func (Q *QueryToBI) LoadSettingsFirstRowFromPgByJobIDByTableName(DB *sql.DB, TableName string) error {

	var argsquery []interface{}
	argsquery = append(argsquery, TableName)

	var LoadValue SettingsJobSliceQueryToBI
	err := DB.QueryRow("SELECT json_byte FROM settings_jobs WHERE table_name = $1", argsquery...).Scan(&LoadValue)
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

type MappingForExcelArray struct {
	NumberField int    `json:"НомерСтроки"`
	Name        string `json:"Имя"`
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
	Description                       string                              `json:"НаименованиеЗадания"`
	TableName                         string                              `json:"ИмяТаблицы"`
	DSNconnection                     string                              `json:"БазаСУБДDSN"`
	HardRemoval                       bool                                `json:"ПолноеУдаление"`
	SelectionFields                   []string                            `json:"ПоляОтбора"`
	ComparionFields                   []string                            `json:"ПоляСравнения"`
	CompareAllFields                  bool                                `json:"СравниватьПоВсемПолям"`
	CompressBody                      bool                                `json:"СжиматьОтвет"`
	AnonymizingFields                 []AnonymizingFields                 `json:"АнонимизацияПолей"`
	DataVisualizationSystemParameters []DataVisualizationSystemParameters `json:"ПараметрыСистемыВизуализацииДанных"`
	//CleaningFieldsBeforeLoading       []CleaningFieldsBeforeLoading       `json:"ПоляОчисткиПередЗагрузкой"`
}

type AnonymizingFields struct {
	Name string `json:"Имя"`
	Type string `json:"Тип"`
}

type DataVisualizationSystemParameters struct {
	Name      string `json:"Name"`
	ValueData string `json:"Value"`
}

type CleaningFieldsBeforeLoading struct {
	Name      string      `json:"Name"`
	Type      string      `json:"Type"`
	ValueData interface{} `json:"Value"`
}
