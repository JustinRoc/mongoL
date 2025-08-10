package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Document 文档接口
type Document interface {
	GetID() primitive.ObjectID
	SetID(id primitive.ObjectID)
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	SetUpdatedAt(t time.Time)
}

// BaseDocument 基础文档结构体
type BaseDocument struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// GetID 获取文档ID
func (d *BaseDocument) GetID() primitive.ObjectID {
	return d.ID
}

// SetID 设置文档ID
func (d *BaseDocument) SetID(id primitive.ObjectID) {
	d.ID = id
}

// GetCreatedAt 获取创建时间
func (d *BaseDocument) GetCreatedAt() time.Time {
	return d.CreatedAt
}

// GetUpdatedAt 获取更新时间
func (d *BaseDocument) GetUpdatedAt() time.Time {
	return d.UpdatedAt
}

// SetUpdatedAt 设置更新时间
func (d *BaseDocument) SetUpdatedAt(t time.Time) {
	d.UpdatedAt = t
}

// BeforeInsert 插入前的钩子函数
func (d *BaseDocument) BeforeInsert() {
	now := time.Now()
	if d.ID.IsZero() {
		d.ID = primitive.NewObjectID()
	}
	if d.CreatedAt.IsZero() {
		d.CreatedAt = now
	}
	d.UpdatedAt = now
}

// BeforeUpdate 更新前的钩子函数
func (d *BaseDocument) BeforeUpdate() {
	d.UpdatedAt = time.Now()
}

// User 用户文档示例
type User struct {
	BaseDocument `bson:",inline"`
	Username     string `bson:"username" json:"username"`
	Email        string `bson:"email" json:"email"`
	Password     string `bson:"password" json:"-"` // 不在JSON中显示密码
	Status       string `bson:"status" json:"status"`
	Profile      struct {
		FirstName string `bson:"first_name" json:"first_name"`
		LastName  string `bson:"last_name" json:"last_name"`
		Avatar    string `bson:"avatar" json:"avatar"`
		Bio       string `bson:"bio" json:"bio"`
	} `bson:"profile" json:"profile"`
}

// Article 文章文档示例
type Article struct {
	BaseDocument `bson:",inline"`
	Title        string               `bson:"title" json:"title"`
	Content      string               `bson:"content" json:"content"`
	AuthorID     primitive.ObjectID   `bson:"author_id" json:"author_id"`
	Tags         []string             `bson:"tags" json:"tags"`
	Status       string               `bson:"status" json:"status"` // draft, published, archived
	ViewCount    int64                `bson:"view_count" json:"view_count"`
	LikeCount    int64                `bson:"like_count" json:"like_count"`
	CategoryID   primitive.ObjectID   `bson:"category_id,omitempty" json:"category_id,omitempty"`
	Comments     []primitive.ObjectID `bson:"comments" json:"comments"`
}

// Category 分类文档示例
type Category struct {
	BaseDocument `bson:",inline"`
	Name         string `bson:"name" json:"name"`
	Description  string `bson:"description" json:"description"`
	ParentID     *primitive.ObjectID `bson:"parent_id,omitempty" json:"parent_id,omitempty"`
	Sort         int    `bson:"sort" json:"sort"`
	IsActive     bool   `bson:"is_active" json:"is_active"`
}