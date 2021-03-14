package common

import (
	"github.com/sta-golang/go-lib-utils/cache"
	"github.com/sta-golang/go-lib-utils/cache/memory"
)

var globalCache cache.Cache

func InitMemoryCache() error {
	globalCache = memory.New(memory.NewConfig(16, 60, 5000, 512))
	return nil
}

func GlobalCache() cache.Cache {
	return globalCache
}
