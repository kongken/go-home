package service

import (
	"context"
	"errors"

	"github.com/kongken/go-home/internal/model"
	"github.com/kongken/go-home/internal/repository"
)

var (
	ErrCommentNotFound = errors.New("comment not found")
)

// CommentService 评论服务接口
type CommentService interface {
	Create(ctx context.Context, userID, targetID, targetType, content, parentID string, attachments []model.MediaAttachment) (*model.Comment, error)
	Get(ctx context.Context, id string) (*model.Comment, error)
	Delete(ctx context.Context, id, userID string) error
	ListByTarget(ctx context.Context, targetID, targetType string, offset, limit int64) ([]*model.Comment, int64, error)
	ListReplies(ctx context.Context, parentID string, offset, limit int64) ([]*model.Comment, int64, error)
	Like(ctx context.Context, id string, delta int32) error
}

// commentService 评论服务实现
type commentService struct {
	commentRepo repository.CommentRepository
}

// NewCommentService 创建评论服务
func NewCommentService(commentRepo repository.CommentRepository) CommentService {
	return &commentService{commentRepo: commentRepo}
}

// Create 创建评论
func (s *commentService) Create(ctx context.Context, userID, targetID, targetType, content, parentID string, attachments []model.MediaAttachment) (*model.Comment, error) {
	comment := &model.Comment{
		UserID:      userID,
		TargetID:    targetID,
		TargetType:  targetType,
		Content:     content,
		ParentID:    parentID,
		Attachments: attachments,
	}

	if err := s.commentRepo.Create(ctx, comment); err != nil {
		return nil, err
	}

	// 如果是回复，增加父评论的回复数
	if parentID != "" {
		s.commentRepo.IncrementReplies(ctx, parentID, 1)
	}

	return comment, nil
}

// Get 获取评论
func (s *commentService) Get(ctx context.Context, id string) (*model.Comment, error) {
	return s.commentRepo.GetByID(ctx, id)
}

// Delete 删除评论
func (s *commentService) Delete(ctx context.Context, id, userID string) error {
	comment, err := s.commentRepo.GetByID(ctx, id)
	if err != nil {
		return ErrCommentNotFound
	}

	// 验证所有权
	if comment.UserID != userID {
		return ErrUnauthorized
	}

	// 如果是回复，减少父评论的回复数
	if comment.ParentID != "" {
		s.commentRepo.IncrementReplies(ctx, comment.ParentID, -1)
	}

	return s.commentRepo.Delete(ctx, id)
}

// ListByTarget 获取目标评论列表
func (s *commentService) ListByTarget(ctx context.Context, targetID, targetType string, offset, limit int64) ([]*model.Comment, int64, error) {
	return s.commentRepo.ListByTarget(ctx, targetID, targetType, offset, limit)
}

// ListReplies 获取回复列表
func (s *commentService) ListReplies(ctx context.Context, parentID string, offset, limit int64) ([]*model.Comment, int64, error) {
	return s.commentRepo.ListReplies(ctx, parentID, offset, limit)
}

// Like 点赞/取消点赞
func (s *commentService) Like(ctx context.Context, id string, delta int32) error {
	return s.commentRepo.IncrementLikes(ctx, id, delta)
}