package repository

import (
	"context"
	"time"

	"github.com/kongken/go-home/internal/model"
	"butterfly.orx.me/core/store/mongo"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// FriendRepositoryButterfly butterfly 版本的好友仓库
type FriendRepositoryButterfly struct {
	friendshipColl *mongodriver.Collection
	requestColl    *mongodriver.Collection
}

// NewFriendRepositoryButterfly 创建 butterfly 好友仓库
func NewFriendRepositoryButterfly() FriendRepository {
	client := mongo.GetClient("primary")
	if client == nil {
		panic("mongo client 'primary' not found")
	}
	db := client.Database("gohome")
	return &FriendRepositoryButterfly{
		friendshipColl: db.Collection("friendships"),
		requestColl:    db.Collection("friend_requests"),
	}
}

// CreateFriendship 创建好友关系
func (r *FriendRepositoryButterfly) CreateFriendship(ctx context.Context, friendship *model.Friendship) error {
	friendship.BeforeInsert()
	_, err := r.friendshipColl.InsertOne(ctx, friendship)
	return err
}

// GetFriendship 获取好友关系
func (r *FriendRepositoryButterfly) GetFriendship(ctx context.Context, userID, friendID string) (*model.Friendship, error) {
	var friendship model.Friendship
	err := r.friendshipColl.FindOne(ctx, bson.M{
		"user_id":   userID,
		"friend_id": friendID,
	}).Decode(&friendship)
	if err != nil {
		return nil, err
	}
	return &friendship, nil
}

// DeleteFriendship 删除好友关系
func (r *FriendRepositoryButterfly) DeleteFriendship(ctx context.Context, userID, friendID string) error {
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
func (r *FriendRepositoryButterfly) ListFriendships(ctx context.Context, userID string, offset, limit int64) ([]*model.Friendship, int64, error) {
	filter := bson.M{"user_id": userID, "status": model.FriendshipAccepted}

	total, err := r.friendshipColl.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().SetSkip(offset).SetLimit(limit).SetSort(bson.D{{Key: "created_at", Value: -1}})
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
func (r *FriendRepositoryButterfly) ListByGroup(ctx context.Context, userID, groupName string, offset, limit int64) ([]*model.Friendship, int64, error) {
	filter := bson.M{"user_id": userID, "group_name": groupName, "status": model.FriendshipAccepted}
	total, err := r.friendshipColl.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	opts := options.Find().SetSkip(offset).SetLimit(limit)
	cursor, err := r.friendshipColl.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	var friendships []*model.Friendship
	cursor.All(ctx, &friendships)
	return friendships, total, nil
}

// UpdateFriendship 更新好友关系
func (r *FriendRepositoryButterfly) UpdateFriendship(ctx context.Context, friendship *model.Friendship) error {
	friendship.UpdatedAt = time.Now()
	_, err := r.friendshipColl.ReplaceOne(ctx, bson.M{"_id": friendship.ID}, friendship)
	return err
}

// CreateRequest 创建好友请求
func (r *FriendRepositoryButterfly) CreateRequest(ctx context.Context, request *model.FriendRequest) error {
	request.BeforeInsert()
	_, err := r.requestColl.InsertOne(ctx, request)
	return err
}

// GetRequest 获取好友请求
func (r *FriendRepositoryButterfly) GetRequest(ctx context.Context, id string) (*model.FriendRequest, error) {
	var request model.FriendRequest
	err := r.requestColl.FindOne(ctx, bson.M{"_id": id}).Decode(&request)
	if err != nil {
		return nil, err
	}
	return &request, nil
}

// GetRequestByUsers 根据用户获取请求
func (r *FriendRepositoryButterfly) GetRequestByUsers(ctx context.Context, fromUserID, toUserID string) (*model.FriendRequest, error) {
	var request model.FriendRequest
	err := r.requestColl.FindOne(ctx, bson.M{
		"from_user_id": fromUserID,
		"to_user_id":   toUserID,
		"status":       model.RequestPending,
	}).Decode(&request)
	if err != nil {
		return nil, err
	}
	return &request, nil
}

// UpdateRequest 更新好友请求
func (r *FriendRepositoryButterfly) UpdateRequest(ctx context.Context, request *model.FriendRequest) error {
	request.UpdatedAt = time.Now()
	_, err := r.requestColl.ReplaceOne(ctx, bson.M{"_id": request.ID}, request)
	return err
}

// ListReceivedRequests 获取收到的好友请求
func (r *FriendRepositoryButterfly) ListReceivedRequests(ctx context.Context, userID string, offset, limit int64) ([]*model.FriendRequest, int64, error) {
	filter := bson.M{"to_user_id": userID, "status": model.RequestPending}
	total, err := r.requestColl.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	opts := options.Find().SetSkip(offset).SetLimit(limit).SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := r.requestColl.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	var requests []*model.FriendRequest
	cursor.All(ctx, &requests)
	return requests, total, nil
}

// ListSentRequests 获取发送的好友请求
func (r *FriendRepositoryButterfly) ListSentRequests(ctx context.Context, userID string, offset, limit int64) ([]*model.FriendRequest, int64, error) {
	filter := bson.M{"from_user_id": userID, "status": model.RequestPending}
	total, err := r.requestColl.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	opts := options.Find().SetSkip(offset).SetLimit(limit)
	cursor, err := r.requestColl.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	var requests []*model.FriendRequest
	cursor.All(ctx, &requests)
	return requests, total, nil
}