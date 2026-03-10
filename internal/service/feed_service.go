package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"butterfly.orx.me/core/log"
	"github.com/kongken/go-home/internal/cache"
	"github.com/kongken/go-home/internal/metrics"
	"github.com/kongken/go-home/internal/model"
	"github.com/kongken/go-home/internal/repository"
	"github.com/kongken/go-home/internal/trace"
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
	ctx, span := trace.FeedCreate(ctx, userID, int32(feedType))
	defer span.End()
	
	logger := log.FromContext(ctx)
	logger.Info("creating feed", "user_id", userID, "type", feedType)
	
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
		logger.Error("failed to create feed", "error", err, "user_id", userID)
		trace.RecordError(span, err)
		return nil, err
	}
	
	trace.SetAttributes(span, trace.Str("feed.id", feed.ID))
	metrics.FeedCreated()
	
	logger.Info("feed created", "feed_id", feed.ID, "user_id", userID)
	return feed, nil
}

// Get 获取动态
func (s *feedService) Get(ctx context.Context, id string) (*model.FeedItem, error) {
	logger := log.FromContext(ctx)
	
	feed, err := s.feedRepo.GetByID(ctx, id)
	if err != nil {
		logger.Warn("feed not found", "feed_id", id)
		return nil, ErrFeedNotFound
	}
	
	logger.Debug("feed fetched", "feed_id", id)
	return feed, nil
}

// Delete 删除动态
func (s *feedService) Delete(ctx context.Context, id, userID string) error {
	logger := log.FromContext(ctx)
	logger.Info("deleting feed", "feed_id", id, "user_id", userID)
	
	feed, err := s.feedRepo.GetByID(ctx, id)
	if err != nil {
		logger.Warn("feed not found for delete", "feed_id", id)
		return ErrFeedNotFound
	}

	if feed.UserID != userID {
		logger.Warn("unauthorized feed delete attempt", "feed_id", id, "user_id", userID, "owner_id", feed.UserID)
		return ErrUnauthorized
	}

	if err := s.feedRepo.Delete(ctx, id); err != nil {
		logger.Error("failed to delete feed", "feed_id", id, "error", err)
		return err
	}
	
	// 删除缓存
	if s.feedCache != nil {
		s.feedCache.DeleteHomeFeed(ctx, userID)
	}
	
	logger.Info("feed deleted", "feed_id", id)
	return nil
}

// ListHome 获取首页动态流 (带缓存)
func (s *feedService) ListHome(ctx context.Context, userIDs []string, offset, limit int64) ([]*model.FeedItem, int64, error) {
	logger := log.FromContext(ctx)
	
	// 生成缓存 key
	cacheKey := strings.Join(userIDs, ",") + ":" + string(rune(offset)) + ":" + string(rune(limit))
	
	// 尝试从缓存获取
	if s.feedCache != nil {
		if feeds, err := s.feedCache.GetHomeFeed(ctx, cacheKey, int(offset/limit)+1); err == nil && feeds != nil {
			metrics.CacheHit("feed_home")
			logger.Debug("feed home cache hit", "key", cacheKey)
			return feeds, int64(len(feeds)), nil
		}
		metrics.CacheMiss("feed_home")
		logger.Debug("feed home cache miss", "key", cacheKey)
	}
	
	// 从数据库获取
	feeds, total, err := s.feedRepo.List(ctx, userIDs, offset, limit)
	if err != nil {
		logger.Error("failed to list feeds", "error", err)
		return nil, 0, err
	}
	
	// 写入缓存
	if s.feedCache != nil {
		if err := s.feedCache.SetHomeFeed(ctx, cacheKey, int(offset/limit)+1, feeds, 5*time.Minute); err != nil {
			logger.Warn("failed to set feed cache", "key", cacheKey, "error", err)
		}
	}
	
	logger.Debug("feeds fetched", "count", len(feeds), "total", total)
	return feeds, total, nil
}

// ListByUser 获取用户动态
func (s *feedService) ListByUser(ctx context.Context, userID string, offset, limit int64) ([]*model.FeedItem, int64, error) {
	logger := log.FromContext(ctx)
	
	feeds, total, err := s.feedRepo.ListByUser(ctx, userID, offset, limit)
	if err != nil {
		logger.Error("failed to list user feeds", "user_id", userID, "error", err)
		return nil, 0, err
	}
	
	logger.Debug("user feeds fetched", "user_id", userID, "count", len(feeds))
	return feeds, total, nil
}

// Like 点赞/取消点赞
func (s *feedService) Like(ctx context.Context, id string, delta int32) error {
	logger := log.FromContext(ctx)
	logger.Info("feed like", "feed_id", id, "delta", delta)
	
	if err := s.feedRepo.IncrementStats(ctx, id, "likes_count", delta); err != nil {
		logger.Error("failed to like feed", "feed_id", id, "error", err)
		return err
	}
	
	return nil
}
