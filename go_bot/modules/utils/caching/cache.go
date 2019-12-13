package caching

import (
	"time"

	"github.com/ATechnoHazard/ginko/go_bot/modules/utils/error_handling"
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
