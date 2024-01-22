package service

import (
	"context"
	"february/common"
	"february/common/consts"
	"february/entity"
	"time"
)

// CacheInterface 缓存接口
type CacheInterface interface {
	Exist(key string) (bool, error)
	Set(key string, value interface{}, ttl time.Duration) error
	Get(key string) (interface{}, error)
	Delete(key string) error
}

// CacheConfigInterface 缓存接口
type CacheConfigInterface interface {
	EnableCache() bool //是否开启缓存
	CacheKey() string  //缓存key
}

// CacheService 缓存服务
type CacheService struct {
	ctx context.Context
	CacheConfigInterface
	CacheInterface
}

// NewCacheService 创建缓存服务
func NewCacheService(ctx context.Context, cacheEntity CacheConfigInterface) *CacheService {
	if cacheEntity.EnableCache() {
		switch common.GlobalConfig.Cache.CacheType {
		case consts.CacheTypeMemory:
			cacheEntity = NewCacheMemoryService(ctx, cacheEntity)
			break
		case consts.CacheTypeRedis:
			cacheEntity = NewCacheRedisService(ctx, cacheEntity)
			break
		}
	}
	return &CacheService{CacheConfigInterface: cacheEntity, ctx: ctx}
}

func (c *CacheService) EnableCache() bool {
	return c.CacheConfigInterface.EnableCache()
}

func (c *CacheService) CacheKey() string {
	return c.CacheConfigInterface.CacheKey()
}

func (c *CacheService) Exist(key string) (bool, error) {
	return c.CacheInterface.Exist(c.CacheKey() + key)
}

func (c *CacheService) Set(key string, value interface{}, ttl time.Duration) error {
	return c.CacheInterface.Set(c.CacheKey()+key, value, ttl)
}

func (c *CacheService) Get(key string) (interface{}, error) {
	return c.CacheInterface.Get(c.CacheKey() + key)
}

func (c *CacheService) Delete(key string) error {
	return c.CacheInterface.Delete(c.CacheKey() + key)
}

// CacheMemoryService 缓存内存服务
type CacheMemoryService struct {
	CacheService
	entity.MemoryCache
}

func NewCacheMemoryService(ctx context.Context, cacheEntity CacheConfigInterface) *CacheMemoryService {
	return &CacheMemoryService{MemoryCache: entity.NewMemoryCache(), CacheService: CacheService{CacheConfigInterface: cacheEntity, ctx: ctx}}
}

func (c *CacheMemoryService) Exist(key string) (bool, error) {
	return c.MemoryCache.Exist(key)
}

func (c *CacheMemoryService) Set(key string, value interface{}, ttl time.Duration) error {
	return c.MemoryCache.Set(key, value, ttl)
}

func (c *CacheMemoryService) Get(key string) (interface{}, error) {
	return c.MemoryCache.Get(key)
}

func (c *CacheMemoryService) Delete(key string) error {
	return c.MemoryCache.Delete(key)
}

// CacheRedisService 缓存redis服务
type CacheRedisService struct {
	CacheService
	entity.Redis
}

func NewCacheRedisService(ctx context.Context, cacheEntity CacheConfigInterface) *CacheRedisService {
	return &CacheRedisService{Redis: common.RedisCache, CacheService: CacheService{CacheConfigInterface: cacheEntity, ctx: ctx}}
}

func (c *CacheRedisService) Exist(key string) (bool, error) {
	res, err := c.Redis.Exists(c.ctx, key).Result()
	if err != nil {
		return false, err
	}
	return res > 0, nil
}

func (c *CacheRedisService) Set(key string, value interface{}, ttl time.Duration) error {
	return c.Redis.Set(c.ctx, key, value, ttl).Err()
}

func (c *CacheRedisService) Get(key string) (interface{}, error) {
	return c.Redis.Get(c.ctx, key).Result()
}

func (c *CacheRedisService) Delete(key string) error {
	return c.Redis.Del(c.ctx, key).Err()
}
