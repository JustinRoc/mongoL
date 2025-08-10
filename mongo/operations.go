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

// Repository 仓储接口
type Repository struct {
	client     *Client
	collection *mongo.Collection
}

// NewRepository 创建新的仓储实例
func NewRepository(client *Client, collectionName string) *Repository {
	return &Repository{
		client:     client,
		collection: client.GetCollection(collectionName),
	}
}

// InsertOne 插入单个文档
func (r *Repository) InsertOne(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error) {
	// 如果文档实现了 BaseDocument，调用 BeforeInsert 钩子
	if doc, ok := document.(*BaseDocument); ok {
		doc.BeforeInsert()
	}

	result, err := r.collection.InsertOne(ctx, document)
	if err != nil {
		return nil, fmt.Errorf("failed to insert document: %w", err)
	}
	return result, nil
}

// InsertMany 插入多个文档
func (r *Repository) InsertMany(ctx context.Context, documents []interface{}) (*mongo.InsertManyResult, error) {
	// 为每个文档调用 BeforeInsert 钩子
	for _, doc := range documents {
		if baseDoc, ok := doc.(*BaseDocument); ok {
			baseDoc.BeforeInsert()
		}
	}

	result, err := r.collection.InsertMany(ctx, documents)
	if err != nil {
		return nil, fmt.Errorf("failed to insert documents: %w", err)
	}
	return result, nil
}

// FindOne 查找单个文档
func (r *Repository) FindOne(ctx context.Context, filter bson.M, result interface{}) error {
	err := r.collection.FindOne(ctx, filter).Decode(result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("document not found")
		}
		return fmt.Errorf("failed to find document: %w", err)
	}
	return nil
}

// FindByID 根据ID查找文档
func (r *Repository) FindByID(ctx context.Context, id primitive.ObjectID, result interface{}) error {
	filter := bson.M{"_id": id}
	return r.FindOne(ctx, filter, result)
}

// Find 查找多个文档
func (r *Repository) Find(ctx context.Context, filter bson.M, results interface{}, opts ...*options.FindOptions) error {
	cursor, err := r.collection.Find(ctx, filter, opts...)
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
func (r *Repository) FindWithPagination(ctx context.Context, filter bson.M, page, pageSize int64, results interface{}) (*PaginationResult, error) {
	// 计算跳过的文档数量
	skip := (page - 1) * pageSize

	// 设置查找选项
	findOptions := options.Find().
		SetSkip(skip).
		SetLimit(pageSize)

	// 执行查找
	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %w", err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, results); err != nil {
		return nil, fmt.Errorf("failed to decode documents: %w", err)
	}

	// 计算总数
	total, err := r.collection.CountDocuments(ctx, filter)
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
func (r *Repository) UpdateOne(ctx context.Context, filter bson.M, update bson.M) (*mongo.UpdateResult, error) {
	// 添加更新时间
	if update["$set"] == nil {
		update["$set"] = bson.M{}
	}
	update["$set"].(bson.M)["updated_at"] = time.Now()

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("failed to update document: %w", err)
	}
	return result, nil
}

// UpdateByID 根据ID更新文档
func (r *Repository) UpdateByID(ctx context.Context, id primitive.ObjectID, update bson.M) (*mongo.UpdateResult, error) {
	filter := bson.M{"_id": id}
	return r.UpdateOne(ctx, filter, update)
}

// UpdateMany 更新多个文档
func (r *Repository) UpdateMany(ctx context.Context, filter bson.M, update bson.M) (*mongo.UpdateResult, error) {
	// 添加更新时间
	if update["$set"] == nil {
		update["$set"] = bson.M{}
	}
	update["$set"].(bson.M)["updated_at"] = time.Now()

	result, err := r.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("failed to update documents: %w", err)
	}
	return result, nil
}

// ReplaceOne 替换单个文档
func (r *Repository) ReplaceOne(ctx context.Context, filter bson.M, replacement interface{}) (*mongo.UpdateResult, error) {
	// 如果替换文档实现了 BaseDocument，调用 BeforeUpdate 钩子
	if doc, ok := replacement.(*BaseDocument); ok {
		doc.BeforeUpdate()
	}

	result, err := r.collection.ReplaceOne(ctx, filter, replacement)
	if err != nil {
		return nil, fmt.Errorf("failed to replace document: %w", err)
	}
	return result, nil
}

// DeleteOne 删除单个文档
func (r *Repository) DeleteOne(ctx context.Context, filter bson.M) (*mongo.DeleteResult, error) {
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to delete document: %w", err)
	}
	return result, nil
}

// DeleteByID 根据ID删除文档
func (r *Repository) DeleteByID(ctx context.Context, id primitive.ObjectID) (*mongo.DeleteResult, error) {
	filter := bson.M{"_id": id}
	return r.DeleteOne(ctx, filter)
}

// DeleteMany 删除多个文档
func (r *Repository) DeleteMany(ctx context.Context, filter bson.M) (*mongo.DeleteResult, error) {
	result, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to delete documents: %w", err)
	}
	return result, nil
}

// Count 计算文档数量
func (r *Repository) Count(ctx context.Context, filter bson.M) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count documents: %w", err)
	}
	return count, nil
}

// Exists 检查文档是否存在
func (r *Repository) Exists(ctx context.Context, filter bson.M) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, filter, options.Count().SetLimit(1))
	if err != nil {
		return false, fmt.Errorf("failed to check document existence: %w", err)
	}
	return count > 0, nil
}

// Aggregate 聚合查询
func (r *Repository) Aggregate(ctx context.Context, pipeline []bson.M, results interface{}) error {
	cursor, err := r.collection.Aggregate(ctx, pipeline)
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