package service

import (
	"context"
	"february/common"
	"february/entity"
)

// CacheInterface 缓存接口
type CacheInterface interface {
	EnableRedis() bool //是否开启缓存
	CacheKey() string  //缓存key
}

type CacheService struct {
	ctx context.Context
	CacheInterface
}

// CacheMemoryService 缓存内存服务
type CacheMemoryService struct {
	CacheService
}

// CacheRedisService 缓存redis服务
type CacheRedisService struct {
	CacheService
	entity.Redis
}

func NewCacheRedisService(ctx context.Context, cacheEntity CacheInterface) *CacheRedisService {
	return &CacheRedisService{Redis: common.Redisx, CacheService: CacheService{CacheInterface: cacheEntity, ctx: ctx}}
}
