package repository

import (
	"context"

	"github.com/kongken/go-home/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// GroupRepository 群组仓库接口
type GroupRepository interface {
	Create(ctx context.Context, group *model.Group) error
	GetByID(ctx context.Context, id string) (*model.Group, error)
	Update(ctx context.Context, group *model.Group) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, category string, offset, limit int64) ([]*model.Group, int64, error)
	Search(ctx context.Context, keyword string, offset, limit int64) ([]*model.Group, int64, error)

	// 成员管理
	AddMember(ctx context.Context, member *model.GroupMember) error
	RemoveMember(ctx context.Context, groupID, userID string) error
	GetMember(ctx context.Context, groupID, userID string) (*model.GroupMember, error)
	UpdateMemberRole(ctx context.Context, groupID, userID string, role model.MemberRole) error
	ListMembers(ctx context.Context, groupID string, offset, limit int64) ([]*model.GroupMember, int64, error)
	CountMembers(ctx context.Context, groupID string) (int64, error)
}

// groupRepository 群组仓库实现
type groupRepository struct {
	collection       *mongo.Collection
	memberCollection *mongo.Collection
}

// NewGroupRepository 创建群组仓库
func NewGroupRepository(db *mongo.Database) GroupRepository {
	return &groupRepository{
		collection:       db.Collection(model.Group{}.CollectionName()),
		memberCollection: db.Collection(model.GroupMember{}.CollectionName()),
	}
}

// Create 创建群组
func (r *groupRepository) Create(ctx context.Context, group *model.Group) error {
	group.BeforeInsert()
	_, err := r.collection.InsertOne(ctx, group)
	return err
}

// GetByID 根据ID获取群组
func (r *groupRepository) GetByID(ctx context.Context, id string) (*model.Group, error) {
	var group model.Group
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&group)
	if err != nil {
		return nil, err
	}
	return &group, nil
}

// Update 更新群组
func (r *groupRepository) Update(ctx context.Context, group *model.Group) error {
	group.BeforeUpdate()
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": group.ID}, group)
	return err
}

// Delete 删除群组
func (r *groupRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	// 删除所有成员
	_, err = r.memberCollection.DeleteMany(ctx, bson.M{"group_id": id})
	return err
}

// List 获取群组列表
func (r *groupRepository) List(ctx context.Context, category string, offset, limit int64) ([]*model.Group, int64, error) {
	filter := bson.M{}
	if category != "" {
		filter["category"] = category
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

	var groups []*model.Group
	if err := cursor.All(ctx, &groups); err != nil {
		return nil, 0, err
	}

	return groups, total, nil
}

// Search 搜索群组
func (r *groupRepository) Search(ctx context.Context, keyword string, offset, limit int64) ([]*model.Group, int64, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"name": bson.M{"$regex": keyword, "$options": "i"}},
			{"description": bson.M{"$regex": keyword, "$options": "i"}},
		},
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSkip(offset).
		SetLimit(limit).
		SetSort(bson.D{{Key: "members_count", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var groups []*model.Group
	if err := cursor.All(ctx, &groups); err != nil {
		return nil, 0, err
	}

	return groups, total, nil
}

// AddMember 添加成员
func (r *groupRepository) AddMember(ctx context.Context, member *model.GroupMember) error {
	member.BeforeInsert()
	_, err := r.memberCollection.InsertOne(ctx, member)
	if err != nil {
		return err
	}
	// 更新成员数
	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": member.GroupID},
		bson.M{"$inc": bson.M{"members_count": 1}},
	)
	return err
}

// RemoveMember 移除成员
func (r *groupRepository) RemoveMember(ctx context.Context, groupID, userID string) error {
	_, err := r.memberCollection.DeleteOne(ctx, bson.M{
		"group_id": groupID,
		"user_id":  userID,
	})
	if err != nil {
		return err
	}
	// 更新成员数
	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": groupID},
		bson.M{"$inc": bson.M{"members_count": -1}},
	)
	return err
}

// GetMember 获取成员
func (r *groupRepository) GetMember(ctx context.Context, groupID, userID string) (*model.GroupMember, error) {
	var member model.GroupMember
	err := r.memberCollection.FindOne(ctx, bson.M{
		"group_id": groupID,
		"user_id":  userID,
	}).Decode(&member)
	if err != nil {
		return nil, err
	}
	return &member, nil
}

// UpdateMemberRole 更新成员角色
func (r *groupRepository) UpdateMemberRole(ctx context.Context, groupID, userID string, role model.MemberRole) error {
	_, err := r.memberCollection.UpdateOne(
		ctx,
		bson.M{"group_id": groupID, "user_id": userID},
		bson.M{"$set": bson.M{"role": role}},
	)
	return err
}

// ListMembers 获取成员列表
func (r *groupRepository) ListMembers(ctx context.Context, groupID string, offset, limit int64) ([]*model.GroupMember, int64, error) {
	filter := bson.M{"group_id": groupID}

	total, err := r.memberCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSkip(offset).
		SetLimit(limit).
		SetSort(bson.D{{Key: "joined_at", Value: -1}})

	cursor, err := r.memberCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var members []*model.GroupMember
	if err := cursor.All(ctx, &members); err != nil {
		return nil, 0, err
	}

	return members, total, nil
}

// CountMembers 统计成员数
func (r *groupRepository) CountMembers(ctx context.Context, groupID string) (int64, error) {
	return r.memberCollection.CountDocuments(ctx, bson.M{"group_id": groupID})
}
