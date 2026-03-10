package repository

import (
	"context"

	"github.com/kongken/go-home/internal/model"
	"butterfly.orx.me/core/store/mongo"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// GroupRepositoryButterfly butterfly 版本的群组仓库
type GroupRepositoryButterfly struct {
	collection       *mongodriver.Collection
	memberCollection *mongodriver.Collection
}

// NewGroupRepositoryButterfly 创建 butterfly 群组仓库
func NewGroupRepositoryButterfly() GroupRepository {
	client := mongo.GetClient("primary")
	if client == nil {
		panic("mongo client 'primary' not found")
	}
	db := client.Database("gohome")
	return &GroupRepositoryButterfly{
		collection:       db.Collection("groups"),
		memberCollection: db.Collection("group_members"),
	}
}

func (r *GroupRepositoryButterfly) Create(ctx context.Context, group *model.Group) error {
	group.BeforeInsert()
	_, err := r.collection.InsertOne(ctx, group)
	return err
}

func (r *GroupRepositoryButterfly) GetByID(ctx context.Context, id string) (*model.Group, error) {
	var group model.Group
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&group)
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *GroupRepositoryButterfly) Update(ctx context.Context, group *model.Group) error {
	group.BeforeUpdate()
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": group.ID}, group)
	return err
}

func (r *GroupRepositoryButterfly) Delete(ctx context.Context, id string) error {
	r.collection.DeleteOne(ctx, bson.M{"_id": id})
	r.memberCollection.DeleteMany(ctx, bson.M{"group_id": id})
	return nil
}

func (r *GroupRepositoryButterfly) List(ctx context.Context, category string, offset, limit int64) ([]*model.Group, int64, error) {
	filter := bson.M{}
	if category != "" {
		filter["category"] = category
	}
	total, _ := r.collection.CountDocuments(ctx, filter)
	opts := options.Find().SetSkip(offset).SetLimit(limit).SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, _ := r.collection.Find(ctx, filter, opts)
	defer cursor.Close(ctx)
	var groups []*model.Group
	cursor.All(ctx, &groups)
	return groups, total, nil
}

func (r *GroupRepositoryButterfly) Search(ctx context.Context, keyword string, offset, limit int64) ([]*model.Group, int64, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"name": bson.M{"$regex": keyword, "$options": "i"}},
			{"description": bson.M{"$regex": keyword, "$options": "i"}},
		},
	}
	total, _ := r.collection.CountDocuments(ctx, filter)
	opts := options.Find().SetSkip(offset).SetLimit(limit)
	cursor, _ := r.collection.Find(ctx, filter, opts)
	defer cursor.Close(ctx)
	var groups []*model.Group
	cursor.All(ctx, &groups)
	return groups, total, nil
}

func (r *GroupRepositoryButterfly) AddMember(ctx context.Context, member *model.GroupMember) error {
	member.BeforeInsert()
	_, err := r.memberCollection.InsertOne(ctx, member)
	if err != nil {
		return err
	}
	r.collection.UpdateOne(ctx, bson.M{"_id": member.GroupID}, bson.M{"$inc": bson.M{"members_count": 1}})
	return nil
}

func (r *GroupRepositoryButterfly) RemoveMember(ctx context.Context, groupID, userID string) error {
	r.memberCollection.DeleteOne(ctx, bson.M{"group_id": groupID, "user_id": userID})
	r.collection.UpdateOne(ctx, bson.M{"_id": groupID}, bson.M{"$inc": bson.M{"members_count": -1}})
	return nil
}

func (r *GroupRepositoryButterfly) GetMember(ctx context.Context, groupID, userID string) (*model.GroupMember, error) {
	var member model.GroupMember
	err := r.memberCollection.FindOne(ctx, bson.M{"group_id": groupID, "user_id": userID}).Decode(&member)
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *GroupRepositoryButterfly) UpdateMemberRole(ctx context.Context, groupID, userID string, role model.MemberRole) error {
	_, err := r.memberCollection.UpdateOne(ctx, bson.M{"group_id": groupID, "user_id": userID}, bson.M{"$set": bson.M{"role": role}})
	return err
}

func (r *GroupRepositoryButterfly) ListMembers(ctx context.Context, groupID string, offset, limit int64) ([]*model.GroupMember, int64, error) {
	filter := bson.M{"group_id": groupID}
	total, _ := r.memberCollection.CountDocuments(ctx, filter)
	opts := options.Find().SetSkip(offset).SetLimit(limit)
	cursor, _ := r.memberCollection.Find(ctx, filter, opts)
	defer cursor.Close(ctx)
	var members []*model.GroupMember
	cursor.All(ctx, &members)
	return members, total, nil
}

func (r *GroupRepositoryButterfly) CountMembers(ctx context.Context, groupID string) (int64, error) {
	return r.memberCollection.CountDocuments(ctx, bson.M{"group_id": groupID})
}