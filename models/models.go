package models

import "time"

type TableDescription struct {
	TableName string   `json:"ИмяТаблицы"`
	Fields    []Fields `json:"СведенияКолонкиТаблицы"`
}

type IndexDescription struct {
	TableName string        `json:"ИмяТаблицы"`
	Fields    []FieldsIndex `json:"СведенияКолонкиТаблицы"`
}

type FieldsIndex struct {
	Name       string `json:"Имя"`
	Definition string `json:"Определение"`
	TypeChange string `json:"ИзменитьВСУБД"`
}

type Fields struct {
	Name       string `json:"Имя"`
	Type       string `json:"Тип"`
	NotNull    bool   `json:"NotNull"`
	PrimaryKey bool   `json:"ПервичныйКлюч"`
	TypeChange string `json:"ИзменитьВСУБД"`
}

type ColumnsStruct struct {
	ColumnName string
	DataType   string
	IsNullable string
	PrimaryKey bool
}

type IndexesStruct struct {
	INDEXNAME string
	INDEXDEF  string
}

type ExchangeJob struct {
	JobID         string `json:"ИдентификаторЗадания"`
	ExchangeJobID string `json:"ИдентификаторЗапроса"`
	Area          string `json:"Область"`
	Event         string `json:"Событие"`
	Priod         string `json:"Дата"`
	Notes         string `json:"Заметки"`
}

type ExchangeJobAllInform struct {
	JobID         string `json:"ИдентификаторЗадания"`
	ExchangeJobID string `json:"ИдентификаторЗапроса"`
	Area          string `json:"Область"`
	Event         string `json:"Событие"`
	Priod         string `json:"Дата"`
	Base          string `json:"База"`
	Notes         string `json:"Заметки"`
}

type ExchangeJobGroup struct {
	Event string `json:"Событие"`
	Count string `json:"Количество"`
}

type ExchangeJobStatusProblem struct {
	Event string `json:"Событие"`
	Area  string `json:"Область"`
}

type AllJobs struct {
	JobID      string
	Status     string
	Priod      string
	JSONString string
}

type QueryToBISimpleID struct {
	JobID string `json:"ИдентификаторЗадания"`
	Time  string
}

type MessageNSQ struct {
	Type string `json:"Тип"`
	Body string `json:"Тело"`
}

type RemoteJob struct {
	JobID        string `json:"ИдентификаторЗадания"`
	RemoteBaseID string `json:"ИдентификаторУдаленнойБазы"`
	JobJSON      string `json:"ЗаданиеJSON"`
}

type EkisAreas struct {
	Ekis             string
	Inn              string
	ShortName        string
	FullName         string
	Mrsd             string
	AreaBgu          string
	AreaZkgu         string
	AddressBgu       string
	AddressZkgu      string
	BaseBgu          string
	BaseZkgu         string
	District         string
	Region           string
	Director         string
	Site             string
	MailDirector     string
	MailOrganization string
	Cgu              string
	Updated          string
	AreaUnf          string
	AddressUnf       string
	BaseUnf          string
}

type OrganizationRegistrationInformation struct {
	IDinDVS    int       `json:"id"`
	Name       string    `json:"name"`
	INN        string    `json:"inn"`
	Area       string    `json:"area"`
	UpdateData time.Time `json:"updated_at"`
}

type OrganizationRegistrationInformationMessage struct {
	Action  string                              `json:"action"`
	Payload OrganizationRegistrationInformation `json:"payload"`
}

type MessageWithPassport struct {
	Srvr              string
	Ref               string
	Pass              string
	Contur            string
	Usr               string
	Mail_timestamp    time.Time
	Date_from_subject time.Time
}

type RowPass struct {
	Date           string `json:"date"`
	Address_server string `json:"address_server"`
	Stage          string `json:"stage"`
	Db_name        string `json:"db_name"`
	Db_username    string `json:"db_username"`
	Db_userpwd     string `json:"db_userpwd"`
}
