package form_data

import (
	"context"
	"fmt"
	"github.com/Lenora-Z/low-code/service/mongodb"
	"github.com/Lenora-Z/low-code/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func getItemDetail(db *mongo.Collection, filter bson.M) (bool, FormDataItem) {
	col := mongodb.NewMongoService(db)
	found, mongoItem := col.First(filter)
	if !found {
		return false, nil
	}
	var item FormDataItem
	_ = mongoItem.Decode(&item)
	return true, item
}

func createItem(db *mongo.Collection, v map[string]interface{}) (error, interface{}) {
	col := mongodb.NewMongoService(db)
	mongoItem, err := col.Insert(v)
	if err != nil {
		return err, nil
	}
	return nil, mongoItem.InsertedID
}

func updateItem(db *mongo.Collection, item map[string]interface{}) (error, int64) {
	col := mongodb.NewMongoService(db)
	where := map[string]interface{}{"_id": item["_id"]}
	delete(item, "_id")
	item["updated_at"] = time.Now().Unix()
	result, err := col.UpdateItem(where, item)
	if err != nil {
		return err, 0
	}
	return nil, result.ModifiedCount

}

type FormDataList []FormDataItem

func getList(db *mongo.Collection, offset, limit uint32, filter map[string]interface{}) (error, uint32, FormDataList) {
	col := mongodb.NewMongoService(db)
	var f = bson.D{}
	f = append(f, bson.E{"is_delete", false})
	for i, x := range filter {
		f = append(f, bson.E{i, x})
	}
	err, count, cur := col.FindAll(f, int64(offset), int64(limit), nil)
	if err != nil {
		return err, 0, nil
	}
	list := make(FormDataList, 0, 50)
	for cur.Next(context.TODO()) {
		var item FormDataItem
		if err := cur.Decode(&item); err != nil {
			return err, 0, nil
		}
		list = append(list, item)
	}
	return nil, count, list
}

func (list FormDataList) UserIds() []uint64 {
	userIds := make([]uint64, 0, cap(list))
	for _, v := range list {
		idStr := fmt.Sprintf("%v", v["user_id"])
		id := utils.NewStr(idStr).Uint64()
		if id == 0 {
			continue
		}
		userIds = append(userIds, id)
	}
	return userIds
}

func (list FormDataList) ObjIds() []primitive.ObjectID {
	ids := make([]primitive.ObjectID, 0, cap(list))
	for _, v := range list {
		id, ok := v["_id"].(primitive.ObjectID)
		if !ok {
			continue
		}
		ids = append(ids, id)
	}
	return ids
}
