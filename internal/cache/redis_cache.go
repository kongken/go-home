package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"butterfly.orx.me/core/store/redis"
	"github.com/kongken/go-home/internal/model"
)

// RedisCache Redis 缓存服务
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache 创建 Redis 缓存服务
func NewRedisCache() *RedisCache {
	client := redis.GetClient("cache")
	if client == nil {
		// 如果没有配置，返回 nil，表示不使用缓存
		return nil
	}
	return &RedisCache{client: client}
}

// UserCache 用户缓存
func (c *RedisCache) UserCache() *UserCache {
	if c == nil {
		return nil
	}
	return &UserCache{client: c.client}
}

// BlogCache 博客缓存
func (c *RedisCache) BlogCache() *BlogCache {
	if c == nil {
		return nil
	}
	return &BlogCache{client: c.client}
}

// FeedCache 动态缓存
func (c *RedisCache) FeedCache() *FeedCache {
	if c == nil {
		return nil
	}
	return &FeedCache{client: c.client}
}

// UserCache 用户缓存
type UserCache struct {
	client *redis.Client
}

// Get 获取用户缓存
func (c *UserCache) Get(ctx context.Context, userID string) (*model.User, error) {
	if c == nil {
		return nil, nil
	}
	val, err := c.client.Get(ctx, fmt.Sprintf("user:%s", userID)).Result()
	if err != nil {
		return nil, err
	}
	var user model.User
	if err := json.Unmarshal([]byte(val), &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// Set 设置用户缓存
func (c *UserCache) Set(ctx context.Context, user *model.User, ttl time.Duration) error {
	if c == nil {
		return nil
	}
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, fmt.Sprintf("user:%s", user.ID), data, ttl).Err()
}

// Delete 删除用户缓存
func (c *UserCache) Delete(ctx context.Context, userID string) error {
	if c == nil {
		return nil
	}
	return c.client.Del(ctx, fmt.Sprintf("user:%s", userID)).Err()
}

// BlogCache 博客缓存
type BlogCache struct {
	client *redis.Client
}

// Get 获取博客缓存
func (c *BlogCache) Get(ctx context.Context, blogID string) (*model.Blog, error) {
	if c == nil {
		return nil, nil
	}
	val, err := c.client.Get(ctx, fmt.Sprintf("blog:%s", blogID)).Result()
	if err != nil {
		return nil, err
	}
	var blog model.Blog
	if err := json.Unmarshal([]byte(val), &blog); err != nil {
		return nil, err
	}
	return &blog, nil
}

// Set 设置博客缓存
func (c *BlogCache) Set(ctx context.Context, blog *model.Blog, ttl time.Duration) error {
	if c == nil {
		return nil
	}
	data, err := json.Marshal(blog)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, fmt.Sprintf("blog:%s", blog.ID), data, ttl).Err()
}

// Delete 删除博客缓存
func (c *BlogCache) Delete(ctx context.Context, blogID string) error {
	if c == nil {
		return nil
	}
	return c.client.Del(ctx, fmt.Sprintf("blog:%s", blogID)).Err()
}

// DeleteUserBlogs 删除用户博客列表缓存
func (c *BlogCache) DeleteUserBlogs(ctx context.Context, userID string) error {
	if c == nil {
		return nil
	}
	// 使用 pattern 删除所有该用户的博客缓存
	// 注意：生产环境应该使用更精确的方式
	return c.client.Del(ctx, fmt.Sprintf("user_blogs:%s", userID)).Err()
}

// FeedCache 动态缓存
type FeedCache struct {
	client *redis.Client
}

// GetHomeFeed 获取首页动态缓存
func (c *FeedCache) GetHomeFeed(ctx context.Context, userID string, page int) ([]*model.FeedItem, error) {
	if c == nil {
		return nil, nil
	}
	val, err := c.client.Get(ctx, fmt.Sprintf("feed:home:%s:%d", userID, page)).Result()
	if err != nil {
		return nil, err
	}
	var feeds []*model.FeedItem
	if err := json.Unmarshal([]byte(val), &feeds); err != nil {
		return nil, err
	}
	return feeds, nil
}

// SetHomeFeed 设置首页动态缓存
func (c *FeedCache) SetHomeFeed(ctx context.Context, userID string, page int, feeds []*model.FeedItem, ttl time.Duration) error {
	if c == nil {
		return nil
	}
	data, err := json.Marshal(feeds)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, fmt.Sprintf("feed:home:%s:%d", userID, page), data, ttl).Err()
}

// DeleteHomeFeed 删除首页动态缓存
func (c *FeedCache) DeleteHomeFeed(ctx context.Context, userID string) error {
	if c == nil {
		return nil
	}
	// 删除该用户的所有首页动态缓存
	pattern := fmt.Sprintf("feed:home:%s:*", userID)
	keys, err := c.client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}
	if len(keys) > 0 {
		return c.client.Del(ctx, keys...).Err()
	}
	return nil
}

// GeneralCache 通用缓存方法
type GeneralCache struct {
	client *redis.Client
}

// General 获取通用缓存
func (c *RedisCache) General() *GeneralCache {
	if c == nil {
		return nil
	}
	return &GeneralCache{client: c.client}
}

// Set 设置缓存
func (c *GeneralCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if c == nil {
		return nil
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, ttl).Err()
}

// Get 获取缓存
func (c *GeneralCache) Get(ctx context.Context, key string, dest interface{}) error {
	if c == nil {
		return nil
	}
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

// Delete 删除缓存
func (c *GeneralCache) Delete(ctx context.Context, key string) error {
	if c == nil {
		return nil
	}
	return c.client.Del(ctx, key).Err()
}

// Exists 检查缓存是否存在
func (c *GeneralCache) Exists(ctx context.Context, key string) (bool, error) {
	if c == nil {
		return false, nil
	}
	n, err := c.client.Exists(ctx, key).Result()
	return n > 0, err
}

// TTL 获取缓存过期时间
func (c *GeneralCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	if c == nil {
		return 0, nil
	}
	return c.client.TTL(ctx, key).Result()
}