package utils

import (
	"context"
	"february/common"
	"february/entity"
	"time"
)

var expire time.Duration

type RedisUtils struct {
	ctx context.Context
	entity.Redis
}

func NewRedisUtils(ctx context.Context) *RedisUtils {
	expire = common.GlobalConfig.Redis.KeyExpire * time.Second
	return &RedisUtils{ctx, common.Redisx}
}

// Lock 分布式锁
func (r *RedisUtils) Lock(key string, value interface{}) (bool, error) {
	return r.Redis.SetNX(r.ctx, key, value, expire).Result()
}

// Unlock 分布式锁
func (r *RedisUtils) Unlock(key string) error {
	return r.Redis.Del(r.ctx, key).Err()
}

func (r *RedisUtils) MustSet(key string, value interface{}) error {
	return r.Redis.Set(r.ctx, key, value, expire).Err()
}

func (r *RedisUtils) MustGet(key string) string {
	return r.Redis.Get(r.ctx, key).Val()
}

func (r *RedisUtils) MustDel(key string) error {
	return r.Redis.Del(r.ctx, key).Err()
}

func (r *RedisUtils) MustExists(key string) bool {
	return r.Redis.Exists(r.ctx, key).Val() == 1
}
