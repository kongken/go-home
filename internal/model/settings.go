package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// PrivacySettings 隐私设置
type PrivacySettings struct {
	DefaultBlogPrivacy   PrivacyLevel `bson:"default_blog_privacy" json:"default_blog_privacy"`
	DefaultAlbumPrivacy  PrivacyLevel `bson:"default_album_privacy" json:"default_album_privacy"`
	DefaultSharePrivacy  PrivacyLevel `bson:"default_share_privacy" json:"default_share_privacy"`
	DefaultFeedPrivacy   PrivacyLevel `bson:"default_feed_privacy" json:"default_feed_privacy"`
}

// NotificationSettings 通知设置
type NotificationSettings struct {
	NotifyFriendRequest bool `bson:"notify_friend_request" json:"notify_friend_request"`
	NotifyComment       bool `bson:"notify_comment" json:"notify_comment"`
	NotifyLike          bool `bson:"notify_like" json:"notify_like"`
	NotifyMention       bool `bson:"notify_mention" json:"notify_mention"`
	NotifyGroup         bool `bson:"notify_group" json:"notify_group"`
	NotifySystem        bool `bson:"notify_system" json:"notify_system"`
}

// Settings 用户设置
type Settings struct {
	ID                   string                 `bson:"_id,omitempty" json:"id"`
	UserID               string                 `bson:"user_id" json:"user_id"`
	Privacy              PrivacySettings        `bson:"privacy" json:"privacy"`
	Notification         NotificationSettings   `bson:"notification" json:"notification"`
	Blacklist            []string               `bson:"blacklist" json:"blacklist"`
	CreatedAt            time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt            time.Time              `bson:"updated_at" json:"updated_at"`
}

// BeforeInsert 插入前钩子
func (s *Settings) BeforeInsert() {
	if s.ID == "" {
		s.ID = bson.NewObjectID().Hex()
	}
	now := time.Now()
	if s.CreatedAt.IsZero() {
		s.CreatedAt = now
	}
	s.UpdatedAt = now
}

// BeforeUpdate 更新前钩子
func (s *Settings) BeforeUpdate() {
	s.UpdatedAt = time.Now()
}

// CollectionName 集合名
func (Settings) CollectionName() string {
	return "settings"
}
