/* Package config определяет структуры конфигурации и некоторые интерфейсы
 */
package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/ometcenter/keeper/env"
)

const (
	LoggerSentry = "Sentry"
	LoggerLorgus = "Lorgus"
	RabbitMQ     = "RabbitMQ"
	NSQ          = "NSQ"
)

var (
	errInvalidHost     = errors.New("invalid host")
	defaultPort        = "8080"
	defaultGRPCPort    = "5300"
	defaultGRPCAddress = "localhost"
	defaultMessagePath = "/"

	defaultPusherAddress   = "http://localhost"
	defaultAddressRabbitMQ = "amqp://localhost:5672"

	defaultNSQTopic      = "go-keeper-messages"
	defaultNSQChannel    = "keeper-agent"
	defaultNSQLookupdPub = "localhost:4150"
	defaultNSQLookupdSub = "localhost:4161"
)

// Pusher содержит настройки для использования сервиса Pusher
type PusherConfig struct {
	// Адрес сервиса принимающего события из Pusher
	Address string
	// Логин и пароль в формате base64
	Token string
}

// PubSubConfig содержит настройки для использования сервиса очередей
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

// LoggerConfig содержит настройки для логгера
type LoggerConfig struct {
	Level int
	Name  string
}

// ServiceConfig основные настроки приложения
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
	LoggerConfig
	Release         string
	SecretKeyJWT    string
	QueueType       string
	AddressRabbitMQ string
	PubSubConfig
	PusherConfig
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

// New возвращает новый ServiceConfig тип
func New() *ServiceConfig {

	pbc := PubSubConfig{
		Topic:         env.GetEnv("NSQ_TOPIC", defaultNSQTopic),
		Channel:       env.GetEnv("NSQ_CHANNEL", defaultNSQChannel),
		NsqLookupdPub: env.GetEnv("NSQ_LOOKUPD_PUB", defaultNSQLookupdPub),
		NsqLookupdSub: env.GetEnv("NSQ_LOOKUPD_SUB", defaultNSQLookupdSub),
		MaxRequeue:    env.GetEnvAsInt("NSQ_MAX_REQUEUE", 10),
		Concurrent:    env.GetEnvAsInt("NSQ_CONCURRENT", 1),
		MaxInFlight:   env.GetEnvAsInt("NSQ_MAX_IN_FLIGHT", 3),
		UseDailyTopic: env.GetEnvAsBool("NSQ_USE_DAILY_TOPIC", false),
	}

	pc := PusherConfig{
		Address: env.GetEnv("PUSHER_ADDRESS", defaultPusherAddress),
		Token:   env.GetEnv("PUSHER_TOKEN", ""),
	}

	l := LoggerConfig{
		Name:  env.GetEnv("LOGGER_DEFAULT", LoggerSentry), //Lorgus
		Level: env.GetEnvAsInt("LOG_LEVEL", 0),
	}
	return &ServiceConfig{
		Port:            env.GetEnv("PORT", defaultPort),
		PortgRPC:        env.GetEnv("PORT_gRPC", defaultGRPCPort),
		AddressPortgRPC: env.GetEnv("ADDRESS_PORT_gRPC", fmt.Sprintf("%s:%s", defaultGRPCAddress, defaultGRPCPort)),
		MessagePath:     env.GetEnv("MESSAGE_PATH", defaultMessagePath),
		MaxWorker:       env.GetEnvAsInt("MAX_WORKERS", 1),
		MaxQueue:        env.GetEnvAsInt("MAX_JOBS_IN_QUEUE", 100),
		MaxLength:       env.GetEnvAsInt("MAX_LENGTH", 1048576),
		DatabaseURL:     env.GetEnv("DB_CONNECTION", ""),
		SentryUrlDSN:    env.GetEnv("SENTRY_URL_DSN", ""),
		Release:         env.GetEnv("RELEASE", "Nope"),
		QueueType:       env.GetEnv("QUEUE_TYPE", "RabbitMQ"), //NSQ
		AddressRabbitMQ: env.GetEnv("ADDRESS_RABBIT_MQ", defaultAddressRabbitMQ),
		SecretKeyJWT:    env.GetEnv("SECRET_KEY_JWT", ""),
		PubSubConfig:    pbc,
		PusherConfig:    pc,
		LoggerConfig:    l,
	}

}
