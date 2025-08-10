package test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/JustinRoc/mongodbL/mongo"
	"github.com/JustinRoc/pkg/slogw"
	"github.com/JustinRoc/pkg/util"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
)

/*
	运行单个测试举例：go test -run "TestSuite/Test_demonstrateUserOperations"
	运行测试套件：go test -run "TestSuite"
*/

type Suite struct {
	suite.Suite
	ctx    context.Context
	client *mongo.Client
}

func TestSuite(t *testing.T) {
	slogw.Init("", "info", nil)
	ctx := context.Background()
	// 加载 .env 文件
	err := godotenv.Load("../.env")
	if err != nil {
		slogw.ErrorContext(ctx, "加载 .env 文件失败", "err", err)
		return
	}
	// 创建 MongoDB 客户端配置
	config := &mongo.Config{
		URI:            os.Getenv("MongoAddress"),
		Database:       "testdb",
		ConnectTimeout: 10 * time.Second,
		MaxPoolSize:    100,
		MinPoolSize:    5,
	}
	slogw.Init("", "info", nil)
	slogw.InfoContext(ctx, "Connecting to MongoDB", "config", util.ToJSONStr(config))

	// 连接到 MongoDB
	client, err := mongo.NewClient(config)
	if err != nil {
		slogw.ErrorContext(ctx, "Failed to connect to MongoDB", "err", err)
		return
	}
	defer func() {
		if err := client.Close(); err != nil {
			slogw.ErrorContext(ctx, "Error closing MongoDB connection", "err", err)
		}
	}()

	suite := &Suite{
		ctx:    ctx,
		client: client,
	}

	if !t.Run("TestChat", suite.Test_demonstrateUserOperations) {
		return
	}
}

func (suite *Suite) Test_demonstrateUserOperations(t *testing.T) {
	err := demonstrateUserOperations(suite.client)
	if err != nil {
		t.Fatalf("demonstrateUserOperations failed: %v", err)
	}
}
