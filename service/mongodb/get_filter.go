package mongodb

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Like(value interface{}) primitive.Regex {
	return primitive.Regex{
		Pattern: fmt.Sprintf("%v", value),
	}
}

func Between(value interface{}) map[string]interface{} {
	internal, ok := value.([]interface{})
	if !ok {
		return nil
	}
	if len(internal) > 2 {
		return nil
	}
	return map[string]interface{}{
		"$gte": internal[0],
		"$lte": internal[1],
	}
}

func In(value interface{}) map[string]interface{} {
	return map[string]interface{}{
		"$in": value,
	}
}

func Where(conditions []bson.E) map[string]interface{} {
	content := make(map[string]interface{}, 0)
	for _, n := range conditions {
		switch n.Key {
		case "=":
			content["$eq"] = n.Value
		case "<":
			content["$lt"] = n.Value
		case "<=":
			content["$lte"] = n.Value
		case ">":
			content["$gt"] = n.Value
		case ">=":
			content["$gte"] = n.Value
		default:
			continue
		}
	}
	return content
}
