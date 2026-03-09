package repository

import (
	"context"
	"time"

	"github.com/kongken/go-home/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// FriendRepository 好友仓库接口
type FriendRepository interface {
	// 好友关系
	CreateFriendship(ctx context.Context, friendship *model.Friendship) error
	GetFriendship(ctx context.Context, userID, friendID string) (*model.Friendship, error)
	DeleteFriendship(ctx context.Context, userID, friendID string) error
	ListFriendships(ctx context.Context, userID string, offset, limit int64) ([]*model.Friendship, int64, error)
	ListByGroup(ctx context.Context, userID, groupName string, offset, limit int64) ([]*model.Friendship, int64, error)
	UpdateFriendship(ctx context.Context, friendship *model.Friendship) error

	// 好友请求
	CreateRequest(ctx context.Context, request *model.FriendRequest) error
	GetRequest(ctx context.Context, id string) (*model.FriendRequest, error)
	GetRequestByUsers(ctx context.Context, fromUserID, toUserID string) (*model.FriendRequest, error)
	UpdateRequest(ctx context.Context, request *model.FriendRequest) error
	ListReceivedRequests(ctx context.Context, userID string, offset, limit int64) ([]*model.FriendRequest, int64, error)
	ListSentRequests(ctx context.Context, userID string, offset, limit int64) ([]*model.FriendRequest, int64, error)
}

// friendRepository 好友仓库实现
type friendRepository struct {
	friendshipColl *mongo.Collection
	requestColl    *mongo.Collection
}

// NewFriendRepository 创建好友仓库
func NewFriendRepository(db *mongo.Database) FriendRepository {
	return &friendRepository{
		friendshipColl: db.Collection(model.Friendship{}.CollectionName()),
		requestColl:    db.Collection(model.FriendRequest{}.CollectionName()),
	}
}

// CreateFriendship 创建好友关系
func (r *friendRepository) CreateFriendship(ctx context.Context, friendship *model.Friendship) error {
	friendship.BeforeInsert()
	_, err := r.friendshipColl.InsertOne(ctx, friendship)
	return err
}

// GetFriendship 获取好友关系
func (r *friendRepository) GetFriendship(ctx context.Context, userID, friendID string) (*model.Friendship, error) {
	var friendship model.Friendship
	filter := bson.M{
		"user_id":   userID,
		"friend_id": friendID,
	}
	err := r.friendshipColl.FindOne(ctx, filter).Decode(&friendship)
	if err != nil {
		return nil, err
	}
	return &friendship, nil
}

// DeleteFriendship 删除好友关系
func (r *friendRepository) DeleteFriendship(ctx context.Context, userID, friendID string) error {
	filter := bson.M{
		"$or": []bson.M{
			{"user_id": userID, "friend_id": friendID},
			{"user_id": friendID, "friend_id": userID},
		},
	}
	_, err := r.friendshipColl.DeleteMany(ctx, filter)
	return err
}

// ListFriendships 获取好友列表
func (r *friendRepository) ListFriendships(ctx context.Context, userID string, offset, limit int64) ([]*model.Friendship, int64, error) {
	filter := bson.M{"user_id": userID, "status": model.FriendshipAccepted}

	total, err := r.friendshipColl.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSkip(offset).
		SetLimit(limit).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.friendshipColl.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var friendships []*model.Friendship
	if err := cursor.All(ctx, &friendships); err != nil {
		return nil, 0, err
	}

	return friendships, total, nil
}

// ListByGroup 按分组获取好友
func (r *friendRepository) ListByGroup(ctx context.Context, userID, groupName string, offset, limit int64) ([]*model.Friendship, int64, error) {
	filter := bson.M{
		"user_id":    userID,
		"group_name": groupName,
		"status":     model.FriendshipAccepted,
	}

	total, err := r.friendshipColl.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSkip(offset).
		SetLimit(limit).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.friendshipColl.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var friendships []*model.Friendship
	if err := cursor.All(ctx, &friendships); err != nil {
		return nil, 0, err
	}

	return friendships, total, nil
}

// UpdateFriendship 更新好友关系
func (r *friendRepository) UpdateFriendship(ctx context.Context, friendship *model.Friendship) error {
	friendship.UpdatedAt = time.Now()
	_, err := r.friendshipColl.ReplaceOne(
		ctx,
		bson.M{"_id": friendship.ID},
		friendship,
	)
	return err
}

// CreateRequest 创建好友请求
func (r *friendRepository) CreateRequest(ctx context.Context, request *model.FriendRequest) error {
	request.BeforeInsert()
	_, err := r.requestColl.InsertOne(ctx, request)
	return err
}

// GetRequest 获取好友请求
func (r *friendRepository) GetRequest(ctx context.Context, id string) (*model.FriendRequest, error) {
	var request model.FriendRequest
	err := r.requestColl.FindOne(ctx, bson.M{"_id": id}).Decode(&request)
	if err != nil {
		return nil, err
	}
	return &request, nil
}

// GetRequestByUsers 根据用户获取请求
func (r *friendRepository) GetRequestByUsers(ctx context.Context, fromUserID, toUserID string) (*model.FriendRequest, error) {
	var request model.FriendRequest
	filter := bson.M{
		"from_user_id": fromUserID,
		"to_user_id":   toUserID,
		"status":       model.RequestPending,
	}
	err := r.requestColl.FindOne(ctx, filter).Decode(&request)
	if err != nil {
		return nil, err
	}
	return &request, nil
}

// UpdateRequest 更新好友请求
func (r *friendRepository) UpdateRequest(ctx context.Context, request *model.FriendRequest) error {
	request.UpdatedAt = time.Now()
	_, err := r.requestColl.ReplaceOne(
		ctx,
		bson.M{"_id": request.ID},
		request,
	)
	return err
}

// ListReceivedRequests 获取收到的好友请求
func (r *friendRepository) ListReceivedRequests(ctx context.Context, userID string, offset, limit int64) ([]*model.FriendRequest, int64, error) {
	filter := bson.M{
		"to_user_id": userID,
		"status":     model.RequestPending,
	}

	total, err := r.requestColl.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSkip(offset).
		SetLimit(limit).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.requestColl.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var requests []*model.FriendRequest
	if err := cursor.All(ctx, &requests); err != nil {
		return nil, 0, err
	}

	return requests, total, nil
}

// ListSentRequests 获取发送的好友请求
func (r *friendRepository) ListSentRequests(ctx context.Context, userID string, offset, limit int64) ([]*model.FriendRequest, int64, error) {
	filter := bson.M{
		"from_user_id": userID,
		"status":       model.RequestPending,
	}

	total, err := r.requestColl.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSkip(offset).
		SetLimit(limit).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.requestColl.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var requests []*model.FriendRequest
	if err := cursor.All(ctx, &requests); err != nil {
		return nil, 0, err
	}

	return requests, total, nil
}
