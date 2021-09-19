package models

type SettingsJobs struct {
	JobID      string `json:"ИдентификаторЗадания"`
	JSONString string `json:"JSONСтрокаНастроек"`
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
