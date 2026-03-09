package repository

import (
	"context"

	"github.com/kongken/go-home/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// CommentRepository 评论仓库接口
type CommentRepository interface {
	Create(ctx context.Context, comment *model.Comment) error
	GetByID(ctx context.Context, id string) (*model.Comment, error)
	Delete(ctx context.Context, id string) error
	ListByTarget(ctx context.Context, targetID, targetType string, offset, limit int64) ([]*model.Comment, int64, error)
	ListReplies(ctx context.Context, parentID string, offset, limit int64) ([]*model.Comment, int64, error)
	IncrementLikes(ctx context.Context, id string, delta int32) error
	IncrementReplies(ctx context.Context, id string, delta int32) error
}

// commentRepository 评论仓库实现
type commentRepository struct {
	collection *mongo.Collection
}

// NewCommentRepository 创建评论仓库
func NewCommentRepository(db *mongo.Database) CommentRepository {
	return &commentRepository{
		collection: db.Collection(model.Comment{}.CollectionName()),
	}
}

// Create 创建评论
func (r *commentRepository) Create(ctx context.Context, comment *model.Comment) error {
	comment.BeforeInsert()
	_, err := r.collection.InsertOne(ctx, comment)
	return err
}

// GetByID 根据ID获取评论
func (r *commentRepository) GetByID(ctx context.Context, id string) (*model.Comment, error) {
	var comment model.Comment
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&comment)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// Delete 删除评论
func (r *commentRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// ListByTarget 获取目标评论列表
func (r *commentRepository) ListByTarget(ctx context.Context, targetID, targetType string, offset, limit int64) ([]*model.Comment, int64, error) {
	filter := bson.M{
		"target_id":   targetID,
		"target_type": targetType,
		"parent_id":   "",
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

	var comments []*model.Comment
	if err := cursor.All(ctx, &comments); err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

// ListReplies 获取回复列表
func (r *commentRepository) ListReplies(ctx context.Context, parentID string, offset, limit int64) ([]*model.Comment, int64, error) {
	filter := bson.M{"parent_id": parentID}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSkip(offset).
		SetLimit(limit).
		SetSort(bson.D{{Key: "created_at", Value: 1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var comments []*model.Comment
	if err := cursor.All(ctx, &comments); err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

// IncrementLikes 增加点赞数
func (r *commentRepository) IncrementLikes(ctx context.Context, id string, delta int32) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$inc": bson.M{"likes_count": delta}},
	)
	return err
}

// IncrementReplies 增加回复数
func (r *commentRepository) IncrementReplies(ctx context.Context, id string, delta int32) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$inc": bson.M{"replies_count": delta}},
	)
	return err
}