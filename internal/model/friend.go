package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// FriendshipStatus 好友关系状态
type FriendshipStatus int32

const (
	FriendshipPending  FriendshipStatus = 0
	FriendshipAccepted FriendshipStatus = 1
	FriendshipBlocked  FriendshipStatus = 2
)

// RequestStatus 请求状态
type RequestStatus int32

const (
	RequestPending   RequestStatus = 0
	RequestAccepted  RequestStatus = 1
	RequestRejected  RequestStatus = 2
	RequestExpired   RequestStatus = 3
)

// Friendship 好友关系
type Friendship struct {
	ID         string           `bson:"_id,omitempty" json:"id"`
	UserID     string           `bson:"user_id" json:"user_id"`
	FriendID   string           `bson:"friend_id" json:"friend_id"`
	Status     FriendshipStatus `bson:"status" json:"status"`
	GroupName  string           `bson:"group_name,omitempty" json:"group_name,omitempty"`
	Remark     string           `bson:"remark,omitempty" json:"remark,omitempty"`
	CreatedAt  time.Time        `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time        `bson:"updated_at" json:"updated_at"`
}

// BeforeInsert 插入前钩子
func (f *Friendship) BeforeInsert() {
	if f.ID == "" {
		f.ID = bson.NewObjectID().Hex()
	}
	now := time.Now()
	if f.CreatedAt.IsZero() {
		f.CreatedAt = now
	}
	f.UpdatedAt = now
}

// CollectionName 集合名
func (Friendship) CollectionName() string {
	return "friendships"
}

// FriendRequest 好友请求
type FriendRequest struct {
	ID          string        `bson:"_id,omitempty" json:"id"`
	FromUserID  string        `bson:"from_user_id" json:"from_user_id"`
	ToUserID    string        `bson:"to_user_id" json:"to_user_id"`
	Message     string        `bson:"message,omitempty" json:"message,omitempty"`
	Status      RequestStatus `bson:"status" json:"status"`
	CreatedAt   time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time     `bson:"updated_at" json:"updated_at"`
}

// BeforeInsert 插入前钩子
func (fr *FriendRequest) BeforeInsert() {
	if fr.ID == "" {
		fr.ID = bson.NewObjectID().Hex()
	}
	now := time.Now()
	if fr.CreatedAt.IsZero() {
		fr.CreatedAt = now
	}
	fr.UpdatedAt = now
}

// CollectionName 集合名
func (FriendRequest) CollectionName() string {
	return "friend_requests"
}
