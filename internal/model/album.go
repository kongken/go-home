package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// Album 相册模型
type Album struct {
	ID          string       `bson:"_id,omitempty" json:"id"`
	UserID      string       `bson:"user_id" json:"user_id"`
	Name        string       `bson:"name" json:"name"`
	Description string       `bson:"description" json:"description"`
	CoverPhoto  string       `bson:"cover_photo" json:"cover_photo"`
	Privacy     PrivacyLevel `bson:"privacy" json:"privacy"`
	CreatedAt   time.Time    `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time    `bson:"updated_at" json:"updated_at"`

	// 统计
	PhotosCount  int32 `bson:"photos_count" json:"photos_count"`
	ViewsCount   int32 `bson:"views_count" json:"views_count"`
	CommentsCount int32 `bson:"comments_count" json:"comments_count"`
}

// BeforeInsert 插入前钩子
func (a *Album) BeforeInsert() {
	if a.ID == "" {
		a.ID = bson.NewObjectID().Hex()
	}
	now := time.Now()
	if a.CreatedAt.IsZero() {
		a.CreatedAt = now
	}
	a.UpdatedAt = now
}

// BeforeUpdate 更新前钩子
func (a *Album) BeforeUpdate() {
	a.UpdatedAt = time.Now()
}

// CollectionName 集合名
func (Album) CollectionName() string {
	return "albums"
}

// Photo 照片模型
type Photo struct {
	ID          string    `bson:"_id,omitempty" json:"id"`
	AlbumID     string    `bson:"album_id" json:"album_id"`
	UserID      string    `bson:"user_id" json:"user_id"`
	URL         string    `bson:"url" json:"url"`
	Thumbnail   string    `bson:"thumbnail" json:"thumbnail"`
	Description string    `bson:"description" json:"description"`
	Width       int32     `bson:"width" json:"width"`
	Height      int32     `bson:"height" json:"height"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
}

// BeforeInsert 插入前钩子
func (p *Photo) BeforeInsert() {
	if p.ID == "" {
		p.ID = bson.NewObjectID().Hex()
	}
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now()
	}
}

// CollectionName 集合名
func (Photo) CollectionName() string {
	return "photos"
}
