package repository

import (
	"context"

	"github.com/kongken/go-home/internal/model"
	"gorm.io/gorm"
)

// BlogRepository 博客仓库接口
type BlogRepository interface {
	Create(ctx context.Context, blog *model.Blog) error
	GetByID(ctx context.Context, id string) (*model.Blog, error)
	Update(ctx context.Context, blog *model.Blog) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, userID string, category string, offset, limit int) ([]*model.Blog, int64, error)
	ListByUser(ctx context.Context, userID string, offset, limit int) ([]*model.Blog, int64, error)
}

// blogRepository 博客仓库实现
type blogRepository struct {
	db *gorm.DB
}

// NewBlogRepository 创建博客仓库
func NewBlogRepository(db *gorm.DB) BlogRepository {
	return &blogRepository{db: db}
}

// Create 创建博客
func (r *blogRepository) Create(ctx context.Context, blog *model.Blog) error {
	return r.db.WithContext(ctx).Create(blog).Error
}

// GetByID 根据ID获取博客
func (r *blogRepository) GetByID(ctx context.Context, id string) (*model.Blog, error) {
	var blog model.Blog
	if err := r.db.WithContext(ctx).Preload("User").First(&blog, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &blog, nil
}

// Update 更新博客
func (r *blogRepository) Update(ctx context.Context, blog *model.Blog) error {
	return r.db.WithContext(ctx).Save(blog).Error
}

// Delete 删除博客
func (r *blogRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.Blog{}, "id = ?", id).Error
}

// List 获取博客列表
func (r *blogRepository) List(ctx context.Context, userID string, category string, offset, limit int) ([]*model.Blog, int64, error) {
	var blogs []*model.Blog
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Blog{}).Where("status = ?", model.BlogStatusPublished)

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("User").Offset(offset).Limit(limit).Order("created_at DESC").Find(&blogs).Error; err != nil {
		return nil, 0, err
	}

	return blogs, total, nil
}

// ListByUser 获取指定用户的博客列表
func (r *blogRepository) ListByUser(ctx context.Context, userID string, offset, limit int) ([]*model.Blog, int64, error) {
	return r.List(ctx, userID, "", offset, limit)
}
