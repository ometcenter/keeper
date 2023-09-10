package utility

import (
	"bytes"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/ometcenter/keeper/config"
	log "github.com/ometcenter/keeper/logging"
	"github.com/ometcenter/keeper/models"
	queue "github.com/ometcenter/keeper/queue"
	tracing "github.com/ometcenter/keeper/tracing/jaeger"
	tracingRabbitMQ "github.com/ometcenter/keeper/tracing/jaeger/rabbitmq"
	"github.com/streadway/amqp"
)

func GetAreasByStasus(DB *sql.DB, JobID, Stasus string) ([]string, error) {

	var argsquery []interface{}
	argsquery = append(argsquery, JobID)
	argsquery = append(argsquery, Stasus)

	QueryText := `select
		area
	from
		public.exchange_jobs
	where
		job_id = $1
		and "event" = $2;`

	rows, err := DB.Query(QueryText, argsquery...)
	if err != nil {
		return nil, err
	}

	var AreasForReturn []string
	for rows.Next() {
		var area string
		err = rows.Scan(&area)
		if err != nil {
			return nil, err
		}
		AreasForReturn = append(AreasForReturn, area)
	}

	defer rows.Close()

	return AreasForReturn, nil

}

func GetAreasNotEqualToStatus(DB *sql.DB, JobID, Stasus string) ([]string, error) {

	var argsquery []interface{}
	argsquery = append(argsquery, JobID)
	argsquery = append(argsquery, Stasus)

	QueryText := `SELECT area
	FROM public.exchange_jobs where job_id = $1 and "event" <> $2;`

	rows, err := DB.Query(QueryText, argsquery...)
	if err != nil {
		return nil, err
	}

	var AreasForReturn []string
	for rows.Next() {
		var area string
		err = rows.Scan(&area)
		if err != nil {
			return nil, err
		}
		AreasForReturn = append(AreasForReturn, area)
	}

	defer rows.Close()

	return AreasForReturn, nil

}

func SendInQueueRabbitMQUniversal(TypeMessage string, DataStruct interface{},
	topicName string, ConnectRabbitMQ *amqp.Connection) error { // RabbitMQchannelMessage *amqp.Channel

	// TODO: А так же вариант прямой передачи через тип интерфейс Переделать под общую структуру сообщения.
	MessageQueueGeneralInterface := models.MessageQueueGeneralInterface{Type: TypeMessage, Body: DataStruct}

	JsonMessageBody, err := json.Marshal(&MessageQueueGeneralInterface)
	if err != nil {
		log.Impl.Errorf("ошибка маршалинга: %s", err)
		return err
	}

	// if RabbitMQchannel == nil {
	// 	err := errors.New("Connection to RabbitMQ not established")
	// 	return err
	// }

	// TODO: При сохранении сообщений возникает ошибка "Exception (505) Reason: "UNEXPECTED_FRAME - expected content body, got non content body frame instead""
	// я создавал два канала для глобальной переменной для Сообщений и Событий, но ошибка сохранилась.
	// возможно попробовать на другой докере не моем развернуть шину, делаю для каждого сообщения отдельный канал
	// так он Gin видимо пораждает отдельную горутину для каждого вызова конечной точки, один канал под одну горутину.
	// Создаю глобальную переменную ConnectRabbitMQ
	RabbitMQchannel, err := ConnectRabbitMQ.Channel()
	if err != nil {
		//Connector.LoggerConn.ErrorLogger.Println("Failed to open a channel")
		//Connector.RabbitMQchannel = nil
		log.Impl.Error("Failed to open a channel :", err)
		log.Impl.Panic()
		panic("Failed to open a channel")
	}

	// Закрываем канал чтобы избежать ошибки Exception (504) Reason: "channel id space exhausted" с наплождением
	// слишком много количества каналов
	defer RabbitMQchannel.Close()

	q, err := RabbitMQchannel.QueueDeclare(
		topicName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return err
	}

	amqpPublishing := amqp.Publishing{
		ContentType: "text/plain",
		Body:        JsonMessageBody,
	}

	if config.Conf.UseTracing {
		// Сделанно для инициализации карты, с nil валиться в паник, внутри библиотеки не предусмотрели.
		amqpPublishing.Headers = map[string]interface{}{}
		//HeaderMaps["Test1"] = "Test1"

		clientSpan := tracing.Tracer.StartSpan("clientspan")
		defer clientSpan.Finish()

		if err := tracingRabbitMQ.Inject(clientSpan, amqpPublishing.Headers); err != nil {
			log.Impl.Error(err)
		}

		//amqpPublishing.Headers["JobID"] = ej.JobID
		//amqpPublishing.Headers["ExchangeJobID"] = ej.ExchangeJobID
		//amqpPublishing.Headers["Area"] = ej.Area

		err = RabbitMQchannel.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqpPublishing)
		if err != nil {
			return err
		}

	} else {
		err = RabbitMQchannel.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqpPublishing)
		if err != nil {
			return err
		}
	}

	return nil

}

// Изменяет статус всего задания "Выполняется"
func ChangeStatusJobsTask(DB *sql.DB, JobID, Status string) error {

	var argsUpdate []interface{}
	argsUpdate = append(argsUpdate, JobID)
	argsUpdate = append(argsUpdate, Status)
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
		argsInsert = append(argsInsert, JobID)
		argsInsert = append(argsInsert, Status)
		argsInsert = append(argsInsert, time.Now().Format("2006-01-02T15:04:05"))

		_, err := DB.Exec(`INSERT INTO jobs (job_id, status, priod)
		VALUES($1, $2, $3);`, argsInsert...)

		if err != nil {
			return err
		}

	}

	return nil

}

// Удалить все задание
func DeleteJobs(DB *sql.DB, JobID string) error {

	var argsquery []interface{}
	argsquery = append(argsquery, JobID)

	QueryString := `DELETE FROM public.exchange_jobs where job_id = $1`

	_, err := DB.Exec(QueryString, argsquery...)
	if err != nil {
		return err
	}

	return nil

}

func GetAllDataFromTables(DB *sql.DB, TableNameParam string, mapAvailableTables map[string]bool, QueryURL url.Values) ([]map[string]interface{}, error) {

	IsAvailableTables, ok := mapAvailableTables[TableNameParam]
	if !ok || !IsAvailableTables {
		err := fmt.Errorf("Problems with getting data")
		return nil, err
	}
	// TODO: Вставить защиту от SQL иньекций, например проверкой таблицы
	//var queryText = `select * from ` + TableName + `;`

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	queryBuilder := psql.Select("*").From(TableNameParam)

	PaginatioRegim := false
	var Limit, Offset int
	//Page := c.Query("page")
	//PerPage := c.Query("per_page")
	Page := QueryURL.Get("page")
	PerPage := QueryURL.Get("per_page")

	if Page != "" && PerPage != "" {
		PaginatioRegim = true
		PageInt, err := strconv.Atoi(Page)
		if err != nil {
			return nil, err
		}

		PerPageInt, err := strconv.Atoi(PerPage)
		if err != nil {
			return nil, err
		}

		Limit = PerPageInt
		Offset = PageInt*PerPageInt - PerPageInt
	}

	getCountOnFieldForPagination := QueryURL.Get("getCountOnFieldForPagination")

	daydeepFilter := QueryURL.Get("daysDeep")
	var daydeep int
	if daydeepFilter == "" {
		daydeep = 2
	} else {
		var err error
		daydeep, err = strconv.Atoi(daydeepFilter)
		if err != nil {
			return nil, err
		}
	}

	dataTimeDeep := QueryURL.Get("dataTimeDeep")

	//fmt.Printf("PaginatioRegim: %t Limit: %d Offset: %d getCountOnFieldForPagination: %s\n", PaginatioRegim, Limit, Offset, getCountOnFieldForPagination)

	param := []interface{}{}

	if daydeepFilter != "" {

		timeNow := time.Now()
		timePast := timeNow.AddDate(0, 0, -daydeep)

		//where
		//updated_at > $1
		//argsquery = append(argsquery, timePast)

		if PaginatioRegim {
			//queryBuilder = psql.Select("*").From(TableNameParam).Where(sq.GtOrEq{"updated_at": timePast}).Limit(uint64(Limit)).Offset(uint64(Offset)).OrderBy("created_at")
			queryBuilder = psql.Select("*").From(TableNameParam).Where(
				sq.Or{sq.GtOrEq{"updated_at": timePast}, sq.GtOrEq{"deleted_at": timePast}, sq.GtOrEq{"created_at": timePast}}).Limit(
				uint64(Limit)).Offset(uint64(Offset)).OrderBy("id ASC")
		} else {
			//queryBuilder = psql.Select("*").From(TableNameParam).Where(sq.GtOrEq{"updated_at": timePast})
			queryBuilder = psql.Select("*").From(TableNameParam).Where(sq.Or{sq.GtOrEq{"updated_at": timePast}, sq.GtOrEq{"deleted_at": timePast}, sq.GtOrEq{"created_at": timePast}})
		}
		param = append(param, timePast)
		param = append(param, timePast)
		param = append(param, timePast)

	} else if dataTimeDeep != "" {

		// timePast, err := time.Parse("2006-01-02T15:04:05", dataTimeDeep)
		// if err != nil {
		// 	return nil, err
		// }

		timePast, err := time.ParseInLocation("2006-01-02T15:04:05", dataTimeDeep, time.Local)
		if err != nil {
			return nil, err
		}

		fmt.Println("timePast - ", timePast)

		if PaginatioRegim {
			//queryBuilder = psql.Select("*").From(TableNameParam).Where(sq.GtOrEq{"updated_at": timePast}).Limit(uint64(Limit)).Offset(uint64(Offset)).OrderBy("created_at")
			queryBuilder = psql.Select("*").From(TableNameParam).Where(
				sq.Or{sq.GtOrEq{"updated_at": timePast}, sq.GtOrEq{"deleted_at": timePast}, sq.GtOrEq{"created_at": timePast}}).Limit(
				uint64(Limit)).Offset(uint64(Offset)).OrderBy("id ASC")
		} else {
			//queryBuilder = psql.Select("*").From(TableNameParam).Where(sq.GtOrEq{"updated_at": timePast})
			queryBuilder = psql.Select("*").From(TableNameParam).Where(sq.Or{sq.GtOrEq{"updated_at": timePast}, sq.GtOrEq{"deleted_at": timePast}, sq.GtOrEq{"created_at": timePast}})
		}
		param = append(param, timePast)
		param = append(param, timePast)
		param = append(param, timePast)

		//return nil, errors.New("Функционал не реализован dataTimeDeep")

	} else if getCountOnFieldForPagination != "" {
		queryBuilder = psql.Select(getCountOnFieldForPagination, "count(1)").From(TableNameParam).GroupBy(getCountOnFieldForPagination).OrderBy(getCountOnFieldForPagination)
		//SELECT area, count(1) FROM collaborators_posle group by area;

	} else if PaginatioRegim {

		queryBuilder = psql.Select("*").From(TableNameParam).Limit(uint64(Limit)).Offset(uint64(Offset)).OrderBy("id ASC")

		if len(QueryURL) > 2 {
			//var conditionSlice []sq.Eq
			//conditionMap := make(map[string]interface{})
			//var conditionSlice sq.Eq

			sqEq := make(sq.Eq)

			for key, value := range QueryURL {

				if key == "total" || key == "page" || key == "per_page" {
					continue
				}
				//conditionSlice = append(conditionSlice, sq.Eq{key: value})
				//conditionMap[key] = value
				if len(value) == 0 {
					param = append(param, "")
				} else {
					param = append(param, value[0])
				}

				sqEq[key] = ""

			}
			queryBuilder = psql.Select("*").From(TableNameParam).Limit(uint64(Limit)).Offset(uint64(Offset)).OrderBy("id ASC").Where(sqEq)
		}

	} else if len(QueryURL) != 0 {
		//var conditionSlice []sq.Eq
		//conditionMap := make(map[string]interface{})
		//var conditionSlice sq.Eq

		sqEq := make(sq.Eq)

		for key, value := range QueryURL {

			if key == "total" {
				continue
			}
			//conditionSlice = append(conditionSlice, sq.Eq{key: value})
			//conditionMap[key] = value
			if len(value) == 0 {
				param = append(param, "")
			} else {
				param = append(param, value[0])
			}

			sqEq[key] = ""

		}
		queryBuilder = psql.Select("*").From(TableNameParam).Where(sqEq)
	}

	queryText, _, err := queryBuilder.ToSql()
	fmt.Println(queryText)
	if err != nil {
		return nil, err
	}

	rows, err := DB.Query(queryText, param...)
	if err != nil {
		return nil, err
	}
	cols, _ := rows.Columns()

	defer rows.Close()

	var MatureDataSlice []map[string]interface{}
	for rows.Next() {

		MatureData := make(map[string]interface{})
		// Create a slice of interface{}'s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}

		// Scan the result into the column pointers...
		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		// Create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		//m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			MatureData[colName] = *val
		}

		MatureDataSlice = append(MatureDataSlice, MatureData)

	}

	// byteResult, err := json.Marshal(MatureDataSlice)
	// if err != nil {
	// 	c.String(http.StatusBadRequest, "Response: %s", err.Error())
	// 	log.Impl.Error(err.Error())
	// 	return
	// }

	//fmt.Println(string(byteResult))
	//c.String(http.StatusOK, string(byteResult))

	//c.Data(http.StatusOK, "application/json", byteResult)

	//c.JSON(http.StatusOK, MatureDataSlice)

	return MatureDataSlice, nil

}

func InTimeSpan(start, end, check time.Time) bool {
	if start.Before(end) {
		return !check.Before(start) && !check.After(end)
	}
	if start.Equal(end) {
		return check.Equal(start)
	}
	return !start.After(check) || !end.Before(check)
}

func Weekday(d time.Weekday) int {
	day := (d - 1) % 7
	if day < 0 {
		day += 7
	}
	return int(day)
}

func StartOfWeek(t time.Time) time.Time {
	// Figure out number of days to back up until Mon:
	// Sun is 0 -> 6, Sat is 6 -> 5, etc.
	toMon := Weekday(t.Weekday())
	y, m, d := t.AddDate(0, 0, -int(toMon)).Date()
	// Result is 00:00:00 on that year, month, day.
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

// func СhunkSliceBy[T any](items []T, chunkSize int) (chunks [][]T) {
//     for chunkSize < len(items) {
//         items, chunks = items[chunkSize:], append(chunks, items[0:chunkSize:chunkSize])
//     }
//     return append(chunks, items)
// }

//Or if you want to manually set the capacity:
// func СhunkSliceBy[T any](items []T, chunkSize int) (chunks [][]T) {
// 	var _chunks = make([][]T, 0, (len(items)/chunkSize)+1)
// 	for chunkSize < len(items) {
// 		items, _chunks = items[chunkSize:], append(_chunks, items[0:chunkSize:chunkSize])
// 	}
// 	return append(_chunks, items)
// }

func ShortDur(d time.Duration) string {
	s := d.String()
	if strings.HasSuffix(s, "m0s") {
		s = s[:len(s)-2]
	}
	if strings.HasSuffix(s, "h0m") {
		s = s[:len(s)-2]
	}
	return s
}

// Если таблица вестит больше 500МБ, а это примерно 3 000 000 записей, то чистим за 6 месяцев метрики.
// 1 month --- 330275
func ShrinkTablesUniversal(DB *sql.DB, TableName string, CounterLimit int, DurationTimeRemaindRows time.Duration,
	DataFieldForCondition string, UseLimitOnly bool) error {

	// queryAllColumns := `SELECT count(*)
	// FROM public.quantity_metrics;`

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	queryBuilder := psql.Select("count(*)").From(TableName)

	queryText, _, err := queryBuilder.ToSql()
	fmt.Printf("ShrinkTables: %s query: %s\n", TableName, queryText)
	if err != nil {
		return err
	}

	var counter int
	err = DB.QueryRow(queryText).Scan(&counter)
	if err != nil {
		return err
	}

	//Если таблица вестит больше 500МБ, а это примерно 3 000 000 записей, то чистим за 6 месяцев метрики.
	//1 month --- 330275
	//TODO: Может быть брать из системных таблиц размер таблицы и тогда чистить? а не ориентироваться по количеству
	if counter > CounterLimit {

		if UseLimitOnly {

			queryTextDelete := fmt.Sprintf("Delete from %s where id in (Select id FROM %s ORDER BY %s DESC OFFSET %d)",
				TableName, TableName, DataFieldForCondition, CounterLimit)
			// queryBuilderDelete := psql.Delete("").From(TableName).Offset(uint64(CounterLimit)).OrderBy(DataFieldForCondition + " DESC")

			// queryTextDelete, _, err := queryBuilderDelete.ToSql()
			// //fmt.Println(queryTextDelete)
			fmt.Printf("ShrinkTables: %s - query: %s - counter %d\n", TableName, queryTextDelete, counter)
			// if err != nil {
			// 	return err
			// }
			_, err = DB.Exec(queryTextDelete)
			if err != nil {
				return err
			}

		} else {
			now := time.Now()
			//after := now.Add(-25 * time.Hour)
			after := now.Add(-1 * DurationTimeRemaindRows)
			//after := now.AddDate(0, -6, 0)
			//fmt.Println("Subtract:", after)

			queryBuilderDelete := psql.Delete("").From(TableName).Where(sq.LtOrEq{DataFieldForCondition: after})

			queryTextDelete, _, err := queryBuilderDelete.ToSql()
			//fmt.Println(queryTextDelete)
			fmt.Printf("ShrinkTables: %s - query: %s - counter %d\n", TableName, queryTextDelete, counter)
			if err != nil {
				return err
			}
			_, err = DB.Exec(queryTextDelete, after)
			if err != nil {
				return err
			}
		}

		log.Impl.Errorf("Обнаруженно переполнение таблицы %s\n Количество записей : %d выполненно усечение больше чем %-8v",
			TableName, counter, DurationTimeRemaindRows)

	}

	return nil

}

func GetCurrentYearAsString() string {
	today := time.Now()
	yearFilterInt := today.Year() //"2022"
	return strconv.Itoa(yearFilterInt)
}

func CloseStatusJob(DB *sql.DB) error {

	var argsquery []interface{}
	argsquery = append(argsquery, "Выполнено")

	//	queryAllColumns := ""

	//	if strings.EqualFold(os.Getenv("USE_SETTINGS_JOB_V2"), "true") {

	queryAllColumns := `select
		jobs.job_id as job_id1,
		coalesce(settings_jobs.code_external, '') as code_external,
		coalesce(settings_jobs.name_external, '') as name_external,
		coalesce(settings_jobs.table_name, '') as table_name,
		coalesce(exchange_jobs."event", '') as status,
		count(exchange_jobs."event") as event_count
	from
		public.jobs as jobs
	left join public.settings_jobs_v2 as settings_jobs on
		jobs.job_id = settings_jobs.job_id
	left join public.exchange_jobs as exchange_jobs on
		jobs.job_id = exchange_jobs.job_id
	where
		status <> $1
	group by
		job_id1,
		code_external,
		name_external,
		table_name,
		exchange_jobs."event"
	order by
		job_id1,
		"event"`

	//	} else {

	// 		// Мы не закрываем задания удаленного сбора, они закрываются в отдельном микросервисе
	// 		queryAllColumns = `select
	// 	jobs.job_id as job_id1,
	// 	coalesce(settings_jobs.code_external, '') as code_external,
	// 	coalesce(settings_jobs.name_external, '') as name_external,
	// 	coalesce(settings_jobs.table_name, '') as table_name,
	// 	coalesce(exchange_jobs."event", '') as status,
	// 	count(exchange_jobs."event") as event_count
	// from
	// 	public.jobs as jobs
	// left join public.settings_jobs as settings_jobs on
	// 	jobs.job_id = settings_jobs.job_id
	// left join public.exchange_jobs as exchange_jobs on
	// 	jobs.job_id = exchange_jobs.job_id
	// where
	// 	status <> $1
	// group by
	// 	job_id1,
	// 	code_external,
	// 	name_external,
	// 	table_name,
	// 	exchange_jobs."event"
	// order by
	// 	job_id1,
	// 	"event"`
	// 		//and coalesce(settings_jobs.use_remote_collection, false) <> true
	// 	}

	rows, err := DB.Query(queryAllColumns, argsquery...)
	if err != nil {
		return err
	}

	defer rows.Close()

	ResultMap := make(map[string][]string)
	ResultSettingsMap := make(map[string]models.JobDescription)

	for rows.Next() {
		var JobID, Status, Event_count string
		var JobDescription models.JobDescription
		err = rows.Scan(&JobID, &JobDescription.Code1C, &JobDescription.Name1C, &JobDescription.TableName, &Status, &Event_count)
		if err != nil {
			return err
		}

		Records, ok := ResultMap[JobID]
		if ok != true {
			var NewRecord []string
			NewRecord = append(NewRecord, Status)
			ResultMap[JobID] = NewRecord
		} else {
			Records = append(Records, Status)
			ResultMap[JobID] = Records
		}

		ResultSettingsMap[JobID] = JobDescription

	}

	// 	var strBuilder strings.Builder
	// 	strBuilder.WriteString(JobID + " ")
	// 	flagExecuteJob := true
	// 	for rows2.Next() {

	// 		var Status string
	// 		var jobs_count int
	// 		err = rows2.Scan(&jobs_count, &Status)
	// 		if err != nil {
	// 			return err
	// 		}

	// 		if Status != "Выполнено" {
	// 			flagExecuteJob = false
	// 		}

	// 		strBuilder.WriteString(fmt.Sprintf("%s:%d\n", Status, jobs_count))

	// 	}

	// Закрываем задание
	for JobID, value := range ResultMap {

		if len(value) == 1 && value[0] == "Выполнено" {

			// err := ChangeStatusJobsTask(DB, JobID, "Выполнено")
			// if err != nil {
			// 	//return nil
			// 	continue
			// }

			var JobV2 models.JobV2
			JobV2.JobID = JobID
			JobV2.Status = "Выполнено"
			var MessageQueueGeneralInterface models.MessageQueueGeneralInterface
			MessageQueueGeneralInterface.Type = "ChangeStatusForJobV2"
			MessageQueueGeneralInterface.Body = JobV2
			JsonMessageBody, err := json.Marshal(&MessageQueueGeneralInterface)
			if err != nil {
				log.Impl.Errorf("ошибка маршалинга: %s", err)
				//return err
				continue
			}

			headers := map[string]interface{}{}
			err = queue.RabbitMQConnectorVb.SendInRabbitMQUniversalV2newChannel(JsonMessageBody, "go-keeper-status", headers)
			if err != nil {
				//return err
				continue
			}

			ResultSettings := ResultSettingsMap[JobID]

			// var QueryToBI models.QueryToBI
			// err = QueryToBI.LoadSettingsFirstRowFromPgByJobID(DB, JobID)
			var SettingsJobsAllV2 models.SettingsJobsAllV2
			err = SettingsJobsAllV2.LoadSettingsFromPgByJobID(DB, JobID)
			if err != nil {
				err = fmt.Errorf("НастройкиМодели not filled in QueryResult для JobId %s", JobID)
				log.Impl.Error(err)
				continue
			}

			configDSN, err := pq.ParseURL(SettingsJobsAllV2.DSNconnection)
			if err != nil {
				log.Impl.Error(err)
				continue
			}

			s := strings.Split(configDSN, " ")
			sSummary := ""
			for _, subS := range s {
				if strings.Contains(subS, "dbname=") {
					sSummary = sSummary + subS + " " //strings.ReplaceAll(subS, "dbname=", "")
				}
				if strings.Contains(subS, "host=") {
					sSummary = sSummary + subS + " " //strings.ReplaceAll(subS, "host=", "")
				}
			}

			// err = SendTextToTelegramChat(fmt.Sprintf("Выполнено задание: %s\nКод в 1С: %s\nИмя таблицы: %s\nCommit микросервиса: %s", ResultSettings.Name1C, ResultSettings.Code1C,
			// 	ResultSettings.TableName, version.Commit))
			err = SendTextToTelegramChat(fmt.Sprintf("Выполнено задание: %s\nКод в админке: %s\nИмя таблицы: %s\nServer: %s", ResultSettings.Name1C, ResultSettings.Code1C,
				ResultSettings.TableName, sSummary))
			if err != nil {
				log.Impl.Error(err)
			}

			if SettingsJobsAllV2.UseHandleAfterLoadAlgorithms {

				var HandleAfterLoad models.HandleAfterLoad
				HandleAfterLoad.JobID = JobID
				HandleAfterLoad.Algorithms = SettingsJobsAllV2.ListHandleAfterLoadAlgorithms

				var MessageQueueGeneralInterface models.MessageQueueGeneralInterface
				MessageQueueGeneralInterface.Type = "HandleAfterLoad"
				MessageQueueGeneralInterface.Body = HandleAfterLoad

				JsonMessageBody, err := json.Marshal(MessageQueueGeneralInterface)
				if err != nil {
					return err
				}

				// var RESTRequestUniversal models.RESTRequestUniversal
				// Headers := make(map[string]string)
				// Headers["TokenBearer"] = config.Conf.TokenBearer
				// RESTRequestUniversal.Headers = Headers
				// RESTRequestUniversal.Method = "POST"
				// RESTRequestUniversal.Body = MessageQueueGeneralInterfaceByte
				// // TODO: Переделать на переменную окружения
				// RESTRequestUniversal.UrlToCall = os.Getenv("ADDRESS_PORT_SERVICE_FRONT") + "/save-event-to-queue"
				// _, err = RESTRequestUniversal.Send()
				// if err != nil {
				// 	log.Impl.Error(err)
				// }

				headers := map[string]interface{}{}
				err = queue.RabbitMQConnectorVb.SendInRabbitMQUniversalV2newChannel(JsonMessageBody, "go-keeper-events", headers)
				if err != nil {
					log.Impl.Error(err)
					continue
				}

			}

			log.Impl.Infof("Обновленно задание: %s", JobID)
		}

	}

	return nil

}

// sendTextToTelegramChat sends a text message to the Telegram chat identified by its chat Id
// func SendTextToTelegramChat(chatId int, text string) (string, error) {
func SendTextToTelegramChat(Text string) error {

	var MessageTelegram models.MessageTelegram
	MessageTelegram.ChatID = os.Getenv("TELEGRAM_CHAT_ID")
	MessageTelegram.Text = Text

	JsonByteMessageBody, err := json.Marshal(&MessageTelegram)
	if err != nil {
		return err
	}

	method := "POST"
	//urlToCall := "https://1d56fe65.proxy.webhookapp.com"
	urlToCall := "https://api.telegram.org"
	tokenBot := os.Getenv("TELEGRAM_TOKEN_BOT")
	urlToCall += "/bot" + tokenBot + "/" + "sendMessage"

	body := bytes.NewBuffer(JsonByteMessageBody)
	//useAuth := true

	req, err := http.NewRequest(method, urlToCall, body)
	req.Header.Set("Content-Type", "application/json")
	// if useAuth {
	// 	req.Header.Set("Authorization", "Basic "+token)
	// }
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Impl.Error(err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyResp, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Impl.Error(err)
			return errors.New("status code not ok: " + strconv.Itoa(resp.StatusCode))
		}

		return errors.New("status code not ok: " + strconv.Itoa(resp.StatusCode) + " body: " + string(bodyResp) + " url: " + urlToCall)
	}

	return nil

}

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}
