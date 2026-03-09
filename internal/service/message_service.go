package service

import (
	"context"
	"errors"

	"github.com/kongken/go-home/internal/model"
	"github.com/kongken/go-home/internal/repository"
)

var (
	ErrMessageNotFound = errors.New("message not found")
)

// MessageService 消息服务接口
type MessageService interface {
	Send(ctx context.Context, fromUserID, toUserID, content string, attachments []model.MediaAttachment) (*model.Message, error)
	Get(ctx context.Context, id string) (*model.Message, error)
	ListByConversation(ctx context.Context, userID, otherUserID string, offset, limit int64) ([]*model.Message, int64, error)
	MarkAsRead(ctx context.Context, userID, otherUserID string) error
	GetUnreadCount(ctx context.Context, userID string) (int64, error)
	ListConversations(ctx context.Context, userID string, offset, limit int64) ([]*model.Conversation, int64, error)
}

// messageService 消息服务实现
type messageService struct {
	messageRepo repository.MessageRepository
}

// NewMessageService 创建消息服务
func NewMessageService(messageRepo repository.MessageRepository) MessageService {
	return &messageService{messageRepo: messageRepo}
}

// Send 发送消息
func (s *messageService) Send(ctx context.Context, fromUserID, toUserID, content string, attachments []model.MediaAttachment) (*model.Message, error) {
	message := &model.Message{
		FromUserID:  fromUserID,
		ToUserID:    toUserID,
		Content:     content,
		Attachments: attachments,
		IsRead:      false,
	}

	if err := s.messageRepo.Create(ctx, message); err != nil {
		return nil, err
	}

	// 更新发送者的会话
	conversation1, err := s.messageRepo.GetConversation(ctx, fromUserID, toUserID)
	if err != nil {
		// 创建新会话
		conversation1 = &model.Conversation{
			UserID:        fromUserID,
			OtherUserID:   toUserID,
			LastMessageID: message.ID,
			UnreadCount:   0,
		}
		s.messageRepo.CreateConversation(ctx, conversation1)
	} else {
		conversation1.LastMessageID = message.ID
		conversation1.UnreadCount = 0
		s.messageRepo.UpdateConversation(ctx, conversation1)
	}

	// 更新接收者的会话
	conversation2, err := s.messageRepo.GetConversation(ctx, toUserID, fromUserID)
	if err != nil {
		// 创建新会话
		conversation2 = &model.Conversation{
			UserID:        toUserID,
			OtherUserID:   fromUserID,
			LastMessageID: message.ID,
			UnreadCount:   1,
		}
		s.messageRepo.CreateConversation(ctx, conversation2)
	} else {
		conversation2.LastMessageID = message.ID
		conversation2.UnreadCount++
		s.messageRepo.UpdateConversation(ctx, conversation2)
	}

	return message, nil
}

// Get 获取消息
func (s *messageService) Get(ctx context.Context, id string) (*model.Message, error) {
	return s.messageRepo.GetByID(ctx, id)
}

// ListByConversation 获取会话消息列表
func (s *messageService) ListByConversation(ctx context.Context, userID, otherUserID string, offset, limit int64) ([]*model.Message, int64, error) {
	return s.messageRepo.ListByConversation(ctx, userID, otherUserID, offset, limit)
}

// MarkAsRead 标记消息已读
func (s *messageService) MarkAsRead(ctx context.Context, userID, otherUserID string) error {
	// 标记消息已读
	if err := s.messageRepo.MarkAsRead(ctx, userID, otherUserID); err != nil {
		return err
	}

	// 更新会话未读数
	conversation, err := s.messageRepo.GetConversation(ctx, userID, otherUserID)
	if err == nil {
		conversation.UnreadCount = 0
		s.messageRepo.UpdateConversation(ctx, conversation)
	}

	return nil
}

// GetUnreadCount 获取未读消息数
func (s *messageService) GetUnreadCount(ctx context.Context, userID string) (int64, error) {
	return s.messageRepo.GetUnreadCount(ctx, userID)
}

// ListConversations 获取会话列表
func (s *messageService) ListConversations(ctx context.Context, userID string, offset, limit int64) ([]*model.Conversation, int64, error) {
	return s.messageRepo.ListConversations(ctx, userID, offset, limit)
}
