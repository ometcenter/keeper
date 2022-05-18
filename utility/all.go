package utility

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/go-redis/redis/v8"
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

//"github.com/go-redis/redis/v8"
func GetDataRedisByParamGoRedislibrary(Param string, RedisDB int, RedisClient *redis.Client) (string, error) {

	var Result string

	_, err := RedisClient.Do(context.Background(), "select", RedisDB).Result()
	if err != nil {
		return Result, err
	}

	val, err := RedisClient.Get(context.Background(), Param).Result()
	if err == redis.Nil {
		//fmt.Println("key2 does not exist")
		//return Result, fmt.Errorf("Не найден ключ для JobId: %s в Redis", InsuranceNumber)
		return Result, nil
	} else if err != nil {
		//panic(err)
		return Result, err
		//return Result, nil
	} else {
		//fmt.Println("key2", val2)
		return val, nil
	}

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

	//queryBuilder := sq.Select("*").From(TableNameParam)
	param := []interface{}{}

	if len(QueryURL) != 0 {
		//var conditionSlice []sq.Eq
		//conditionMap := make(map[string]interface{})
		//var conditionSlice sq.Eq

		sqEq := make(sq.Eq)

		for key, value := range QueryURL {
			//conditionSlice = append(conditionSlice, sq.Eq{key: value})
			//conditionMap[key] = value
			if len(value) == 0 {
				param = append(param, "")
			} else {
				param = append(param, value[0])
			}

			sqEq[key] = ""

		}
		//queryBuilder = sq.Select("*").From(TableNameParam).Where(sqEq)
		queryBuilder = psql.Select("*").From(TableNameParam).Where(sqEq)
	}

	queryText, _, err := queryBuilder.ToSql()
	//fmt.Println(queryText)
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
