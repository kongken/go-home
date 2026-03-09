package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// BlogStatus 博客状态
type BlogStatus int32

const (
	BlogStatusDraft     BlogStatus = 0
	BlogStatusPublished BlogStatus = 1
	BlogStatusDeleted   BlogStatus = 2
)

// PrivacyLevel 隐私级别
type PrivacyLevel int32

const (
	PrivacyPublic   PrivacyLevel = 0
	PrivacyFriends  PrivacyLevel = 1
	PrivacyPrivate  PrivacyLevel = 2
)

// Blog 博客模型
type Blog struct {
	ID          string       `bson:"_id,omitempty" json:"id"`
	UserID      string       `bson:"user_id" json:"user_id"`
	Title       string       `bson:"title" json:"title"`
	Content     string       `bson:"content" json:"content"`
	Summary     string       `bson:"summary" json:"summary"`
	CoverImage  string       `bson:"cover_image" json:"cover_image"`
	Tags        string       `bson:"tags" json:"tags"`
	Category    string       `bson:"category" json:"category"`
	Privacy     PrivacyLevel `bson:"privacy" json:"privacy"`
	Status      BlogStatus   `bson:"status" json:"status"`
	CreatedAt   time.Time    `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time    `bson:"updated_at" json:"updated_at"`

	// 统计信息
	ViewsCount     int32 `bson:"views_count" json:"views_count"`
	LikesCount     int32 `bson:"likes_count" json:"likes_count"`
	CommentsCount  int32 `bson:"comments_count" json:"comments_count"`
	FavoritesCount int32 `bson:"favorites_count" json:"favorites_count"`
}

// BeforeInsert 插入前钩子
func (b *Blog) BeforeInsert(tx interface{}) {
	if b.ID == "" {
		b.ID = bson.NewObjectID().Hex()
	}
	now := time.Now()
	if b.CreatedAt.IsZero() {
		b.CreatedAt = now
	}
	b.UpdatedAt = now
}

// BeforeUpdate 更新前钩子
func (b *Blog) BeforeUpdate() {
	b.UpdatedAt = time.Now()
}

// CollectionName 集合名
func (Blog) CollectionName() string {
	return "blogs"
}
