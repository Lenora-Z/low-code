package apply

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service interface {
	//获取记录by process_id
	GetLogByProcessId(id string) (bool, *TaskLog)
	// NewDealLog 新增处理记录
	NewDealLog(instanceId, processId string, assignee uint64, status int8) (error, interface{})
}

type service struct {
	db *mongo.Collection
}

func NewService(db *mongo.Database) Service {
	u := new(service)
	u.db = db.Collection(collectionName)
	return u
}

func (srv *service) NewDealLog(instanceId, processId string, assignee uint64, status int8) (error, interface{}) {
	return createDealLog(srv.db, TaskLog{
		InstanceId: instanceId,
		Process_Id: processId,
		Assignee:   assignee,
		Status:     status,
	})
}

func (srv *service) GetLogByProcessId(id string) (bool, *TaskLog) {
	f := bson.M{"process_id": id}
	return getLogItem(srv.db, f)
}
