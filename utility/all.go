package utility

import (
	"database/sql"
	"encoding/json"

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
