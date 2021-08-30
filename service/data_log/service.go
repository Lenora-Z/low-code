package data_log

import (
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Service interface {
	NewLog(userId, formId uint64, action string, data map[string]interface{}) (error, interface{})
}

type service struct {
	db *mongo.Collection
}

func NewService(db *mongo.Database) Service {
	u := new(service)
	u.db = db.Collection(collectionName)
	return u
}

func (srv *service) NewLog(userId, formId uint64, action string, data map[string]interface{}) (error, interface{}) {
	data["user_id"] = userId
	data["form_id"] = formId
	data["action"] = action
	data["created_at"] = time.Now().Unix()
	delete(data, "updated_at")
	delete(data, "is_delete")
	return createLog(srv.db, data)
}
