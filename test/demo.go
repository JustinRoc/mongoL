package test

import (
	"context"
	"fmt"
	"time"

	"github.com/JustinRoc/mongodbL/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

// demonstrateUserOperations 演示用户相关操作
func demonstrateUserOperations(client *mongo.Client) error {
	ctx := context.Background()
	userCol := mongo.NewCollection(client, "users")

	// 创建用户
	user := &mongo.User{
		Username: "john_doe",
		Email:    "john@example.com",
		Password: "hashed_password_here",
		Status:   "active",
	}
	user.Profile.FirstName = "John"
	user.Profile.LastName = "Doe"
	user.Profile.Bio = "Software Developer"

	// 插入用户
	result, err := userCol.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}
	fmt.Printf("✅ 插入用户成功，ID: %v\n", result.InsertedID)

	// 查找用户
	var foundUser mongo.User
	filter := bson.M{"username": "john_doe"}
	if err := userCol.FindOne(ctx, filter, &foundUser); err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}
	fmt.Printf("✅ 找到用户: %s (%s)\n", foundUser.Username, foundUser.Email)

	// 更新用户
	update := bson.M{
		"$set": bson.M{
			"profile.bio": "Senior Software Developer",
			"status":      "premium",
		},
	}
	updateResult, err := userCol.UpdateByID(ctx, foundUser.ID, update)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	fmt.Printf("✅ 更新用户成功，修改了 %d 个文档\n", updateResult.ModifiedCount)

	// 计算用户数量
	count, err := userCol.Count(ctx, bson.M{"status": "premium"})
	if err != nil {
		return fmt.Errorf("failed to count users: %w", err)
	}
	fmt.Printf("✅ Premium 用户数量: %d\n", count)

	return nil
}

// demonstrateArticleOperations 演示文章相关操作
func demonstrateArticleOperations(client *mongo.Client) error {
	ctx := context.Background()
	articleRepo := mongo.NewCollection(client, "articles")

	// 创建文章
	article := &mongo.Article{
		Title:     "Go MongoDB 教程",
		Content:   "这是一篇关于如何在 Go 中使用 MongoDB 的详细教程...",
		AuthorID:  primitive.NewObjectID(),
		Tags:      []string{"golang", "mongodb", "tutorial", "database"},
		Status:    "published",
		ViewCount: 0,
		LikeCount: 0,
	}

	// 插入文章
	result, err := articleRepo.InsertOne(ctx, article)
	if err != nil {
		return fmt.Errorf("failed to insert article: %w", err)
	}
	fmt.Printf("✅ 插入文章成功，ID: %v\n", result.InsertedID)

	// 批量插入文章
	articles := []interface{}{
		&mongo.Article{
			Title:     "MongoDB 索引优化",
			Content:   "深入了解 MongoDB 索引的最佳实践...",
			AuthorID:  primitive.NewObjectID(),
			Tags:      []string{"mongodb", "performance", "index"},
			Status:    "published",
			ViewCount: 150,
			LikeCount: 25,
		},
		&mongo.Article{
			Title:     "Go 并发编程",
			Content:   "掌握 Go 语言的并发编程模式...",
			AuthorID:  primitive.NewObjectID(),
			Tags:      []string{"golang", "concurrency", "goroutine"},
			Status:    "draft",
			ViewCount: 0,
			LikeCount: 0,
		},
	}

	batchResult, err := articleRepo.InsertMany(ctx, articles)
	if err != nil {
		return fmt.Errorf("failed to insert articles: %w", err)
	}
	fmt.Printf("✅ 批量插入文章成功，插入了 %d 篇文章\n", len(batchResult.InsertedIDs))

	// 分页查询文章
	var paginatedArticles []mongo.Article
	pagination, err := articleRepo.FindWithPagination(
		ctx,
		bson.M{"status": "published"},
		1, // 第一页
		2, // 每页2条
		&paginatedArticles,
	)
	if err != nil {
		return fmt.Errorf("failed to find articles with pagination: %w", err)
	}
	fmt.Printf("✅ 分页查询结果: 第%d页，共%d页，总计%d篇文章\n",
		pagination.Page, pagination.TotalPage, pagination.Total)

	// 聚合查询 - 按标签统计文章数量
	pipeline := []bson.M{
		{"$match": bson.M{"status": "published"}},
		{"$unwind": "$tags"},
		{"$group": bson.M{
			"_id":   "$tags",
			"count": bson.M{"$sum": 1},
		}},
		{"$sort": bson.M{"count": -1}},
	}

	var tagStats []bson.M
	if err := articleRepo.Aggregate(ctx, pipeline, &tagStats); err != nil {
		return fmt.Errorf("failed to aggregate tag stats: %w", err)
	}
	fmt.Printf("✅ 标签统计结果: %+v\n", tagStats)

	return nil
}

// demonstrateIndexOperations 演示索引相关操作
func demonstrateIndexOperations(client *mongo.Client) error {
	ctx := context.Background()

	// 为用户集合创建索引
	userIndexes := mongo.NewCommonIndexes(client, "users")
	if err := userIndexes.CreateUserIndexes(ctx); err != nil {
		return fmt.Errorf("failed to create user indexes: %w", err)
	}
	fmt.Println("✅ 用户集合索引创建成功")

	// 为文章集合创建索引
	articleIndexes := mongo.NewCommonIndexes(client, "articles")
	if err := articleIndexes.CreateArticleIndexes(ctx); err != nil {
		return fmt.Errorf("failed to create article indexes: %w", err)
	}
	fmt.Println("✅ 文章集合索引创建成功")

	// 创建自定义索引
	indexManager := mongo.NewIndexManager(client, "articles")

	// 创建TTL索引（30天后过期）
	_, err := indexManager.CreateTTLIndex(ctx, "created_at", 30*24*time.Hour, nil)
	if err != nil {
		return fmt.Errorf("failed to create TTL index: %w", err)
	}
	fmt.Println("✅ TTL索引创建成功")

	// 列出所有索引
	indexes, err := indexManager.ListIndexes(ctx)
	if err != nil {
		return fmt.Errorf("failed to list indexes: %w", err)
	}
	fmt.Printf("✅ 文章集合共有 %d 个索引\n", len(indexes))

	return nil
}

// demonstrateTransactionOperations 演示事务相关操作
func demonstrateTransactionOperations(client *mongo.Client) error {
	ctx := context.Background()
	txnManager := mongo.NewTransactionManager(client)

	// 在事务中执行多个操作
	err := txnManager.WithTransaction(ctx, func(sessCtx mongodriver.SessionContext) error {
		userRepo := mongo.NewCollection(client, "users")
		articleRepo := mongo.NewCollection(client, "articles")

		// 创建用户
		user := &mongo.User{
			Username: "transaction_user",
			Email:    "txn@example.com",
			Status:   "active",
		}

		userResult, err := userRepo.InsertOne(sessCtx, user)
		if err != nil {
			return fmt.Errorf("failed to insert user in transaction: %w", err)
		}

		// 创建该用户的文章
		article := &mongo.Article{
			Title:    "事务测试文章",
			Content:  "这是在事务中创建的文章",
			AuthorID: userResult.InsertedID.(primitive.ObjectID),
			Tags:     []string{"transaction", "test"},
			Status:   "published",
		}

		_, err = articleRepo.InsertOne(sessCtx, article)
		if err != nil {
			return fmt.Errorf("failed to insert article in transaction: %w", err)
		}

		fmt.Println("✅ 事务中的操作执行成功")
		return nil
	})

	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}

	fmt.Println("✅ 事务提交成功")
	return nil
}
