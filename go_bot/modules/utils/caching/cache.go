package caching

import (
	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/error_handling"
	"github.com/allegro/bigcache"
	"time"
)

var CACHE *bigcache.BigCache

func InitCache() {
	config := bigcache.Config{Shards: 1024,
		LifeWindow:         2 * time.Hour,
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
