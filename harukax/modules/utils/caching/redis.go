/*
 *    Copyright © 2020 Haruka Network Development
 *    This file is part of Haruka X.
 *
 *    Haruka X is free software: you can redistribute it and/or modify
 *    it under the terms of the Raphielscape Public License as published by
 *    the Devscapes Open Source Holding GmbH., version 1.d
 *
 *    Haruka X is distributed in the hope that it will be useful,
 *    but WITHOUT ANY WARRANTY; without even the implied warranty of
 *    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *    Devscapes Raphielscape Public License for more details.
 *
 *    You should have received a copy of the Devscapes Raphielscape Public License
 */

package caching

import (
	"time"

	"github.com/HarukaNetwork/HarukaX/harukax"
	"github.com/go-redis/redis"
)

var REDIS *redis.Client

func InitRedis() {
	client := redis.NewClient(&redis.Options{
		Addr:         harukax.BotConfig.RedisAddress,
		Password:     harukax.BotConfig.RedisPassword,
		DB:           0,
		DialTimeout:  time.Second,
		MinIdleConns: 0,
	})
	REDIS = client
}
