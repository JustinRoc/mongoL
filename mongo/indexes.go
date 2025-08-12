package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DocumentIndexes 文档索引配置
type DocumentIndexes struct {
	client *Client
}

// NewDocumentIndexes 创建文档索引配置实例
func NewDocumentIndexes(client *Client) *DocumentIndexes {
	return &DocumentIndexes{
		client: client,
	}
}

// CreateUserIndexes 为 User 集合创建索引
func (di *DocumentIndexes) CreateUserIndexes(ctx context.Context) error {
	indexManager := NewIndexManager(di.client, "users")
	
	indexes := []mongo.IndexModel{
		// 1. 用户名唯一索引 - 用于登录和用户查找
		{
			Keys:    bson.D{{"username", 1}},
			Options: options.Index().SetUnique(true).SetName("idx_username_unique"),
		},
		// 2. 邮箱唯一索引 - 用于登录和用户查找
		{
			Keys:    bson.D{{"email", 1}},
			Options: options.Index().SetUnique(true).SetName("idx_email_unique"),
		},
		// 3. 状态索引 - 用于查询活跃用户等
		{
			Keys:    bson.D{{"status", 1}},
			Options: options.Index().SetName("idx_status"),
		},
		// 4. 创建时间索引 - 用于按时间排序和范围查询
		{
			Keys:    bson.D{{"created_at", -1}},
			Options: options.Index().SetName("idx_created_at_desc"),
		},
		// 5. 更新时间索引 - 用于查找最近更新的用户
		{
			Keys:    bson.D{{"updated_at", -1}},
			Options: options.Index().SetName("idx_updated_at_desc"),
		},
		// 6. 用户姓名复合索引 - 用于按姓名搜索
		{
			Keys: bson.D{
				{"profile.first_name", 1},
				{"profile.last_name", 1},
			},
			Options: options.Index().SetName("idx_profile_name"),
		},
		// 7. 状态+创建时间复合索引 - 用于查询特定状态的用户并按时间排序
		{
			Keys: bson.D{
				{"status", 1},
				{"created_at", -1},
			},
			Options: options.Index().SetName("idx_status_created_at"),
		},
		// 8. 用户简介文本索引 - 用于全文搜索
		{
			Keys:    bson.D{{"profile.bio", "text"}},
			Options: options.Index().SetName("idx_profile_bio_text"),
		},
	}
	
	_, err := indexManager.CreateIndexes(ctx, indexes)
	return err
}

// CreateArticleIndexes 为 Article 集合创建索引
func (di *DocumentIndexes) CreateArticleIndexes(ctx context.Context) error {
	indexManager := NewIndexManager(di.client, "articles")
	
	indexes := []mongo.IndexModel{
		// 1. 作者ID索引 - 用于查询某个用户的所有文章
		{
			Keys:    bson.D{{"author_id", 1}},
			Options: options.Index().SetName("idx_author_id"),
		},
		// 2. 状态索引 - 用于查询已发布、草稿等状态的文章
		{
			Keys:    bson.D{{"status", 1}},
			Options: options.Index().SetName("idx_status"),
		},
		// 3. 分类ID索引 - 用于按分类查询文章
		{
			Keys:    bson.D{{"category_id", 1}},
			Options: options.Index().SetName("idx_category_id"),
		},
		// 4. 创建时间索引 - 用于按发布时间排序
		{
			Keys:    bson.D{{"created_at", -1}},
			Options: options.Index().SetName("idx_created_at_desc"),
		},
		// 5. 浏览量索引 - 用于热门文章排序
		{
			Keys:    bson.D{{"view_count", -1}},
			Options: options.Index().SetName("idx_view_count_desc"),
		},
		// 6. 点赞数索引 - 用于按点赞数排序
		{
			Keys:    bson.D{{"like_count", -1}},
			Options: options.Index().SetName("idx_like_count_desc"),
		},
		// 7. 标签索引 - 用于按标签查询文章
		{
			Keys:    bson.D{{"tags", 1}},
			Options: options.Index().SetName("idx_tags"),
		},
		// 8. 状态+创建时间复合索引 - 用于查询已发布文章并按时间排序
		{
			Keys: bson.D{
				{"status", 1},
				{"created_at", -1},
			},
			Options: options.Index().SetName("idx_status_created_at"),
		},
		// 9. 作者+状态复合索引 - 用于查询某作者的特定状态文章
		{
			Keys: bson.D{
				{"author_id", 1},
				{"status", 1},
			},
			Options: options.Index().SetName("idx_author_status"),
		},
		// 10. 分类+状态+创建时间复合索引 - 用于分类页面的文章列表
		{
			Keys: bson.D{
				{"category_id", 1},
				{"status", 1},
				{"created_at", -1},
			},
			Options: options.Index().SetName("idx_category_status_created_at"),
		},
		// 11. 标题和内容文本索引 - 用于全文搜索
		{
			Keys: bson.D{
				{"title", "text"},
				{"content", "text"},
			},
			Options: options.Index().SetName("idx_title_content_text"),
		},
	}
	
	_, err := indexManager.CreateIndexes(ctx, indexes)
	return err
}

// CreateCategoryIndexes 为 Category 集合创建索引
func (di *DocumentIndexes) CreateCategoryIndexes(ctx context.Context) error {
	indexManager := NewIndexManager(di.client, "categories")
	
	indexes := []mongo.IndexModel{
		// 1. 分类名称唯一索引 - 确保分类名称不重复
		{
			Keys:    bson.D{{"name", 1}},
			Options: options.Index().SetUnique(true).SetName("idx_name_unique"),
		},
		// 2. 父分类ID索引 - 用于查询子分类
		{
			Keys:    bson.D{{"parent_id", 1}},
			Options: options.Index().SetSparse(true).SetName("idx_parent_id"), // 稀疏索引，因为根分类没有parent_id
		},
		// 3. 激活状态索引 - 用于查询激活的分类
		{
			Keys:    bson.D{{"is_active", 1}},
			Options: options.Index().SetName("idx_is_active"),
		},
		// 4. 排序索引 - 用于分类排序显示
		{
			Keys:    bson.D{{"sort", 1}},
			Options: options.Index().SetName("idx_sort"),
		},
		// 5. 父分类+排序复合索引 - 用于获取某个父分类下的子分类并排序
		{
			Keys: bson.D{
				{"parent_id", 1},
				{"sort", 1},
			},
			Options: options.Index().SetName("idx_parent_sort"),
		},
		// 6. 激活状态+排序复合索引 - 用于获取激活的分类并排序
		{
			Keys: bson.D{
				{"is_active", 1},
				{"sort", 1},
			},
			Options: options.Index().SetName("idx_active_sort"),
		},
		// 7. 分类描述文本索引 - 用于分类搜索
		{
			Keys:    bson.D{{"description", "text"}},
			Options: options.Index().SetName("idx_description_text"),
		},
	}
	
	_, err := indexManager.CreateIndexes(ctx, indexes)
	return err
}

// CreateAllDocumentIndexes 为所有文档类型创建索引
func (di *DocumentIndexes) CreateAllDocumentIndexes(ctx context.Context) error {
	// 创建用户索引
	if err := di.CreateUserIndexes(ctx); err != nil {
		return err
	}
	
	// 创建文章索引
	if err := di.CreateArticleIndexes(ctx); err != nil {
		return err
	}
	
	// 创建分类索引
	if err := di.CreateCategoryIndexes(ctx); err != nil {
		return err
	}
	
	return nil
}

// CreateBaseDocumentIndexes 为所有继承BaseDocument的集合创建基础索引
func (di *DocumentIndexes) CreateBaseDocumentIndexes(ctx context.Context, collectionName string) error {
	indexManager := NewIndexManager(di.client, collectionName)
	
	indexes := []mongo.IndexModel{
		// 创建时间索引
		{
			Keys:    bson.D{{"created_at", -1}},
			Options: options.Index().SetName("idx_created_at_desc"),
		},
		// 更新时间索引
		{
			Keys:    bson.D{{"updated_at", -1}},
			Options: options.Index().SetName("idx_updated_at_desc"),
		},
	}
	
	_, err := indexManager.CreateIndexes(ctx, indexes)
	return err
}

// DropAllDocumentIndexes 删除所有文档索引（谨慎使用）
func (di *DocumentIndexes) DropAllDocumentIndexes(ctx context.Context) error {
	collections := []string{"users", "articles", "categories"}
	
	for _, collectionName := range collections {
		indexManager := NewIndexManager(di.client, collectionName)
		if err := indexManager.DropAllIndexes(ctx); err != nil {
			return err
		}
	}
	
	return nil
}

// GetIndexUsageStats 获取索引使用统计
func (di *DocumentIndexes) GetIndexUsageStats(ctx context.Context, collectionName string) ([]bson.M, error) {
	indexManager := NewIndexManager(di.client, collectionName)
	return indexManager.GetIndexStats(ctx)
}