package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

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

	ConnetRabbitMQ()

	return nil

}

func ConnetRabbitMQ() {

	c := make(chan *amqp.Error)
	go func() {
		err := <-c
		log.Impl.Error("reconnect RabbitMQ: " + err.Error())
		ConnetRabbitMQ()
	}()

	var err error
	ConnectRabbitMQ, err = amqp.Dial(config.Conf.AddressRabbitMQ)
	if err != nil {
		log.Impl.Error("Can not connect to RabbitMQ :", err)
		log.Impl.Panic()
		panic("Can not connect to RabbitMQ")
	}
	ConnectRabbitMQ.NotifyClose(c)

}

func SendInRabbitMQUniversalV2newChannel(messageBody []byte, topicName string, ConnectRabbitMQ *amqp.Connection,
	Headers map[string]interface{}) error {

	// TODO: При сохранении сообщений возникает ошибка "Exception (505) Reason: "UNEXPECTED_FRAME - expected content body, got non content body frame instead""
	// я создавал два канала для глобальной переменной для Сообщений и Событий, но ошибка сохранилась.
	// возможно попробовать на другой докере не моем развернуть шину, делаю для каждого сообщения отдельный канал
	// так он Gin видимо пораждает отдельную горутину для каждого вызова конечной точки, один канал под одну горутину.
	// Создаю глобальную переменную ConnectRabbitMQ
	RabbitMQchannel, err := ConnectRabbitMQ.Channel()
	if err != nil {
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

func SendInRabbitMQUniversalV2currentChannel(messageBody []byte, topicName string, Channel *amqp.Channel, Headers map[string]interface{}) error {

	// TODO: При сохранении сообщений возникает ошибка "Exception (505) Reason: "UNEXPECTED_FRAME - expected content body, got non content body frame instead""
	// я создавал два канала для глобальной переменной для Сообщений и Событий, но ошибка сохранилась.
	// возможно попробовать на другой докере не моем развернуть шину, делаю для каждого сообщения отдельный канал
	// так он Gin видимо пораждает отдельную горутину для каждого вызова конечной точки, один канал под одну горутину.
	// Создаю глобальную переменную ConnectRabbitMQ

	q, err := Channel.QueueDeclare(
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

	err = Channel.Publish(
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

type RabbitMQConnector struct {
	commandChannel           chan string
	connectRabbitMQ          *amqp.Connection
	rabbitMQErrorsChan       chan *amqp.Error
	activeTokens             map[string]string
	activeTokensMu           sync.RWMutex
	systems                  map[string]string
	ctx                      context.Context
	ctxCancelFn              func()
	saveMapToExternalStorage func()
}

var RabbitMQConnectorVb *RabbitMQConnector

func NewRabbitMQConnector() *RabbitMQConnector {
	ctx, cancel := context.WithCancel(context.Background())

	systems := make(map[string]string)
	systems["keeper"] = "keeper"
	systems["ekis"] = "ekis"

	var sayHelloWorld = func() {
		//fmt.Println("Hello World !")
	}

	return &RabbitMQConnector{
		commandChannel: make(chan string),
		// out:            make(chan interface{}, 10),
		rabbitMQErrorsChan:       make(chan *amqp.Error),
		activeTokens:             make(map[string]string),
		systems:                  systems,
		ctx:                      ctx,
		ctxCancelFn:              cancel,
		saveMapToExternalStorage: sayHelloWorld,
	}
}

func (t *RabbitMQConnector) Run() error {

	t.ConnetRabbitMQ()

	// for key, _ := range t.systems {
	// 	if key == "keeper" {
	// 		err := t.CreateSessionKeeper()
	// 		if err != nil {
	// 			log.Impl.Error(err)
	// 		}
	// 	}
	// }

	tkExpiring := time.NewTicker(time.Second * 300)
	defer tkExpiring.Stop()

	for {
		select {
		case <-t.ctx.Done():
			fmt.Printf("ServerTokenStore STOP!\n")
			return errors.New("STOP")
		case <-tkExpiring.C:

			// for key, _ := range t.systems {
			// 	// if key == "keeper" {
			// 	// 	err := t.validateSessionKeeper()
			// 	// 	if err != nil {
			// 	// 		log.Impl.Error(err)
			// 	// 	}
			// 	// }
			// }

			t.saveMapToExternalStorage()

			//default:
		}
		//time.Sleep(time.Second * 10)

		//t.Stop()
	}
}

func (t *RabbitMQConnector) Stop() {
	t.ctxCancelFn()

	//RabbitMQchannelConsumer.Close()
	//RabbitMQchannelPublic.Close()

	// close(w.out)
}

func (t *RabbitMQConnector) ConnetRabbitMQ() {

	//c := make(chan *amqp.Error)
	go func() {
		err := <-t.rabbitMQErrorsChan
		log.Impl.Error("reconnect RabbitMQ: " + err.Error())
		time.Sleep(time.Second * 10)
		t.ConnetRabbitMQ()
	}()

	var err error
	t.connectRabbitMQ, err = amqp.Dial(config.Conf.AddressRabbitMQ)
	if err != nil {
		log.Impl.Error("Can not connect to RabbitMQ :", err)
		log.Impl.Panic()
		panic("Can not connect to RabbitMQ")
	}

	t.connectRabbitMQ.NotifyClose(t.rabbitMQErrorsChan)
}

func (t *RabbitMQConnector) SendInRabbitMQUniversalV2newChannel(messageBody []byte, topicName string, Headers map[string]interface{}) error {

	// TODO: При сохранении сообщений возникает ошибка "Exception (505) Reason: "UNEXPECTED_FRAME - expected content body, got non content body frame instead""
	// я создавал два канала для глобальной переменной для Сообщений и Событий, но ошибка сохранилась.
	// возможно попробовать на другой докере не моем развернуть шину, делаю для каждого сообщения отдельный канал
	// так он Gin видимо пораждает отдельную горутину для каждого вызова конечной точки, один канал под одну горутину.
	// Создаю глобальную переменную ConnectRabbitMQ
	RabbitMQchannel, err := t.connectRabbitMQ.Channel()
	if err != nil {
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
