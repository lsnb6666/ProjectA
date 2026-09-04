package config

import (
	"ProjectA/global"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

func initRedis() {
	Logger.Info("Initializing Redis connection...")

	RedisClient := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		DB:       0,
		Password: "",
	})
	_, err := RedisClient.Ping().Result()

	if err != nil {
		Logger.Fatal("Failed to connect to Redis", zap.Error(err), zap.String("addr", "localhost:6379"))
	}

	global.RedisDB = RedisClient
	Logger.Info("Redis connected successfully")
}
