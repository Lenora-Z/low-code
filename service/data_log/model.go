package data_log

import (
	"github.com/Lenora-Z/low-code/service/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

func createLog(db *mongo.Collection, v map[string]interface{}) (error, interface{}) {
	col := mongodb.NewMongoService(db)
	mongoItem, err := col.Insert(v)
	if err != nil {
		return err, nil
	}
	return nil, mongoItem.InsertedID
}
