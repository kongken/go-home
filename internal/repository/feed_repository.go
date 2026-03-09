package repository

import (
	"context"

	"github.com/kongken/go-home/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// FeedRepository 动态仓库接口
type FeedRepository interface {
	Create(ctx context.Context, feed *model.FeedItem) error
	GetByID(ctx context.Context, id string) (*model.FeedItem, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, userIDs []string, offset, limit int64) ([]*model.FeedItem, int64, error)
	ListByUser(ctx context.Context, userID string, offset, limit int64) ([]*model.FeedItem, int64, error)
	IncrementStats(ctx context.Context, id string, field string, delta int32) error
}

// feedRepository 动态仓库实现
type feedRepository struct {
	collection *mongo.Collection
}

// NewFeedRepository 创建动态仓库
func NewFeedRepository(db *mongo.Database) FeedRepository {
	return &feedRepository{
		collection: db.Collection(model.FeedItem{}.CollectionName()),
	}
}

// Create 创建动态
func (r *feedRepository) Create(ctx context.Context, feed *model.FeedItem) error {
	feed.BeforeInsert()
	_, err := r.collection.InsertOne(ctx, feed)
	return err
}

// GetByID 根据ID获取动态
func (r *feedRepository) GetByID(ctx context.Context, id string) (*model.FeedItem, error) {
	var feed model.FeedItem
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&feed)
	if err != nil {
		return nil, err
	}
	return &feed, nil
}

// Delete 删除动态
func (r *feedRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// List 获取动态列表（多用户）
func (r *feedRepository) List(ctx context.Context, userIDs []string, offset, limit int64) ([]*model.FeedItem, int64, error) {
	filter := bson.M{
		"user_id": bson.M{"$in": userIDs},
		"privacy": bson.M{"$in": []model.PrivacyLevel{model.PrivacyPublic, model.PrivacyFriends}},
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

	var feeds []*model.FeedItem
	if err := cursor.All(ctx, &feeds); err != nil {
		return nil, 0, err
	}

	return feeds, total, nil
}

// ListByUser 获取指定用户的动态列表
func (r *feedRepository) ListByUser(ctx context.Context, userID string, offset, limit int64) ([]*model.FeedItem, int64, error) {
	filter := bson.M{"user_id": userID}

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

	var feeds []*model.FeedItem
	if err := cursor.All(ctx, &feeds); err != nil {
		return nil, 0, err
	}

	return feeds, total, nil
}

// IncrementStats 增加统计字段
func (r *feedRepository) IncrementStats(ctx context.Context, id string, field string, delta int32) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$inc": bson.M{field: delta}},
	)
	return err
}
