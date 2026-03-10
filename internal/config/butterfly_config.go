package config

import (
	"fmt"
)

// ButterflyConfig butterfly 框架配置
type ButterflyConfig struct {
	Service string          `yaml:"service"`
	Store   StoreConfig     `yaml:"store"`
	OTel    OTelConfig      `yaml:"otel"`
	JWT     JWTConfig       `yaml:"jwt"`
	Server  ServerConfig    `yaml:"server"`
}

// StoreConfig 存储配置
type StoreConfig struct {
	Mongo MongoConfig `yaml:"mongo"`
	Redis RedisConfig `yaml:"redis"`
	DB    DBConfig    `yaml:"db"`
}

// MongoConfig MongoDB 配置
type MongoConfig struct {
	Primary   MongoConnection `yaml:"primary"`
	Secondary MongoConnection `yaml:"secondary"`
}

// MongoConnection MongoDB 连接配置
type MongoConnection struct {
	URI string `yaml:"uri"`
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Cache   RedisConnection `yaml:"cache"`
	Session RedisConnection `yaml:"session"`
}

// RedisConnection Redis 连接配置
type RedisConnection struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// DBConfig 关系型数据库配置
type DBConfig struct {
	Main DBConnection `yaml:"main"`
}

// DBConnection 数据库连接配置
type DBConnection struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name"`
}

// OTelConfig OpenTelemetry 配置
type OTelConfig struct {
	Tracing TracingConfig `yaml:"tracing"`
}

// TracingConfig 链路追踪配置
type TracingConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Endpoint string `yaml:"endpoint"`
	Provider string `yaml:"provider"`
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret        string `yaml:"secret"`
	AccessExpiry  int64  `yaml:"access_expiry"`
	RefreshExpiry int64  `yaml:"refresh_expiry"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

// Print 打印配置信息
func (c *ButterflyConfig) Print() {
	fmt.Printf("Service: %s\n", c.Service)
	fmt.Printf("Server: %s:%d (mode: %s)\n", c.Server.Host, c.Server.Port, c.Server.Mode)
	fmt.Printf("MongoDB: %s\n", c.Store.Mongo.Primary.URI)
	fmt.Printf("Redis Cache: %s\n", c.Store.Redis.Cache.Addr)
	fmt.Printf("Tracing: enabled=%v, endpoint=%s\n", c.OTel.Tracing.Enabled, c.OTel.Tracing.Endpoint)
}

// GetDSN 获取数据库 DSN
func (c *DBConnection) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.DBName)
}