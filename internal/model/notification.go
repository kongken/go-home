package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// NotificationType 通知类型
type NotificationType int32

const (
	NotificationTypeSystem        NotificationType = 0
	NotificationTypeFriendRequest NotificationType = 1
	NotificationTypeFriendAccepted NotificationType = 2
	NotificationTypeLike          NotificationType = 3
	NotificationTypeComment       NotificationType = 4
	NotificationTypeMention       NotificationType = 5
	NotificationTypeGroupInvite   NotificationType = 6
	NotificationTypeGroupApproved NotificationType = 7
)

// Notification 通知模型
type Notification struct {
	ID         string           `bson:"_id,omitempty" json:"id"`
	UserID     string           `bson:"user_id" json:"user_id"`       // 接收者
	Type       NotificationType `bson:"type" json:"type"`
	Title      string           `bson:"title" json:"title"`
	Content    string           `bson:"content" json:"content"`
	ActorID    string           `bson:"actor_id" json:"actor_id"`     // 触发者
	TargetID   string           `bson:"target_id" json:"target_id"`   // 关联对象
	TargetType string           `bson:"target_type" json:"target_type"` // 关联对象类型
	IsRead     bool             `bson:"is_read" json:"is_read"`
	CreatedAt  time.Time        `bson:"created_at" json:"created_at"`
}

// BeforeInsert 插入前钩子
func (n *Notification) BeforeInsert() {
	if n.ID == "" {
		n.ID = bson.NewObjectID().Hex()
	}
	if n.CreatedAt.IsZero() {
		n.CreatedAt = time.Now()
	}
}

// CollectionName 集合名
func (Notification) CollectionName() string {
	return "notifications"
}
