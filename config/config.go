package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// Conf ...
var Conf *Config

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
type Config struct {
	Port            string
	MessagePath     string
	MaxWorker       int
	MaxQueue        int
	MaxLength       int
	DatabaseURL     string
	SentryUrlDSN    string
	LoggerDefault   string
	LogLevel        int
	Release         string
	SecretKeyJWT    string
	QueueType       string
	AddressRabbitMQ string
	PubSubConfig
	Pusher
}

// New returns a new Config struct
func New() *Config {

	if tz := getEnv("TIMEZONE", "Local"); tz != "" {
		var err error
		time.Local, err = time.LoadLocation(tz)
		if err != nil {
			log.Printf("[ERROR] loading location '%s': %v\n", tz, err)
		}
		nameLocation, _ := time.Now().Zone()
		log.Printf("[INFO] текущая таймзона %s", nameLocation)
	}

	c := &Config{
		Port:            getEnv("PORT", "8080"),
		MessagePath:     getEnv("MESSAGE_PATH", "/"),
		MaxWorker:       getEnvAsInt("MAX_WORKERS", 1),
		MaxQueue:        getEnvAsInt("MAX_JOBS_IN_QUEUE", 100),
		MaxLength:       getEnvAsInt("MAX_LENGTH", 1048576),
		DatabaseURL:     getEnv("DB_CONNECTION", ""),
		SentryUrlDSN:    getEnv("SENTRY_URL_DSN", "http://ded6d3a6b5c64c0d9e38c042a365fa39:0aed9c2ef0994bf39f40e7227174bfa2@localhost:9000/2"),
		LoggerDefault:   getEnv("LOGGER_DEFAULT", "Sentry"),
		LogLevel:        getEnvAsInt("LOG_LEVEL", 0),
		Release:         getEnv("RELEASE", "Nope"),
		QueueType:       getEnv("QUEUE_TYPE", "RabbitMQ"), //NSQ
		AddressRabbitMQ: getEnv("ADDRESS_RABBIT_MQ", "amqp://localhost:5672"),
		SecretKeyJWT:    getEnv("SECRET_KEY_JWT", ""),
		PubSubConfig: PubSubConfig{
			Topic:         getEnv("NSQ_TOPIC", "go-keeper-messages"),
			Channel:       getEnv("NSQ_CHANNEL", "keeper-agent"),
			NsqLookupdPub: getEnv("NSQ_LOOKUPD_PUB", "localhost:4150"),
			NsqLookupdSub: getEnv("NSQ_LOOKUPD_SUB", "localhost:4161"),
			MaxRequeue:    getEnvAsInt("NSQ_MAX_REQUEUE", 10),
			Concurrent:    getEnvAsInt("NSQ_CONCURRENT", 1),
			MaxInFlight:   getEnvAsInt("NSQ_MAX_IN_FLIGHT", 3),
			UseDailyTopic: getEnvAsBool("NSQ_USE_DAILY_TOPIC", false),
		},
		Pusher: Pusher{
			Address: getEnv("PUSHER_ADDRESS", "http://localhost"),
			Token:   getEnv("PUSHER_TOKEN", ""),
		},
	}

	Conf = c
	return c
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

// Simple helper function to read an environment variable into integer or return a default value
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

// Helper to read an environment variable into a bool or return default value
func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}

// Helper to read an environment variable into a string slice or return default value
func getEnvAsSlice(name string, defaultVal []string, sep string) []string {
	valStr := getEnv(name, "")

	if valStr == "" {
		return defaultVal
	}

	val := strings.Split(valStr, sep)

	return val
}
