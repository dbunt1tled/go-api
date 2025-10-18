package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/dbunt1tled/go-api/internal/config/env"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	rdb *redis.Client
}

var (
	redisCache *RedisCache //nolint:gochecknoglobals // singleton
	m          sync.Once   //nolint:gochecknoglobals // singleton
)

func GetRedisCache() *RedisCache {
	if redisCache == nil {
		var err error
		cfg := env.GetConfigInstance()
		redisCache, err = NewRedisCache(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password, cfg.Redis.DB)
		if err != nil {
			panic(err)
		}
	}
	return redisCache
}

func NewRedisCache(host, port, password string, db int) (*RedisCache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:        fmt.Sprintf("%v:%v", host, port),
		Password:    password,
		DB:          db,
		ReadTimeout: 0, // 3 second
	})

	return &RedisCache{
		rdb: rdb,
	}, nil
}

func (c *RedisCache) Set(key string, value string, expSecond int32) error {
	ctx := context.Background()

	return c.rdb.Set(ctx, key, value, time.Duration(expSecond)*time.Second).Err()
}

func (c *RedisCache) Get(key string) (string, error) {
	ctx := context.Background()

	result, err := c.rdb.Get(ctx, key).Result()
	if err == nil {
		return result, nil
	}

	if errors.Is(err, redis.Nil) {
		return "", errors.New("key not found")
	}

	return "", err
}

func (c *RedisCache) Delete(key string) error {
	ctx := context.Background()

	err := c.rdb.Del(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *RedisCache) Close() error {
	return c.rdb.Close()
}
