package redis

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	//libraryGoRedis "github.com/go-redis/redis/v8"
	libraryGoRedis "github.com/go-redis/redis/v8"
	libraryRediGo "github.com/gomodule/redigo/redis"
	"github.com/ometcenter/keeper/config"
	log "github.com/ometcenter/keeper/logging"
	"github.com/ometcenter/keeper/models"
)

// var PoolRedisRediGolibrary *libraryRediGo.Pool

// var ctxGoRedisLibrary = context.Background()
// var RedisClientGoRedisLibrary *libraryGoRedis.Client

// // "github.com/gomodule/redigo/redis"
// func IntiClientLibraryRediGo(AddressPort string) (*libraryRediGo.Pool, error) {

// 	Pool := &libraryRediGo.Pool{

// 		MaxIdle:     10,
// 		IdleTimeout: 240 * time.Second,

// 		Dial: func() (libraryRediGo.Conn, error) {
// 			c, err := libraryRediGo.Dial("tcp", AddressPort)
// 			if err != nil {
// 				return nil, err
// 			}
// 			return c, err
// 		},

// 		TestOnBorrow: func(c libraryRediGo.Conn, t time.Time) error {
// 			_, err := c.Do("PING")
// 			return err
// 		},
// 	}

// 	PoolRedisRediGolibrary = Pool

// 	return Pool, nil
// }

// // "github.com/gomodule/redigo/redis"
// func SelectLibraryRediGo(Pool *libraryRediGo.Pool, RedisDB int) error {

// 	conn := Pool.Get()
// 	defer conn.Close()

// 	//_, err := conn.Do("SET", key, value)
// 	_, err := conn.Do("SELECT", RedisDB) // 10 секунд
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// // "github.com/gomodule/redigo/redis"
// func SetLibraryRediGo(Pool *libraryRediGo.Pool, key string, value interface{}, RedisDB int, TTL int64) error {

// 	conn := Pool.Get()
// 	defer conn.Close()

// 	_, err := conn.Do("SELECT", RedisDB) // 10 секунд
// 	if err != nil {
// 		fmt.Println("Auth err --- _, err := conn.Do(SELECT, RedisDB) // 10 секунд ----", err)
// 		return err
// 	}

// 	//_, err := conn.Do("SET", key, value) --- _, err = conn.Do("SET", key, value, "EX", "100") // 10 секунд
// 	_, err = conn.Do("SET", key, value)
// 	if err != nil {
// 		// v := string(value)
// 		// if len(v) > 15 {
// 		// 	v = v[0:12] + "..."
// 		// }
// 		// return fmt.Errorf("error setting key %s to %s: %v", key, v, err)
// 		fmt.Println("Auth err --- _, err = conn.Do(SET, key, value) ----", err)
// 		return err
// 	}

// 	// // Установить время истечения 24 часа
// 	// //n, _ := conn.Do("EXPIRE", key, 24*3600)
// 	if TTL != 0 {
// 		_, err := conn.Do("EXPIRE", key, TTL)
// 		// if n == int64(1) {
// 		// 	fmt.Println("success: ", n)
// 		// }
// 		if err != nil {
// 			// v := string(value)
// 			// if len(v) > 15 {
// 			// 	v = v[0:12] + "..."
// 			// }
// 			// return fmt.Errorf("error EXPIRE key %s to %s: %v", key, v, err)
// 			fmt.Println("Auth err --- _, err := conn.Do(EXPIRE, key, TTL) ----", err)
// 			return err
// 		}
// 	}

// 	return nil
// }

// // "github.com/gomodule/redigo/redis"
// func GetLibraryRediGo(Pool *libraryRediGo.Pool, key string, RedisDB int) (string, error) {

// 	conn := Pool.Get()
// 	defer conn.Close()

// 	_, err := conn.Do("SELECT", RedisDB) // 10 секунд
// 	if err != nil {
// 		return "", err
// 	}

// 	//var data []byte
// 	//data, err = libraryRediGo.Bytes(conn.Do("GET", key))
// 	data, err := libraryRediGo.String(conn.Do("GET", key))
// 	if err != nil {
// 		ErrReturn := fmt.Errorf("error getting key %s: %v", key, err)
// 		fmt.Println(ErrReturn)
// 		return "", ErrReturn
// 	}

// 	return data, nil
// }

// // "github.com/go-redis/redis/v8"
// func IntiClientLibraryGoRedis(AddressPort string) (*libraryGoRedis.Client, error) {

// 	rdb := libraryGoRedis.NewClient(&libraryGoRedis.Options{
// 		Addr:     AddressPort,
// 		Password: "", // no password set
// 		DB:       0,  // use default DB
// 	})

// 	//_, err := rdb.Ping(ctxRedis).Result()
// 	_, err := rdb.Ping(context.Background()).Result()
// 	if err != nil {
// 		return nil, err
// 	}

// 	RedisClientGoRedisLibrary = rdb

// 	return rdb, err
// }

// // "github.com/go-redis/redis/v8"
// func GetLibraryGoRedis(RedisClient *libraryGoRedis.Client, Key string, RedisDB int) (string, error) {

// 	var Result string

// 	_, err := RedisClient.Do(context.Background(), "select", RedisDB).Result()
// 	if err != nil {
// 		return Result, err
// 	}

// 	val, err := RedisClient.Get(context.Background(), Key).Result()
// 	if err == libraryGoRedis.Nil {
// 		//fmt.Println("key2 does not exist")
// 		//return Result, fmt.Errorf("Не найден ключ для JobId: %s в Redis", InsuranceNumber)
// 		return Result, nil
// 	} else if err != nil {
// 		//panic(err)
// 		return Result, err
// 		//return Result, nil
// 	} else {
// 		//fmt.Println("key2", val2)
// 		return val, nil
// 	}

// }

// // "github.com/go-redis/redis/v8"
// func SetLibraryGoRedis(RedisClient *libraryGoRedis.Client, Key string, Value interface{}, RedisDB int, TTLsec int64) error {

// 	_, err := RedisClient.Do(context.Background(), "select", RedisDB).Result()
// 	if err != nil {
// 		return err
// 	}

// 	if TTLsec == 0 {
// 		err = RedisClient.Set(context.Background(), Key, Value, 0).Err()
// 	} else {
// 		err = RedisClient.Set(context.Background(), Key, Value, time.Second*time.Duration(TTLsec)).Err()
// 	}
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// // "github.com/go-redis/redis/v8"
// func FlushdbLibraryGoRedis(RedisClient *libraryGoRedis.Client, RedisDB int) error {

// 	_, err := RedisClient.Do(context.Background(), "select", RedisDB).Result()
// 	if err != nil {
// 		return err
// 	}
// 	_, err = RedisClient.Do(context.Background(), "flushdb").Result()
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// // "github.com/go-redis/redis/v8"
// func SelectLibraryGoRedis(RedisClient *libraryGoRedis.Client, RedisDB int) error {

// 	_, err := RedisClient.Do(context.Background(), "select", RedisDB).Result()
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type RedisConnector struct {
	commandChannel                       chan string
	connectPoolRedisRediGolibrary        *libraryRediGo.Pool
	connectRedisClientGoRedisLibrary     *libraryGoRedis.Client
	currentLibary                        string
	connectPool                          map[string]interface{}
	connectPoolRedisClientGoRedisLibrary map[int]*libraryGoRedis.Client
	activeAreas                          []int
	activeTokensMu                       sync.RWMutex
	redislibraries                       map[string]string
	ctx                                  context.Context
	ctxCancelFn                          func()
	saveMapToExternalStorage             func()
}

var RedisConnectorVb *RedisConnector

func NewRedisConnector(redislibraries map[string]string, currentLibary string, activeAreas []int) *RedisConnector {
	ctx, cancel := context.WithCancel(context.Background())

	// redislibraries := make(map[string]string)
	// redislibraries["LibraryRediGo"] = "LibraryRediGo"
	// redislibraries["LibraryGoRedis"] = "LibraryGoRedis"

	//currentLibary:            "LibraryRediGo",

	var sayHelloWorld = func() {
		//fmt.Println("Hello World !")
	}

	return &RedisConnector{
		commandChannel: make(chan string),
		// out:            make(chan interface{}, 10),
		connectPool:                          make(map[string]interface{}),
		connectPoolRedisClientGoRedisLibrary: make(map[int]*libraryGoRedis.Client),
		activeAreas:                          activeAreas,
		currentLibary:                        currentLibary,
		redislibraries:                       redislibraries,
		ctx:                                  ctx,
		ctxCancelFn:                          cancel,
		saveMapToExternalStorage:             sayHelloWorld,
	}
}

func (t *RedisConnector) Run() error {

	err := t.Connect()
	if err != nil {
		return err
	}

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

func (r *RedisConnector) ChangeCurrentLibrary(currentLibary string) error {

	r.currentLibary = currentLibary

	return nil

}

func (r *RedisConnector) Connect() error {

	for key, _ := range r.redislibraries {
		if key == "LibraryRediGo" {
			err := r.IntiClientLibraryRediGo(config.Conf.RedisAddressPort)
			if err != nil {
				return err
			}
		}
		if key == "LibraryGoRedis" {
			err := r.IntiClientLibraryGoRedis(config.Conf.RedisAddressPort)
			if err != nil {
				return err
			}
		}
	}

	return nil

}

func (r *RedisConnector) IntiClientLibraryRediGo(AddressPort string) error {

	Pool := &libraryRediGo.Pool{

		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,

		Dial: func() (libraryRediGo.Conn, error) {
			c, err := libraryRediGo.Dial("tcp", AddressPort)
			if err != nil {
				return nil, err
			}
			return c, err
		},

		TestOnBorrow: func(c libraryRediGo.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	r.connectPoolRedisRediGolibrary = Pool
	r.connectPool["LibraryRediGo"] = Pool

	return nil
}

// "github.com/go-redis/redis/v8"
func (r *RedisConnector) IntiClientLibraryGoRedis(AddressPort string) error {

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

	r.connectRedisClientGoRedisLibrary = rdb
	r.connectPool["LibraryGoRedis"] = rdb

	if len(r.activeAreas) > 0 {
		for _, value := range r.activeAreas {
			AreaValue := value

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

			_, err = rdb.Do(context.Background(), "select", AreaValue).Result()
			if err != nil {
				return err
			}

			r.connectPoolRedisClientGoRedisLibrary[AreaValue] = rdb

		}
	}

	return err
}

func (r *RedisConnector) Select(RedisDB int) error {

	if r.currentLibary == "LibraryRediGo" {
		err := r.SelectLibraryRediGo(RedisDB)
		if err != nil {
			return err
		}
	} else if r.currentLibary == "LibraryGoRedis" {
		err := r.SelectLibraryGoRedis(RedisDB)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *RedisConnector) Set(Key string, Value interface{}, RedisDB int, TTLsec int64) error {

	if r.currentLibary == "LibraryRediGo" {
		err := r.SetLibraryRediGo(Key, Value, RedisDB, TTLsec)
		if err != nil {
			return err
		}
	} else if r.currentLibary == "LibraryGoRedis" {
		err := r.SetLibraryGoRedis(Key, Value, RedisDB, TTLsec)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *RedisConnector) GetCurrentConnection() interface{} {

	if r.currentLibary == "LibraryRediGo" {
		return r.connectPoolRedisRediGolibrary
	} else if r.currentLibary == "LibraryGoRedis" {
		return r.connectRedisClientGoRedisLibrary
	}

	return nil
}

func (r *RedisConnector) GetCurrentConnectionByArea(Area int) (*libraryGoRedis.Client, error) {

	if r.currentLibary == "LibraryRediGo" {
		// not implemented
		return nil, fmt.Errorf("not implemented for LibraryRediGo")
	} else if r.currentLibary == "LibraryGoRedis" {

		value, ok := r.connectPoolRedisClientGoRedisLibrary[Area]
		if !ok {
			return nil, fmt.Errorf("connection not found for %s", value)
		}
		return value, nil
	}

	return nil, fmt.Errorf("not implemented for _____")
}

func (r *RedisConnector) Get(Key string, RedisDB int) (string, error) {

	if r.currentLibary == "LibraryRediGo" {
		resultStr, err := r.GetLibraryRediGo(Key, RedisDB)
		if err != nil {
			return "", err
		}
		return resultStr, nil

	} else if r.currentLibary == "LibraryGoRedis" {
		resultStr, err := r.GetLibraryGoRedis(Key, RedisDB)
		if err != nil {
			return "", err
		}
		return resultStr, nil
	}

	return "", nil
}

func (r *RedisConnector) Flushdb(RedisDB int) error {

	if r.currentLibary == "LibraryRediGo" {
		// TODO: Implement this function
		// err := r.FlushdbLibraryGoRedis(RedisDB)
		// if err != nil {
		// 	return err
		// }
	} else if r.currentLibary == "LibraryGoRedis" {
		err := r.FlushdbLibraryGoRedis(RedisDB)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *RedisConnector) Del(RedisDB int, keyForDel string) error {

	if r.currentLibary == "LibraryRediGo" {
		// TODO: Implement this function
		// err := r.DelGoRedis(RedisDB)
		// if err != nil {
		// 	return err
		// }
	} else if r.currentLibary == "LibraryGoRedis" {
		err := r.DelLibraryGoRedis(RedisDB, keyForDel)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *RedisConnector) DelLibraryGoRedis(RedisDB int, keyForDel string) error {

	_, err := r.connectRedisClientGoRedisLibrary.Do(context.Background(), "select", RedisDB).Result()
	if err != nil {
		return err
	}

	_, err = r.connectRedisClientGoRedisLibrary.Del(context.Background(), keyForDel).Result()
	if err != nil {
		return err
	}

	return nil

}

func (r *RedisConnector) SelectLibraryRediGo(RedisDB int) error {

	conn := r.connectPoolRedisRediGolibrary.Get()
	defer conn.Close()

	//_, err := conn.Do("SET", key, value)
	_, err := conn.Do("SELECT", RedisDB) // 10 секунд
	if err != nil {
		return err
	}

	return nil
}

// "github.com/go-redis/redis/v8"
func (r *RedisConnector) SelectLibraryGoRedis(RedisDB int) error {

	_, err := r.connectRedisClientGoRedisLibrary.Do(context.Background(), "select", RedisDB).Result()
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisConnector) FlushdbLibraryGoRedis(RedisDB int) error {

	_, err := r.connectRedisClientGoRedisLibrary.Do(context.Background(), "select", RedisDB).Result()
	if err != nil {
		return err
	}
	_, err = r.connectRedisClientGoRedisLibrary.Do(context.Background(), "flushdb").Result()
	if err != nil {
		return err
	}

	return nil
}

// "github.com/gomodule/redigo/redis"
func (r *RedisConnector) SetLibraryRediGo(key string, value interface{}, RedisDB int, TTL int64) error {

	conn := r.connectPoolRedisRediGolibrary.Get()
	defer conn.Close()

	_, err := conn.Do("SELECT", RedisDB) // 10 секунд
	if err != nil {
		fmt.Println("Auth err --- _, err := conn.Do(SELECT, RedisDB) // 10 секунд ----", err)
		return err
	}

	//_, err := conn.Do("SET", key, value) --- _, err = conn.Do("SET", key, value, "EX", "100") // 10 секунд
	_, err = conn.Do("SET", key, value)
	if err != nil {
		// v := string(value)
		// if len(v) > 15 {
		// 	v = v[0:12] + "..."
		// }
		// return fmt.Errorf("error setting key %s to %s: %v", key, v, err)
		fmt.Println("Auth err --- _, err = conn.Do(SET, key, value) ----", err)
		return err
	}

	// // Установить время истечения 24 часа
	// //n, _ := conn.Do("EXPIRE", key, 24*3600)
	if TTL != 0 {
		_, err := conn.Do("EXPIRE", key, TTL)
		// if n == int64(1) {
		// 	fmt.Println("success: ", n)
		// }
		if err != nil {
			// v := string(value)
			// if len(v) > 15 {
			// 	v = v[0:12] + "..."
			// }
			// return fmt.Errorf("error EXPIRE key %s to %s: %v", key, v, err)
			fmt.Println("Auth err --- _, err := conn.Do(EXPIRE, key, TTL) ----", err)
			return err
		}
	}

	return nil
}

func (r *RedisConnector) PipeSet(RedisDB int, arr []models.PipeArr) error {
	_, err := r.connectRedisClientGoRedisLibrary.Do(context.Background(), "select", RedisDB).Result()
	if err != nil {
		return err
	}

	pipe := r.connectRedisClientGoRedisLibrary.Pipeline()
	ctx := context.Background()
	for i := 0; i < len(arr); i++ {
		cmd := pipe.Set(ctx, arr[i].Key, arr[i].Value, 0)
		if cmd.Err() != nil {
			return cmd.Err()
		}
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

// "github.com/go-redis/redis/v8"
func (r *RedisConnector) SetLibraryGoRedis(Key string, Value interface{}, RedisDB int, TTLsec int64) error {

	_, err := r.connectRedisClientGoRedisLibrary.Do(context.Background(), "select", RedisDB).Result()
	if err != nil {
		return err
	}

	if TTLsec == 0 {
		err = r.connectRedisClientGoRedisLibrary.Set(context.Background(), Key, Value, 0).Err()
	} else {
		err = r.connectRedisClientGoRedisLibrary.Set(context.Background(), Key, Value, time.Second*time.Duration(TTLsec)).Err()
	}
	if err != nil {
		return err
	}

	return nil
}

// "github.com/gomodule/redigo/redis"
func (r *RedisConnector) GetLibraryRediGo(key string, RedisDB int) (string, error) {

	conn := r.connectPoolRedisRediGolibrary.Get()
	defer conn.Close()

	_, err := conn.Do("SELECT", RedisDB) // 10 секунд
	if err != nil {
		return "", err
	}

	//var data []byte
	//data, err = libraryRediGo.Bytes(conn.Do("GET", key))
	data, err := libraryRediGo.String(conn.Do("GET", key))
	if err != nil {
		ErrReturn := fmt.Errorf("error getting key %s: %v", key, err)
		fmt.Println(ErrReturn)
		return "", ErrReturn
	}

	return data, nil
}

// "github.com/go-redis/redis/v8"
func (r *RedisConnector) GetLibraryGoRedis(Key string, RedisDB int) (string, error) {

	var Result string

	_, err := r.connectRedisClientGoRedisLibrary.Do(context.Background(), "select", RedisDB).Result()
	if err != nil {
		return Result, err
	}

	val, err := r.connectRedisClientGoRedisLibrary.Get(context.Background(), Key).Result()
	if err == libraryGoRedis.Nil {
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

func (r *RedisConnector) Stop() {

	for key, value := range r.connectPool {

		if key == "LibraryRediGo" {

			RedisClient, ok := value.(*libraryRediGo.Pool)
			if !ok {
				log.Impl.Error(fmt.Errorf("Internal error convert type failed"))
			}

			err := RedisClient.Close()
			if err != nil {
				log.Impl.Error(err)
			}
		}

		if key == "LibraryGoRedis" {
			RedisClient, ok := value.(*libraryGoRedis.Client)
			if !ok {
				log.Impl.Error(fmt.Errorf("Internal error convert type failed"))
			}

			err := RedisClient.Close()
			if err != nil {
				log.Impl.Error(err)
			}
		}

	}

	r.ctxCancelFn()

	// close(w.out)
}
