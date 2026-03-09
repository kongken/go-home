package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// Gender 性别
type Gender int32

const (
	GenderUnknown Gender = 0
	GenderMale    Gender = 1
	GenderFemale  Gender = 2
	GenderOther   Gender = 3
)

// UserStatus 用户状态
type UserStatus int32

const (
	UserStatusNormal  UserStatus = 0
	UserStatusFrozen  UserStatus = 1
	UserStatusBanned  UserStatus = 2
)

// User 用户模型
type User struct {
	ID        string     `bson:"_id,omitempty" json:"id"`
	Username  string     `bson:"username" json:"username"`
	Password  string     `bson:"password" json:"-"`
	Nickname  string     `bson:"nickname" json:"nickname"`
	Avatar    string     `bson:"avatar" json:"avatar"`
	Bio       string     `bson:"bio" json:"bio"`
	Gender    Gender     `bson:"gender" json:"gender"`
	Birthday  *time.Time `bson:"birthday,omitempty" json:"birthday,omitempty"`
	Email     string     `bson:"email" json:"email"`
	Phone     string     `bson:"phone" json:"phone"`
	Location  string     `bson:"location" json:"location"`
	Status    UserStatus `bson:"status" json:"status"`
	CreatedAt time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time  `bson:"updated_at" json:"updated_at"`

	// 统计信息
	FriendsCount   int32 `bson:"friends_count" json:"friends_count"`
	FollowersCount int32 `bson:"followers_count" json:"followers_count"`
	FollowingCount int32 `bson:"following_count" json:"following_count"`
	BlogsCount     int32 `bson:"blogs_count" json:"blogs_count"`
	AlbumsCount    int32 `bson:"albums_count" json:"albums_count"`
	SharesCount    int32 `bson:"shares_count" json:"shares_count"`
}

// BeforeInsert 插入前钩子
func (u *User) BeforeInsert(tx interface{}) {
	if u.ID == "" {
		u.ID = bson.NewObjectID().Hex()
	}
	now := time.Now()
	if u.CreatedAt.IsZero() {
		u.CreatedAt = now
	}
	u.UpdatedAt = now
}

// CollectionName 集合名
func (User) CollectionName() string {
	return "users"
}
