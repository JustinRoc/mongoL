# MongoDB Go 包

这是一个基于 Go 语言的 MongoDB 标准服务端代码包，提供了完整的 MongoDB 操作封装，包括连接管理、CRUD 操作、索引管理、事务支持等功能。

## 功能特性

- ✅ **连接管理**: 支持连接池配置、自动重连、健康检查
- ✅ **文档模型**: 提供基础文档结构体和常用文档类型
- ✅ **CRUD 操作**: 完整的增删改查操作，支持批量操作
- ✅ **分页查询**: 内置分页支持，自动计算总页数
- ✅ **聚合查询**: 支持复杂的聚合管道操作
- ✅ **索引管理**: 提供索引创建、删除、统计等功能
- ✅ **事务支持**: 支持 MongoDB 事务操作
- ✅ **工具函数**: 提供常用的辅助函数和类型转换

## 安装

```bash
go mod init your-project
go get go.mongodb.org/mongo-driver
```

## 快速开始

### 1. 连接 MongoDB

```go
package main

import (
    "context"
    "log"
    "time"
    
    "your-project/mongo"
)

func main() {
    // 创建配置
    config := &mongo.Config{
        URI:            "mongodb://localhost:27017",
        Database:       "myapp",
        ConnectTimeout: 10 * time.Second,
        MaxPoolSize:    100,
        MinPoolSize:    5,
    }
    
    // 连接数据库
    client, err := mongo.NewClient(config)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()
    
    // 测试连接
    if err := client.Ping(); err != nil {
        log.Fatal(err)
    }
    
    log.Println("Connected to MongoDB!")
}
```

### 2. 定义文档结构

```go
// 使用内置的用户文档
type User struct {
    mongo.BaseDocument `bson:",inline"`
    Username           string `bson:"username" json:"username"`
    Email              string `bson:"email" json:"email"`
    Status             string `bson:"status" json:"status"`
}

// 或者自定义文档
type Product struct {
    mongo.BaseDocument `bson:",inline"`
    Name               string  `bson:"name" json:"name"`
    Price              float64 `bson:"price" json:"price"`
    Category           string  `bson:"category" json:"category"`
}
```

### 3. 基础 CRUD 操作

```go
func crudExample(client *mongo.Client) {
    ctx := context.Background()
    userRepo := mongo.NewRepository(client, "users")
    
    // 创建用户
    user := &mongo.User{
        Username: "john_doe",
        Email:    "john@example.com",
        Status:   "active",
    }
    
    // 插入
    result, err := userRepo.InsertOne(ctx, user)
    if err != nil {
        log.Fatal(err)
    }
    
    // 查询
    var foundUser mongo.User
    filter := bson.M{"username": "john_doe"}
    err = userRepo.FindOne(ctx, filter, &foundUser)
    if err != nil {
        log.Fatal(err)
    }
    
    // 更新
    update := bson.M{"$set": bson.M{"status": "premium"}}
    _, err = userRepo.UpdateByID(ctx, foundUser.ID, update)
    if err != nil {
        log.Fatal(err)
    }
    
    // 删除
    _, err = userRepo.DeleteByID(ctx, foundUser.ID)
    if err != nil {
        log.Fatal(err)
    }
}
```

### 4. 分页查询

```go
func paginationExample(client *mongo.Client) {
    ctx := context.Background()
    userRepo := mongo.NewRepository(client, "users")
    
    var users []mongo.User
    pagination, err := userRepo.FindWithPagination(
        ctx,
        bson.M{"status": "active"}, // 过滤条件
        1,  // 页码
        10, // 每页数量
        &users,
    )
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("第 %d 页，共 %d 页，总计 %d 个用户\n", 
        pagination.Page, pagination.TotalPage, pagination.Total)
}
```

### 5. 聚合查询

```go
func aggregationExample(client *mongo.Client) {
    ctx := context.Background()
    userRepo := mongo.NewRepository(client, "users")
    
    // 按状态统计用户数量
    pipeline := []bson.M{
        {"$group": bson.M{
            "_id":   "$status",
            "count": bson.M{"$sum": 1},
        }},
        {"$sort": bson.M{"count": -1}},
    }
    
    var results []bson.M
    err := userRepo.Aggregate(ctx, pipeline, &results)
    if err != nil {
        log.Fatal(err)
    }
    
    for _, result := range results {
        fmt.Printf("状态: %s, 数量: %v\n", result["_id"], result["count"])
    }
}
```

### 6. 索引管理

```go
func indexExample(client *mongo.Client) {
    ctx := context.Background()
    
    // 创建常用索引
    userIndexes := mongo.NewCommonIndexes(client, "users")
    err := userIndexes.CreateUserIndexes(ctx)
    if err != nil {
        log.Fatal(err)
    }
    
    // 创建自定义索引
    indexManager := mongo.NewIndexManager(client, "products")
    
    // 创建唯一索引
    _, err = indexManager.CreateUniqueIndex(ctx, "sku", nil)
    if err != nil {
        log.Fatal(err)
    }
    
    // 创建文本索引
    _, err = indexManager.CreateTextIndex(ctx, []string{"name", "description"}, nil)
    if err != nil {
        log.Fatal(err)
    }
    
    // 创建TTL索引（30天过期）
    _, err = indexManager.CreateTTLIndex(ctx, "expires_at", 30*24*time.Hour, nil)
    if err != nil {
        log.Fatal(err)
    }
}
```

### 7. 事务操作

```go
func transactionExample(client *mongo.Client) {
    ctx := context.Background()
    txnManager := mongo.NewTransactionManager(client)
    
    err := txnManager.WithTransaction(ctx, func(sessCtx mongodriver.SessionContext) error {
        userRepo := mongo.NewRepository(client, "users")
        orderRepo := mongo.NewRepository(client, "orders")
        
        // 在事务中创建用户
        user := &mongo.User{
            Username: "buyer",
            Email:    "buyer@example.com",
            Status:   "active",
        }
        userResult, err := userRepo.InsertOne(sessCtx, user)
        if err != nil {
            return err
        }
        
        // 在事务中创建订单
        order := map[string]interface{}{
            "user_id": userResult.InsertedID,
            "amount":  99.99,
            "status":  "pending",
        }
        _, err = orderRepo.InsertOne(sessCtx, order)
        return err
    })
    
    if err != nil {
        log.Fatal(err)
    }
}
```

## 工具函数

包中提供了许多有用的工具函数：

```go
// ObjectID 转换
id, err := mongo.ObjectIDFromString("507f1f77bcf86cd799439011")
str := mongo.StringFromObjectID(id)

// 构建查询过滤器
filter := mongo.BuildFilter(map[string]interface{}{
    "status": "active",
    "age":    mongo.BuildRangeFilter("age", 18, 65),
})

// 构建正则表达式查询
regexFilter := mongo.BuildRegexFilter("name", "^John", "i")

// 构建文本搜索
textFilter := mongo.BuildTextSearchFilter("golang mongodb")
```

## 配置选项

```go
config := &mongo.Config{
    URI:            "mongodb://localhost:27017", // MongoDB 连接字符串
    Database:       "myapp",                     // 数据库名称
    ConnectTimeout: 10 * time.Second,           // 连接超时时间
    MaxPoolSize:    100,                        // 最大连接池大小
    MinPoolSize:    5,                          // 最小连接池大小
}
```

## 运行示例

1. 确保 MongoDB 服务正在运行
2. 运行示例程序：

```bash
go run main.go
```

## 项目结构

```
mongo/
├── client.go      # 客户端连接管理
├── document.go    # 文档结构体定义
├── operations.go  # CRUD 操作
├── index.go       # 索引管理
├── transaction.go # 事务支持
└── utils.go       # 工具函数
```

## 最佳实践

1. **连接管理**: 在应用启动时创建客户端，在应用关闭时关闭连接
2. **错误处理**: 始终检查和处理错误
3. **上下文使用**: 为所有操作传递适当的上下文
4. **索引优化**: 根据查询模式创建合适的索引
5. **事务使用**: 只在需要原子性操作时使用事务
6. **分页查询**: 对大量数据使用分页避免内存问题

## 许可证

MIT License