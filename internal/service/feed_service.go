package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/kongken/go-home/internal/cache"
	"github.com/kongken/go-home/internal/model"
	"github.com/kongken/go-home/internal/repository"
)

var (
	ErrFeedNotFound = errors.New("feed not found")
)

// FeedService 动态服务接口
type FeedService interface {
	Create(ctx context.Context, userID string, feedType model.FeedType, content, targetID, targetType string, attachments []model.MediaAttachment, privacy model.PrivacyLevel) (*model.FeedItem, error)
	Get(ctx context.Context, id string) (*model.FeedItem, error)
	Delete(ctx context.Context, id, userID string) error
	ListHome(ctx context.Context, userIDs []string, offset, limit int64) ([]*model.FeedItem, int64, error)
	ListByUser(ctx context.Context, userID string, offset, limit int64) ([]*model.FeedItem, int64, error)
	Like(ctx context.Context, id string, delta int32) error
}

// feedService 动态服务实现
type feedService struct {
	feedRepo  repository.FeedRepository
	feedCache *cache.FeedCache
}

// NewFeedService 创建动态服务
func NewFeedService(feedRepo repository.FeedRepository, redisCache *cache.RedisCache) FeedService {
	return &feedService{
		feedRepo:  feedRepo,
		feedCache: redisCache.FeedCache(),
	}
}

// Create 创建动态
func (s *feedService) Create(ctx context.Context, userID string, feedType model.FeedType, content, targetID, targetType string, attachments []model.MediaAttachment, privacy model.PrivacyLevel) (*model.FeedItem, error) {
	feed := &model.FeedItem{
		UserID:      userID,
		Type:        feedType,
		Content:     content,
		TargetID:    targetID,
		TargetType:  targetType,
		Attachments: attachments,
		Privacy:     privacy,
	}

	if err := s.feedRepo.Create(ctx, feed); err != nil {
		return nil, err
	}

	return feed, nil
}

// Get 获取动态
func (s *feedService) Get(ctx context.Context, id string) (*model.FeedItem, error) {
	return s.feedRepo.GetByID(ctx, id)
}

// Delete 删除动态
func (s *feedService) Delete(ctx context.Context, id, userID string) error {
	feed, err := s.feedRepo.GetByID(ctx, id)
	if err != nil {
		return ErrFeedNotFound
	}

	if feed.UserID != userID {
		return ErrUnauthorized
	}

	return s.feedRepo.Delete(ctx, id)
}

// ListHome 获取首页动态流 (带缓存)
func (s *feedService) ListHome(ctx context.Context, userIDs []string, offset, limit int64) ([]*model.FeedItem, int64, error) {
	// 生成缓存 key
	cacheKey := strings.Join(userIDs, ",") + ":" + string(rune(offset)) + ":" + string(rune(limit))
	
	// 尝试从缓存获取
	if s.feedCache != nil {
		if feeds, err := s.feedCache.GetHomeFeed(ctx, cacheKey, int(offset/limit)+1); err == nil && feeds != nil {
			return feeds, int64(len(feeds)), nil
		}
	}
	
	// 从数据库获取
	feeds, total, err := s.feedRepo.List(ctx, userIDs, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	
	// 写入缓存
	if s.feedCache != nil {
		s.feedCache.SetHomeFeed(ctx, cacheKey, int(offset/limit)+1, feeds, 5*time.Minute)
	}
	
	return feeds, total, nil
}

// ListByUser 获取用户动态
func (s *feedService) ListByUser(ctx context.Context, userID string, offset, limit int64) ([]*model.FeedItem, int64, error) {
	return s.feedRepo.ListByUser(ctx, userID, offset, limit)
}

// Like 点赞/取消点赞
func (s *feedService) Like(ctx context.Context, id string, delta int32) error {
	return s.feedRepo.IncrementStats(ctx, id, "likes_count", delta)
}
