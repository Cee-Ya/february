package entity

import (
	"fmt"
	"sync"
	"time"
)

type MemoryCacheItem struct {
	Value      interface{}
	Expiration int64
}

type MemoryCache struct {
	items     map[string]MemoryCacheItem
	mu        sync.RWMutex
	expiryChs map[string]chan struct{}
	onExpired func(string)
}

func NewCache(onExpired func(string)) *MemoryCache {
	return &MemoryCache{
		items:     make(map[string]MemoryCacheItem),
		expiryChs: make(map[string]chan struct{}),
		onExpired: onExpired,
	}
}

func (c *MemoryCache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	expiration := time.Now().Add(ttl).UnixNano()
	c.items[key] = MemoryCacheItem{
		Value:      value,
		Expiration: expiration,
	}
	if ch, found := c.expiryChs[key]; found {
		close(ch)
	}
	ch := make(chan struct{})
	c.expiryChs[key] = ch
	go c.expiryListener(key, ch, expiration)
}

func (c *MemoryCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, found := c.items[key]
	if !found || item.Expiration < time.Now().UnixNano() {
		return nil, false
	}
	return item.Value, true
}

func (c *MemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
	if ch, found := c.expiryChs[key]; found {
		close(ch)
		delete(c.expiryChs, key)
	}
}

func (c *MemoryCache) expiryListener(key string, ch chan struct{}, expiration int64) {
	<-ch
	if time.Now().UnixNano() > expiration {
		c.Delete(key)
		fmt.Println("delete key:", key)
		fmt.Println("onExpired:", c.onExpired != nil)
		if c.onExpired != nil {
			c.onExpired(key)
		}
	}
}
