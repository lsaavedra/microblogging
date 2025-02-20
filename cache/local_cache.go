package cache

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/pkg/errors"
)

var (
	errorKeyNotFound = errors.New("key not found")
)

type CacheItem struct {
	Value      []byte
	Expiration int64
}

type LocalCache struct {
	items map[string]CacheItem
	mutex sync.RWMutex
}

func NewLocalCache() Cache {
	return &LocalCache{
		items: make(map[string]CacheItem),
	}
}

func (c *LocalCache) Set(key string, value []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.items[key] = CacheItem{
		Value: value,
	}
	return nil
}

func (c *LocalCache) Get(key string) ([]byte, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, found := c.items[key]
	if !found {
		return nil, errorKeyNotFound
	}
	return item.Value, nil
}

func (c *LocalCache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.items, key)
}

func (c *LocalCache) Cleanup() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	now := time.Now().UnixNano()
	for key, item := range c.items {
		if now > item.Expiration {
			delete(c.items, key)
		}
	}
}

func DecodeCacheValue[T any](data []byte, v *T) error {
	return json.Unmarshal(data, v)
}

func EncodeCacheValue[T any](v T) ([]byte, error) {
	return json.Marshal(v)
}
