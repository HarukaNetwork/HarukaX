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

	"github.com/HarukaNetwork/HarukaX/harukax/modules/utils/error_handling"
	"github.com/allegro/bigcache"
)

var CACHE *bigcache.BigCache

func InitCache() {
	config := bigcache.Config{
		Shards:             1024,
		LifeWindow:         30 * 24 * time.Hour, // one month
		CleanWindow:        5 * time.Minute,
		MaxEntriesInWindow: 1000 * 10 * 60,
		MaxEntrySize:       500,
		HardMaxCacheSize:   512,
		OnRemove:           nil,
		OnRemoveWithReason: nil,
	}
	cache, err := bigcache.NewBigCache(config)
	error_handling.HandleErr(err)
	CACHE = cache
}
