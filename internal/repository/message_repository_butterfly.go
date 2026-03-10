package repository

import (
	"context"

	"github.com/kongken/go-home/internal/model"
	"butterfly.orx.me/core/store/mongo"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// MessageRepositoryButterfly butterfly 版本的消息仓库
type MessageRepositoryButterfly struct {
	collection            *mongodriver.Collection
	conversationCollection *mongodriver.Collection
}

// NewMessageRepositoryButterfly 创建 butterfly 消息仓库
func NewMessageRepositoryButterfly() MessageRepository {
	client := mongo.GetClient("primary")
	if client == nil {
		panic("mongo client 'primary' not found")
	}
	db := client.Database("gohome")
	return &MessageRepositoryButterfly{
		collection:             db.Collection("messages"),
		conversationCollection: db.Collection("conversations"),
	}
}

func (r *MessageRepositoryButterfly) Create(ctx context.Context, message *model.Message) error {
	message.BeforeInsert()
	_, err := r.collection.InsertOne(ctx, message)
	return err
}

func (r *MessageRepositoryButterfly) GetByID(ctx context.Context, id string) (*model.Message, error) {
	var message model.Message
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&message)
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (r *MessageRepositoryButterfly) ListByConversation(ctx context.Context, userID, otherUserID string, offset, limit int64) ([]*model.Message, int64, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"from_user_id": userID, "to_user_id": otherUserID},
			{"from_user_id": otherUserID, "to_user_id": userID},
		},
	}
	total, _ := r.collection.CountDocuments(ctx, filter)
	opts := options.Find().SetSkip(offset).SetLimit(limit).SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, _ := r.collection.Find(ctx, filter, opts)
	defer cursor.Close(ctx)
	var messages []*model.Message
	cursor.All(ctx, &messages)
	return messages, total, nil
}

func (r *MessageRepositoryButterfly) MarkAsRead(ctx context.Context, userID, otherUserID string) error {
	_, err := r.collection.UpdateMany(ctx,
		bson.M{"from_user_id": otherUserID, "to_user_id": userID, "is_read": false},
		bson.M{"$set": bson.M{"is_read": true}},
	)
	return err
}

func (r *MessageRepositoryButterfly) GetUnreadCount(ctx context.Context, userID string) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{"to_user_id": userID, "is_read": false})
}

func (r *MessageRepositoryButterfly) GetConversation(ctx context.Context, userID, otherUserID string) (*model.Conversation, error) {
	var conversation model.Conversation
	err := r.conversationCollection.FindOne(ctx, bson.M{"user_id": userID, "other_user_id": otherUserID}).Decode(&conversation)
	if err != nil {
		return nil, err
	}
	return &conversation, nil
}

func (r *MessageRepositoryButterfly) CreateConversation(ctx context.Context, conversation *model.Conversation) error {
	conversation.BeforeInsert()
	_, err := r.conversationCollection.InsertOne(ctx, conversation)
	return err
}

func (r *MessageRepositoryButterfly) UpdateConversation(ctx context.Context, conversation *model.Conversation) error {
	conversation.BeforeUpdate()
	_, err := r.conversationCollection.ReplaceOne(ctx, bson.M{"_id": conversation.ID}, conversation)
	return err
}

func (r *MessageRepositoryButterfly) ListConversations(ctx context.Context, userID string, offset, limit int64) ([]*model.Conversation, int64, error) {
	filter := bson.M{"user_id": userID}
	total, _ := r.conversationCollection.CountDocuments(ctx, filter)
	opts := options.Find().SetSkip(offset).SetLimit(limit).SetSort(bson.D{{Key: "updated_at", Value: -1}})
	cursor, _ := r.conversationCollection.Find(ctx, filter, opts)
	defer cursor.Close(ctx)
	var conversations []*model.Conversation
	cursor.All(ctx, &conversations)
	return conversations, total, nil
}