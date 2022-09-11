package queue

import (
	"github.com/ometcenter/keeper/config"
	log "github.com/ometcenter/keeper/logging"
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
