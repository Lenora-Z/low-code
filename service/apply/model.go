package apply

import (
	"github.com/Lenora-Z/low-code/service/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type TaskLog struct {
	InstanceId string    `json:"instance_id"` //实例id
	Process_Id string    `json:"process_id"`  //任务id
	Assignee   uint64    `json:"assignee"`    //处理人
	Updated_At time.Time `json:"updated_at"`  //处理时间
	Status     int8      `json:"status"`      //处理状态
}

func createDealLog(db *mongo.Collection, log TaskLog) (error, interface{}) {
	col := mongodb.NewMongoService(db)
	log.Updated_At = time.Now()
	item, err := col.Insert(log)
	if err != nil {
		return err, nil
	}
	return nil, item.InsertedID
}

func getLogItem(db *mongo.Collection, filter bson.M) (bool, *TaskLog) {
	col := mongodb.NewMongoService(db)
	found, mongoItem := col.First(filter)
	if !found {
		return false, nil
	}
	var item TaskLog
	_ = mongoItem.Decode(&item)
	return true, &item
}
