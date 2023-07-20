package main

import (
	"os"

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

func NewRedisClient() *redisClient {
	rKey = "keyss"

	rcl = &redisClient{
		redisCL: redis.NewClient(&redis.Options{
			Addr:     redisHost + ":" + redisPort,
			DB:       0,
			Password: "",
		}),
	}

	return rcl
}
