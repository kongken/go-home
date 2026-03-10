package repository

import (
	"context"

	"github.com/kongken/go-home/internal/model"
	"butterfly.orx.me/core/store/mongo"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// SettingsRepositoryButterfly butterfly 版本的设置仓库
type SettingsRepositoryButterfly struct {
	collection *mongodriver.Collection
}

// NewSettingsRepositoryButterfly 创建 butterfly 设置仓库
func NewSettingsRepositoryButterfly() SettingsRepository {
	client := mongo.GetClient("primary")
	if client == nil {
		panic("mongo client 'primary' not found")
	}
	return &SettingsRepositoryButterfly{
		collection: client.Database("gohome").Collection("settings"),
	}
}

func (r *SettingsRepositoryButterfly) GetByUserID(ctx context.Context, userID string) (*model.Settings, error) {
	var settings model.Settings
	err := r.collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&settings)
	if err != nil {
		// 返回默认设置
		return &model.Settings{
			UserID: userID,
			Privacy: model.PrivacySettings{
				DefaultBlogPrivacy:  model.PrivacyPublic,
				DefaultAlbumPrivacy: model.PrivacyPublic,
				DefaultSharePrivacy: model.PrivacyPublic,
				DefaultFeedPrivacy:  model.PrivacyPublic,
			},
			Notification: model.NotificationSettings{
				NotifyFriendRequest: true,
				NotifyComment:       true,
				NotifyLike:          true,
				NotifyMention:       true,
				NotifyGroup:         true,
				NotifySystem:        true,
			},
			Blacklist: []string{},
		}, nil
	}
	return &settings, nil
}

func (r *SettingsRepositoryButterfly) Create(ctx context.Context, settings *model.Settings) error {
	settings.BeforeInsert()
	_, err := r.collection.InsertOne(ctx, settings)
	return err
}

func (r *SettingsRepositoryButterfly) Update(ctx context.Context, settings *model.Settings) error {
	settings.BeforeUpdate()
	_, err := r.collection.ReplaceOne(ctx, bson.M{"user_id": settings.UserID}, settings)
	return err
}

func (r *SettingsRepositoryButterfly) AddToBlacklist(ctx context.Context, userID, blockedUserID string) error {
	_, err := r.collection.UpdateOne(ctx,
		bson.M{"user_id": userID},
		bson.M{"$addToSet": bson.M{"blacklist": blockedUserID}},
	)
	return err
}

func (r *SettingsRepositoryButterfly) RemoveFromBlacklist(ctx context.Context, userID, blockedUserID string) error {
	_, err := r.collection.UpdateOne(ctx,
		bson.M{"user_id": userID},
		bson.M{"$pull": bson.M{"blacklist": blockedUserID}},
	)
	return err
}