package repository

import (
	"context"

	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// NewUserRepositoryWithClient creates a UserRepository with provided MongoDB client
// This is useful for testing with a real database connection
func NewUserRepositoryWithClient(client *mongodriver.Client, dbName string) UserRepository {
	if dbName == "" {
		dbName = "gohome"
	}
	return &UserRepositoryButterfly{
		collection: client.Database(dbName).Collection("users"),
	}
}

// NewBlogRepositoryWithClient creates a BlogRepository with provided MongoDB client
func NewBlogRepositoryWithClient(client *mongodriver.Client, dbName string) BlogRepository {
	if dbName == "" {
		dbName = "gohome"
	}
	return &BlogRepositoryButterfly{
		collection: client.Database(dbName).Collection("blogs"),
	}
}

// NewFeedRepositoryWithClient creates a FeedRepository with provided MongoDB client
func NewFeedRepositoryWithClient(client *mongodriver.Client, dbName string) FeedRepository {
	if dbName == "" {
		dbName = "gohome"
	}
	return &FeedRepositoryButterfly{
		collection: client.Database(dbName).Collection("feeds"),
	}
}

// NewFriendRepositoryWithClient creates a FriendRepository with provided MongoDB client
func NewFriendRepositoryWithClient(client *mongodriver.Client, dbName string) FriendRepository {
	if dbName == "" {
		dbName = "gohome"
	}
	db := client.Database(dbName)
	return &FriendRepositoryButterfly{
		friendshipColl: db.Collection("friendships"),
		requestColl:    db.Collection("friend_requests"),
	}
}

// NewGroupRepositoryWithClient creates a GroupRepository with provided MongoDB client
func NewGroupRepositoryWithClient(client *mongodriver.Client, dbName string) GroupRepository {
	if dbName == "" {
		dbName = "gohome"
	}
	return &GroupRepositoryButterfly{
		collection: client.Database(dbName).Collection("groups"),
	}
}

// NewMessageRepositoryWithClient creates a MessageRepository with provided MongoDB client
func NewMessageRepositoryWithClient(client *mongodriver.Client, dbName string) MessageRepository {
	if dbName == "" {
		dbName = "gohome"
	}
	return &MessageRepositoryButterfly{
		collection: client.Database(dbName).Collection("messages"),
	}
}

// NewCommentRepositoryWithClient creates a CommentRepository with provided MongoDB client
func NewCommentRepositoryWithClient(client *mongodriver.Client, dbName string) CommentRepository {
	if dbName == "" {
		dbName = "gohome"
	}
	return &CommentRepositoryButterfly{
		collection: client.Database(dbName).Collection("comments"),
	}
}

// NewAlbumRepositoryWithClient creates an AlbumRepository with provided MongoDB client
func NewAlbumRepositoryWithClient(client *mongodriver.Client, dbName string) AlbumRepository {
	if dbName == "" {
		dbName = "gohome"
	}
	return &AlbumRepositoryButterfly{
		collection: client.Database(dbName).Collection("albums"),
	}
}

// NewNotificationRepositoryWithClient creates a NotificationRepository with provided MongoDB client
func NewNotificationRepositoryWithClient(client *mongodriver.Client, dbName string) NotificationRepository {
	if dbName == "" {
		dbName = "gohome"
	}
	return &NotificationRepositoryButterfly{
		collection: client.Database(dbName).Collection("notifications"),
	}
}

// NewSettingsRepositoryWithClient creates a SettingsRepository with provided MongoDB client
func NewSettingsRepositoryWithClient(client *mongodriver.Client, dbName string) SettingsRepository {
	if dbName == "" {
		dbName = "gohome"
	}
	return &SettingsRepositoryButterfly{
		collection: client.Database(dbName).Collection("settings"),
	}
}

// TestCollectionGetter interface for getting collection for testing
type TestCollectionGetter interface {
	GetTestCollection() *mongodriver.Collection
}

// GetTestCollection returns the collection for testing (UserRepositoryButterfly)
func (r *UserRepositoryButterfly) GetTestCollection() *mongodriver.Collection {
	return r.collection
}

// GetTestCollection returns the collection for testing (BlogRepositoryButterfly)
func (r *BlogRepositoryButterfly) GetTestCollection() *mongodriver.Collection {
	return r.collection
}

// GetTestCollection returns the collection for testing (FeedRepositoryButterfly)
func (r *FeedRepositoryButterfly) GetTestCollection() *mongodriver.Collection {
	return r.collection
}

// GetTestCollection returns the collection for testing (GroupRepositoryButterfly)
func (r *GroupRepositoryButterfly) GetTestCollection() *mongodriver.Collection {
	return r.collection
}

// GetTestCollection returns the collection for testing (MessageRepositoryButterfly)
func (r *MessageRepositoryButterfly) GetTestCollection() *mongodriver.Collection {
	return r.collection
}

// GetTestCollection returns the collection for testing (CommentRepositoryButterfly)
func (r *CommentRepositoryButterfly) GetTestCollection() *mongodriver.Collection {
	return r.collection
}

// GetTestCollection returns the collection for testing (AlbumRepositoryButterfly)
func (r *AlbumRepositoryButterfly) GetTestCollection() *mongodriver.Collection {
	return r.collection
}

// GetTestCollection returns the collection for testing (NotificationRepositoryButterfly)
func (r *NotificationRepositoryButterfly) GetTestCollection() *mongodriver.Collection {
	return r.collection
}

// GetTestCollection returns the collection for testing (SettingsRepositoryButterfly)
func (r *SettingsRepositoryButterfly) GetTestCollection() *mongodriver.Collection {
	return r.collection
}

// DeleteAllDocuments deletes all documents from a collection (for test cleanup)
func DeleteAllDocuments(ctx context.Context, repo interface{}) error {
	var coll *mongodriver.Collection

	switch r := repo.(type) {
	case *UserRepositoryButterfly:
		coll = r.collection
	case *BlogRepositoryButterfly:
		coll = r.collection
	case *FeedRepositoryButterfly:
		coll = r.collection
	case *GroupRepositoryButterfly:
		coll = r.collection
	case *MessageRepositoryButterfly:
		coll = r.collection
	case *CommentRepositoryButterfly:
		coll = r.collection
	case *AlbumRepositoryButterfly:
		coll = r.collection
	case *NotificationRepositoryButterfly:
		coll = r.collection
	case *SettingsRepositoryButterfly:
		coll = r.collection
	case TestCollectionGetter:
		coll = r.GetTestCollection()
	default:
		return nil
	}

	_, err := coll.DeleteMany(ctx, bson.M{})
	return err
}

// FindOne finds a single document by filter (for testing)
func FindOne(ctx context.Context, repo interface{}, filter bson.M, result interface{}) error {
	var coll *mongodriver.Collection

	switch r := repo.(type) {
	case *UserRepositoryButterfly:
		coll = r.collection
	case *BlogRepositoryButterfly:
		coll = r.collection
	case *FeedRepositoryButterfly:
		coll = r.collection
	case *GroupRepositoryButterfly:
		coll = r.collection
	case *MessageRepositoryButterfly:
		coll = r.collection
	case *CommentRepositoryButterfly:
		coll = r.collection
	case *AlbumRepositoryButterfly:
		coll = r.collection
	case *NotificationRepositoryButterfly:
		coll = r.collection
	case *SettingsRepositoryButterfly:
		coll = r.collection
	case TestCollectionGetter:
		coll = r.GetTestCollection()
	default:
		return nil
	}

	return coll.FindOne(ctx, filter).Decode(result)
}

// FindAll finds all documents matching filter (for testing)
func FindAll(ctx context.Context, repo interface{}, filter bson.M, results interface{}) error {
	var coll *mongodriver.Collection

	switch r := repo.(type) {
	case *UserRepositoryButterfly:
		coll = r.collection
	case *BlogRepositoryButterfly:
		coll = r.collection
	case *FeedRepositoryButterfly:
		coll = r.collection
	case *GroupRepositoryButterfly:
		coll = r.collection
	case *MessageRepositoryButterfly:
		coll = r.collection
	case *CommentRepositoryButterfly:
		coll = r.collection
	case *AlbumRepositoryButterfly:
		coll = r.collection
	case *NotificationRepositoryButterfly:
		coll = r.collection
	case *SettingsRepositoryButterfly:
		coll = r.collection
	case TestCollectionGetter:
		coll = r.GetTestCollection()
	default:
		return nil
	}

	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	return cursor.All(ctx, results)
}

// FindWithOptions finds documents with options (for testing)
func FindWithOptions(ctx context.Context, repo interface{}, filter bson.M, results interface{}, skip int64, limit int64) error {
	var coll *mongodriver.Collection

	switch r := repo.(type) {
	case *UserRepositoryButterfly:
		coll = r.collection
	case *BlogRepositoryButterfly:
		coll = r.collection
	case *FeedRepositoryButterfly:
		coll = r.collection
	case *GroupRepositoryButterfly:
		coll = r.collection
	case *MessageRepositoryButterfly:
		coll = r.collection
	case *CommentRepositoryButterfly:
		coll = r.collection
	case *AlbumRepositoryButterfly:
		coll = r.collection
	case *NotificationRepositoryButterfly:
		coll = r.collection
	case *SettingsRepositoryButterfly:
		coll = r.collection
	case TestCollectionGetter:
		coll = r.GetTestCollection()
	default:
		return nil
	}

	findOpts := options.Find().SetSkip(skip).SetLimit(limit)
	cursor, err := coll.Find(ctx, filter, findOpts)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	return cursor.All(ctx, results)
}