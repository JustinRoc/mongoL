package biz

import (
	"context"
	"log"

	"github.com/JustinRoc/mongodbL/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CompoundIndexExample 演示正确使用复合索引的方法
func CompoundIndexExample(client *mongo.Client) error {
	ctx := context.Background()
	indexManager := mongo.NewIndexManager(client, "articles")

	log.Println("演示复合索引的正确创建方法...")

	// ✅ 正确方式：使用 bson.D 保证字段顺序
	log.Println("1. 创建分类+状态+时间复合索引（字段顺序很重要）")
	_, err := indexManager.CreateIndex(ctx, bson.D{
		{"category_id", 1},   // 第一个字段：分类ID
		{"status", 1},        // 第二个字段：状态
		{"created_at", -1},   // 第三个字段：创建时间（降序）
	}, options.Index().SetName("idx_category_status_created_at"))
	if err != nil {
		log.Printf("创建复合索引失败: %v", err)
	} else {
		log.Println("✅ 复合索引创建成功")
	}

	// ✅ 另一个复合索引示例：作者+状态
	log.Println("2. 创建作者+状态复合索引")
	_, err = indexManager.CreateIndex(ctx, bson.D{
		{"author_id", 1},
		{"status", 1},
	}, options.Index().SetName("idx_author_status"))
	if err != nil {
		log.Printf("创建复合索引失败: %v", err)
	} else {
		log.Println("✅ 作者+状态复合索引创建成功")
	}
	return nil
}

// QueryWithCompoundIndex 演示如何正确使用复合索引进行查询
func QueryWithCompoundIndex(client *mongo.Client) error {
	ctx := context.Background()
	articleCollection := mongo.NewCollection(client, "articles")

	log.Println("演示复合索引查询的最佳实践...")

	// 假设我们有索引：{category_id: 1, status: 1, created_at: -1}
	
	// 重要说明：
	// 1. 创建索引时使用 bson.D 保证字段顺序
	// 2. 查询时使用 bson.M，MongoDB 会自动优化查询计划
	// 3. 复合索引的效率遵循"最左前缀"原则

	// ✅ 高效查询：遵循最左前缀原则
	log.Println("1. 使用第一个字段查询（高效）")
	var articles []mongo.Article
	err := articleCollection.Find(ctx, bson.M{
		"category_id": "507f1f77bcf86cd799439011",
	}, &articles)
	if err != nil {
		log.Printf("查询失败: %v", err)
	}

	log.Println("2. 使用前两个字段查询（高效）")
	err = articleCollection.Find(ctx, bson.M{
		"category_id": "507f1f77bcf86cd799439011",
		"status":      "published",
	}, &articles)
	if err != nil {
		log.Printf("查询失败: %v", err)
	}

	log.Println("3. 使用所有字段查询（最高效）")
	err = articleCollection.Find(ctx, bson.M{
		"category_id": "507f1f77bcf86cd799439011",
		"status":      "published",
		"created_at":  bson.M{"$gte": "2024-01-01"},
	}, &articles)
	if err != nil {
		log.Printf("查询失败: %v", err)
	}

	// ❌ 低效查询：跳过了前面的字段
	log.Println("4. 低效查询示例（跳过第一个字段）")
	err = articleCollection.Find(ctx, bson.M{
		"status": "published", // 跳过了 category_id
	}, &articles)
	if err != nil {
		log.Printf("查询失败: %v", err)
	}

	log.Println("复合索引查询演示完成")
	return nil
}

// IndexOrderComparison 对比不同字段顺序的索引效果
func IndexOrderComparison(client *mongo.Client) {
	log.Println("复合索引字段顺序的重要性:")
	log.Println("")
	
	log.Println("假设有以下两个不同的复合索引:")
	log.Println("索引A: {category_id: 1, status: 1, created_at: -1}")
	log.Println("索引B: {status: 1, category_id: 1, created_at: -1}")
	log.Println("")
	
	log.Println("对于查询: {category_id: 'xxx', status: 'published'}")
	log.Println("✅ 索引A: 高效 - 可以完全利用索引")
	log.Println("⚠️ 索引B: 低效 - 只能部分利用索引")
	log.Println("")
	
	log.Println("对于查询: {status: 'published'}")
	log.Println("⚠️ 索引A: 低效 - 跳过了第一个字段")
	log.Println("✅ 索引B: 高效 - 可以利用索引")
	log.Println("")
	
	log.Println("结论: 索引字段的顺序应该根据实际查询模式来设计")
	log.Println("最常用的查询字段应该放在前面")
}