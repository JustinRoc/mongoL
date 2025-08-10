package biz

import (
	"context"

	"github.com/JustinRoc/mongodbL/mongo"
	"github.com/JustinRoc/pkg/slogw"
	"github.com/JustinRoc/pkg/util"
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
		return err
	}
	slogw.Info("insert user success", "user", util.ToJSONStr(user))
	return nil
}
