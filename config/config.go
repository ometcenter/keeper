package config

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
	"time"

	consul "github.com/hashicorp/consul/api"
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

	//TODO: Что делать с ошибками и вызывать ли паник, сейчас вызываю паник
	Conf.LoadSettingsFromConsul()
	Conf.LoadSettingsFromDockerSecrets()
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

type LoadSettings struct {
	LoadSettingsFromConsul       bool
	ConsulServerAddres           string
	LoadSettingsFromDockerSecret bool
	ArrayDockerSecretKey         string
	ConsulUpdateIndex            uint64
}

// Config ...
type ServiceConfig struct {
	Port                   string
	PortgRPC               string
	AddressPortgRPC        string
	MessagePath            string
	MaxWorker              int
	MaxQueue               int
	MaxLength              int
	DatabaseURL            string
	SentryUrlDSN           string
	Release                string
	SecretKeyJWT           string
	QueueType              string
	AddressRabbitMQ        string
	UseRedis               bool
	RedisAddressPort       string
	TokenBearer            string
	AddressPortCRONService string
	AddressPostJaeger      string
	UseTracing             bool
	PubSubConfig
	Pusher
	LoggerConfig
	LoadSettings
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
		Port:                   env.GetEnv("PORT", "8080"),
		PortgRPC:               env.GetEnv("PORT_gRPC", defaultGRPCPort),
		AddressPortgRPC:        env.GetEnv("ADDRESS_PORT_gRPC", fmt.Sprintf("%s:%s", defaultGRPCAddress, defaultGRPCPort)),
		MessagePath:            env.GetEnv("MESSAGE_PATH", "/"),
		MaxWorker:              env.GetEnvAsInt("MAX_WORKERS", 1),
		MaxQueue:               env.GetEnvAsInt("MAX_JOBS_IN_QUEUE", 100),
		MaxLength:              env.GetEnvAsInt("MAX_LENGTH", 1048576),
		DatabaseURL:            env.GetEnv("DB_CONNECTION", ""),
		Release:                env.GetEnv("RELEASE", "Nope"),
		QueueType:              env.GetEnv("QUEUE_TYPE", "RabbitMQ"), //NSQ
		AddressRabbitMQ:        env.GetEnv("ADDRESS_RABBIT_MQ", "amqp://localhost:5672"),
		SecretKeyJWT:           env.GetEnv("SECRET_KEY_JWT", ""),
		SentryUrlDSN:           env.GetEnv("SENTRY_URL_DSN", ""),
		UseRedis:               env.GetEnvAsBool("USE_REDIS", false),
		RedisAddressPort:       env.GetEnv("REDIS_ADDRESS_PORT", "localhost:6379"),
		AddressPortCRONService: env.GetEnv("ADDRESS_PORT_CRON_SERVICE", "localhost:8087"),
		AddressPostJaeger:      env.GetEnv("ADDRESS_PORT_JAEGER", "localhost:6831"),
		UseTracing:             env.GetEnvAsBool("USE_TRACING", false),
		TokenBearer:            env.GetEnv("TOKEN_BEARER", ""),
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
		LoadSettings: LoadSettings{
			LoadSettingsFromConsul:       env.GetEnvAsBool("LOAD_SETTINGS_FROM_CONSUL", false),
			ConsulServerAddres:           env.GetEnv("CONSUL_SERVER_ADDRESS", ""),
			LoadSettingsFromDockerSecret: env.GetEnvAsBool("LOAD_SETTINGS_FROM_DOCKER_SECRET", false),
			ArrayDockerSecretKey:         env.GetEnv("ARRAY_DOCKER_SECRET_KEY", ""),
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

func (s *ServiceConfig) LoadSettingsFromConsul() {
	if s.LoadSettings.LoadSettingsFromConsul {
		err := s.GetSettingsFromConsul()
		if err != nil {
			panic(err)
		}
		go func(s *ServiceConfig) {
			for {
				time.Sleep(time.Minute * 30)
				err := s.GetSettingsFromConsul()
				if err != nil {
					//panic(err)
					fmt.Printf("Ошибка загрузки настроек конфига из консула %s", err.Error())
				}
				//fmt.Printf("Conf.LoadSettings.ConsulUpdateIndex = %d\n", Conf.LoadSettings.ConsulUpdateIndex)
			}
		}(s)

	}
}

func (s *ServiceConfig) LoadSettingsFromDockerSecrets() {

	if s.LoadSettings.LoadSettingsFromDockerSecret {
		err := s.GetSettingsFromDockerSecrets()
		if err != nil {
			panic(err)
		}
	}

}

func (s *ServiceConfig) GetSettingsFromConsul() error {

	var consulClient *consul.Client

	consulConf := consul.DefaultConfig()
	consulConf.Address = s.LoadSettings.ConsulServerAddres

	var err error
	consulClient, err = consul.NewClient(consulConf)
	if err != nil {
		return err
	}

	qo := &consul.QueryOptions{
		WaitIndex: s.LoadSettings.ConsulUpdateIndex,
		WaitTime:  time.Second * 3,
	}

	// TODO: Тут иногда происходит зависание параметр WaitTime:  time.Second * 10, спасает, но ключи приходят корректно
	// понять причину. Возможно если WaitIndex устанавливать большим чем текущий в консуле идет эта ошибка
	kvPairs, qm, err := consulClient.KV().List("GoKeeper", qo)
	//kvPairs, qm, err := consulClient.KV().List("", qo)
	if err != nil {
		return err
	}

	//_ = qm

	//fmt.Println("remoute consul last index", qm.LastIndex)
	if s.LoadSettings.ConsulUpdateIndex == qm.LastIndex {
		return nil
	}

	//fmt.Printf("qm: %v\n", qm)

	//newConfig := make(map[string]string)

	//GlobalSettingsReturn := Global_settings{}

	PrifixGroup := "GoKeeper/"

	reflectGlobalSettings := reflect.ValueOf(s)
	reflectElem := reflectGlobalSettings.Elem()

	//TODO: Алгоритм устанавливает для полей структура в структуре, но что будет если поля в основной и вложенной структтуре совпадают?

	for _, item := range kvPairs {
		if item.Key == PrifixGroup {
			continue
		}
		//fmt.Println(string(item.Key), string(item.Value))
		res := strings.ReplaceAll(string(item.Key), PrifixGroup, "")
		//fmt.Println("res:", res)

		field := reflectElem.FieldByName(res)

		if field.IsValid() {

			switch field.Kind() {
			case reflect.Int:
				{
					ParseIntVariable, _ := strconv.Atoi(string(item.Value))
					field.SetInt(int64(ParseIntVariable))
				}
			case reflect.Bool:
				{
					ParseBoolVariable, _ := strconv.ParseBool(string(item.Value))
					field.SetBool(ParseBoolVariable)
				}
			case reflect.String:
				{
					field.SetString(string(item.Value))
				}
			}

		}
	}

	// Обновляем индекс консула

	s.LoadSettings.ConsulUpdateIndex = qm.LastIndex
	fmt.Printf("Обновлены настройки консула LastIndex: %d | Адрес консула: %s\n", s.LoadSettings.ConsulUpdateIndex, s.LoadSettings.ConsulServerAddres)

	return nil

}

func (s *ServiceConfig) GetSettingsFromDockerSecrets() error {

	files, err := ioutil.ReadDir("/run/secrets")
	if err != nil {
		fmt.Println("Secret error :", err.Error())
		return fmt.Errorf("Secret error :%s/n", err.Error())
	}

	var mapSecrets map[string]string
	mapSecrets = make(map[string]string)

	for _, file := range files {
		//fmt.Println("file.Name() : ", file.Name())

		if file.IsDir() == true {
			fmt.Println("Secret error :", "IsDir")
			continue
		}
		buf, err := ioutil.ReadFile("/run/secrets/" + file.Name())
		if err != nil {
			fmt.Println("Secret error :", err.Error())
			continue
		}

		mapSecrets[file.Name()] = strings.TrimSpace(string(buf))
		//fmt.Println("value : ", strings.TrimSpace(string(buf)))

	}

	reflectGlobalSettings := reflect.ValueOf(s)
	reflectElem := reflectGlobalSettings.Elem()

	Keys := strings.Split(s.LoadSettings.ArrayDockerSecretKey, ",")
	for _, Key := range Keys {
		//fmt.Println("444")

		valueSecret, ok := mapSecrets[strings.TrimSpace(Key)]
		if !ok {
			fmt.Println("Secret error :")
			return fmt.Errorf("Secret error :%s\n", "Нет такого ключа")
		}

		field := reflectElem.FieldByName(strings.TrimSpace(Key))

		if field.IsValid() {

			switch field.Kind() {
			case reflect.Int:
				{
					ParseIntVariable, _ := strconv.Atoi(string(valueSecret))
					field.SetInt(int64(ParseIntVariable))
				}
			case reflect.Bool:
				{
					ParseBoolVariable, _ := strconv.ParseBool(string(valueSecret))
					field.SetBool(ParseBoolVariable)
				}
			case reflect.String:
				{
					field.SetString(string(valueSecret))
				}
			}

		}

	}

	return nil
}
