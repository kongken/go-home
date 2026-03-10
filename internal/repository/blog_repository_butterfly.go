package repository

import (
	"context"

	"github.com/kongken/go-home/internal/model"
	"butterfly.orx.me/core/store/mongo"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// BlogRepositoryButterfly butterfly 版本的博客仓库
type BlogRepositoryButterfly struct {
	collection *mongodriver.Collection
}

// NewBlogRepositoryButterfly 创建 butterfly 博客仓库
func NewBlogRepositoryButterfly() BlogRepository {
	client := mongo.GetClient("primary")
	if client == nil {
		panic("mongo client 'primary' not found")
	}
	
	return &BlogRepositoryButterfly{
		collection: client.Database("gohome").Collection("blogs"),
	}
}

// Create 创建博客
func (r *BlogRepositoryButterfly) Create(ctx context.Context, blog *model.Blog) error {
	blog.BeforeInsert(nil)
	_, err := r.collection.InsertOne(ctx, blog)
	return err
}

// GetByID 根据ID获取博客
func (r *BlogRepositoryButterfly) GetByID(ctx context.Context, id string) (*model.Blog, error) {
	var blog model.Blog
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&blog)
	if err != nil {
		return nil, err
	}
	return &blog, nil
}

// Update 更新博客
func (r *BlogRepositoryButterfly) Update(ctx context.Context, blog *model.Blog) error {
	blog.BeforeUpdate()
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": blog.ID}, blog)
	return err
}

// Delete 删除博客
func (r *BlogRepositoryButterfly) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// List 获取博客列表
func (r *BlogRepositoryButterfly) List(ctx context.Context, userID string, category string, offset, limit int) ([]*model.Blog, int64, error) {
	filter := bson.M{"status": model.BlogStatusPublished}

	if userID != "" {
		filter["user_id"] = userID
	}
	if category != "" {
		filter["category"] = category
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSkip(int64(offset)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var blogs []*model.Blog
	if err := cursor.All(ctx, &blogs); err != nil {
		return nil, 0, err
	}

	return blogs, total, nil
}

// ListByUser 获取指定用户的博客列表
func (r *BlogRepositoryButterfly) ListByUser(ctx context.Context, userID string, offset, limit int) ([]*model.Blog, int64, error) {
	return r.List(ctx, userID, "", offset, limit)
}