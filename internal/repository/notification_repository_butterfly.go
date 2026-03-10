package repository

import (
	"context"

	"github.com/kongken/go-home/internal/model"
	"butterfly.orx.me/core/store/mongo"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// NotificationRepositoryButterfly butterfly 版本的通知仓库
type NotificationRepositoryButterfly struct {
	collection *mongodriver.Collection
}

// NewNotificationRepositoryButterfly 创建 butterfly 通知仓库
func NewNotificationRepositoryButterfly() NotificationRepository {
	client := mongo.GetClient("primary")
	if client == nil {
		panic("mongo client 'primary' not found")
	}
	return &NotificationRepositoryButterfly{
		collection: client.Database("gohome").Collection("notifications"),
	}
}

func (r *NotificationRepositoryButterfly) Create(ctx context.Context, notification *model.Notification) error {
	notification.BeforeInsert()
	_, err := r.collection.InsertOne(ctx, notification)
	return err
}

func (r *NotificationRepositoryButterfly) GetByID(ctx context.Context, id string) (*model.Notification, error) {
	var notification model.Notification
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&notification)
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

func (r *NotificationRepositoryButterfly) List(ctx context.Context, userID string, unreadOnly bool, offset, limit int64) ([]*model.Notification, int64, error) {
	filter := bson.M{"user_id": userID}
	if unreadOnly {
		filter["is_read"] = false
	}
	total, _ := r.collection.CountDocuments(ctx, filter)
	opts := options.Find().SetSkip(offset).SetLimit(limit).SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, _ := r.collection.Find(ctx, filter, opts)
	defer cursor.Close(ctx)
	var notifications []*model.Notification
	cursor.All(ctx, &notifications)
	return notifications, total, nil
}

func (r *NotificationRepositoryButterfly) MarkAsRead(ctx context.Context, id string) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"is_read": true}})
	return err
}

func (r *NotificationRepositoryButterfly) MarkAllAsRead(ctx context.Context, userID string) error {
	_, err := r.collection.UpdateMany(ctx,
		bson.M{"user_id": userID, "is_read": false},
		bson.M{"$set": bson.M{"is_read": true}},
	)
	return err
}

func (r *NotificationRepositoryButterfly) GetUnreadCount(ctx context.Context, userID string) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{"user_id": userID, "is_read": false})
}

func (r *NotificationRepositoryButterfly) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}