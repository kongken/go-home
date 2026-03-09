package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// Comment 评论模型
type Comment struct {
	ID          string            `bson:"_id,omitempty" json:"id"`
	UserID      string            `bson:"user_id" json:"user_id"`
	TargetID    string            `bson:"target_id" json:"target_id"`
	TargetType  string            `bson:"target_type" json:"target_type"`
	Content     string            `bson:"content" json:"content"`
	ParentID    string            `bson:"parent_id,omitempty" json:"parent_id,omitempty"`
	Attachments []MediaAttachment `bson:"attachments,omitempty" json:"attachments,omitempty"`
	CreatedAt   time.Time         `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time         `bson:"updated_at" json:"updated_at"`

	// 统计
	LikesCount   int32 `bson:"likes_count" json:"likes_count"`
	RepliesCount int32 `bson:"replies_count" json:"replies_count"`
}

// BeforeInsert 插入前钩子
func (c *Comment) BeforeInsert() {
	if c.ID == "" {
		c.ID = bson.NewObjectID().Hex()
	}
	now := time.Now()
	if c.CreatedAt.IsZero() {
		c.CreatedAt = now
	}
	c.UpdatedAt = now
}

// CollectionName 集合名
func (Comment) CollectionName() string {
	return "comments"
}
