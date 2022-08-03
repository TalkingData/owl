package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

const (
	defaultKeyPrefix = "/td/owl_v6"
)

type RedisTool struct {
	client      *redis.Client
	serviceName string
}

// NewRedisTool 新建RedisTool
func NewRedisTool(address, password, srvName string, db int) *RedisTool {
	return &RedisTool{
		client: redis.NewClient(&redis.Options{
			Addr:     address,
			Password: password,
			DB:       db,
		}),
		serviceName: srvName,
	}
}

// Close 关闭Redis
func (rt *RedisTool) Close() {
	if rt.client == nil {
		return
	}
	_ = rt.client.Close()
}

func (rt *RedisTool) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return rt.client.Set(ctx, rt.getFinalKey(key), value, expiration).Err()
}

func (rt *RedisTool) Del(ctx context.Context, key string) error {
	return rt.client.Del(ctx, rt.getFinalKey(key)).Err()
}

func (rt *RedisTool) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return rt.client.PExpire(ctx, rt.getFinalKey(key), expiration).Err()
}

func (rt *RedisTool) Get(ctx context.Context, key string) (string, error) {
	return rt.client.Get(ctx, rt.getFinalKey(key)).Result()
}

func (rt *RedisTool) HasKey(ctx context.Context, key string) bool {
	if err := rt.client.Get(ctx, rt.getFinalKey(key)).Err(); err == nil {
		return true
	}
	return false
}

// Deprecated: DirectGet 直接获取redis中的值
func (rt *RedisTool) DirectGet(ctx context.Context, key string) (string, error) {
	return rt.client.Get(ctx, key).Result()
}

// Deprecated: DirectHasKey 直接判断redis中是否存在key
func (rt *RedisTool) DirectHasKey(ctx context.Context, key string) bool {
	if err := rt.client.Get(ctx, key).Err(); err == nil {
		return true
	}
	return false
}

// getFinalKey 生成KeyName
func (rt *RedisTool) getFinalKey(key string) string {
	return defaultKeyPrefix + "/" + rt.serviceName + "/" + key
}
