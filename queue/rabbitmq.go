package queue

import (
	"encoding/json"

	"github.com/ometcenter/keeper/config"
	log "github.com/ometcenter/keeper/logging"
	tracing "github.com/ometcenter/keeper/tracing/jaeger"
	tracingRabbitMQ "github.com/ometcenter/keeper/tracing/jaeger/rabbitmq"
	"github.com/streadway/amqp"
)

var ConnectRabbitMQ *amqp.Connection

func StartQueue() error {

	switch config.Conf.QueueType {
	case "RabbitMQ":
		err := InitRabbitMQ()
		if err != nil {
			return err
		}
	default:
		// err := StartNSQ()
		// if err != nil {
		// 	return err
		// }
	}

	return nil
}

func InitRabbitMQ() error {

	// TODO: При сохранении сообщений возникает ошибка "Exception (505) Reason: "UNEXPECTED_FRAME - expected content body, got non content body frame instead""
	// я создавал два канала для глобальной переменной для Сообщений и Событий, но ошибка сохранилась.
	// возможно попробовать на другой докере не моем развернуть шину, делаю для каждого сообщения отдельный канал
	// так он Gin видимо пораждает отдельную горутину для каждого вызова конечной точки, один канал под одну горутину.
	// Создаю глобальную переменную ConnectRabbitMQ
	ConnetRabbitMQ()

	// conn, err := amqp.Dial(config.Conf.AddressRabbitMQ) //5672
	// if err != nil {
	// 	//RabbitMQchannel = nil
	// 	return err
	// }
	// //defer conn.Close()

	// ch, err := conn.Channel()
	// if err != nil {
	// 	//Connector.LoggerConn.ErrorLogger.Println("Failed to open a channel")
	// 	//Connector.RabbitMQchannel = nil
	// 	return err
	// }
	// //defer ch.Close()

	// RabbitMQchannel = ch

	return nil

}

func ConnetRabbitMQ() {

	c := make(chan *amqp.Error)
	go func() {
		err := <-c
		log.Impl.Error("reconnect RabbitMQ: " + err.Error())
		ConnetRabbitMQ()
	}()

	// TODO: При сохранении сообщений возникает ошибка "Exception (505) Reason: "UNEXPECTED_FRAME - expected content body, got non content body frame instead""
	// я создавал два канала для глобальной переменной для Сообщений и Событий, но ошибка сохранилась.
	// возможно попробовать на другой докере не моем развернуть шину, делаю для каждого сообщения отдельный канал
	// так он Gin видимо пораждает отдельную горутину для каждого вызова конечной точки, один канал под одну горутину.
	// Создаю глобальную переменную ConnectRabbitMQ

	//conn, err := amqp.Dial(config.Conf.AddressRabbitMQ)
	var err error
	ConnectRabbitMQ, err = amqp.Dial(config.Conf.AddressRabbitMQ)
	if err != nil {
		log.Impl.Error("Can not connect to RabbitMQ :", err)
		log.Impl.Panic()
		panic("Can not connect to RabbitMQ")
	}
	ConnectRabbitMQ.NotifyClose(c)

	// chMessage, err := ConnectRabbitMQ.Channel()
	// if err != nil {
	// 	//Connector.LoggerConn.ErrorLogger.Println("Failed to open a channel")
	// 	//Connector.RabbitMQchannel = nil
	// 	log.Impl.Error("Failed to open a channel :", err)
	// 	log.Impl.Panic()
	// 	panic("Failed to open a channel")
	// }
	// //defer ch.Close()

	// RabbitMQchannelMessage = chMessage

	// // create topology

	// chEvent, err := ConnectRabbitMQ.Channel()
	// if err != nil {
	// 	//Connector.LoggerConn.ErrorLogger.Println("Failed to open a channel")
	// 	//Connector.RabbitMQchannel = nil
	// 	log.Impl.Error("Failed to open a channel :", err)
	// 	log.Impl.Panic()
	// 	panic("Failed to open a channel")
	// }
	// //defer ch.Close()

	// RabbitMQchannelEvent = chEvent
}

func SendInRabbitMQUniversalV2newChannel(messageBody []byte, topicName string, ConnectRabbitMQ *amqp.Connection, Headers map[string]interface{}) error {

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
		Body:        messageBody,
	}

	if len(Headers) > 0 {

		// Сделанно для инициализации карты, с nil валиться в паник, внутри библиотеки не предусмотрели.
		amqpPublishing.Headers = map[string]interface{}{}

		for key, value := range Headers {
			amqpPublishing.Headers[key] = value
		}

		if config.Conf.UseTracing {

			clientSpan := tracing.Tracer.StartSpan("clientspan")
			defer clientSpan.Finish()

			if err := tracingRabbitMQ.Inject(clientSpan, amqpPublishing.Headers); err != nil {
				log.Impl.Error(err)
			}
		}
	}

	err = RabbitMQchannel.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqpPublishing)
	if err != nil {
		return err
	}

	return nil

}

// TODO: Replace this
type MessageNSQ struct {
	Type string `json:"Тип"`
	Body string `json:"Тело"`
}

func SaveStatus(BodyString string, TypeMessage string, ConnectRabbitMQ *amqp.Connection) error {

	MessageNSQMarshal := MessageNSQ{Type: TypeMessage, Body: BodyString}

	JsonMessage, err := json.Marshal(&MessageNSQMarshal)
	if err != nil {
		log.Impl.Errorf("ошибка маршалинга: %s", err)
		return err
	}

	switch config.Conf.QueueType {
	case "RabbitMQ":
		headers := map[string]interface{}{}
		err = SendInRabbitMQUniversalV2newChannel(JsonMessage, "go-keeper-status", ConnectRabbitMQ, headers)
		if err != nil {
			log.Impl.Errorf("ошибка отправки в RabbitMQ: %s", err)
			return err
		}

	default:
		// err = SendInQueueNSQ(JsonMessage, "go-keeper-status")
		// if err != nil {
		// 	log.Impl.Errorf("ошибка отправки в RabbitMQ: %s", err)
		// 	return err
		// }
	}

	return nil

}
