package models

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

type SettingsJobs struct {
	JobID      string `json:"ИдентификаторЗадания"`
	JSONString string `json:"JSONСтрокаНастроек"`
}

type AllJobs struct {
	JobID      string
	Status     string
	Priod      string
	JSONString string
}

type QueryToBISimpleID struct {
	JobID string `json:"ИдентификаторЗадания,JobID"`
	Time  string
}

type SettingsJobSliceQueryToBI struct {
	JobID          string      `json:"ИдентификаторЗадания"`
	SliceQueryToBI []QueryToBI `json:"Настройки"`
}

type QueryToBI struct {
	JobID                   string        `json:"ИдентификаторЗадания"`
	Portions                int           `json:"Порции"`
	Query                   []Query       `json:"Запросы"`
	AddParam                AdditionParam `json:"ДополнительныеПараметрыJSON"`
	AddParamJSNOString      string        `json:"JSONСтрокаДополнительныеПараметры"`
	Connect                 Connect       `json:"ПараметрыПодключения"`
	ConnectContur           ConnectContur `json:"ПараметрыПодключенияКонтура"`
	ConnectBI1C             ConnectBI1C   `json:"ПараметрыПодключенияBI1C"`
	ConnectConturJSNOString string        `json:"JSONСтрокаПараметрыПодключенияКонтура"`
}

type Query struct {
	QueryText                            string `json:"Запрос"`
	Base                                 string `json:"База"`
	Area                                 string `json:"Области"`
	ExchangeJobID                        string `json:"ИдентификаторЗапроса"`
	PText                                string `json:"ПараметрыЗапроса"`
	UsedCalculatedFieldsInQueryParametrs bool   `json:"ИспользуетсяВычисляемыеПоляВПараметрахЗапроса"`
}

type AdditionParam struct {
	ZipAnswer  bool     `json:"СжиматьОтвет"`
	HashAnswer bool     `json:"ХешироватьРезультат"`
	Options    *Options `json:"НастройкиМоделиДанных"`
	Connect    Connect  `json:"ПараметрыПодключенияHTTPОтвета"`
}

type Connect struct {
	AddressServer   string `json:"АдресСервиса"`
	PortServer      int    `json:"Порт"`
	Resource        string `json:"Ресурс"`
	SecureConnetion bool   `json:"ЗащищенноеСоединение"`
	LoginConnetion  string `json:"Логин"`
	Password        string `json:"Пароль"`
	Headers         string `json:"Заголовки"`
}

type ConnectBI1C struct {
	AddressServer   string `json:"АдресСервиса"`
	PortServer      int    `json:"Порт"`
	Resource        string `json:"Ресурс"`
	SecureConnetion bool   `json:"ЗащищенноеСоединение"`
	LoginConnetion  string `json:"Логин"`
	Password        string `json:"Пароль"`
	Headers         string `json:"Заголовки"`
}

type ConnectContur struct {
	AddressServer   string `json:"АдресСервисаПриемника"`
	PortServer      int    `json:"Порт"`
	Resource        string `json:"Ресурс"`
	SecureConnetion bool   `json:"ЗащищенноеСоединение"`
	LoginConnetion  string `json:"Логин"`
	Password        string `json:"Пароль"`
	Headers         string `json:"Заголовки"`
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

type MessageNSQ struct {
	Type string `json:"Тип"`
	Body string `json:"Тело"`
}

type RemoteJob struct {
	JobID        string `json:"ИдентификаторЗадания"`
	RemoteBaseID string `json:"ИдентификаторУдаленнойБазы"`
	JobJSON      string `json:"ЗаданиеJSON"`
}
