package service

import (
	"context"
	"errors"

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
	feedRepo repository.FeedRepository
}

// NewFeedService 创建动态服务
func NewFeedService(feedRepo repository.FeedRepository) FeedService {
	return &feedService{feedRepo: feedRepo}
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

// ListHome 获取首页动态流
func (s *feedService) ListHome(ctx context.Context, userIDs []string, offset, limit int64) ([]*model.FeedItem, int64, error) {
	return s.feedRepo.List(ctx, userIDs, offset, limit)
}

// ListByUser 获取用户动态
func (s *feedService) ListByUser(ctx context.Context, userID string, offset, limit int64) ([]*model.FeedItem, int64, error) {
	return s.feedRepo.ListByUser(ctx, userID, offset, limit)
}

// Like 点赞/取消点赞
func (s *feedService) Like(ctx context.Context, id string, delta int32) error {
	return s.feedRepo.IncrementStats(ctx, id, "likes_count", delta)
}
