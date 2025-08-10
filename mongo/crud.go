package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Collection 集合操作
type Collection struct {
	cli        *Client
	collection *mongo.Collection
}

// NewCollection 创建新的集合实例
func NewCollection(client *Client, collectionName string) *Collection {
	return &Collection{
		cli:        client,
		collection: client.GetCollection(collectionName),
	}
}

// InsertOne 插入单个文档
func (c *Collection) InsertOne(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error) {
	if doc, ok := document.(Document); ok {
		doc.BeforeInsert()
	}

	result, err := c.collection.InsertOne(ctx, document)
	if err != nil {
		return nil, fmt.Errorf("failed to insert document: %w", err)
	}

	// 将生成的 ID 写入到 document 中
	if doc, ok := document.(Document); ok {
		if insertedID, ok := result.InsertedID.(primitive.ObjectID); ok {
			doc.SetID(insertedID)
		} else {
			return nil, fmt.Errorf("insertedID is not ObjectID")
		}
	}
	return result, nil
}

// InsertMany 插入多个文档
func (c *Collection) InsertMany(ctx context.Context, documents []interface{}) (*mongo.InsertManyResult, error) {
	// 为每个文档调用 BeforeInsert 钩子
	for _, doc := range documents {
		if baseDoc, ok := doc.(*BaseDocument); ok {
			baseDoc.BeforeInsert()
		}
	}

	result, err := c.collection.InsertMany(ctx, documents)
	if err != nil {
		return nil, fmt.Errorf("failed to insert documents: %w", err)
	}
	return result, nil
}

// FindOne 查找单个文档
func (c *Collection) FindOne(ctx context.Context, filter bson.M, result interface{}) error {
	err := c.collection.FindOne(ctx, filter).Decode(result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("document not found")
		}
		return fmt.Errorf("failed to find document: %w", err)
	}
	return nil
}

// FindByID 根据ID查找文档
func (c *Collection) FindByID(ctx context.Context, id primitive.ObjectID, result interface{}) error {
	filter := bson.M{"_id": id}
	return c.FindOne(ctx, filter, result)
}

// Find 查找多个文档
func (c *Collection) Find(ctx context.Context, filter bson.M, results interface{}, opts ...*options.FindOptions) error {
	cursor, err := c.collection.Find(ctx, filter, opts...)
	if err != nil {
		return fmt.Errorf("failed to find documents: %w", err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, results); err != nil {
		return fmt.Errorf("failed to decode documents: %w", err)
	}
	return nil
}

// FindWithPagination 分页查找文档
func (c *Collection) FindWithPagination(ctx context.Context, filter bson.M, page, pageSize int64, results interface{}) (*PaginationResult, error) {
	// 计算跳过的文档数量
	skip := (page - 1) * pageSize

	// 设置查找选项
	findOptions := options.Find().
		SetSkip(skip).
		SetLimit(pageSize)

	// 执行查找
	cursor, err := c.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %w", err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, results); err != nil {
		return nil, fmt.Errorf("failed to decode documents: %w", err)
	}

	// 计算总数
	total, err := c.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to count documents: %w", err)
	}

	return &PaginationResult{
		Page:      page,
		PageSize:  pageSize,
		Total:     total,
		TotalPage: (total + pageSize - 1) / pageSize,
	}, nil
}

// UpdateOne 更新单个文档
func (c *Collection) UpdateOne(ctx context.Context, filter bson.M, update bson.M, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	// 添加更新时间
	if update["$set"] == nil {
		update["$set"] = bson.M{}
	}
	update["$set"].(bson.M)["updated_at"] = time.Now()

	result, err := c.collection.UpdateOne(ctx, filter, update, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to update document: %w", err)
	}
	return result, nil
}

// UpdateByID 根据ID更新文档
func (c *Collection) UpdateByID(ctx context.Context, id primitive.ObjectID, update bson.M, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	filter := bson.M{"_id": id}
	return c.UpdateOne(ctx, filter, update, opts...)
}

// UpdateMany 更新多个文档
func (c *Collection) UpdateMany(ctx context.Context, filter bson.M, update bson.M) (*mongo.UpdateResult, error) {
	// 添加更新时间
	if update["$set"] == nil {
		update["$set"] = bson.M{}
	}
	update["$set"].(bson.M)["updated_at"] = time.Now()

	result, err := c.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("failed to update documents: %w", err)
	}
	return result, nil
}

// ReplaceOne 替换单个文档
func (c *Collection) ReplaceOne(ctx context.Context, filter bson.M, replacement interface{}) (*mongo.UpdateResult, error) {
	// 如果替换文档实现了 BaseDocument，调用 BeforeUpdate 钩子
	if doc, ok := replacement.(*BaseDocument); ok {
		doc.BeforeUpdate()
	}

	result, err := c.collection.ReplaceOne(ctx, filter, replacement)
	if err != nil {
		return nil, fmt.Errorf("failed to replace document: %w", err)
	}
	return result, nil
}

// DeleteOne 删除单个文档
func (c *Collection) DeleteOne(ctx context.Context, filter bson.M) (*mongo.DeleteResult, error) {
	result, err := c.collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to delete document: %w", err)
	}
	return result, nil
}

// DeleteByID 根据ID删除文档
func (c *Collection) DeleteByID(ctx context.Context, id primitive.ObjectID) (*mongo.DeleteResult, error) {
	filter := bson.M{"_id": id}
	return c.DeleteOne(ctx, filter)
}

// DeleteMany 删除多个文档
func (c *Collection) DeleteMany(ctx context.Context, filter bson.M) (*mongo.DeleteResult, error) {
	result, err := c.collection.DeleteMany(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to delete documents: %w", err)
	}
	return result, nil
}

// Count 计算文档数量
func (c *Collection) Count(ctx context.Context, filter bson.M) (int64, error) {
	count, err := c.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count documents: %w", err)
	}
	return count, nil
}

// Exists 检查文档是否存在
func (c *Collection) Exists(ctx context.Context, filter bson.M) (bool, error) {
	count, err := c.collection.CountDocuments(ctx, filter, options.Count().SetLimit(1))
	if err != nil {
		return false, fmt.Errorf("failed to check document existence: %w", err)
	}
	return count > 0, nil
}

// Aggregate 聚合查询
func (c *Collection) Aggregate(ctx context.Context, pipeline []bson.M, results interface{}) error {
	cursor, err := c.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return fmt.Errorf("failed to aggregate: %w", err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, results); err != nil {
		return fmt.Errorf("failed to decode aggregation results: %w", err)
	}
	return nil
}

// PaginationResult 分页结果
type PaginationResult struct {
	Page      int64 `json:"page"`
	PageSize  int64 `json:"page_size"`
	Total     int64 `json:"total"`
	TotalPage int64 `json:"total_page"`
}
