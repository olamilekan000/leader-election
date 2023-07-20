package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

var redisHost string = os.Getenv("REDIS_HOST")
var redisPort string = os.Getenv("REDIS_PORT")

var rcl *redisClient
var rKey string

type redisClient struct {
	redisCL *redis.Client
}

func (rc *redisClient) Set(ctx context.Context, key, value string) error {
	if err := rc.redisCL.Set(ctx, key, value, 0).Err(); err != nil {
		return err
	}

	return nil
}

func (rc *redisClient) Get(ctx context.Context, key string) (*string, error) {
	val, err := rc.redisCL.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	return &val, err
}

func (rc *redisClient) Ping(ctx context.Context) (string, error) {
	return rcl.redisCL.Ping(ctx).Result()
}

func NewRedisClient() *redisClient {
	rKey = "process_count"

	rcl = &redisClient{
		redisCL: redis.NewClient(&redis.Options{
			Addr:     redisHost + ":" + redisPort,
			DB:       0,
			Password: "",
		}),
	}

	return rcl
}

func processData(ctx context.Context, processor string) {
	var rVal *string
	var num int
	var err error

	for {
		time.Sleep(10 * time.Second)

		fmt.Println("Pod " + processor + " " + " is processing workload")

		rVal, err = rcl.Get(ctx, rKey)
		if errors.Is(err, redis.Nil) {
			if err := rcl.Set(ctx, rKey, "0"); err != nil {
				log.Fatalln(err.Error())
			}
			continue
		}

		if err != nil {
			log.Fatalln("error fetching data from redis: ", err.Error())
		}

		if rVal != nil {
			num, err = strconv.Atoi(*rVal)
			if err != nil {
				fmt.Println("Error: convert", err)
			}
		}

		num += 1
		fmt.Println("NUM IS: ", num)

		err := rcl.Set(ctx, rKey, strconv.Itoa(num))
		if err != nil {
			log.Fatalln("error persisting data: ", err.Error())
		}
	}
}
