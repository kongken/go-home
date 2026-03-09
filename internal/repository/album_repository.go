package repository

import (
	"context"

	"github.com/kongken/go-home/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// AlbumRepository 相册仓库接口
type AlbumRepository interface {
	Create(ctx context.Context, album *model.Album) error
	GetByID(ctx context.Context, id string) (*model.Album, error)
	Update(ctx context.Context, album *model.Album) error
	Delete(ctx context.Context, id string) error
	ListByUser(ctx context.Context, userID string, offset, limit int64) ([]*model.Album, int64, error)

	// 照片管理
	AddPhoto(ctx context.Context, photo *model.Photo) error
	DeletePhoto(ctx context.Context, id string) error
	GetPhotosByAlbum(ctx context.Context, albumID string, offset, limit int64) ([]*model.Photo, int64, error)
	IncrementPhotoCount(ctx context.Context, albumID string, delta int32) error
}

// albumRepository 相册仓库实现
type albumRepository struct {
	collection      *mongo.Collection
	photoCollection *mongo.Collection
}

// NewAlbumRepository 创建相册仓库
func NewAlbumRepository(db *mongo.Database) AlbumRepository {
	return &albumRepository{
		collection:      db.Collection(model.Album{}.CollectionName()),
		photoCollection: db.Collection(model.Photo{}.CollectionName()),
	}
}

// Create 创建相册
func (r *albumRepository) Create(ctx context.Context, album *model.Album) error {
	album.BeforeInsert()
	_, err := r.collection.InsertOne(ctx, album)
	return err
}

// GetByID 根据ID获取相册
func (r *albumRepository) GetByID(ctx context.Context, id string) (*model.Album, error) {
	var album model.Album
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&album)
	if err != nil {
		return nil, err
	}
	return &album, nil
}

// Update 更新相册
func (r *albumRepository) Update(ctx context.Context, album *model.Album) error {
	album.BeforeUpdate()
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": album.ID}, album)
	return err
}

// Delete 删除相册
func (r *albumRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	// 删除所有照片
	_, err = r.photoCollection.DeleteMany(ctx, bson.M{"album_id": id})
	return err
}

// ListByUser 获取用户相册列表
func (r *albumRepository) ListByUser(ctx context.Context, userID string, offset, limit int64) ([]*model.Album, int64, error) {
	filter := bson.M{"user_id": userID}

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

	var albums []*model.Album
	if err := cursor.All(ctx, &albums); err != nil {
		return nil, 0, err
	}

	return albums, total, nil
}

// AddPhoto 添加照片
func (r *albumRepository) AddPhoto(ctx context.Context, photo *model.Photo) error {
	photo.BeforeInsert()
	_, err := r.photoCollection.InsertOne(ctx, photo)
	return err
}

// DeletePhoto 删除照片
func (r *albumRepository) DeletePhoto(ctx context.Context, id string) error {
	_, err := r.photoCollection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// GetPhotosByAlbum 获取相册照片
func (r *albumRepository) GetPhotosByAlbum(ctx context.Context, albumID string, offset, limit int64) ([]*model.Photo, int64, error) {
	filter := bson.M{"album_id": albumID}

	total, err := r.photoCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSkip(offset).
		SetLimit(limit).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.photoCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var photos []*model.Photo
	if err := cursor.All(ctx, &photos); err != nil {
		return nil, 0, err
	}

	return photos, total, nil
}

// IncrementPhotoCount 增加照片数
func (r *albumRepository) IncrementPhotoCount(ctx context.Context, albumID string, delta int32) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": albumID},
		bson.M{"$inc": bson.M{"photos_count": delta}},
	)
	return err
}