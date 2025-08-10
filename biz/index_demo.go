package biz

import (
	"context"
	"fmt"
	"log"

	"github.com/JustinRoc/mongodbL/mongo"
)

// InitializeIndexes 初始化所有文档索引
func InitializeIndexes(client *mongo.Client) error {
	ctx := context.Background()
	
	// 创建文档索引配置实例
	docIndexes := mongo.NewDocumentIndexes(client)
	
	log.Println("开始创建数据库索引...")
	
	// 创建所有文档索引
	if err := docIndexes.CreateAllDocumentIndexes(ctx); err != nil {
		return fmt.Errorf("创建文档索引失败: %w", err)
	}
	
	log.Println("所有索引创建完成")
	return nil
}

// InitializeUserIndexes 只初始化用户相关索引
func InitializeUserIndexes(client *mongo.Client) error {
	ctx := context.Background()
	
	docIndexes := mongo.NewDocumentIndexes(client)
	
	log.Println("创建用户索引...")
	if err := docIndexes.CreateUserIndexes(ctx); err != nil {
		return fmt.Errorf("创建用户索引失败: %w", err)
	}
	
	log.Println("用户索引创建完成")
	return nil
}

// CheckIndexUsage 检查索引使用情况
func CheckIndexUsage(client *mongo.Client) error {
	ctx := context.Background()
	
	docIndexes := mongo.NewDocumentIndexes(client)
	
	// 检查用户集合的索引使用统计
	stats, err := docIndexes.GetIndexUsageStats(ctx, "users")
	if err != nil {
		return fmt.Errorf("获取索引统计失败: %w", err)
	}
	
	log.Println("用户集合索引使用统计:")
	for _, stat := range stats {
		if name, ok := stat["name"].(string); ok {
			if ops, ok := stat["accesses"].(map[string]interface{}); ok {
				if opsCount, ok := ops["ops"].(int64); ok {
					log.Printf("索引 %s 使用次数: %d", name, opsCount)
				}
			}
		}
	}
	
	return nil
}

// DemoIndexQueries 演示使用索引的查询
func DemoIndexQueries(client *mongo.Client) error {
	ctx := context.Background()
	userCollection := mongo.NewCollection(client, "users")
	
	log.Println("演示索引查询...")
	
	// 1. 使用用户名索引查询（唯一索引）
	log.Println("1. 按用户名查询（使用唯一索引）")
	var user mongo.User
	err := userCollection.FindOne(ctx, map[string]interface{}{
		"username": "john_doe",
	}, &user)
	if err != nil {
		log.Printf("查询失败: %v", err)
	}
	
	// 2. 使用邮箱索引查询（唯一索引）
	log.Println("2. 按邮箱查询（使用唯一索引）")
	err = userCollection.FindOne(ctx, map[string]interface{}{
		"email": "john@example.com",
	}, &user)
	if err != nil {
		log.Printf("查询失败: %v", err)
	}
	
	// 3. 使用状态索引查询
	log.Println("3. 按状态查询（使用单字段索引）")
	var users []mongo.User
	err = userCollection.Find(ctx, map[string]interface{}{
		"status": "active",
	}, &users)
	if err != nil {
		log.Printf("查询失败: %v", err)
	}
	
	// 4. 使用复合索引查询（状态+创建时间）
	log.Println("4. 按状态和创建时间查询（使用复合索引）")
	err = userCollection.Find(ctx, map[string]interface{}{
		"status": "active",
	}, &users)
	if err != nil {
		log.Printf("查询失败: %v", err)
	}
	
	// 5. 使用嵌套字段索引查询
	log.Println("5. 按用户姓名查询（使用嵌套字段复合索引）")
	err = userCollection.Find(ctx, map[string]interface{}{
		"profile.first_name": "John",
		"profile.last_name":  "Doe",
	}, &users)
	if err != nil {
		log.Printf("查询失败: %v", err)
	}
	
	log.Println("索引查询演示完成")
	return nil
}

// DemoArticleQueries 演示文章相关的索引查询
func DemoArticleQueries(client *mongo.Client) error {
	ctx := context.Background()
	articleCollection := mongo.NewCollection(client, "articles")
	
	log.Println("演示文章索引查询...")
	
	// 1. 按作者查询文章（使用作者ID索引）
	log.Println("1. 按作者查询文章")
	var articles []mongo.Article
	err := articleCollection.Find(ctx, map[string]interface{}{
		"author_id": "507f1f77bcf86cd799439011", // 示例ObjectID
	}, &articles)
	if err != nil {
		log.Printf("查询失败: %v", err)
	}
	
	// 2. 按状态查询文章（使用状态索引）
	log.Println("2. 按状态查询文章")
	err = articleCollection.Find(ctx, map[string]interface{}{
		"status": "published",
	}, &articles)
	if err != nil {
		log.Printf("查询失败: %v", err)
	}
	
	// 3. 按分类查询文章（使用分类ID索引）
	log.Println("3. 按分类查询文章")
	err = articleCollection.Find(ctx, map[string]interface{}{
		"category_id": "507f1f77bcf86cd799439012", // 示例ObjectID
	}, &articles)
	if err != nil {
		log.Printf("查询失败: %v", err)
	}
	
	// 4. 按标签查询文章（使用标签索引）
	log.Println("4. 按标签查询文章")
	err = articleCollection.Find(ctx, map[string]interface{}{
		"tags": "golang",
	}, &articles)
	if err != nil {
		log.Printf("查询失败: %v", err)
	}
	
	// 5. 使用复合索引查询（分类+状态+时间排序）
	log.Println("5. 分类页面文章列表（使用复合索引）")
	err = articleCollection.Find(ctx, map[string]interface{}{
		"category_id": "507f1f77bcf86cd799439012",
		"status":      "published",
	}, &articles)
	if err != nil {
		log.Printf("查询失败: %v", err)
	}
	
	log.Println("文章索引查询演示完成")
	return nil
}

// OptimizeIndexes 索引优化建议
func OptimizeIndexes(client *mongo.Client) {
	log.Println("索引优化建议:")
	log.Println("1. 定期检查索引使用统计，删除未使用的索引")
	log.Println("2. 根据实际查询模式调整复合索引的字段顺序")
	log.Println("3. 对于大集合，考虑使用部分索引减少索引大小")
	log.Println("4. 监控索引对写入性能的影响")
	log.Println("5. 使用 explain() 分析查询计划，确保索引被正确使用")
}