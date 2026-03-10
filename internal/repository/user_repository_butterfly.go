package repository

import (
	"context"

	"github.com/kongken/go-home/internal/model"
	"butterfly.orx.me/core/store/mongo"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// UserRepositoryButterfly butterfly 版本的用户仓库
type UserRepositoryButterfly struct {
	collection *mongodriver.Collection
}

// NewUserRepositoryButterfly 创建 butterfly 用户仓库
func NewUserRepositoryButterfly() UserRepository {
	// 使用 butterfly mongo store 获取 client
	client := mongo.GetClient("primary")
	if client == nil {
		panic("mongo client 'primary' not found, please check config.yaml")
	}
	
	return &UserRepositoryButterfly{
		collection: client.Database("gohome").Collection("users"),
	}
}

// Create 创建用户
func (r *UserRepositoryButterfly) Create(ctx context.Context, user *model.User) error {
	user.BeforeInsert(nil)
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

// GetByID 根据ID获取用户
func (r *UserRepositoryButterfly) GetByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (r *UserRepositoryButterfly) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (r *UserRepositoryButterfly) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户
func (r *UserRepositoryButterfly) Update(ctx context.Context, user *model.User) error {
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": user.ID}, user)
	return err
}

// Delete 删除用户
func (r *UserRepositoryButterfly) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// List 获取用户列表
func (r *UserRepositoryButterfly) List(ctx context.Context, offset, limit int) ([]*model.User, int64, error) {
	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSkip(int64(offset)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var users []*model.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// Ping 检查连接
func (r *UserRepositoryButterfly) Ping(ctx context.Context) error {
	return r.collection.Database().Client().Ping(ctx, nil)
}