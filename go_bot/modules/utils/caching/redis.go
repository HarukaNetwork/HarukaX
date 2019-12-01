package caching

import (
	"github.com/ATechnoHazard/ginko/go_bot"
	"github.com/go-redis/redis"
	"time"
)

var REDIS *redis.Client

func InitRedis() {
	client := redis.NewClient(&redis.Options{
		Addr:         go_bot.BotConfig.RedisAddress,
		Password:     go_bot.BotConfig.RedisPassword,
		DB:           0,
		DialTimeout:  time.Second,
		MinIdleConns: 0,
	})
	REDIS = client
}
