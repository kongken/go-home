package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"butterfly.orx.me/core/log"
	"github.com/kongken/go-home/internal/cache"
	"github.com/kongken/go-home/internal/metrics"
	"github.com/kongken/go-home/internal/model"
	"github.com/kongken/go-home/internal/repository"
	"github.com/kongken/go-home/internal/trace"
)

var (
	ErrBlogNotFound = errors.New("blog not found")
	ErrUnauthorized = errors.New("unauthorized")
)

// BlogService 博客服务接口
type BlogService interface {
	Create(ctx context.Context, userID string, title, content, summary, coverImage string, tags []string, category string, privacy model.PrivacyLevel, status model.BlogStatus) (*model.Blog, error)
	Get(ctx context.Context, id string) (*model.Blog, error)
	Update(ctx context.Context, id, userID string, updates map[string]interface{}) (*model.Blog, error)
	Delete(ctx context.Context, id, userID string) error
	List(ctx context.Context, userID, category string, tags []string, offset, limit int) ([]*model.Blog, int64, error)
	ListByUser(ctx context.Context, userID string, offset, limit int) ([]*model.Blog, int64, error)
}

// blogService 博客服务实现
type blogService struct {
	blogRepo  repository.BlogRepository
	blogCache *cache.BlogCache
}

// NewBlogService 创建博客服务
func NewBlogService(blogRepo repository.BlogRepository, redisCache *cache.RedisCache) BlogService {
	return &blogService{
		blogRepo:  blogRepo,
		blogCache: redisCache.BlogCache(),
	}
}

// Create 创建博客
func (s *blogService) Create(ctx context.Context, userID string, title, content, summary, coverImage string, tags []string, category string, privacy model.PrivacyLevel, status model.BlogStatus) (*model.Blog, error) {
	ctx, span := trace.BlogCreate(ctx, userID, "")
	defer span.End()
	
	logger := log.FromContext(ctx)
	logger.Info("creating blog", "user_id", userID, "title", title)
	
	tagsJSON, err := json.Marshal(tags)
	if err != nil {
		logger.Error("failed to marshal tags", "error", err)
		trace.RecordError(span, err)
		return nil, err
	}

	blog := &model.Blog{
		UserID:     userID,
		Title:      title,
		Content:    content,
		Summary:    summary,
		CoverImage: coverImage,
		Tags:       string(tagsJSON),
		Category:   category,
		Privacy:    privacy,
		Status:     status,
	}

	if err := s.blogRepo.Create(ctx, blog); err != nil {
		logger.Error("failed to create blog", "error", err, "user_id", userID)
		trace.RecordError(span, err)
		return nil, err
	}
	
	trace.SetAttributes(span, trace.Str("blog.id", blog.ID))
	metrics.BlogCreated()
	
	logger.Info("blog created", "blog_id", blog.ID, "user_id", userID)
	return blog, nil
}

// Get 获取博客 (带缓存)
func (s *blogService) Get(ctx context.Context, id string) (*model.Blog, error) {
	logger := log.FromContext(ctx)
	
	// 先查缓存
	if s.blogCache != nil {
		if blog, err := s.blogCache.Get(ctx, id); err == nil && blog != nil {
			metrics.CacheHit("blog")
			logger.Debug("blog cache hit", "blog_id", id)
			return blog, nil
		}
		metrics.CacheMiss("blog")
		logger.Debug("blog cache miss", "blog_id", id)
	}
	
	// 查数据库
	blog, err := s.blogRepo.GetByID(ctx, id)
	if err != nil {
		logger.Warn("blog not found", "blog_id", id)
		return nil, ErrBlogNotFound
	}
	
	// 写入缓存
	if s.blogCache != nil {
		if err := s.blogCache.Set(ctx, blog, time.Hour); err != nil {
			logger.Warn("failed to set blog cache", "blog_id", id, "error", err)
		}
	}
	
	logger.Debug("blog fetched", "blog_id", id)
	return blog, nil
}

// Update 更新博客
func (s *blogService) Update(ctx context.Context, id, userID string, updates map[string]interface{}) (*model.Blog, error) {
	logger := log.FromContext(ctx)
	logger.Info("updating blog", "blog_id", id, "user_id", userID)
	
	blog, err := s.blogRepo.GetByID(ctx, id)
	if err != nil {
		logger.Warn("blog not found for update", "blog_id", id)
		return nil, ErrBlogNotFound
	}

	// 检查权限
	if blog.UserID != userID {
		logger.Warn("unauthorized blog update attempt", "blog_id", id, "user_id", userID, "owner_id", blog.UserID)
		return nil, ErrUnauthorized
	}

	// 更新字段
	if title, ok := updates["title"].(string); ok {
		blog.Title = title
	}
	if content, ok := updates["content"].(string); ok {
		blog.Content = content
	}
	if summary, ok := updates["summary"].(string); ok {
		blog.Summary = summary
	}
	if coverImage, ok := updates["cover_image"].(string); ok {
		blog.CoverImage = coverImage
	}
	if category, ok := updates["category"].(string); ok {
		blog.Category = category
	}
	if tags, ok := updates["tags"].([]string); ok {
		tagsJSON, err := json.Marshal(tags)
		if err != nil {
			logger.Error("failed to marshal tags", "error", err)
			return nil, err
		}
		blog.Tags = string(tagsJSON)
	}
	if privacy, ok := updates["privacy"].(model.PrivacyLevel); ok {
		blog.Privacy = privacy
	}
	if status, ok := updates["status"].(model.BlogStatus); ok {
		blog.Status = status
	}

	if err := s.blogRepo.Update(ctx, blog); err != nil {
		logger.Error("failed to update blog", "blog_id", id, "error", err)
		return nil, err
	}
	
	// 删除缓存
	if s.blogCache != nil {
		s.blogCache.Delete(ctx, id)
	}
	
	logger.Info("blog updated", "blog_id", id)
	return blog, nil
}

// Delete 删除博客
func (s *blogService) Delete(ctx context.Context, id, userID string) error {
	logger := log.FromContext(ctx)
	logger.Info("deleting blog", "blog_id", id, "user_id", userID)
	
	blog, err := s.blogRepo.GetByID(ctx, id)
	if err != nil {
		logger.Warn("blog not found for delete", "blog_id", id)
		return ErrBlogNotFound
	}

	// 检查权限
	if blog.UserID != userID {
		logger.Warn("unauthorized blog delete attempt", "blog_id", id, "user_id", userID, "owner_id", blog.UserID)
		return ErrUnauthorized
	}

	if err := s.blogRepo.Delete(ctx, id); err != nil {
		logger.Error("failed to delete blog", "blog_id", id, "error", err)
		return err
	}
	
	// 删除缓存
	if s.blogCache != nil {
		s.blogCache.Delete(ctx, id)
	}
	
	logger.Info("blog deleted", "blog_id", id)
	return nil
}

// List 获取博客列表
func (s *blogService) List(ctx context.Context, userID, category string, tags []string, offset, limit int) ([]*model.Blog, int64, error) {
	return s.blogRepo.List(ctx, userID, category, offset, limit)
}

// ListByUser 获取指定用户的博客列表
func (s *blogService) ListByUser(ctx context.Context, userID string, offset, limit int) ([]*model.Blog, int64, error) {
	return s.blogRepo.ListByUser(ctx, userID, offset, limit)
}
