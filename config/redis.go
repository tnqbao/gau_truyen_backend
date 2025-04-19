package config

import (
	"context"
	"os"
	"strconv"
	"sync"

	"github.com/go-redis/redis/v8"
)

var (
	redisClient *redis.Client
	once        sync.Once
)

func InitRedis() {
	once.Do(func() {
		redisAddr := os.Getenv("REDIS_ADDRESS")
		if redisAddr == "" {
			redisAddr = "127.0.0.1:6379"
		}
		redisPassword := os.Getenv("REDIS_PASSWORD")
		redisDB := os.Getenv("REDIS_DB")
		if redisDB == "" {
			redisDB = "1"
		}
		db, err := strconv.Atoi(redisDB)
		if err != nil {
			panic("Invalid REDIS_DB value: " + redisDB)
		}
		client := redis.NewClient(&redis.Options{
			Addr:     redisAddr,
			Password: redisPassword,
			DB:       db,
		})

		ctx := context.Background()
		_, err = client.Ping(ctx).Result()
		if err != nil {
			panic("Failed to connect to Redis: " + err.Error())
		}

		redisClient = client
	})
}

func GetRedisClient() *redis.Client {
	return redisClient
}
