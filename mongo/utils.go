package mongo

import (
	"fmt"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ObjectIDFromString 从字符串创建 ObjectID
func ObjectIDFromString(s string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(s)
}

// ObjectIDsFromStrings 从字符串数组创建 ObjectID 数组
func ObjectIDsFromStrings(strs []string) ([]primitive.ObjectID, error) {
	var ids []primitive.ObjectID
	for _, str := range strs {
		id, err := ObjectIDFromString(str)
		if err != nil {
			return nil, fmt.Errorf("invalid ObjectID: %s", str)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// StringFromObjectID 将 ObjectID 转换为字符串
func StringFromObjectID(id primitive.ObjectID) string {
	return id.Hex()
}

// StringsFromObjectIDs 将 ObjectID 数组转换为字符串数组
func StringsFromObjectIDs(ids []primitive.ObjectID) []string {
	var strs []string
	for _, id := range ids {
		strs = append(strs, id.Hex())
	}
	return strs
}

// BuildUpdateSet 构建更新操作的 $set 部分
func BuildUpdateSet(data interface{}) bson.M {
	update := bson.M{}
	setValue := reflect.ValueOf(data)
	setType := reflect.TypeOf(data)

	// 如果是指针，获取其指向的值
	if setValue.Kind() == reflect.Ptr {
		setValue = setValue.Elem()
		setType = setType.Elem()
	}

	if setValue.Kind() != reflect.Struct {
		return update
	}

	setFields := bson.M{}
	for i := 0; i < setValue.NumField(); i++ {
		field := setValue.Field(i)
		fieldType := setType.Field(i)

		// 跳过未导出的字段
		if !field.CanInterface() {
			continue
		}

		// 获取 bson 标签
		bsonTag := fieldType.Tag.Get("bson")
		if bsonTag == "" || bsonTag == "-" {
			continue
		}

		// 解析 bson 标签
		tagParts := strings.Split(bsonTag, ",")
		fieldName := tagParts[0]
		if fieldName == "" {
			fieldName = strings.ToLower(fieldType.Name)
		}

		// 跳过 omitempty 字段如果值为零值
		if len(tagParts) > 1 && contains(tagParts[1:], "omitempty") && field.IsZero() {
			continue
		}

		// 跳过 _id 字段
		if fieldName == "_id" {
			continue
		}

		setFields[fieldName] = field.Interface()
	}

	if len(setFields) > 0 {
		update["$set"] = setFields
	}

	return update
}

// BuildFilter 构建查询过滤器
func BuildFilter(conditions map[string]interface{}) bson.M {
	filter := bson.M{}
	for key, value := range conditions {
		if value != nil {
			filter[key] = value
		}
	}
	return filter
}

// BuildSort 构建排序条件
func BuildSort(sorts map[string]int) bson.D {
	var sortDoc bson.D
	for field, order := range sorts {
		sortDoc = append(sortDoc, bson.E{Key: field, Value: order})
	}
	return sortDoc
}

// contains 检查字符串切片是否包含指定字符串
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// ValidateObjectID 验证 ObjectID 字符串是否有效
func ValidateObjectID(id string) bool {
	_, err := primitive.ObjectIDFromHex(id)
	return err == nil
}

// NewObjectID 生成新的 ObjectID
func NewObjectID() primitive.ObjectID {
	return primitive.NewObjectID()
}

// IsZeroObjectID 检查 ObjectID 是否为零值
func IsZeroObjectID(id primitive.ObjectID) bool {
	return id.IsZero()
}

// MergeBsonM 合并多个 bson.M
func MergeBsonM(docs ...bson.M) bson.M {
	result := bson.M{}
	for _, doc := range docs {
		for key, value := range doc {
			result[key] = value
		}
	}
	return result
}

// ToObjectID 尝试将 interface{} 转换为 ObjectID
func ToObjectID(v interface{}) (primitive.ObjectID, error) {
	switch val := v.(type) {
	case primitive.ObjectID:
		return val, nil
	case string:
		return ObjectIDFromString(val)
	default:
		return primitive.NilObjectID, fmt.Errorf("cannot convert %T to ObjectID", v)
	}
}

// BuildRegexFilter 构建正则表达式过滤器
func BuildRegexFilter(field, pattern string, options ...string) bson.M {
	regex := bson.M{"$regex": pattern}
	if len(options) > 0 {
		regex["$options"] = strings.Join(options, "")
	}
	return bson.M{field: regex}
}

// BuildInFilter 构建 $in 过滤器
func BuildInFilter(field string, values []interface{}) bson.M {
	return bson.M{field: bson.M{"$in": values}}
}

// BuildRangeFilter 构建范围过滤器
func BuildRangeFilter(field string, min, max interface{}) bson.M {
	filter := bson.M{}
	if min != nil {
		filter["$gte"] = min
	}
	if max != nil {
		filter["$lte"] = max
	}
	return bson.M{field: filter}
}

// BuildTextSearchFilter 构建文本搜索过滤器
func BuildTextSearchFilter(text string) bson.M {
	return bson.M{"$text": bson.M{"$search": text}}
}