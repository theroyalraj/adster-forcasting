package db

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/theroyalraj/adster-forcasting/internal/config"
	"github.com/theroyalraj/adster-forcasting/internal/utils"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func InitRedis(cfg *config.Config) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		utils.Log.Fatal("Failed to connect to Redis:", err)
	}

	utils.Log.Info("Connected to Redis")
}
