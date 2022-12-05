package utility

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/ometcenter/keeper/config"
	log "github.com/ometcenter/keeper/logging"
	"github.com/ometcenter/keeper/models"
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
				uint64(Limit)).Offset(uint64(Offset)).OrderBy("created_at")
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
				uint64(Limit)).Offset(uint64(Offset)).OrderBy("created_at")
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

		queryBuilder = psql.Select("*").From(TableNameParam).Limit(uint64(Limit)).Offset(uint64(Offset)).OrderBy("created_at")

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
