package globalvariables

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	libraryGoRedis "github.com/go-redis/redis/v8"
	"github.com/ometcenter/keeper/config"
	log "github.com/ometcenter/keeper/logging"
)

type GlobalVariablesConnector struct {
	commandChannel                   chan string
	globalVariablesMapMu             sync.RWMutex
	globalVariablesMap               map[string]interface{}
	connectRedisClientGoRedisLibrary *libraryGoRedis.Client
	ctx                              context.Context
	durationTicker                   time.Duration
	ctxCancelFn                      func()
	saveMapToExternalStorage         func()
}

var GlobalVariablesConnectorVb *GlobalVariablesConnector

func NewGlobalVariablesConnector(durationTicker time.Duration) *GlobalVariablesConnector {
	ctx, cancel := context.WithCancel(context.Background())

	// redislibraries := make(map[string]string)
	// redislibraries["LibraryRediGo"] = "LibraryRediGo"
	// redislibraries["LibraryGoRedis"] = "LibraryGoRedis"

	//currentLibary:            "LibraryRediGo",

	var sayHelloWorld = func() {
		//fmt.Println("Hello World !")
	}

	if durationTicker == 0 {
		durationTicker = time.Second * 30
	}

	return &GlobalVariablesConnector{
		commandChannel: make(chan string),
		// out:            make(chan interface{}, 10),
		globalVariablesMap:       make(map[string]interface{}),
		ctx:                      ctx,
		ctxCancelFn:              cancel,
		saveMapToExternalStorage: sayHelloWorld,
		durationTicker:           durationTicker,
	}
}

func (t *GlobalVariablesConnector) Run() error {

	err := t.ConnectRedis(config.Conf.RedisAddressPort)
	if err != nil {
		return err
	}

	tkExpiring := time.NewTicker(t.durationTicker)
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
			t.RefreshAllGlobalVariables()

			//fmt.Printf("Global variables: %v\n", t.globalVariablesMap)

			//t.SetValueForGlobalVariable("currentTime", time.Now())
		}
		//time.Sleep(time.Second * 10)

		//t.Stop()
	}
}

func (r *GlobalVariablesConnector) ConnectRedis(AddressPort string) error {

	rdb := libraryGoRedis.NewClient(&libraryGoRedis.Options{
		Addr:     AddressPort,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	//_, err := rdb.Ping(ctxRedis).Result()
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return err
	}

	_, err = rdb.Do(context.Background(), "select", 10).Result()
	if err != nil {
		return err
	}

	r.connectRedisClientGoRedisLibrary = rdb

	return nil
}

func (r *GlobalVariablesConnector) RefreshAllGlobalVariables() error {

	var keysAnswer []string

	RedisClient := r.connectRedisClientGoRedisLibrary

	// СПОСОБ №1
	var cursor uint64
	for {
		var keys []string
		var err error
		keys, cursor, err = RedisClient.Scan(context.Background(), cursor, "", 0).Result() // Scan(ctx, cursor, "prefix:*", 0).Result()
		if err != nil {
			return err
		}

		for _, key := range keys {
			//fmt.Println("key", key)
			keysAnswer = append(keysAnswer, key)
		}

		if cursor == 0 { // no more keys
			break
		}
		//fmt.Println("new iterate --- ")
	}

	for _, value := range keysAnswer {

		//var Result string
		val, err := r.connectRedisClientGoRedisLibrary.Get(context.Background(), value).Result()
		if err == libraryGoRedis.Nil {
			//fmt.Println("key2 does not exist")
			//return Result, fmt.Errorf("Не найден ключ для JobId: %s в Redis", InsuranceNumber)
			continue
		} else if err != nil {
			//panic(err)
			log.Impl.Error(err)
			continue
			//return Result, nil
		} else {
			//fmt.Println("key2", val2)
			//return val, nil
		}

		r.globalVariablesMapMu.Lock()
		r.globalVariablesMap[value] = val
		r.globalVariablesMapMu.Unlock()
	}

	return nil

}

func (r *GlobalVariablesConnector) SetValueForGlobalVariable(Key string, Value interface{}) error {

	RedisClient := r.connectRedisClientGoRedisLibrary

	err := RedisClient.Set(context.Background(), Key, Value, 0).Err()
	if err != nil {
		return err
	}

	r.globalVariablesMapMu.Lock()
	r.globalVariablesMap[Key] = Value
	r.globalVariablesMapMu.Unlock()

	return nil

}
