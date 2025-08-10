package biz

import (
	"context"

	"github.com/JustinRoc/mongodbL/mongo"
	"github.com/JustinRoc/pkg/slogw"
	"github.com/JustinRoc/pkg/util"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserBiz struct {
	col *mongo.Collection
}

func NewUserBiz(client *mongo.Client) *UserBiz {
	return &UserBiz{
		col: mongo.NewCollection(client, "users"),
	}
}

func (u *UserBiz) InsertUser(ctx context.Context, user *mongo.User) error {
	_, err := u.col.InsertOne(ctx, user)
	if err != nil {
		return errors.Wrap(err, "insert user")
	}
	slogw.Info("insert user success", "user", util.ToJSONStr(user))
	return nil
}

func (u *UserBiz) FindOne(ctx context.Context) error {
	// 查找用户
	var foundUser mongo.User
	filter := bson.M{"username": "john_doe"}
	if err := u.col.FindOne(ctx, filter, &foundUser); err != nil {
		return errors.Wrap(err, "find user")
	}
	slogw.Info("find user success", "user", util.ToJSONStr(foundUser))
	return nil
}


func (u *UserBiz) UpdateByID(ctx context.Context, id primitive.ObjectID) error {
	update := bson.M{
		"$set": bson.M{
			"profile.bio": "Senior Software Developer",
			"status":      "premium",
		},
	}
	// 演示更新一个文档
	// updateResult, err := u.col.UpdateByID(ctx, id, update)
	// if err != nil {
	// 	return errors.Wrap(err, "update user by id")
	// }
	// slogw.Info("update user by id success", "updateResult", util.ToJSONStr(updateResult))
	// return nil

	// 演示执行upsert操作
	opts := options.Update().SetUpsert(true)
	updateResult, err := u.col.UpdateByID(ctx, id, update, opts)
	if err != nil {
		return errors.Wrap(err, "update user by id")
	}
	slogw.Info("update user by id success", "updateResult", util.ToJSONStr(updateResult))
	return nil
}


func (u *UserBiz) Count(ctx context.Context) error {
	count, err := u.col.Count(ctx, bson.M{"status": "premium"})
	if err != nil {
		return errors.Wrap(err, "count user by id")
	}
	slogw.Info("count user by id success", "count", count)
	return nil
}