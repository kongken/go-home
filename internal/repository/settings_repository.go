package repository

import (
	"context"

	"github.com/kongken/go-home/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// SettingsRepository 设置仓库接口
type SettingsRepository interface {
	GetByUserID(ctx context.Context, userID string) (*model.Settings, error)
	Create(ctx context.Context, settings *model.Settings) error
	Update(ctx context.Context, settings *model.Settings) error
	AddToBlacklist(ctx context.Context, userID, blockedUserID string) error
	RemoveFromBlacklist(ctx context.Context, userID, blockedUserID string) error
}

// settingsRepository 设置仓库实现
type settingsRepository struct {
	collection *mongo.Collection
}

// NewSettingsRepository 创建设置仓库
func NewSettingsRepository(db *mongo.Database) SettingsRepository {
	return &settingsRepository{
		collection: db.Collection(model.Settings{}.CollectionName()),
	}
}

// GetByUserID 根据用户ID获取设置
func (r *settingsRepository) GetByUserID(ctx context.Context, userID string) (*model.Settings, error) {
	var settings model.Settings
	err := r.collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&settings)
	if err != nil {
		if err == mongo.ErrNoDocuments {
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
		return nil, err
	}
	return &settings, nil
}

// Create 创建设置
func (r *settingsRepository) Create(ctx context.Context, settings *model.Settings) error {
	settings.BeforeInsert()
	_, err := r.collection.InsertOne(ctx, settings)
	return err
}

// Update 更新设置
func (r *settingsRepository) Update(ctx context.Context, settings *model.Settings) error {
	settings.BeforeUpdate()
	_, err := r.collection.ReplaceOne(
		ctx,
		bson.M{"user_id": settings.UserID},
		settings,
	)
	return err
}

// AddToBlacklist 添加到黑名单
func (r *settingsRepository) AddToBlacklist(ctx context.Context, userID, blockedUserID string) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"user_id": userID},
		bson.M{"$addToSet": bson.M{"blacklist": blockedUserID}},
	)
	return err
}

// RemoveFromBlacklist 从黑名单移除
func (r *settingsRepository) RemoveFromBlacklist(ctx context.Context, userID, blockedUserID string) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"user_id": userID},
		bson.M{"$pull": bson.M{"blacklist": blockedUserID}},
	)
	return err
}