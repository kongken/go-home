package repository

import (
	"context"

	"github.com/kongken/go-home/internal/model"
	"butterfly.orx.me/core/store/mongo"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// CommentRepositoryButterfly butterfly 版本的评论仓库
type CommentRepositoryButterfly struct {
	collection *mongodriver.Collection
}

// NewCommentRepositoryButterfly 创建 butterfly 评论仓库
func NewCommentRepositoryButterfly() CommentRepository {
	client := mongo.GetClient("primary")
	if client == nil {
		panic("mongo client 'primary' not found")
	}
	return &CommentRepositoryButterfly{
		collection: client.Database("gohome").Collection("comments"),
	}
}

func (r *CommentRepositoryButterfly) Create(ctx context.Context, comment *model.Comment) error {
	comment.BeforeInsert()
	_, err := r.collection.InsertOne(ctx, comment)
	return err
}

func (r *CommentRepositoryButterfly) GetByID(ctx context.Context, id string) (*model.Comment, error) {
	var comment model.Comment
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&comment)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *CommentRepositoryButterfly) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *CommentRepositoryButterfly) ListByTarget(ctx context.Context, targetID, targetType string, offset, limit int64) ([]*model.Comment, int64, error) {
	filter := bson.M{"target_id": targetID, "target_type": targetType, "parent_id": ""}
	total, _ := r.collection.CountDocuments(ctx, filter)
	opts := options.Find().SetSkip(offset).SetLimit(limit).SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, _ := r.collection.Find(ctx, filter, opts)
	defer cursor.Close(ctx)
	var comments []*model.Comment
	cursor.All(ctx, &comments)
	return comments, total, nil
}

func (r *CommentRepositoryButterfly) ListReplies(ctx context.Context, parentID string, offset, limit int64) ([]*model.Comment, int64, error) {
	filter := bson.M{"parent_id": parentID}
	total, _ := r.collection.CountDocuments(ctx, filter)
	opts := options.Find().SetSkip(offset).SetLimit(limit)
	cursor, _ := r.collection.Find(ctx, filter, opts)
	defer cursor.Close(ctx)
	var comments []*model.Comment
	cursor.All(ctx, &comments)
	return comments, total, nil
}

func (r *CommentRepositoryButterfly) IncrementLikes(ctx context.Context, id string, delta int32) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$inc": bson.M{"likes_count": delta}})
	return err
}

func (r *CommentRepositoryButterfly) IncrementReplies(ctx context.Context, id string, delta int32) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$inc": bson.M{"replies_count": delta}})
	return err
}