package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// MemberRole 成员角色
type MemberRole int32

const (
	MemberRoleMember MemberRole = 0
	MemberRoleAdmin  MemberRole = 1
	MemberRoleOwner  MemberRole = 2
)

// GroupType 群组类型
type GroupType int32

const (
	GroupTypePublic  GroupType = 0
	GroupTypePrivate GroupType = 1
)

// JoinMode 加入模式
type JoinMode int32

const (
	JoinModeFree     JoinMode = 0
	JoinModeApproval JoinMode = 1
	JoinModeInvite   JoinMode = 2
)

// Group 群组模型
type Group struct {
	ID          string     `bson:"_id,omitempty" json:"id"`
	Name        string     `bson:"name" json:"name"`
	Description string     `bson:"description" json:"description"`
	Avatar      string     `bson:"avatar" json:"avatar"`
	Category    string     `bson:"category" json:"category"`
	OwnerID     string     `bson:"owner_id" json:"owner_id"`
	Type        GroupType  `bson:"type" json:"type"`
	JoinMode    JoinMode   `bson:"join_mode" json:"join_mode"`
	MemberLimit int32      `bson:"member_limit" json:"member_limit"`
	CreatedAt   time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `bson:"updated_at" json:"updated_at"`

	// 统计
	MembersCount int32 `bson:"members_count" json:"members_count"`
	PostsCount   int32 `bson:"posts_count" json:"posts_count"`
	TopicsCount  int32 `bson:"topics_count" json:"topics_count"`
}

// BeforeInsert 插入前钩子
func (g *Group) BeforeInsert() {
	if g.ID == "" {
		g.ID = bson.NewObjectID().Hex()
	}
	now := time.Now()
	if g.CreatedAt.IsZero() {
		g.CreatedAt = now
	}
	g.UpdatedAt = now
}

// BeforeUpdate 更新前钩子
func (g *Group) BeforeUpdate() {
	g.UpdatedAt = time.Now()
}

// CollectionName 集合名
func (Group) CollectionName() string {
	return "groups"
}

// GroupMember 群组成员模型
type GroupMember struct {
	ID        string     `bson:"_id,omitempty" json:"id"`
	GroupID   string     `bson:"group_id" json:"group_id"`
	UserID    string     `bson:"user_id" json:"user_id"`
	Role      MemberRole `bson:"role" json:"role"`
	JoinedAt  time.Time  `bson:"joined_at" json:"joined_at"`
}

// BeforeInsert 插入前钩子
func (gm *GroupMember) BeforeInsert() {
	if gm.ID == "" {
		gm.ID = bson.NewObjectID().Hex()
	}
	if gm.JoinedAt.IsZero() {
		gm.JoinedAt = time.Now()
	}
}

// CollectionName 集合名
func (GroupMember) CollectionName() string {
	return "group_members"
}
