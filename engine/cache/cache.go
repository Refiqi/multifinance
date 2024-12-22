package cache

import (
	"fmt"
	lru "github.com/hashicorp/golang-lru"

	"sync"
	"time"
)

type Adder func() (data interface{}, err error)

type DoubleBufferCache interface {
	Get(key string, adder Adder) (value interface{}, err error)
}

type DoubleBuffer struct {
	cache            *lru.Cache
	cacheRefreshMSec int64
	cacheExpiryMSec  int64
	lock             sync.RWMutex
}

type DoubleBufferLruConfig struct {
	CacheSize        int
	CacheRefreshMSec int64
	CacheExpiryMSec  int64
}

type cacheEntry struct {
	data         interface{}
	insertedTime time.Time
}

func NewDoubleBufferLru(config DoubleBufferLruConfig) DoubleBufferCache {
	cache, err := lru.New(config.CacheSize)
	if err != nil {
		panic(fmt.Errorf("error Cause: %v\n", err))
	}

	return &DoubleBuffer{
		cache:            cache,
		cacheRefreshMSec: config.CacheRefreshMSec,
		cacheExpiryMSec:  config.CacheExpiryMSec}
}

func (ths *DoubleBuffer) Get(key string, adder Adder) (value interface{}, err error) {
	if cacheAge, ok, cacheData := getCache(ths.cache, key); ok {
		if cacheAge < ths.cacheExpiryMSec {
			if cacheAge > ths.cacheRefreshMSec {
				go func() {
					ths.lock.Lock()
					if isCacheExpired(ths.cache, key, ths.cacheRefreshMSec) {
						_, err = ths.refreshCache(adder, key)
						if err != nil {
							fmt.Printf("refresh error. Cause: %v", err)
						}
					}
					ths.lock.Unlock()
				}()
			}
			return cacheData, nil
		} else {
			ths.cache.Remove(key)
		}
	}

	return ths.refreshCache(adder, key)
}

func getCache(cache *lru.Cache, key string) (cacheAge int64, ok bool, data interface{}) {
	cacheVal, ok := cache.Get(key)
	if ok {
		cache := cacheVal.(cacheEntry)
		cacheAge = int64(time.Since(cache.insertedTime) / time.Millisecond)

		data = cache.data
		return
	}

	return
}

func isCacheExpired(cache *lru.Cache, key string, cacheRefreshMSec int64) bool {
	if cacheAge, ok, _ := getCache(cache, key); ok {
		return cacheAge > cacheRefreshMSec
	}
	return false
}

func (ths *DoubleBuffer) refreshCache(adder Adder, key string) (value interface{}, err error) {
	data, err := adder()
	if err != nil {
		return nil, err
	}

	ths.add(key, cacheEntry{
		data:         data,
		insertedTime: time.Now()})

	return data, nil
}

func (ths *DoubleBuffer) add(key string, value interface{}) {
	ths.cache.Add(key, value)
}