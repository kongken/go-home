package service

import (
	"context"
	"errors"

	"github.com/kongken/go-home/internal/model"
	"github.com/kongken/go-home/internal/repository"
)

var (
	ErrGroupNotFound      = errors.New("group not found")
	ErrGroupMemberExists  = errors.New("already a member")
	ErrGroupMemberNotFound = errors.New("not a member")
	ErrGroupFull          = errors.New("group is full")
	ErrNotGroupOwner      = errors.New("not group owner")
	ErrNotGroupAdmin      = errors.New("not group admin")
)

// GroupService 群组服务接口
type GroupService interface {
	Create(ctx context.Context, ownerID, name, description, avatar, category string, groupType model.GroupType, joinMode model.JoinMode, memberLimit int32) (*model.Group, error)
	Get(ctx context.Context, id string) (*model.Group, error)
	Update(ctx context.Context, id, userID string, updates map[string]interface{}) (*model.Group, error)
	Delete(ctx context.Context, id, userID string) error
	List(ctx context.Context, category string, offset, limit int64) ([]*model.Group, int64, error)
	Search(ctx context.Context, keyword string, offset, limit int64) ([]*model.Group, int64, error)

	// 成员管理
	Join(ctx context.Context, groupID, userID string) error
	Leave(ctx context.Context, groupID, userID string) error
	KickMember(ctx context.Context, groupID, userID, targetUserID string) error
	UpdateMemberRole(ctx context.Context, groupID, userID, targetUserID string, role model.MemberRole) error
	ListMembers(ctx context.Context, groupID string, offset, limit int64) ([]*model.GroupMember, int64, error)
	GetMember(ctx context.Context, groupID, userID string) (*model.GroupMember, error)
}

// groupService 群组服务实现
type groupService struct {
	groupRepo repository.GroupRepository
}

// NewGroupService 创建群组服务
func NewGroupService(groupRepo repository.GroupRepository) GroupService {
	return &groupService{groupRepo: groupRepo}
}

// Create 创建群组
func (s *groupService) Create(ctx context.Context, ownerID, name, description, avatar, category string, groupType model.GroupType, joinMode model.JoinMode, memberLimit int32) (*model.Group, error) {
	group := &model.Group{
		Name:        name,
		Description: description,
		Avatar:      avatar,
		Category:    category,
		OwnerID:     ownerID,
		Type:        groupType,
		JoinMode:    joinMode,
		MemberLimit: memberLimit,
		MembersCount: 1, // 群主
	}

	if err := s.groupRepo.Create(ctx, group); err != nil {
		return nil, err
	}

	// 添加群主为成员
	member := &model.GroupMember{
		GroupID: group.ID,
		UserID:  ownerID,
		Role:    model.MemberRoleOwner,
	}
	if err := s.groupRepo.AddMember(ctx, member); err != nil {
		return nil, err
	}

	return group, nil
}

// Get 获取群组
func (s *groupService) Get(ctx context.Context, id string) (*model.Group, error) {
	return s.groupRepo.GetByID(ctx, id)
}

// Update 更新群组
func (s *groupService) Update(ctx context.Context, id, userID string, updates map[string]interface{}) (*model.Group, error) {
	group, err := s.groupRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrGroupNotFound
	}

	// 检查权限（只有群主或管理员可以更新）
	member, err := s.groupRepo.GetMember(ctx, id, userID)
	if err != nil || (member.Role != model.MemberRoleOwner && member.Role != model.MemberRoleAdmin) {
		return nil, ErrNotGroupAdmin
	}

	// 更新字段
	if name, ok := updates["name"].(string); ok {
		group.Name = name
	}
	if description, ok := updates["description"].(string); ok {
		group.Description = description
	}
	if avatar, ok := updates["avatar"].(string); ok {
		group.Avatar = avatar
	}
	if category, ok := updates["category"].(string); ok {
		group.Category = category
	}
	if joinMode, ok := updates["join_mode"].(model.JoinMode); ok {
		group.JoinMode = joinMode
	}
	if memberLimit, ok := updates["member_limit"].(int32); ok {
		group.MemberLimit = memberLimit
	}

	if err := s.groupRepo.Update(ctx, group); err != nil {
		return nil, err
	}

	return group, nil
}

// Delete 删除群组
func (s *groupService) Delete(ctx context.Context, id, userID string) error {
	group, err := s.groupRepo.GetByID(ctx, id)
	if err != nil {
		return ErrGroupNotFound
	}

	// 只有群主可以删除
	if group.OwnerID != userID {
		return ErrNotGroupOwner
	}

	return s.groupRepo.Delete(ctx, id)
}

// List 获取群组列表
func (s *groupService) List(ctx context.Context, category string, offset, limit int64) ([]*model.Group, int64, error) {
	return s.groupRepo.List(ctx, category, offset, limit)
}

// Search 搜索群组
func (s *groupService) Search(ctx context.Context, keyword string, offset, limit int64) ([]*model.Group, int64, error) {
	return s.groupRepo.Search(ctx, keyword, offset, limit)
}

// Join 加入群组
func (s *groupService) Join(ctx context.Context, groupID, userID string) error {
	group, err := s.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return ErrGroupNotFound
	}

	// 检查是否已满
	if group.MemberLimit > 0 && group.MembersCount >= group.MemberLimit {
		return ErrGroupFull
	}

	// 检查是否已是成员
	_, err = s.groupRepo.GetMember(ctx, groupID, userID)
	if err == nil {
		return ErrGroupMemberExists
	}

	member := &model.GroupMember{
		GroupID: groupID,
		UserID:  userID,
		Role:    model.MemberRoleMember,
	}

	return s.groupRepo.AddMember(ctx, member)
}

// Leave 离开群组
func (s *groupService) Leave(ctx context.Context, groupID, userID string) error {
	// 检查是否是群主
	group, err := s.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return ErrGroupNotFound
	}
	if group.OwnerID == userID {
		return errors.New("owner cannot leave group")
	}

	return s.groupRepo.RemoveMember(ctx, groupID, userID)
}

// KickMember 踢出成员
func (s *groupService) KickMember(ctx context.Context, groupID, userID, targetUserID string) error {
	// 检查操作者权限
	member, err := s.groupRepo.GetMember(ctx, groupID, userID)
	if err != nil {
		return ErrGroupMemberNotFound
	}
	if member.Role != model.MemberRoleOwner && member.Role != model.MemberRoleAdmin {
		return ErrNotGroupAdmin
	}

	// 不能踢出群主
	group, err := s.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return ErrGroupNotFound
	}
	if group.OwnerID == targetUserID {
		return errors.New("cannot kick owner")
	}

	return s.groupRepo.RemoveMember(ctx, groupID, targetUserID)
}

// UpdateMemberRole 更新成员角色
func (s *groupService) UpdateMemberRole(ctx context.Context, groupID, userID, targetUserID string, role model.MemberRole) error {
	// 只有群主可以设置管理员
	group, err := s.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return ErrGroupNotFound
	}
	if group.OwnerID != userID {
		return ErrNotGroupOwner
	}

	return s.groupRepo.UpdateMemberRole(ctx, groupID, targetUserID, role)
}

// ListMembers 获取成员列表
func (s *groupService) ListMembers(ctx context.Context, groupID string, offset, limit int64) ([]*model.GroupMember, int64, error) {
	return s.groupRepo.ListMembers(ctx, groupID, offset, limit)
}

// GetMember 获取成员信息
func (s *groupService) GetMember(ctx context.Context, groupID, userID string) (*model.GroupMember, error) {
	return s.groupRepo.GetMember(ctx, groupID, userID)
}
