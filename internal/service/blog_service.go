package service

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/kongken/go-home/internal/model"
	"github.com/kongken/go-home/internal/repository"
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
	blogRepo repository.BlogRepository
}

// NewBlogService 创建博客服务
func NewBlogService(blogRepo repository.BlogRepository) BlogService {
	return &blogService{blogRepo: blogRepo}
}

// Create 创建博客
func (s *blogService) Create(ctx context.Context, userID string, title, content, summary, coverImage string, tags []string, category string, privacy model.PrivacyLevel, status model.BlogStatus) (*model.Blog, error) {
	tagsJSON, _ := json.Marshal(tags)

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
		return nil, err
	}

	return blog, nil
}

// Get 获取博客
func (s *blogService) Get(ctx context.Context, id string) (*model.Blog, error) {
	return s.blogRepo.GetByID(ctx, id)
}

// Update 更新博客
func (s *blogService) Update(ctx context.Context, id, userID string, updates map[string]interface{}) (*model.Blog, error) {
	blog, err := s.blogRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrBlogNotFound
	}

	// 检查权限
	if blog.UserID != userID {
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
		tagsJSON, _ := json.Marshal(tags)
		blog.Tags = string(tagsJSON)
	}
	if privacy, ok := updates["privacy"].(model.PrivacyLevel); ok {
		blog.Privacy = privacy
	}
	if status, ok := updates["status"].(model.BlogStatus); ok {
		blog.Status = status
	}

	if err := s.blogRepo.Update(ctx, blog); err != nil {
		return nil, err
	}

	return blog, nil
}

// Delete 删除博客
func (s *blogService) Delete(ctx context.Context, id, userID string) error {
	blog, err := s.blogRepo.GetByID(ctx, id)
	if err != nil {
		return ErrBlogNotFound
	}

	// 检查权限
	if blog.UserID != userID {
		return ErrUnauthorized
	}

	return s.blogRepo.Delete(ctx, id)
}

// List 获取博客列表
func (s *blogService) List(ctx context.Context, userID, category string, tags []string, offset, limit int) ([]*model.Blog, int64, error) {
	return s.blogRepo.List(ctx, userID, category, offset, limit)
}

// ListByUser 获取指定用户的博客列表
func (s *blogService) ListByUser(ctx context.Context, userID string, offset, limit int) ([]*model.Blog, int64, error) {
	return s.blogRepo.ListByUser(ctx, userID, offset, limit)
}
