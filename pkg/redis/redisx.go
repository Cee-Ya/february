package redisx

import (
	"context"
	"february/common"
	"february/common/consts"
	"february/entity"
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
	"strings"
)

//var Redis interface {
//	Exists(ctx context.Context, keys ...string) *redis.IntCmd
//	Del(ctx context.Context, keys ...string) *redis.IntCmd
//	Get(ctx context.Context, key string) *redis.StringCmd
//	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
//	HGetAll(ctx context.Context, key string) *redis.MapStringStringCmd
//	HExists(ctx context.Context, key, field string) *redis.BoolCmd
//	HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
//	HGet(ctx context.Context, key, fields string) *redis.StringCmd
//	HDel(ctx context.Context, key string, fields ...string) *redis.IntCmd
//	Close() error
//	Ping(ctx context.Context) *redis.StatusCmd
//	Publish(ctx context.Context, channel string, message interface{}) *redis.IntCmd
//}

func InitRedis(cfg entity.RedisConfig) error {
	var redisClient entity.Redis
	fmt.Println("redis mode:", cfg.RedisType)
	switch cfg.RedisType {
	case "standalone", "":
		redisOptions := &redis.Options{
			Addr:     cfg.Addr,
			Username: cfg.Username,
			Password: cfg.Password,
			DB:       cfg.DB,
		}

		if cfg.UseTLS {
			tlsConfig, err := cfg.TLSConfig()
			if err != nil {
				fmt.Println("failed to init redis tls config:", err)
				os.Exit(1)
			}
			redisOptions.TLSConfig = tlsConfig
		}

		redisClient = redis.NewClient(redisOptions)
	case "cluster":
		redisOptions := &redis.ClusterOptions{
			Addrs:    strings.Split(cfg.Addr, consts.COMMA),
			Username: cfg.Username,
			Password: cfg.Password,
		}

		if cfg.UseTLS {
			tlsConfig, err := cfg.TLSConfig()
			if err != nil {
				fmt.Println("failed to init redis tls config:", err)
				os.Exit(1)
			}
			redisOptions.TLSConfig = tlsConfig
		}

		redisClient = redis.NewClusterClient(redisOptions)

	case "sentinel":
		redisOptions := &redis.FailoverOptions{
			MasterName:       cfg.MasterName,
			SentinelAddrs:    strings.Split(cfg.Addr, consts.COMMA),
			Username:         cfg.Username,
			Password:         cfg.Password,
			DB:               cfg.DB,
			SentinelUsername: cfg.SentinelUsername,
			SentinelPassword: cfg.SentinelPassword,
		}

		if cfg.UseTLS {
			tlsConfig, err := cfg.TLSConfig()
			if err != nil {
				fmt.Println("failed to init redis tls config:", err)
				os.Exit(1)
			}
			redisOptions.TLSConfig = tlsConfig
		}

		redisClient = redis.NewFailoverClient(redisOptions)

	default:
		fmt.Println("failed to init redis , redis type is illegal:", cfg.RedisType)
		os.Exit(1)
	}

	err := redisClient.Ping(context.Background()).Err()
	if err != nil {
		fmt.Println("failed to ping redis:", err)
		os.Exit(1)
	}
	common.Redisx = redisClient
	return nil
}
