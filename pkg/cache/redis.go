package cache

import (
	"context"
	"fmt"

	"github.com/kongken/go-home/internal/config"
	"github.com/redis/go-redis/v9"
)

// RedisClient Redis客户端
type RedisClient struct {
	client *redis.Client
}

// NewRedisClient 创建 Redis 客户端
func NewRedisClient(cfg *config.RedisConfig) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisClient{client: client}, nil
}

// GetClient 获取底层 Redis 客户端
func (r *RedisClient) GetClient() *redis.Client {
	return r.client
}

// Close 关闭连接
func (r *RedisClient) Close() error {
	return r.client.Close()
}
