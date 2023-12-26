package db

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"

	"exampleApi/helpers"
)

var REDIS_HOST = helpers.GetEnv("REDIS_HOST")
var REDIS_PORT = helpers.GetEnv("REDIS_PORT")
var REDIS_PASSWORD = helpers.GetEnv("REDIS_PASSWORD")

func ConnectRedis() *redis.Client {
	adress := fmt.Sprintf("%v:%v", REDIS_HOST, REDIS_PORT)

	client := redis.NewClient(&redis.Options{
		Addr:     adress, //redis port
		Password: REDIS_PASSWORD,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Panicf("Redis connection: %v", err)
	}

	return client
}
