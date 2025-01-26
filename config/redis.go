package config

import (
	"log"

	"github.com/go-redis/redis"
)

func initRedis() {
	RedisClient := redis.NewClient(&redis.Options{
		Addr:     AppConfig.Database.Dsn,
		Password: "",
		DB:       0,
	})

	_, err := RedisClient.Ping().Result()

	if err != nil {
		log.Fatalf("failed to connect redis")
	}

}
