package db

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"

	"exampleApi/helpers"
)

func ConnectRedis() *redis.Client {
	host := helpers.GetEnv("REDIS_HOST")
	port := helpers.GetEnv("REDIS_PORT")
	password := helpers.GetEnv("REDIS_PASSWORD")

	adress := fmt.Sprintf("%v:%v", host, port)

	client := redis.NewClient(&redis.Options{
		Addr:     adress, //redis port
		Password: password,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Panicf("Redis connection: %v", err)
	}

	return client
}
