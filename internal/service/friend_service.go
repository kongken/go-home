package service

import (
	"context"
	"errors"

	"github.com/kongken/go-home/internal/model"
	"github.com/kongken/go-home/internal/repository"
)

var (
	ErrFriendshipNotFound = errors.New("friendship not found")
	ErrRequestNotFound    = errors.New("friend request not found")
	ErrAlreadyFriends     = errors.New("already friends")
	ErrRequestPending     = errors.New("friend request already pending")
)

// FriendService 好友服务接口
type FriendService interface {
	// 好友请求
	SendRequest(ctx context.Context, fromUserID, toUserID, message string) (*model.FriendRequest, error)
	HandleRequest(ctx context.Context, requestID string, accept bool, groupName string) error
	ListReceivedRequests(ctx context.Context, userID string, offset, limit int64) ([]*model.FriendRequest, int64, error)
	ListSentRequests(ctx context.Context, userID string, offset, limit int64) ([]*model.FriendRequest, int64, error)

	// 好友关系
	GetFriendship(ctx context.Context, userID, friendID string) (*model.Friendship, error)
	ListFriends(ctx context.Context, userID string, offset, limit int64) ([]*model.Friendship, int64, error)
	DeleteFriend(ctx context.Context, userID, friendID string) error
	UpdateFriendGroup(ctx context.Context, userID, friendID, groupName string) error
}

// friendService 好友服务实现
type friendService struct {
	friendRepo repository.FriendRepository
}

// NewFriendService 创建好友服务
func NewFriendService(friendRepo repository.FriendRepository) FriendService {
	return &friendService{friendRepo: friendRepo}
}

// SendRequest 发送好友请求
func (s *friendService) SendRequest(ctx context.Context, fromUserID, toUserID, message string) (*model.FriendRequest, error) {
	// 检查是否已经是好友
	_, err := s.friendRepo.GetFriendship(ctx, fromUserID, toUserID)
	if err == nil {
		return nil, ErrAlreadyFriends
	}

	// 检查是否已有待处理的请求
	_, err = s.friendRepo.GetRequestByUsers(ctx, fromUserID, toUserID)
	if err == nil {
		return nil, ErrRequestPending
	}

	request := &model.FriendRequest{
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Message:    message,
		Status:     model.RequestPending,
	}

	if err := s.friendRepo.CreateRequest(ctx, request); err != nil {
		return nil, err
	}

	return request, nil
}

// HandleRequest 处理好友请求
func (s *friendService) HandleRequest(ctx context.Context, requestID string, accept bool, groupName string) error {
	request, err := s.friendRepo.GetRequest(ctx, requestID)
	if err != nil {
		return ErrRequestNotFound
	}

	if !accept {
		request.Status = model.RequestRejected
		return s.friendRepo.UpdateRequest(ctx, request)
	}

	// 接受请求
	request.Status = model.RequestAccepted
	if err := s.friendRepo.UpdateRequest(ctx, request); err != nil {
		return err
	}

	// 创建双向好友关系
	friendship1 := &model.Friendship{
		UserID:    request.FromUserID,
		FriendID:  request.ToUserID,
		Status:    model.FriendshipAccepted,
		GroupName: groupName,
	}
	if err := s.friendRepo.CreateFriendship(ctx, friendship1); err != nil {
		return err
	}

	friendship2 := &model.Friendship{
		UserID:   request.ToUserID,
		FriendID: request.FromUserID,
		Status:   model.FriendshipAccepted,
	}
	if err := s.friendRepo.CreateFriendship(ctx, friendship2); err != nil {
		return err
	}

	return nil
}

// ListReceivedRequests 获取收到的好友请求
func (s *friendService) ListReceivedRequests(ctx context.Context, userID string, offset, limit int64) ([]*model.FriendRequest, int64, error) {
	return s.friendRepo.ListReceivedRequests(ctx, userID, offset, limit)
}

// ListSentRequests 获取发送的好友请求
func (s *friendService) ListSentRequests(ctx context.Context, userID string, offset, limit int64) ([]*model.FriendRequest, int64, error) {
	return s.friendRepo.ListSentRequests(ctx, userID, offset, limit)
}

// GetFriendship 获取好友关系
func (s *friendService) GetFriendship(ctx context.Context, userID, friendID string) (*model.Friendship, error) {
	return s.friendRepo.GetFriendship(ctx, userID, friendID)
}

// ListFriends 获取好友列表
func (s *friendService) ListFriends(ctx context.Context, userID string, offset, limit int64) ([]*model.Friendship, int64, error) {
	return s.friendRepo.ListFriendships(ctx, userID, offset, limit)
}

// DeleteFriend 删除好友
func (s *friendService) DeleteFriend(ctx context.Context, userID, friendID string) error {
	return s.friendRepo.DeleteFriendship(ctx, userID, friendID)
}

// UpdateFriendGroup 更新好友分组
func (s *friendService) UpdateFriendGroup(ctx context.Context, userID, friendID, groupName string) error {
	friendship, err := s.friendRepo.GetFriendship(ctx, userID, friendID)
	if err != nil {
		return ErrFriendshipNotFound
	}

	friendship.GroupName = groupName
	return s.friendRepo.UpdateFriendship(ctx, friendship)
}
