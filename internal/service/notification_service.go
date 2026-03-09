package service

import (
	"context"
	"errors"

	"github.com/kongken/go-home/internal/model"
	"github.com/kongken/go-home/internal/repository"
)

var (
	ErrNotificationNotFound = errors.New("notification not found")
)

// NotificationService 通知服务接口
type NotificationService interface {
	Create(ctx context.Context, userID string, notifType model.NotificationType, title, content, actorID, targetID, targetType string) (*model.Notification, error)
	Get(ctx context.Context, id string) (*model.Notification, error)
	List(ctx context.Context, userID string, unreadOnly bool, offset, limit int64) ([]*model.Notification, int64, error)
	MarkAsRead(ctx context.Context, id string) error
	MarkAllAsRead(ctx context.Context, userID string) error
	GetUnreadCount(ctx context.Context, userID string) (int64, error)
	Delete(ctx context.Context, id, userID string) error
}

// notificationService 通知服务实现
type notificationService struct {
	notifRepo repository.NotificationRepository
}

// NewNotificationService 创建通知服务
func NewNotificationService(notifRepo repository.NotificationRepository) NotificationService {
	return &notificationService{notifRepo: notifRepo}
}

// Create 创建通知
func (s *notificationService) Create(ctx context.Context, userID string, notifType model.NotificationType, title, content, actorID, targetID, targetType string) (*model.Notification, error) {
	notification := &model.Notification{
		UserID:     userID,
		Type:       notifType,
		Title:      title,
		Content:    content,
		ActorID:    actorID,
		TargetID:   targetID,
		TargetType: targetType,
		IsRead:     false,
	}

	if err := s.notifRepo.Create(ctx, notification); err != nil {
		return nil, err
	}

	return notification, nil
}

// Get 获取通知
func (s *notificationService) Get(ctx context.Context, id string) (*model.Notification, error) {
	return s.notifRepo.GetByID(ctx, id)
}

// List 获取通知列表
func (s *notificationService) List(ctx context.Context, userID string, unreadOnly bool, offset, limit int64) ([]*model.Notification, int64, error) {
	return s.notifRepo.List(ctx, userID, unreadOnly, offset, limit)
}

// MarkAsRead 标记通知已读
func (s *notificationService) MarkAsRead(ctx context.Context, id string) error {
	return s.notifRepo.MarkAsRead(ctx, id)
}

// MarkAllAsRead 标记所有通知已读
func (s *notificationService) MarkAllAsRead(ctx context.Context, userID string) error {
	return s.notifRepo.MarkAllAsRead(ctx, userID)
}

// GetUnreadCount 获取未读通知数
func (s *notificationService) GetUnreadCount(ctx context.Context, userID string) (int64, error) {
	return s.notifRepo.GetUnreadCount(ctx, userID)
}

// Delete 删除通知
func (s *notificationService) Delete(ctx context.Context, id, userID string) error {
	notification, err := s.notifRepo.GetByID(ctx, id)
	if err != nil {
		return ErrNotificationNotFound
	}

	// 验证所有权
	if notification.UserID != userID {
		return errors.New("unauthorized")
	}

	return s.notifRepo.Delete(ctx, id)
}