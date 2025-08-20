package biz

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/JustinRoc/mongodbL/mongo"
	"github.com/JustinRoc/pkg/slogw"
	"github.com/JustinRoc/pkg/util"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
)


type Suite struct {
	suite.Suite
	ctx    context.Context
	client *mongo.Client
	userbiz *UserBiz
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
	if config.URI != "mongodb://localhost:27017" {
		slogw.WarnContext(ctx, "don't test in corporate DB")
		return
	}

	suite := &Suite{
		ctx:    ctx,
		client: client,
		userbiz: NewUserBiz(client),
	}

	if !t.Run("TestChat", suite.TestInsertUser) {
		return
	}
}

func (suite *Suite) TestInsertUser(t *testing.T) {
	now := time.Now().Unix()
	user := &mongo.User{
		Username: fmt.Sprintf("John_%d", now),
		Email:    "john@example.com",
		Password: "hashed_password_here",
		Status:   "active",
	}
	user.Profile.FirstName = "John"
	user.Profile.LastName = "Doe"
	user.Profile.Bio = "Software Developer"
	if err := suite.userbiz.InsertUser(suite.ctx, user); err != nil {
		t.Fatalf("InsertUser failed: %v", err)
	}
}
