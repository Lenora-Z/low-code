package form_data

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

func TestMongoService_Database(t *testing.T) {
	col := db.Collection("form_data")
	ash := bson.M{"form_id": 84}
	ires, err := col.Find(context.TODO(), ash)
	if err != nil {
		fmt.Errorf(err.Error())
		return
	}
	for ires.Next(context.TODO()) {
		var item map[string]interface{}
		if err := ires.Decode(&item); err != nil {
			fmt.Errorf(err.Error())
			return
		}
		fmt.Println(item)
	}
}

func TestMongoService_GetByInstanceId(t *testing.T) {
	//instance := "75e18cd5-9088-11eb-9853-e24146da355"
	//col := NewMongoService(db, "form_data")
	//found, item := col.GetByInstanceId(instance)
	//t.Log(found, item)
}

func TestFormDataService_GetByObjId(t *testing.T) {
	_id := "60a1de95b93970efd1aa814a"
	col := NewFormDataService(db)
	t.Log(col.GetByObjId(_id))
}

func TestFormDataService_NewItem(t *testing.T) {
	col := NewFormDataService(db)
	col.NewItem(1, 1, 1, "222", map[string]interface{}{})
}
