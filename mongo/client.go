package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Client MongoDB 客户端封装
type Client struct {
	client   *mongo.Client
	database *mongo.Database
	dbName   string
}

// Config MongoDB 连接配置
type Config struct {
	URI            string        `json:"uri"`
	Database       string        `json:"database"`
	ConnectTimeout time.Duration `json:"connect_timeout"`
	MaxPoolSize    uint64        `json:"max_pool_size"`
	MinPoolSize    uint64        `json:"min_pool_size"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		URI:            "mongodb://localhost:27017",
		Database:       "test",
		ConnectTimeout: 10 * time.Second,
		MaxPoolSize:    100,
		MinPoolSize:    5,
	}
}

// NewClient 创建新的 MongoDB 客户端
func NewClient(config *Config) (*Client, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// 设置客户端选项
	clientOptions := options.Client().
		ApplyURI(config.URI).
		SetConnectTimeout(config.ConnectTimeout).
		SetMaxPoolSize(config.MaxPoolSize).
		SetMinPoolSize(config.MinPoolSize)

	// 连接到 MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Printf("Successfully connected to MongoDB: %s", config.URI)

	return &Client{
		client:   client,
		database: client.Database(config.Database),
		dbName:   config.Database,
	}, nil
}

// GetDatabase 获取数据库实例
func (c *Client) GetDatabase() *mongo.Database {
	return c.database
}

// GetCollection 获取集合实例
func (c *Client) GetCollection(name string) *mongo.Collection {
	return c.database.Collection(name)
}

// Close 关闭客户端连接
func (c *Client) Close() error {
	if c.client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return c.client.Disconnect(ctx)
	}
	return nil
}

// Ping 测试连接
func (c *Client) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return c.client.Ping(ctx, readpref.Primary())
}

// GetDatabaseName 获取数据库名称
func (c *Client) GetDatabaseName() string {
	return c.dbName
}