package apply

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

var db *mongo.Database

func getMongoConnection() (*mongo.Client, error) {
	format := "mongodb://%s:%s@%s:%s/?maxPoolSize=%s&minPoolSize=%s"
	uri := fmt.Sprintf(format, "dev-bpmn", "123456", "192.168.3.48", "27017", "100", "10")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	db, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func init() {
	cl, err := getMongoConnection()
	if err != nil {
		fmt.Errorf("%+v", err)
		return
	}
	db = cl.Database("dev-bpmn")
}

func Test_createDealLog(t *testing.T) {
	col := db.Collection(collectionName)
	createDealLog(col, TaskLog{
		InstanceId: "111",
		Process_Id:  "333",
		Assignee:    1,
		Updated_At:  time.Now(),
		Status:      3,
	})
}

func Test_GetItem(t *testing.T) {
	col := db.Collection(collectionName)
	status, item := getLogItem(col, bson.M{
		"instance_id": "111",
		"process_id":  "333",
	})
	t.Log(status, item)
}
