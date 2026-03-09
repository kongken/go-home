package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// FeedType 动态类型
type FeedType int32

const (
	FeedTypeStatus     FeedType = 0
	FeedTypeBlog       FeedType = 1
	FeedTypeAlbum      FeedType = 2
	FeedTypeShare      FeedType = 3
	FeedTypeCheckin    FeedType = 4
	FeedTypeFriendship FeedType = 5
	FeedTypeJoinGroup  FeedType = 6
	FeedTypeActivity   FeedType = 7
)

// FeedItem 动态项
type FeedItem struct {
	ID          string         `bson:"_id,omitempty" json:"id"`
	UserID      string         `bson:"user_id" json:"user_id"`
	Type        FeedType       `bson:"type" json:"type"`
	Content     string         `bson:"content" json:"content"`
	TargetID    string         `bson:"target_id,omitempty" json:"target_id,omitempty"`
	TargetType  string         `bson:"target_type,omitempty" json:"target_type,omitempty"`
	Attachments []MediaAttachment `bson:"attachments,omitempty" json:"attachments,omitempty"`
	Location    *Location      `bson:"location,omitempty" json:"location,omitempty"`
	Privacy     PrivacyLevel   `bson:"privacy" json:"privacy"`
	CreatedAt   time.Time      `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time      `bson:"updated_at" json:"updated_at"`

	// 统计信息
	LikesCount    int32 `bson:"likes_count" json:"likes_count"`
	CommentsCount int32 `bson:"comments_count" json:"comments_count"`
	SharesCount   int32 `bson:"shares_count" json:"shares_count"`
	ViewsCount    int32 `bson:"views_count" json:"views_count"`
}

// MediaAttachment 媒体附件
type MediaAttachment struct {
	Type      string `bson:"type" json:"type"` // image, video, audio
	URL       string `bson:"url" json:"url"`
	Thumbnail string `bson:"thumbnail,omitempty" json:"thumbnail,omitempty"`
	Width     int32  `bson:"width,omitempty" json:"width,omitempty"`
	Height    int32  `bson:"height,omitempty" json:"height,omitempty"`
}

// Location 位置
type Location struct {
	Latitude  float64 `bson:"latitude" json:"latitude"`
	Longitude float64 `bson:"longitude" json:"longitude"`
	Address   string  `bson:"address,omitempty" json:"address,omitempty"`
	Name      string  `bson:"name,omitempty" json:"name,omitempty"`
}

// BeforeInsert 插入前钩子
func (f *FeedItem) BeforeInsert() {
	if f.ID == "" {
		f.ID = bson.NewObjectID().Hex()
	}
	now := time.Now()
	if f.CreatedAt.IsZero() {
		f.CreatedAt = now
	}
	f.UpdatedAt = now
}

// BeforeUpdate 更新前钩子
func (f *FeedItem) BeforeUpdate() {
	f.UpdatedAt = time.Now()
}

// CollectionName 集合名
func (FeedItem) CollectionName() string {
	return "feeds"
}
