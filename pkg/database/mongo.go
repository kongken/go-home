package database

import (
	"context"
	"fmt"
	"time"

	"github.com/kongken/go-home/internal/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// MongoDB MongoDB客户端
type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

// NewMongoDB 创建 MongoDB 连接
func NewMongoDB(cfg *config.DatabaseConfig) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?ssl=false",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	
	if cfg.User == "" && cfg.Password == "" {
		uri = fmt.Sprintf("mongodb://%s:%d/%s", cfg.Host, cfg.Port, cfg.DBName)
	}

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	db := client.Database(cfg.DBName)

	return &MongoDB{
		Client:   client,
		Database: db,
	}, nil
}

// Close 关闭连接
func (m *MongoDB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return m.Client.Disconnect(ctx)
}

// Collection 获取集合
func (m *MongoDB) Collection(name string) *mongo.Collection {
	return m.Database.Collection(name)
}
