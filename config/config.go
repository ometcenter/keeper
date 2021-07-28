package config

import (
	"fmt"
	"time"

	"github.com/ometcenter/keeper/env"
)

type LoggerType int

const (
	LoggerSentry LoggerType = iota
	LoggerLorgus
	LoggerRabbitMQ
	LoggerNSQ
)

// const (
// 	LoggerSentry   = "Sentry"
// 	LoggerLorgus   = "Lorgus"
// 	LoggerRabbitMQ = "RabbitMQ"
// 	LoggerNSQ      = "NSQ"
// )

// Conf ...
var Conf *ServiceConfig

// TODO: Почему не в коде, а через переменные?
var (
	defaultGRPCPort    = "5300"
	defaultGRPCAddress = "localhost"
)

func InitConfig() {
	Conf = New()
	_ = Conf.InitTimezone()
}

// Pusher ...
type Pusher struct {
	Address string
	Token   string
}

// PubSubConfig ...
type PubSubConfig struct {
	Topic         string
	Channel       string
	NsqLookupdPub string
	NsqLookupdSub string
	MaxRequeue    int
	Concurrent    int
	MaxInFlight   int
	UseDailyTopic bool
}

// Config ...
type ServiceConfig struct {
	Port            string
	PortgRPC        string
	AddressPortgRPC string
	MessagePath     string
	MaxWorker       int
	MaxQueue        int
	MaxLength       int
	DatabaseURL     string
	SentryUrlDSN    string
	Release         string
	SecretKeyJWT    string
	QueueType       string
	AddressRabbitMQ string
	PubSubConfig
	Pusher
	LoggerConfig
}

// LoggerConfig содержит настройки для логгера
type LoggerConfig struct {
	Level int
	Name  string
	//TODO: Встроить логи через константы например LoggerType[1]
	LoggerType LoggerType
}

// New returns a new Config struct
func New() *ServiceConfig {

	// if tz := env.GetEnv("TIMEZONE", "Local"); tz != "" {
	// 	var err error
	// 	time.Local, err = time.LoadLocation(tz)
	// 	if err != nil {
	// 		log.Printf("[ERROR] loading location '%s': %v\n", tz, err)
	// 	}
	// 	nameLocation, _ := time.Now().Zone()
	// 	log.Printf("[INFO] текущая таймзона %s", nameLocation)
	// }

	ServiceConfig := &ServiceConfig{
		Port:            env.GetEnv("PORT", "8080"),
		PortgRPC:        env.GetEnv("PORT_gRPC", defaultGRPCPort),
		AddressPortgRPC: env.GetEnv("ADDRESS_PORT_gRPC", fmt.Sprintf("%s:%s", defaultGRPCAddress, defaultGRPCPort)),
		MessagePath:     env.GetEnv("MESSAGE_PATH", "/"),
		MaxWorker:       env.GetEnvAsInt("MAX_WORKERS", 1),
		MaxQueue:        env.GetEnvAsInt("MAX_JOBS_IN_QUEUE", 100),
		MaxLength:       env.GetEnvAsInt("MAX_LENGTH", 1048576),
		DatabaseURL:     env.GetEnv("DB_CONNECTION", ""),
		Release:         env.GetEnv("RELEASE", "Nope"),
		QueueType:       env.GetEnv("QUEUE_TYPE", "RabbitMQ"), //NSQ
		AddressRabbitMQ: env.GetEnv("ADDRESS_RABBIT_MQ", "amqp://localhost:5672"),
		SecretKeyJWT:    env.GetEnv("SECRET_KEY_JWT", ""),
		SentryUrlDSN:    env.GetEnv("SENTRY_URL_DSN", ""),
		PubSubConfig: PubSubConfig{
			Topic:         env.GetEnv("NSQ_TOPIC", "go-keeper-messages"),
			Channel:       env.GetEnv("NSQ_CHANNEL", "keeper-agent"),
			NsqLookupdPub: env.GetEnv("NSQ_LOOKUPD_PUB", "localhost:4150"),
			NsqLookupdSub: env.GetEnv("NSQ_LOOKUPD_SUB", "localhost:4161"),
			MaxRequeue:    env.GetEnvAsInt("NSQ_MAX_REQUEUE", 10),
			Concurrent:    env.GetEnvAsInt("NSQ_CONCURRENT", 1),
			MaxInFlight:   env.GetEnvAsInt("NSQ_MAX_IN_FLIGHT", 3),
			UseDailyTopic: env.GetEnvAsBool("NSQ_USE_DAILY_TOPIC", false),
		},
		Pusher: Pusher{
			Address: env.GetEnv("PUSHER_ADDRESS", "http://localhost"),
			Token:   env.GetEnv("PUSHER_TOKEN", ""),
		},
		LoggerConfig: LoggerConfig{
			Name:  env.GetEnv("LOGGER_DEFAULT", "Sentry"), //Lorgus
			Level: env.GetEnvAsInt("LOG_LEVEL", 0),
		},
	}

	return ServiceConfig
}

// InitTimezone устанавливает временную зону в time.Local
func (s ServiceConfig) InitTimezone() string {
	str := ""
	var err error
	if tz := env.GetEnv("TIMEZONE", "Local"); tz != "" {
		time.Local, err = time.LoadLocation(tz)
		if err != nil {
			str = fmt.Sprintf("[ERROR] loading location '%s': %v\n", tz, err)
		}
		nameLocation, _ := time.Now().Zone()
		str = fmt.Sprintf("[INFO] текущая таймзона %s\n", nameLocation)
	}
	return str
}
