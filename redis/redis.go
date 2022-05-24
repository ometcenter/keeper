package redis

import (
	"context"
	"fmt"
	"time"

	libraryGoRedis "github.com/go-redis/redis/v8"
	libraryRedigo "github.com/gomodule/redigo/redis"
)

var PoolRedisRediGolibrary *libraryRedigo.Pool

var ctxGoRedisLibrary = context.Background()
var RedisClientGoRedisLibrary *libraryGoRedis.Client

func IntiRedisClientRediGo(AddressPort string) (*libraryRedigo.Pool, error) {

	Pool := &libraryRedigo.Pool{

		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,

		Dial: func() (libraryRedigo.Conn, error) {
			c, err := libraryRedigo.Dial("tcp", AddressPort)
			if err != nil {
				return nil, err
			}
			return c, err
		},

		TestOnBorrow: func(c libraryRedigo.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	PoolRedisRediGolibrary = Pool

	return Pool, nil
}

func IntiRedisClient(AddressPort string) (*libraryGoRedis.Client, error) {

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

//"github.com/gomodule/redigo/redis"
func SelectLibraryGoRedis(Pool *libraryRedigo.Pool, BaseNumber int) error {

	conn := Pool.Get()
	defer conn.Close()

	//_, err := conn.Do("SET", key, value)
	_, err := conn.Do("SELECT", BaseNumber) // 10 секунд
	if err != nil {
		return err
	}

	return nil
}

//"github.com/gomodule/redigo/redis"
func SetRedigo(Pool *libraryRedigo.Pool, key string, value []byte, TTL int) error {

	conn := Pool.Get()
	defer conn.Close()

	//_, err := conn.Do("SET", key, value)
	_, err := conn.Do("SET", key, value, "EX", "100") // 10 секунд
	if err != nil {
		v := string(value)
		if len(v) > 15 {
			v = v[0:12] + "..."
		}
		return fmt.Errorf("error setting key %s to %s: %v", key, v, err)
	}

	// // Установить время истечения 24 часа
	// //n, _ := conn.Do("EXPIRE", key, 24*3600)
	if TTL != 0 {
		_, err := conn.Do("EXPIRE", key, TTL)
		// if n == int64(1) {
		// 	fmt.Println("success: ", n)
		// }
		if err != nil {
			v := string(value)
			if len(v) > 15 {
				v = v[0:12] + "..."
			}
			return fmt.Errorf("error EXPIRE key %s to %s: %v", key, v, err)
		}
	}

	return nil
}

//"github.com/gomodule/redigo/redis"
func GetRedigo(Pool *libraryRedigo.Pool, key string) error {

	conn := Pool.Get()
	defer conn.Close()

	var data []byte
	data, err := libraryRedigo.Bytes(conn.Do("GET", key))
	if err != nil {
		ErrReturn := fmt.Errorf("error getting key %s: %v", key, err)
		fmt.Println(ErrReturn)
		return ErrReturn
	}
	fmt.Println(string(data))
	return nil
}

//"github.com/go-redis/redis/v8"
func GetLibraryGoRedis(Key string, RedisDB int, RedisClient *libraryGoRedis.Client) (string, error) {

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
func SetLibraryGoRedis(Key string, RedisDB int, RedisClient *libraryGoRedis.Client) error {

	_, err := RedisClient.Do(context.Background(), "select", 12).Result()
	if err != nil {
		return err
	}

	err = RedisClient.Set(context.Background(), "TestSetGoRedis", "777", time.Second*5).Err()
	if err != nil {
		return err
	}

	return nil
}

//"github.com/go-redis/redis/v8"
func flushdbLibraryGoRedis(Key string, RedisDB int, RedisClient *libraryGoRedis.Client) error {

	_, err := RedisClient.Do(context.Background(), "select", 12).Result()
	if err != nil {
		return err
	}
	_, err = RedisClient.Do(context.Background(), "flushdb").Result()
	if err != nil {
		return err
	}

	return nil
}
