package repository

import (
	"context"

	"github.com/kongken/go-home/internal/model"
	"butterfly.orx.me/core/store/mongo"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// AlbumRepositoryButterfly butterfly 版本的相册仓库
type AlbumRepositoryButterfly struct {
	collection      *mongodriver.Collection
	photoCollection *mongodriver.Collection
}

// NewAlbumRepositoryButterfly 创建 butterfly 相册仓库
func NewAlbumRepositoryButterfly() AlbumRepository {
	client := mongo.GetClient("primary")
	if client == nil {
		panic("mongo client 'primary' not found")
	}
	db := client.Database("gohome")
	return &AlbumRepositoryButterfly{
		collection:      db.Collection("albums"),
		photoCollection: db.Collection("photos"),
	}
}

func (r *AlbumRepositoryButterfly) Create(ctx context.Context, album *model.Album) error {
	album.BeforeInsert()
	_, err := r.collection.InsertOne(ctx, album)
	return err
}

func (r *AlbumRepositoryButterfly) GetByID(ctx context.Context, id string) (*model.Album, error) {
	var album model.Album
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&album)
	if err != nil {
		return nil, err
	}
	return &album, nil
}

func (r *AlbumRepositoryButterfly) Update(ctx context.Context, album *model.Album) error {
	album.BeforeUpdate()
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": album.ID}, album)
	return err
}

func (r *AlbumRepositoryButterfly) Delete(ctx context.Context, id string) error {
	r.collection.DeleteOne(ctx, bson.M{"_id": id})
	r.photoCollection.DeleteMany(ctx, bson.M{"album_id": id})
	return nil
}

func (r *AlbumRepositoryButterfly) ListByUser(ctx context.Context, userID string, offset, limit int64) ([]*model.Album, int64, error) {
	filter := bson.M{"user_id": userID}
	total, _ := r.collection.CountDocuments(ctx, filter)
	opts := options.Find().SetSkip(offset).SetLimit(limit).SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, _ := r.collection.Find(ctx, filter, opts)
	defer cursor.Close(ctx)
	var albums []*model.Album
	cursor.All(ctx, &albums)
	return albums, total, nil
}

func (r *AlbumRepositoryButterfly) AddPhoto(ctx context.Context, photo *model.Photo) error {
	photo.BeforeInsert()
	_, err := r.photoCollection.InsertOne(ctx, photo)
	return err
}

func (r *AlbumRepositoryButterfly) DeletePhoto(ctx context.Context, id string) error {
	_, err := r.photoCollection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *AlbumRepositoryButterfly) GetPhotosByAlbum(ctx context.Context, albumID string, offset, limit int64) ([]*model.Photo, int64, error) {
	filter := bson.M{"album_id": albumID}
	total, _ := r.photoCollection.CountDocuments(ctx, filter)
	opts := options.Find().SetSkip(offset).SetLimit(limit).SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, _ := r.photoCollection.Find(ctx, filter, opts)
	defer cursor.Close(ctx)
	var photos []*model.Photo
	cursor.All(ctx, &photos)
	return photos, total, nil
}

func (r *AlbumRepositoryButterfly) IncrementPhotoCount(ctx context.Context, albumID string, delta int32) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": albumID}, bson.M{"$inc": bson.M{"photos_count": delta}})
	return err
}