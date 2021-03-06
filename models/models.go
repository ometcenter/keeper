package models

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	log "github.com/ometcenter/keeper/logging"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

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

// Job структура задания
type Job struct {
	gorm.Model
	JobID  string `json:"ИдентификаторЗадания"`
	Status string `json:"Состояние"`
	Priod  string `json:"Дата"`
}

func (J *Job) GetJobStatus(DB *sql.DB) error {

	var argsquery []interface{}
	argsquery = append(argsquery, J.JobID)

	// queryAllColumns := `SELECT status, priod from jobs where job_id = $1`
	// rows, err := DB.Query(queryAllColumns, argsquery...)
	// if err != nil {
	// 	return err
	// }

	// defer rows.Close()

	// Job := Job{}
	// for rows.Next() {
	// 	err = rows.Scan(&Job.Status, &Job.Priod)
	// 	if err != nil {
	// 		return nil
	// 	}
	// }

	// *J = Job

	var NullTimeCreatedAt sql.NullTime
	var NullTimeUpdatedAt sql.NullTime

	var Job Job
	err := DB.QueryRow("SELECT id, created_at, updated_at, deleted_at, job_id, status, priod FROM jobs WHERE job_id = $1", argsquery...).Scan(&Job.ID, &NullTimeCreatedAt,
		&NullTimeUpdatedAt, &Job.DeletedAt, &Job.JobID, &Job.Status, &Job.Priod)
	if err != nil {
		return err
	}

	if NullTimeCreatedAt.Valid {
		Job.CreatedAt = NullTimeCreatedAt.Time
	}

	if NullTimeUpdatedAt.Valid {
		Job.UpdatedAt = NullTimeUpdatedAt.Time
	}

	*J = Job

	return nil
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

type JobsCronBackGroundTask struct {
	ID   int
	Job  string
	Next time.Time
	Prev time.Time
}

type QuantityMetric struct {
	ID             int
	DateMetric     time.Time // Дата метрики
	Area           string    // Область
	TableName      string    // Имя таблицы
	DataBaseID     string    // Идентификатор базы данных
	Value          int       // Значение метрики
	Hash           int64     // Строка хеш суммы
	SizeBody       int       // Размер сообщения в байтах
	SpeedUnzipping float64   // time.Duration Скорость распаковки в секундах
	SaveSpeed      float64   // time.Duration Скорость сохранения пакета в базу в секундах
	CountRecords   int       // Количество строк в запросе
	DateBeginQuery string    // Дата/время начала выполнения запроса в 1С
	DataEndQuery   string    // Дата/время окончания запроса в 1С
	DataSendQuery  string    // Дата/время отправки ответа в 1С.
}

type DeleteDataForArea struct {
	JobID      string
	TableName  string
	Area       string
	DataBaseID string
}

type ErrorOnBI struct {
	JobID            string `json:"УИД_Пакета"`
	ExchangeJobID    string `json:"УИД"`
	Area             string `json:"НомерОбласти"`
	EmptyQuery       bool   `json:"ПустойЗапрос"`
	ErrorDescription bool   `json:"ОшибкаВыполнения"`
	ResultQueryJSON  string `json:"РезультатЗапроса"`
}

type AllAreasSourses struct {
	gorm.Model
	Area                string
	ExternalID          string
	ShortName           string
	FullName            string
	INN                 string
	TypeSource          string
	BaseURL             string
	BaseName            string
	Notes               string
	AdditionInformation datatypes.JSON
}

type FileAndBinary struct {
	gorm.Model
	Name      string
	Type      string
	TextStore string
	ByteStore []byte `gorm:"type:bytea"`
	Notes     string
}

type RESTRequestUniversal struct {
	Body           []byte
	UrlToCall      string
	Method         string
	Headers        map[string]string
	UseAuth        bool
	AuthUserName   string
	AuthPassword   string
	TimeoutSeconds time.Duration
}

func (requestUniversal *RESTRequestUniversal) Send() ([]byte, error) {

	body := bytes.NewBuffer([]byte(requestUniversal.Body))

	req, err := http.NewRequest(requestUniversal.Method, requestUniversal.UrlToCall, body)
	req.Header.Set("Content-Type", "application/json")
	if requestUniversal.UseAuth {
		req.SetBasicAuth(requestUniversal.AuthUserName, requestUniversal.AuthPassword)
	}
	//if useAuth {
	//req.Header.Set("Authorization", "Basic "+requestUniversal.Headers["TokenBearer"])
	//}

	for key, value := range requestUniversal.Headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	if requestUniversal.TimeoutSeconds != 0 {
		client.Timeout = requestUniversal.TimeoutSeconds * time.Second
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		//return errors.New("status code not ok: " + strconv.Itoa(resp.StatusCode))
		return nil, err
	}

	bodyRespons, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		errReturn := errors.New("status code not ok: " + strconv.Itoa(resp.StatusCode) +
			" service: " + requestUniversal.UrlToCall)
		return bodyRespons, errReturn
	}

	return bodyRespons, nil

}

// TODO: Выкидываем все что можем взять из настроек
type QueryResult struct {
	Area           int             `json:"НомерОбласти"`
	ResultRequest  json.RawMessage `json:"РезультатЗапроса"`
	ErrorExecution bool            `json:"ОшибкаВыполнения"`
	EmptyRequest   bool            `json:"ПустойЗапрос"`
	ExchangeJobID  string          `json:"УИД"`
	JobID          string          `json:"УИД_Пакета"`
	// TODO: Возможно логичнее сжатые данные складывать в отдельное поле, чтобы не было путанницы извращения
	ResultRequestBase64          string                   `json:"РезультатЗапросаBase64"`
	HashSum                      int64                    `json:"ХешСумма"`
	SizeBody                     int                      // Размера тела в байтах, вычисляется при вычитывании сообщения из шины.
	Metrics                      Metrics                  `json:"Метрики"`
	MatureData                   []map[string]interface{} //[]map[string]interface{}
	Settings                     QueryToBI
	ElapsedSpeedUnzippingFloat64 float64
	lenPayload                   int
	BeginTimeMetric              time.Time
	ElapsedSaveSpeed             time.Duration
}

type Metrics struct {
	CountRecords   int    `json:"КоличествоЗаписей"`
	DateBeginQuery string `json:"ДатаНачалаЗапроса"`
	DataEndQuery   string `json:"ДатаОкончанияЗапроса"`
	DataSendQuery  string `json:"ДатаОтправкиОтвета"`
}

// TODO: Старый вариант переделать
type ChangingStatusJob struct {
	Priod string                `json:"Дата"`
	Event string                `json:"Событие"`
	Date  DataChangingStatusJob `json:"Данные"`
}

// TODO: Старый вариант переделать
type DataChangingStatusJob struct {
	JobID         string   `json:"ИдентификаторЗадания"`
	ExchangeJobID string   `json:"ИдентификаторЗапроса"`
	Status        string   `json:"Состояние"`
	Areas         []string `json:"Области"`
}

type RequestHistoryAPI struct {
	gorm.Model
	User   string
	Method string
	Amount int
}

func (RequestHistoryAPI *RequestHistoryAPI) CheckLimit(DB *sql.DB, firstday, lastday time.Time) (int, error) {

	var argsquery []interface{}
	argsquery = append(argsquery, RequestHistoryAPI.User)
	argsquery = append(argsquery, RequestHistoryAPI.Method)
	argsquery = append(argsquery, firstday)
	argsquery = append(argsquery, lastday)

	//fmt.Println(firstday, " - ", lastday)

	queryText := `select
		count(request_history_apis.amount) as amount
	from
		request_history_apis as request_history_apis
	where
		"user" = $1 and method = $2 and created_at >=$3 and created_at <=$4`

	rows, err := DB.Query(queryText, argsquery...)
	if err != nil {
		return 0, err
	}

	defer rows.Close()

	var amount int
	for rows.Next() {
		err = rows.Scan(&amount)
		if err != nil {
			return 0, err
		}

	}

	return amount, nil

}

func (RequestHistoryAPI *RequestHistoryAPI) AddNewRecordLimit(DB *sql.DB) error {

	var argsInsert []interface{}
	argsInsert = append(argsInsert, RequestHistoryAPI.Method)
	argsInsert = append(argsInsert, RequestHistoryAPI.User)
	argsInsert = append(argsInsert, 1)
	argsInsert = append(argsInsert, time.Now())

	_, err := DB.Exec(`INSERT INTO request_history_apis (method, "user", amount, created_at)
		VALUES($1, $2, $3, $4);`, argsInsert...)

	if err != nil {
		log.Impl.Error(err.Error())
		return err
	}

	return nil

}
