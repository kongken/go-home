package repository

import (
	"context"

	"github.com/kongken/go-home/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// MessageRepository 消息仓库接口
type MessageRepository interface {
	Create(ctx context.Context, message *model.Message) error
	GetByID(ctx context.Context, id string) (*model.Message, error)
	ListByConversation(ctx context.Context, userID, otherUserID string, offset, limit int64) ([]*model.Message, int64, error)
	MarkAsRead(ctx context.Context, userID, otherUserID string) error
	GetUnreadCount(ctx context.Context, userID string) (int64, error)

	// 会话管理
	GetConversation(ctx context.Context, userID, otherUserID string) (*model.Conversation, error)
	CreateConversation(ctx context.Context, conversation *model.Conversation) error
	UpdateConversation(ctx context.Context, conversation *model.Conversation) error
	ListConversations(ctx context.Context, userID string, offset, limit int64) ([]*model.Conversation, int64, error)
}

// messageRepository 消息仓库实现
type messageRepository struct {
	collection            *mongo.Collection
	conversationCollection *mongo.Collection
}

// NewMessageRepository 创建消息仓库
func NewMessageRepository(db *mongo.Database) MessageRepository {
	return &messageRepository{
		collection:            db.Collection(model.Message{}.CollectionName()),
		conversationCollection: db.Collection(model.Conversation{}.CollectionName()),
	}
}

// Create 创建消息
func (r *messageRepository) Create(ctx context.Context, message *model.Message) error {
	message.BeforeInsert()
	_, err := r.collection.InsertOne(ctx, message)
	return err
}

// GetByID 根据ID获取消息
func (r *messageRepository) GetByID(ctx context.Context, id string) (*model.Message, error) {
	var message model.Message
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&message)
	if err != nil {
		return nil, err
	}
	return &message, nil
}

// ListByConversation 获取会话消息列表
func (r *messageRepository) ListByConversation(ctx context.Context, userID, otherUserID string, offset, limit int64) ([]*model.Message, int64, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"from_user_id": userID, "to_user_id": otherUserID},
			{"from_user_id": otherUserID, "to_user_id": userID},
		},
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

	var messages []*model.Message
	if err := cursor.All(ctx, &messages); err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}

// MarkAsRead 标记消息已读
func (r *messageRepository) MarkAsRead(ctx context.Context, userID, otherUserID string) error {
	_, err := r.collection.UpdateMany(
		ctx,
		bson.M{
			"from_user_id": otherUserID,
			"to_user_id":   userID,
			"is_read":      false,
		},
		bson.M{"$set": bson.M{"is_read": true}},
	)
	return err
}

// GetUnreadCount 获取未读消息数
func (r *messageRepository) GetUnreadCount(ctx context.Context, userID string) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{
		"to_user_id": userID,
		"is_read":    false,
	})
}

// GetConversation 获取会话
func (r *messageRepository) GetConversation(ctx context.Context, userID, otherUserID string) (*model.Conversation, error) {
	var conversation model.Conversation
	filter := bson.M{
		"user_id":       userID,
		"other_user_id": otherUserID,
	}
	err := r.conversationCollection.FindOne(ctx, filter).Decode(&conversation)
	if err != nil {
		return nil, err
	}
	return &conversation, nil
}

// CreateConversation 创建会话
func (r *messageRepository) CreateConversation(ctx context.Context, conversation *model.Conversation) error {
	conversation.BeforeInsert()
	_, err := r.conversationCollection.InsertOne(ctx, conversation)
	return err
}

// UpdateConversation 更新会话
func (r *messageRepository) UpdateConversation(ctx context.Context, conversation *model.Conversation) error {
	conversation.BeforeUpdate()
	_, err := r.conversationCollection.ReplaceOne(
		ctx,
		bson.M{"_id": conversation.ID},
		conversation,
	)
	return err
}

// ListConversations 获取会话列表
func (r *messageRepository) ListConversations(ctx context.Context, userID string, offset, limit int64) ([]*model.Conversation, int64, error) {
	filter := bson.M{"user_id": userID}

	total, err := r.conversationCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSkip(offset).
		SetLimit(limit).
		SetSort(bson.D{{Key: "updated_at", Value: -1}})

	cursor, err := r.conversationCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var conversations []*model.Conversation
	if err := cursor.All(ctx, &conversations); err != nil {
		return nil, 0, err
	}

	return conversations, total, nil
}
