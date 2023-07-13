package models

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/ometcenter/keeper/config"
	log "github.com/ometcenter/keeper/logging"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type TableDescription struct {
	TableName string   `json:"ИмяТаблицы"`
	Fields    []Fields `json:"СведенияКолонкиТаблицы"`
}

// type IndexDescription struct {
// 	TableName string        `json:"ИмяТаблицы"`
// 	Fields    []FieldsIndex `json:"СведенияКолонкиТаблицы"`
// }

// type FieldsIndex struct {
// 	Name       string `json:"Имя"`
// 	Definition string `json:"Определение"`
// 	TypeChange string `json:"ИзменитьВСУБД"`
// }

// type Fields struct {
// 	Name       string `json:"Имя"`
// 	Type       string `json:"Тип"`
// 	NotNull    bool   `json:"NotNull"`
// 	PrimaryKey bool   `json:"ПервичныйКлюч"`
// 	TypeChange string `json:"ИзменитьВСУБД"`
// }

// type ColumnsStruct struct {
// 	ColumnName string
// 	DataType   string
// 	IsNullable string
// 	PrimaryKey bool
// }

type IndexesStruct struct {
	INDEXNAME string
	INDEXDEF  string
}

type CheckVersionFile struct {
	CheckResult   bool   `json:"checkResult"`
	LatestVersion string `json:"latestVersion"`
	Notes         string `json:"notes"`
}

type ExchangeJob struct {
	JobID         string `json:"ИдентификаторЗадания"`
	ExchangeJobID string `json:"ИдентификаторЗапроса"`
	Area          string `json:"Область"`
	Event         string `json:"Событие"`
	Priod         string `json:"Дата"`
	Notes         string `json:"Заметки"`
}

func (E *ExchangeJob) SendStatusCreateExchangeJobIDThroughREST(UrlToCall string) error {

	ExchangeJobByte, err := json.Marshal(E)
	if err != nil {
		//cCp.JSON(http.StatusBadRequest, err.Error())
		//log.Impl.Error(fmt.Errorf("Ошибка по области %s задание %s: Текст ошибки --- %s\n", Area, JobIDParam, err.Error()))
		return err
	}

	var RESTRequestUniversal3 RESTRequestUniversal
	Headers := make(map[string]string)
	Headers["TokenBearer"] = config.Conf.TokenBearer
	RESTRequestUniversal3.Headers = Headers
	RESTRequestUniversal3.Method = "POST"
	RESTRequestUniversal3.Body = ExchangeJobByte
	// TODO: Переделать на переменную окружения
	RESTRequestUniversal3.UrlToCall = UrlToCall //os.Getenv("ADDRESS_PORT_SERVICE_FRONT") + "/changingstatussimple"
	_, err = RESTRequestUniversal3.Send()
	if err != nil {
		return err
		//log.Impl.Error(fmt.Errorf("Ошибка по области %s задание %s: Текст ошибки --- %s\n", Area, JobIDParam, err.Error()))
	}

	return nil

}

type ExchangeJobV2 struct {
	JobID         string `json:"jobID"`
	ExchangeJobID string `json:"exchangeJobID"`
	Area          string `json:"area"`
	Event         string `json:"status"`
	Priod         string `json:"period"`
	Notes         string `json:"notes"`
}

func (E *ExchangeJobV2) SendStatusCreateExchangeJobIDThroughREST(UrlToCall, TokenBearer string) error {

	ExchangeJobByte, err := json.Marshal(E)
	if err != nil {
		return err
	}

	var RESTRequestUniversal3 RESTRequestUniversal
	Headers := make(map[string]string)
	Headers["TokenBearer"] = TokenBearer
	RESTRequestUniversal3.Headers = Headers
	RESTRequestUniversal3.Method = "POST"
	RESTRequestUniversal3.Body = ExchangeJobByte
	RESTRequestUniversal3.UrlToCall = UrlToCall //os.Getenv("ADDRESS_PORT_SERVICE_FRONT") + "/changingstatussimple"
	_, err = RESTRequestUniversal3.Send()
	if err != nil {
		return err
	}

	return nil
}

func (E *ExchangeJobV2) SaveDirectSQL(DB *sql.DB) error {

	area := string(E.Area)
	if area == "" {
		area = "0"
	}

	var argsquery []interface{}
	argsquery = append(argsquery, E.JobID)
	argsquery = append(argsquery, E.ExchangeJobID)
	argsquery = append(argsquery, area)

	query := `SELECT * FROM exchange_jobs WHERE job_id = $1 AND exchange_job_id = $2 AND area = $3`

	rows, err := DB.Query(query, argsquery...)
	if err != nil {
		return err
	}

	defer rows.Close()

	flag := false
	for rows.Next() {
		flag = true
		break
	}

	if flag == true {

		var argsUpdate []interface{}
		argsUpdate = append(argsUpdate, E.JobID)
		argsUpdate = append(argsUpdate, E.ExchangeJobID)
		argsUpdate = append(argsUpdate, area)
		argsUpdate = append(argsUpdate, E.Event)
		argsUpdate = append(argsUpdate, E.Priod)
		argsUpdate = append(argsUpdate, E.Notes)

		_, err := DB.Exec(`UPDATE exchange_jobs SET job_id=$1, exchange_job_id=$2, area=$3, "event"=$4, priod=$5, notes=$6
		WHERE job_id = $1 AND exchange_job_id = $2 AND area = $3;`, argsUpdate...)
		if err != nil {
			return err
		}

		// `UPDATE exchange_jobs SET created_at='', updated_at='', deleted_at='', job_id='', exchange_job_id='', area='', "event"='', priod=''
		// WHERE id=nextval('exchange_jobs_id_seq'::regclass);`

	} else {

		//  `INSERT INTO exchange_jobs (job_id, exchange_job_id, area, "event", priod)
		//  VALUES('$1, $2, $3, $4, $5);`

		var argsInsert []interface{}
		argsInsert = append(argsInsert, E.JobID)
		argsInsert = append(argsInsert, E.ExchangeJobID)
		argsInsert = append(argsInsert, area)
		argsInsert = append(argsInsert, E.Event)
		argsInsert = append(argsInsert, E.Priod)
		argsInsert = append(argsInsert, E.Notes)

		_, err := DB.Exec(`INSERT INTO exchange_jobs (job_id, exchange_job_id, area, event, priod, notes)
		VALUES($1, $2, $3, $4, $5, $6);`, argsInsert...)

		if err != nil {
			return err
		}

	}

	return nil
}

// func (E *models.ExchangeJobV2) SaveToRabbitMQ(ConnectRabbitMQ *amqp.Connection) error {
// 	var MessageQueueGeneralInterface models.MessageQueueGeneralInterface
// 	MessageQueueGeneralInterface.Type = "ChangeStatusForExchangeJobV2"
// 	MessageQueueGeneralInterface.Body = E

// 	err := SendInQueueRabbitMQUniversal(MessageQueueGeneralInterface.Type, MessageQueueGeneralInterface.Body, "go-keeper-status", ConnectRabbitMQ)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

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

type JobV2 struct {
	JobID  string `json:"jobID"`
	Status string `json:"status"`
	Priod  string `json:"period"`
	Notes  string `json:"notes"`
}

func (J *JobV2) SaveDirectSQL(DB *sql.DB) error {

	var argsUpdate []interface{}
	argsUpdate = append(argsUpdate, J.JobID)
	argsUpdate = append(argsUpdate, J.Status)
	argsUpdate = append(argsUpdate, time.Now().Format("2006-01-02T15:04:05"))

	result, err := DB.Exec(`UPDATE jobs SET status=$2, priod=$3
		WHERE job_id = $1;`, argsUpdate...)
	if err != nil {
		return err
	}

	//LastInsertId, _ := result.LastInsertId()
	RowsAffected, _ := result.RowsAffected()

	//fmt.Println("LastInsertId: ", LastInsertId)
	//fmt.Println("RowsAffected: ", RowsAffected)

	// Если не обновленно не одной записи, значит это новая запись и ее надо добавить
	if RowsAffected == 0 {

		var argsInsert []interface{}
		argsInsert = append(argsInsert, J.JobID)
		argsInsert = append(argsInsert, J.Status)
		argsInsert = append(argsInsert, time.Now().Format("2006-01-02T15:04:05"))

		_, err := DB.Exec(`INSERT INTO jobs (job_id, status, priod)
			VALUES($1, $2, $3);`, argsInsert...)

		if err != nil {
			return err
		}

	}

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

type JobDescription struct {
	Name1C    string
	Code1C    string
	TableName string
}

type MessageTelegram struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

type QuantityMetric struct {
	ID             int
	DateMetric     time.Time // Дата метрики
	Area           string    // Область
	TableName      string    // Имя таблицы
	DataBaseID     string    // Идентификатор базы данных
	Value          int       // Значение метрики
	Hash           string    // Строка хеш суммы
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

type ErrorOrEmptyQuery struct {
	JobID            string `json:"jobID"`
	ExchangeJobID    string `json:"exchangeJobID"`
	Area             string `json:"area"`
	EmptyQuery       bool   `json:"emptyQuery"`
	ErrorDescription string `json:"errorDescription"`
}

func (E *ErrorOrEmptyQuery) SendResultThroughREST(TokenBearer, UrlToCall string) error {

	byteResult, err := json.Marshal(*E)
	if err != nil {
		return err
	}

	var RESTRequestUniversal2 RESTRequestUniversal
	Headers := make(map[string]string)
	Headers["TokenBearer"] = TokenBearer
	RESTRequestUniversal2.Headers = Headers
	RESTRequestUniversal2.Method = "POST"
	RESTRequestUniversal2.Body = byteResult
	RESTRequestUniversal2.UrlToCall = UrlToCall

	//for i := 0; i < 1000; i++ {

	_, err = RESTRequestUniversal2.Send()
	if err != nil {
		return err
	}
	//}

	return nil
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

	// if resp.StatusCode != http.StatusOK {
	// 	//return errors.New("status code not ok: " + strconv.Itoa(resp.StatusCode))
	// 	return nil, err
	// }

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

// TODO: Перевести на type DataToETL struct {
type QueryResultShort struct {
	Area                int                      `json:"НомерОбласти"`
	ResultRequest       []map[string]interface{} `json:"РезультатЗапроса"`
	ErrorExecution      bool                     `json:"ОшибкаВыполнения"`
	EmptyRequest        bool                     `json:"ПустойЗапрос"`
	ExchangeJobID       string                   `json:"УИД"`
	JobID               string                   `json:"УИД_Пакета"`
	ResultRequestBase64 string                   `json:"РезультатЗапросаBase64"`
	HashSum             int64                    `json:"ХешСумма"`
	//Metrics                      Metrics                  `json:"Метрики"`
	CleaningFieldsBeforeLoading []CleaningFieldsBeforeLoading `json:"cleaningFieldsBeforeLoading"`
}

func (Q *QueryResultShort) ZipAnswerGzip() error {

	byteValue, err := json.Marshal(Q.ResultRequest)
	if err != nil {
		return err
	}

	//fmt.Println("Json: ", string(byteValue))

	var buf bytes.Buffer
	g := gzip.NewWriter(&buf)
	if _, err = g.Write(byteValue); err != nil {
		return err
	}
	if err = g.Close(); err != nil {
		return err
	}

	sDec := base64.StdEncoding.EncodeToString(buf.Bytes())
	if err != nil {
		return err
	}

	//fmt.Println(sDec)

	//QueryResult.ResultRequest = nil
	var NilMap []map[string]interface{}
	Q.ResultRequest = NilMap
	Q.ResultRequestBase64 = sDec

	return nil

}

func (Q *QueryResultShort) SendResultThroughREST(TableName, UrlToCall string) error {

	var QueryResultSlice []QueryResultShort //[]modelsShare.QueryResultShort

	QueryResultSlice = append(QueryResultSlice, *Q)
	byteResult, err := json.Marshal(QueryResultSlice)
	if err != nil {
		return err
	}

	AreaString := strconv.Itoa(Q.Area)

	var RESTRequestUniversal2 RESTRequestUniversal
	Headers := make(map[string]string)
	Headers["TokenBearer"] = config.Conf.TokenBearer
	Headers["JobID"] = Q.JobID
	Headers["Area"] = AreaString
	Headers["TableName"] = TableName
	Headers["ExchangeJobID"] = Q.ExchangeJobID
	//layoutISO := "2006-01-02"
	//Headers["NoteCOD"] = fmt.Sprintf(" - статут в ЦОД: дата начала = %s, дата окончания = %s", item.DataStart.Format(layoutISO), item.DataEnd.Format(layoutISO))
	//fmt.Println(Headers["NoteCOD"])
	//fmt.Printf("%v\n", item)

	RESTRequestUniversal2.Headers = Headers
	RESTRequestUniversal2.Method = "POST"
	RESTRequestUniversal2.Body = byteResult
	RESTRequestUniversal2.UrlToCall = UrlToCall

	//for i := 0; i < 1000; i++ {

	_, err = RESTRequestUniversal2.Send()
	if err != nil {
		return err
	}
	//}

	return nil
}

type DataToETL struct {
	Area                        int                           `json:"area"`
	Data                        []map[string]interface{}      `json:"data"`
	ExchangeJobID               string                        `json:"exchangeJobID"`
	JobID                       string                        `json:"jobID"`
	DataBase64                  string                        `json:"dataBase64"`
	HashSum                     string                        `json:"hashSum"`
	CleaningFieldsBeforeLoading []CleaningFieldsBeforeLoading `json:"cleaningFieldsBeforeLoading"`
}

func (D *DataToETL) SendResultThroughREST(TokenBearer, UrlToCall string) error {

	byteResult, err := json.Marshal(*D)
	if err != nil {
		return err
	}

	var RESTRequestUniversal2 RESTRequestUniversal
	Headers := make(map[string]string)
	Headers["TokenBearer"] = TokenBearer
	Headers["JobID"] = D.JobID
	Headers["Area"] = D.GetAreaString()
	//Headers["TableName"] = TableName
	Headers["ExchangeJobID"] = D.ExchangeJobID

	RESTRequestUniversal2.Headers = Headers
	RESTRequestUniversal2.Method = "POST"
	RESTRequestUniversal2.Body = byteResult
	RESTRequestUniversal2.UrlToCall = UrlToCall

	//for i := 0; i < 1000; i++ {

	_, err = RESTRequestUniversal2.Send()
	if err != nil {
		return err
	}
	//}

	return nil
}

func (d *DataToETL) ZipAnswerGzip() error {

	byteValue, err := json.Marshal(d.Data)
	if err != nil {
		return err
	}

	//fmt.Println("Json: ", string(byteValue))

	var buf bytes.Buffer
	g := gzip.NewWriter(&buf)
	if _, err = g.Write(byteValue); err != nil {
		return err
	}
	if err = g.Close(); err != nil {
		return err
	}

	sDec := base64.StdEncoding.EncodeToString(buf.Bytes())
	if err != nil {
		return err
	}

	//fmt.Println(sDec)

	//QueryResult.ResultRequest = nil
	var NilMap []map[string]interface{}
	d.Data = NilMap
	d.DataBase64 = sDec

	return nil

}

func (d *DataToETL) HashDataSha256() error {

	byteValue, err := json.Marshal(d.Data)
	if err != nil {
		return err
	}

	h := sha256.New()
	h.Write(byteValue)
	// Calculate and print the hash
	d.HashSum = fmt.Sprintf("%x", h.Sum(nil))

	return nil

}

func (d *DataToETL) GetAreaString() string {
	AreaString := strconv.Itoa(d.Area)
	return AreaString
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

type AllAreaFromCOD struct {
	Status         string           `json:"status"`
	Date           string           `json:"date"`
	Query          string           `json:"query"`
	Number         string           `json:"number"`
	ID             string           `json:"id"`
	RegionsFromCOD []RegionsFromCOD `json:"regions"`
}

type RegionsFromCOD struct {
	Error  bool   `json:"error"`
	Empty  bool   `json:"empty"`
	Status string `json:"status"`
	Code   int    `json:"code"`
	// DataStart time.Time "date_start"
	// DataEnd   time.Time "date_end"
	DataStart string "date_start"
	DataEnd   string "date_end"
}

type ChangeStatusJobSimple struct {
	JobID  string `json:"jobID"`
	Status string `json:"status"`
}

type RemoteJobs struct {
	JobId        string
	RemoteBaseId string
}

type HistoryReceivedMessages struct {
	gorm.Model
	Area          string `json:"Область"`
	TableName     string `json:"ИмяТаблицы"`
	DateRecord    time.Time
	MessageResult datatypes.JSON
	Settings      datatypes.JSON
}
