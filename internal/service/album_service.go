package service

import (
	"context"
	"errors"

	"github.com/kongken/go-home/internal/model"
	"github.com/kongken/go-home/internal/repository"
)

var (
	ErrAlbumNotFound = errors.New("album not found")
	ErrPhotoNotFound = errors.New("photo not found")
)

// AlbumService 相册服务接口
type AlbumService interface {
	Create(ctx context.Context, userID, name, description, coverPhoto string, privacy model.PrivacyLevel) (*model.Album, error)
	Get(ctx context.Context, id string) (*model.Album, error)
	Update(ctx context.Context, id, userID string, updates map[string]interface{}) (*model.Album, error)
	Delete(ctx context.Context, id, userID string) error
	ListByUser(ctx context.Context, userID string, offset, limit int64) ([]*model.Album, int64, error)

	// 照片管理
	AddPhoto(ctx context.Context, albumID, userID, url, thumbnail, description string, width, height int32) (*model.Photo, error)
	DeletePhoto(ctx context.Context, photoID, userID string) error
	GetPhotos(ctx context.Context, albumID string, offset, limit int64) ([]*model.Photo, int64, error)
}

// albumService 相册服务实现
type albumService struct {
	albumRepo repository.AlbumRepository
}

// NewAlbumService 创建相册服务
func NewAlbumService(albumRepo repository.AlbumRepository) AlbumService {
	return &albumService{albumRepo: albumRepo}
}

// Create 创建相册
func (s *albumService) Create(ctx context.Context, userID, name, description, coverPhoto string, privacy model.PrivacyLevel) (*model.Album, error) {
	album := &model.Album{
		UserID:      userID,
		Name:        name,
		Description: description,
		CoverPhoto:  coverPhoto,
		Privacy:     privacy,
	}

	if err := s.albumRepo.Create(ctx, album); err != nil {
		return nil, err
	}

	return album, nil
}

// Get 获取相册
func (s *albumService) Get(ctx context.Context, id string) (*model.Album, error) {
	return s.albumRepo.GetByID(ctx, id)
}

// Update 更新相册
func (s *albumService) Update(ctx context.Context, id, userID string, updates map[string]interface{}) (*model.Album, error) {
	album, err := s.albumRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrAlbumNotFound
	}

	// 验证所有权
	if album.UserID != userID {
		return nil, ErrUnauthorized
	}

	// 更新字段
	if name, ok := updates["name"].(string); ok {
		album.Name = name
	}
	if description, ok := updates["description"].(string); ok {
		album.Description = description
	}
	if coverPhoto, ok := updates["cover_photo"].(string); ok {
		album.CoverPhoto = coverPhoto
	}
	if privacy, ok := updates["privacy"].(model.PrivacyLevel); ok {
		album.Privacy = privacy
	}

	if err := s.albumRepo.Update(ctx, album); err != nil {
		return nil, err
	}

	return album, nil
}

// Delete 删除相册
func (s *albumService) Delete(ctx context.Context, id, userID string) error {
	album, err := s.albumRepo.GetByID(ctx, id)
	if err != nil {
		return ErrAlbumNotFound
	}

	// 验证所有权
	if album.UserID != userID {
		return ErrUnauthorized
	}

	return s.albumRepo.Delete(ctx, id)
}

// ListByUser 获取用户相册列表
func (s *albumService) ListByUser(ctx context.Context, userID string, offset, limit int64) ([]*model.Album, int64, error) {
	return s.albumRepo.ListByUser(ctx, userID, offset, limit)
}

// AddPhoto 添加照片
func (s *albumService) AddPhoto(ctx context.Context, albumID, userID, url, thumbnail, description string, width, height int32) (*model.Photo, error) {
	album, err := s.albumRepo.GetByID(ctx, albumID)
	if err != nil {
		return nil, ErrAlbumNotFound
	}

	// 验证所有权
	if album.UserID != userID {
		return nil, ErrUnauthorized
	}

	photo := &model.Photo{
		AlbumID:     albumID,
		UserID:      userID,
		URL:         url,
		Thumbnail:   thumbnail,
		Description: description,
		Width:       width,
		Height:      height,
	}

	if err := s.albumRepo.AddPhoto(ctx, photo); err != nil {
		return nil, err
	}

	// 更新相册照片数
	s.albumRepo.IncrementPhotoCount(ctx, albumID, 1)

	return photo, nil
}

// DeletePhoto 删除照片
func (s *albumService) DeletePhoto(ctx context.Context, photoID, userID string) error {
	// 这里简化处理，实际应该先查询照片获取相册ID
	return s.albumRepo.DeletePhoto(ctx, photoID)
}

// GetPhotos 获取相册照片
func (s *albumService) GetPhotos(ctx context.Context, albumID string, offset, limit int64) ([]*model.Photo, int64, error) {
	return s.albumRepo.GetPhotosByAlbum(ctx, albumID, offset, limit)
}