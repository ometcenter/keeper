package redis

import (
	"context"
	"fmt"
	"time"

	libraryGoRedis "github.com/go-redis/redis/v8"
	libraryRediGo "github.com/gomodule/redigo/redis"
)

var PoolRedisRediGolibrary *libraryRediGo.Pool

var ctxGoRedisLibrary = context.Background()
var RedisClientGoRedisLibrary *libraryGoRedis.Client

//"github.com/gomodule/redigo/redis"
func IntiClientLibraryRediGo(AddressPort string) (*libraryRediGo.Pool, error) {

	Pool := &libraryRediGo.Pool{

		MaxIdle:     3,
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

	PoolRedisRediGolibrary = Pool

	return Pool, nil
}

//"github.com/gomodule/redigo/redis"
func SelectLibraryRediGo(Pool *libraryRediGo.Pool, RedisDB int) error {

	conn := Pool.Get()
	defer conn.Close()

	//_, err := conn.Do("SET", key, value)
	_, err := conn.Do("SELECT", RedisDB) // 10 секунд
	if err != nil {
		return err
	}

	return nil
}

//"github.com/gomodule/redigo/redis"
func SetLibraryRediGo(Pool *libraryRediGo.Pool, key string, value interface{}, RedisDB int, TTL int) error {

	conn := Pool.Get()
	defer conn.Close()

	_, err := conn.Do("SELECT", RedisDB) // 10 секунд
	if err != nil {
		return err
	}

	//_, err := conn.Do("SET", key, value)
	_, err = conn.Do("SET", key, value, "EX", "100") // 10 секунд
	if err != nil {
		// v := string(value)
		// if len(v) > 15 {
		// 	v = v[0:12] + "..."
		// }
		// return fmt.Errorf("error setting key %s to %s: %v", key, v, err)
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
			return err
		}
	}

	return nil
}

//"github.com/gomodule/redigo/redis"
func GetLibraryRediGo(Pool *libraryRediGo.Pool, key string, RedisDB int) (string, error) {

	conn := Pool.Get()
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

//"github.com/go-redis/redis/v8"
func IntiClientLibraryGoRedis(AddressPort string) (*libraryGoRedis.Client, error) {

	rdb := libraryGoRedis.NewClient(&libraryGoRedis.Options{
		Addr:     AddressPort,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	//_, err := rdb.Ping(ctxRedis).Result()
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	RedisClientGoRedisLibrary = rdb

	return rdb, err
}

//"github.com/go-redis/redis/v8"
func GetLibraryGoRedis(RedisClient *libraryGoRedis.Client, Key string, RedisDB int) (string, error) {

	var Result string

	_, err := RedisClient.Do(context.Background(), "select", RedisDB).Result()
	if err != nil {
		return Result, err
	}

	val, err := RedisClient.Get(context.Background(), Key).Result()
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

//"github.com/go-redis/redis/v8"
func SetLibraryGoRedis(RedisClient *libraryGoRedis.Client, Key string, Value interface{}, RedisDB int, TTLsec int) error {

	_, err := RedisClient.Do(context.Background(), "select", RedisDB).Result()
	if err != nil {
		return err
	}

	if TTLsec == 0 {
		err = RedisClient.Set(context.Background(), Key, Value, 0).Err()
	} else {
		err = RedisClient.Set(context.Background(), Key, Value, time.Second*5).Err()
	}
	if err != nil {
		return err
	}

	return nil
}

//"github.com/go-redis/redis/v8"
func FlushdbLibraryGoRedis(RedisClient *libraryGoRedis.Client, RedisDB int) error {

	_, err := RedisClient.Do(context.Background(), "select", RedisDB).Result()
	if err != nil {
		return err
	}
	_, err = RedisClient.Do(context.Background(), "flushdb").Result()
	if err != nil {
		return err
	}

	return nil
}

//"github.com/go-redis/redis/v8"
func SelectLibraryGoRedis(RedisClient *libraryGoRedis.Client, RedisDB int) error {

	_, err := RedisClient.Do(context.Background(), "select", RedisDB).Result()
	if err != nil {
		return err
	}

	return nil
}
