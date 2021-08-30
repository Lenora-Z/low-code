//Created by Goland
//@User: lenora
//@Date: 2021/3/22
//@Time: 11:00 上午
package mongodb

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoService interface {
	First(filter bson.M) (bool, *mongo.SingleResult)
	FindAll(filter bson.D, offset, limit int64, sort map[string]int8) (error, uint32, *mongo.Cursor)
	Insert(items interface{}) (*mongo.InsertOneResult, error)
	UpdateItem(where map[string]interface{}, update map[string]interface{}) (*mongo.UpdateResult, error)
}

type mongoService struct {
	*mongo.Collection
}

func NewMongoService(col *mongo.Collection) MongoService {
	u := new(mongoService)
	u.Collection = col
	return u
}

func (col *mongoService) First(filter bson.M) (bool, *mongo.SingleResult) {
	mongoItem := col.FindOne(context.TODO(), filter)
	if mongoItem.Err() != nil {
		return false, nil
	}
	return true, mongoItem
}

func (col *mongoService) FindAll(filter bson.D, offset, limit int64, sort map[string]int8) (error, uint32, *mongo.Cursor) {
	var SORT = bson.M{}
	option := options.Find()
	if len(sort) > 0 {
		for i, x := range sort {
			SORT[i] = x
		}
	} else {
		SORT["_id"] = -1
	}
	option.SetSort(SORT)
	option.SetLimit(limit)
	option.SetSkip(offset)

	count, err := col.CountDocuments(context.TODO(), filter)
	if err != nil {
		return err, 0, nil
	}

	cur, err := col.Find(context.TODO(), filter, option)
	return err, uint32(count), cur
}

func (col *mongoService) Insert(items interface{}) (*mongo.InsertOneResult, error) {
	item := make(map[string]interface{})
	bytes, _ := json.Marshal(items)
	json.Unmarshal(bytes, &item)
	return col.InsertOne(context.TODO(), item)
}

func (col *mongoService) UpdateItem(where map[string]interface{}, update map[string]interface{}) (*mongo.UpdateResult, error) {
	filter := bson.D{}
	for k, v := range where {
		filter = append(filter, bson.E{k, v})
	}
	return col.UpdateOne(context.TODO(), filter, bson.D{
		{"$set", update},
	})
}
