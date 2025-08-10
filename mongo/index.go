package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// IndexManager 索引管理器
type IndexManager struct {
	collection *mongo.Collection
}

// NewIndexManager 创建新的索引管理器
func NewIndexManager(client *Client, collectionName string) *IndexManager {
	return &IndexManager{
		collection: client.GetCollection(collectionName),
	}
}

// IndexModel 索引模型
type IndexModel struct {
	Keys    bson.D
	Options *options.IndexOptions
}

// CreateIndex 创建单个索引
func (im *IndexManager) CreateIndex(ctx context.Context, keys bson.D, opts *options.IndexOptions) (string, error) {
	indexModel := mongo.IndexModel{
		Keys:    keys,
		Options: opts,
	}

	name, err := im.collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return "", fmt.Errorf("failed to create index: %w", err)
	}

	log.Printf("Created index: %s", name)
	return name, nil
}

// CreateIndexes 创建多个索引
func (im *IndexManager) CreateIndexes(ctx context.Context, models []IndexModel) ([]string, error) {
	var indexModels []mongo.IndexModel
	for _, model := range models {
		indexModels = append(indexModels, mongo.IndexModel{
			Keys:    model.Keys,
			Options: model.Options,
		})
	}

	names, err := im.collection.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		return nil, fmt.Errorf("failed to create indexes: %w", err)
	}

	log.Printf("Created indexes: %v", names)
	return names, nil
}

// DropIndex 删除索引
func (im *IndexManager) DropIndex(ctx context.Context, name string) error {
	_, err := im.collection.Indexes().DropOne(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to drop index %s: %w", name, err)
	}

	log.Printf("Dropped index: %s", name)
	return nil
}

// DropAllIndexes 删除所有索引（除了_id索引）
func (im *IndexManager) DropAllIndexes(ctx context.Context) error {
	_, err := im.collection.Indexes().DropAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to drop all indexes: %w", err)
	}

	log.Println("Dropped all indexes")
	return nil
}

// ListIndexes 列出所有索引
func (im *IndexManager) ListIndexes(ctx context.Context) ([]bson.M, error) {
	cursor, err := im.collection.Indexes().List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list indexes: %w", err)
	}
	defer cursor.Close(ctx)

	var indexes []bson.M
	if err := cursor.All(ctx, &indexes); err != nil {
		return nil, fmt.Errorf("failed to decode indexes: %w", err)
	}

	return indexes, nil
}

// CreateTextIndex 创建文本索引
func (im *IndexManager) CreateTextIndex(ctx context.Context, fields []string, opts *options.IndexOptions) (string, error) {
	keys := bson.D{}
	for _, field := range fields {
		keys = append(keys, bson.E{Key: field, Value: "text"})
	}

	if opts == nil {
		opts = options.Index()
	}

	return im.CreateIndex(ctx, keys, opts)
}

// CreateCompoundIndex 创建复合索引
func (im *IndexManager) CreateCompoundIndex(ctx context.Context, fields map[string]int, opts *options.IndexOptions) (string, error) {
	keys := bson.D{}
	for field, order := range fields {
		keys = append(keys, bson.E{Key: field, Value: order})
	}

	return im.CreateIndex(ctx, keys, opts)
}

// CreateUniqueIndex 创建唯一索引
func (im *IndexManager) CreateUniqueIndex(ctx context.Context, field string, opts *options.IndexOptions) (string, error) {
	keys := bson.D{{Key: field, Value: 1}}

	if opts == nil {
		opts = options.Index()
	}
	opts.SetUnique(true)

	return im.CreateIndex(ctx, keys, opts)
}

// CreateSparseIndex 创建稀疏索引
func (im *IndexManager) CreateSparseIndex(ctx context.Context, field string, opts *options.IndexOptions) (string, error) {
	keys := bson.D{{Key: field, Value: 1}}

	if opts == nil {
		opts = options.Index()
	}
	opts.SetSparse(true)

	return im.CreateIndex(ctx, keys, opts)
}

// CreateTTLIndex 创建TTL索引（生存时间索引）
func (im *IndexManager) CreateTTLIndex(ctx context.Context, field string, expireAfter time.Duration, opts *options.IndexOptions) (string, error) {
	keys := bson.D{{Key: field, Value: 1}}

	if opts == nil {
		opts = options.Index()
	}
	opts.SetExpireAfterSeconds(int32(expireAfter.Seconds()))

	return im.CreateIndex(ctx, keys, opts)
}

// CreatePartialIndex 创建部分索引
func (im *IndexManager) CreatePartialIndex(ctx context.Context, field string, filter bson.M, opts *options.IndexOptions) (string, error) {
	keys := bson.D{{Key: field, Value: 1}}

	if opts == nil {
		opts = options.Index()
	}
	opts.SetPartialFilterExpression(filter)

	return im.CreateIndex(ctx, keys, opts)
}

// IndexExists 检查索引是否存在
func (im *IndexManager) IndexExists(ctx context.Context, name string) (bool, error) {
	indexes, err := im.ListIndexes(ctx)
	if err != nil {
		return false, err
	}

	for _, index := range indexes {
		if indexName, ok := index["name"].(string); ok && indexName == name {
			return true, nil
		}
	}

	return false, nil
}

// GetIndexStats 获取索引统计信息
func (im *IndexManager) GetIndexStats(ctx context.Context) ([]bson.M, error) {
	pipeline := []bson.M{
		{"$indexStats": bson.M{}},
	}

	cursor, err := im.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to get index stats: %w", err)
	}
	defer cursor.Close(ctx)

	var stats []bson.M
	if err := cursor.All(ctx, &stats); err != nil {
		return nil, fmt.Errorf("failed to decode index stats: %w", err)
	}

	return stats, nil
}

// CommonIndexes 常用索引创建函数
type CommonIndexes struct {
	indexManager *IndexManager
}

// NewCommonIndexes 创建常用索引管理器
func NewCommonIndexes(client *Client, collectionName string) *CommonIndexes {
	return &CommonIndexes{
		indexManager: NewIndexManager(client, collectionName),
	}
}

// CreateUserIndexes 为用户集合创建常用索引
func (ci *CommonIndexes) CreateUserIndexes(ctx context.Context) error {
	indexes := []IndexModel{
		{
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_username"),
		},
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_email"),
		},
		{
			Keys:    bson.D{{Key: "status", Value: 1}},
			Options: options.Index().SetName("idx_status"),
		},
		{
			Keys:    bson.D{{Key: "created_at", Value: -1}},
			Options: options.Index().SetName("idx_created_at"),
		},
	}

	_, err := ci.indexManager.CreateIndexes(ctx, indexes)
	return err
}

// CreateArticleIndexes 为文章集合创建常用索引
func (ci *CommonIndexes) CreateArticleIndexes(ctx context.Context) error {
	indexes := []IndexModel{
		{
			Keys:    bson.D{{Key: "title", Value: "text"}, {Key: "content", Value: "text"}},
			Options: options.Index().SetName("idx_text_search"),
		},
		{
			Keys:    bson.D{{Key: "author_id", Value: 1}},
			Options: options.Index().SetName("idx_author_id"),
		},
		{
			Keys:    bson.D{{Key: "status", Value: 1}},
			Options: options.Index().SetName("idx_status"),
		},
		{
			Keys:    bson.D{{Key: "tags", Value: 1}},
			Options: options.Index().SetName("idx_tags"),
		},
		{
			Keys:    bson.D{{Key: "created_at", Value: -1}},
			Options: options.Index().SetName("idx_created_at"),
		},
		{
			Keys:    bson.D{{Key: "category_id", Value: 1}, {Key: "status", Value: 1}},
			Options: options.Index().SetName("idx_category_status"),
		},
	}

	_, err := ci.indexManager.CreateIndexes(ctx, indexes)
	return err
}