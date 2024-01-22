package entity

import (
	"fmt"
	"github.com/pkg/errors"
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
}

func NewMemoryCache() MemoryCache {
	return MemoryCache{
		items:     make(map[string]MemoryCacheItem),
		expiryChs: make(map[string]chan struct{}),
		mu:        sync.RWMutex{},
	}
}

func (c *MemoryCache) Exist(key string) (bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, found := c.items[key]
	if !found || item.Expiration < time.Now().UnixNano() {
		return false, nil
	}
	return true, nil
}

func (c *MemoryCache) Set(key string, value interface{}, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	expiration := time.Now().Add(ttl).UnixNano()
	// 数据序列化为json字符串
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
	return nil
}

func (c *MemoryCache) Get(key string) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, found := c.items[key]
	if !found || item.Expiration < time.Now().UnixNano() {
		return nil, errors.New("not record")
	}
	return item.Value, nil
}

func (c *MemoryCache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
	if ch, found := c.expiryChs[key]; found {
		close(ch)
		delete(c.expiryChs, key)
	}
	return nil
}

func (c *MemoryCache) expiryListener(key string, ch chan struct{}, expiration int64) {
	<-ch
	if time.Now().UnixNano() > expiration {
		c.Delete(key)
		fmt.Println("delete key:", key)
		c.onExpired(key)
	}
}

func (c *MemoryCache) onExpired(key string) {
	fmt.Println("onExpired:", key)
}
