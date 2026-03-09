package service

import (
	"context"

	"github.com/kongken/go-home/internal/model"
	"github.com/kongken/go-home/internal/repository"
)

// SettingsService 设置服务接口
type SettingsService interface {
	Get(ctx context.Context, userID string) (*model.Settings, error)
	UpdatePrivacy(ctx context.Context, userID string, privacy model.PrivacySettings) error
	UpdateNotification(ctx context.Context, userID string, notification model.NotificationSettings) error
	AddToBlacklist(ctx context.Context, userID, blockedUserID string) error
	RemoveFromBlacklist(ctx context.Context, userID, blockedUserID string) error
}

// settingsService 设置服务实现
type settingsService struct {
	settingsRepo repository.SettingsRepository
}

// NewSettingsService 创建设置服务
func NewSettingsService(settingsRepo repository.SettingsRepository) SettingsService {
	return &settingsService{settingsRepo: settingsRepo}
}

// Get 获取用户设置
func (s *settingsService) Get(ctx context.Context, userID string) (*model.Settings, error) {
	settings, err := s.settingsRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 如果是默认设置，保存到数据库
	if settings.ID == "" {
		settings.UserID = userID
		s.settingsRepo.Create(ctx, settings)
	}

	return settings, nil
}

// UpdatePrivacy 更新隐私设置
func (s *settingsService) UpdatePrivacy(ctx context.Context, userID string, privacy model.PrivacySettings) error {
	settings, err := s.settingsRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	settings.Privacy = privacy

	if settings.ID == "" {
		return s.settingsRepo.Create(ctx, settings)
	}
	return s.settingsRepo.Update(ctx, settings)
}

// UpdateNotification 更新通知设置
func (s *settingsService) UpdateNotification(ctx context.Context, userID string, notification model.NotificationSettings) error {
	settings, err := s.settingsRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	settings.Notification = notification

	if settings.ID == "" {
		return s.settingsRepo.Create(ctx, settings)
	}
	return s.settingsRepo.Update(ctx, settings)
}

// AddToBlacklist 添加到黑名单
func (s *settingsService) AddToBlacklist(ctx context.Context, userID, blockedUserID string) error {
	return s.settingsRepo.AddToBlacklist(ctx, userID, blockedUserID)
}

// RemoveFromBlacklist 从黑名单移除
func (s *settingsService) RemoveFromBlacklist(ctx context.Context, userID, blockedUserID string) error {
	return s.settingsRepo.RemoveFromBlacklist(ctx, userID, blockedUserID)
}