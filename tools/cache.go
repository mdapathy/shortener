package tools

import (
	"github.com/allegro/bigcache"
	"log"
	"time"
)

type cache struct {
	Storage *bigcache.BigCache
}

var (
	config = bigcache.Config{
		Shards: 1024,

		LifeWindow:         10 * time.Minute,
		CleanWindow:        5 * time.Minute,
		MaxEntriesInWindow: 1000 * 10 * 60,
		MaxEntrySize:       500,
		Verbose:            true,
		HardMaxCacheSize:   8192,
		OnRemove:           nil,
		OnRemoveWithReason: nil,
	}
)

type CacheStorage interface {
	GetValue(string) (string, error)
	SetValue(string, string) error
	DeleteValue(string) error
}

func NewCacheStorage() CacheStorage {
	storage, err := bigcache.NewBigCache(config)
	if err != nil {
		log.Fatalf("Error allocating cache storage:%s", err)
	}
	return &cache{
		Storage: storage,
	}
}

func (c *cache) GetValue(s string) (string, error) {
	res, err := c.Storage.Get(s)
	return string(res), err
}

//Used after both the redirect and create requests
func (c *cache) SetValue(s string, s2 string) error {
	return c.Storage.Set(s, []byte(s2))
}

func (c *cache) DeleteValue(s string) error {
	return c.Storage.Delete(s)
}
