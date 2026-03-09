package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// Message 私信模型
type Message struct {
	ID          string             `bson:"_id,omitempty" json:"id"`
	FromUserID  string             `bson:"from_user_id" json:"from_user_id"`
	ToUserID    string             `bson:"to_user_id" json:"to_user_id"`
	Content     string             `bson:"content" json:"content"`
	Attachments []MediaAttachment  `bson:"attachments,omitempty" json:"attachments,omitempty"`
	IsRead      bool               `bson:"is_read" json:"is_read"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
}

// BeforeInsert 插入前钩子
func (m *Message) BeforeInsert() {
	if m.ID == "" {
		m.ID = bson.NewObjectID().Hex()
	}
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}
}

// CollectionName 集合名
func (Message) CollectionName() string {
	return "messages"
}

// Conversation 会话模型
type Conversation struct {
	ID           string    `bson:"_id,omitempty" json:"id"`
	UserID       string    `bson:"user_id" json:"user_id"`             // 当前用户
	OtherUserID  string    `bson:"other_user_id" json:"other_user_id"` // 对方用户
	LastMessageID string   `bson:"last_message_id" json:"last_message_id"`
	UnreadCount  int32     `bson:"unread_count" json:"unread_count"`
	UpdatedAt    time.Time `bson:"updated_at" json:"updated_at"`
}

// BeforeInsert 插入前钩子
func (c *Conversation) BeforeInsert() {
	if c.ID == "" {
		c.ID = bson.NewObjectID().Hex()
	}
	c.UpdatedAt = time.Now()
}

// BeforeUpdate 更新前钩子
func (c *Conversation) BeforeUpdate() {
	c.UpdatedAt = time.Now()
}

// CollectionName 集合名
func (Conversation) CollectionName() string {
	return "conversations"
}
