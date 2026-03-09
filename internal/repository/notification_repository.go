package repository

import (
	"context"

	"github.com/kongken/go-home/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// NotificationRepository 通知仓库接口
type NotificationRepository interface {
	Create(ctx context.Context, notification *model.Notification) error
	GetByID(ctx context.Context, id string) (*model.Notification, error)
	List(ctx context.Context, userID string, unreadOnly bool, offset, limit int64) ([]*model.Notification, int64, error)
	MarkAsRead(ctx context.Context, id string) error
	MarkAllAsRead(ctx context.Context, userID string) error
	GetUnreadCount(ctx context.Context, userID string) (int64, error)
	Delete(ctx context.Context, id string) error
}

// notificationRepository 通知仓库实现
type notificationRepository struct {
	collection *mongo.Collection
}

// NewNotificationRepository 创建通知仓库
func NewNotificationRepository(db *mongo.Database) NotificationRepository {
	return &notificationRepository{
		collection: db.Collection(model.Notification{}.CollectionName()),
	}
}

// Create 创建通知
func (r *notificationRepository) Create(ctx context.Context, notification *model.Notification) error {
	notification.BeforeInsert()
	_, err := r.collection.InsertOne(ctx, notification)
	return err
}

// GetByID 根据ID获取通知
func (r *notificationRepository) GetByID(ctx context.Context, id string) (*model.Notification, error) {
	var notification model.Notification
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&notification)
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

// List 获取通知列表
func (r *notificationRepository) List(ctx context.Context, userID string, unreadOnly bool, offset, limit int64) ([]*model.Notification, int64, error) {
	filter := bson.M{"user_id": userID}
	if unreadOnly {
		filter["is_read"] = false
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSkip(offset).
		SetLimit(limit).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var notifications []*model.Notification
	if err := cursor.All(ctx, &notifications); err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

// MarkAsRead 标记通知已读
func (r *notificationRepository) MarkAsRead(ctx context.Context, id string) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"is_read": true}},
	)
	return err
}

// MarkAllAsRead 标记所有通知已读
func (r *notificationRepository) MarkAllAsRead(ctx context.Context, userID string) error {
	_, err := r.collection.UpdateMany(
		ctx,
		bson.M{"user_id": userID, "is_read": false},
		bson.M{"$set": bson.M{"is_read": true}},
	)
	return err
}

// GetUnreadCount 获取未读通知数
func (r *notificationRepository) GetUnreadCount(ctx context.Context, userID string) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{
		"user_id": userID,
		"is_read": false,
	})
}

// Delete 删除通知
func (r *notificationRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
