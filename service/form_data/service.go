package form_data

import (
	"github.com/Lenora-Z/low-code/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type FormDataService interface {
	DeleteItem(id string) error
	// FilterList 列表筛选
	FilterList(page, limit uint32, formId uint64, filter map[string]interface{}) (error, uint32, FormDataList)
	//获取记录formId&key
	GetByFormAndKey(formId uint64, keyName string, value interface{}) (bool, FormDataItem)
	//获取记录by内部id
	GetByObjId(id string) (bool, FormDataItem)
	//获取记录by实例id
	GetByInstanceId(instanceId string) (bool, FormDataItem)
	//获取某表单下的全部记录
	GetFullFormList(formId uint64) (error, uint32, FormDataList)
	//分页获取表单下的记录
	GetListByFormId(page, limit uint32, formId uint64) (error, uint32, FormDataList)
	// GetListByParentId 根据父id获取记录列表
	GetListByParentId(id interface{}) (error, uint32, FormDataList)
	//新增记录
	NewItem(userId, formId, flowId uint64, instance string, data map[string]interface{}) (error, interface{})
	// UpdateItem 修改记录
	UpdateItem(id string, data map[string]interface{}) (error, int64)
}

type formDataService struct {
	db *mongo.Collection
}

func NewFormDataService(db *mongo.Database) FormDataService {
	u := new(formDataService)
	u.db = db.Collection(collectionName)
	return u
}
func (srv *formDataService) GetByObjId(id string) (bool, FormDataItem) {
	objId, _ := primitive.ObjectIDFromHex(id)
	f := bson.M{"_id": objId}
	return getItemDetail(srv.db, f)
}

func (srv *formDataService) GetByInstanceId(instanceId string) (bool, FormDataItem) {
	f := bson.M{"instance_id": instanceId}
	return getItemDetail(srv.db, f)
}

func (srv *formDataService) GetByFormAndKey(formId uint64, keyName string, value interface{}) (bool, FormDataItem) {
	f := bson.M{
		"form_id": formId,
		keyName:   value,
	}
	return getItemDetail(srv.db, f)
}

func (srv *formDataService) NewItem(userId, formId, flowId uint64, instance string, data map[string]interface{}) (error, interface{}) {
	data["instance_id"] = instance
	data["user_id"] = userId
	data["form_id"] = formId
	data["flow_id"] = flowId
	data["created_at"] = time.Now().Unix()
	data["is_delete"] = false
	return createItem(srv.db, data)
}

func (srv *formDataService) GetListByFormId(page, limit uint32, formId uint64) (error, uint32, FormDataList) {
	offset := (page - 1) * limit
	filter := make(map[string]interface{})
	filter["form_id"] = formId
	return getList(srv.db, offset, limit, filter)
}

func (srv *formDataService) GetListByParentId(id interface{}) (error, uint32, FormDataList) {
	return getList(srv.db, 0, utils.MAX_LIMIT, map[string]interface{}{
		"parent_id": id,
	})
}

func (srv *formDataService) GetFullFormList(formId uint64) (error, uint32, FormDataList) {
	return srv.GetListByFormId(1, utils.MAX_LIMIT, formId)
}

func (srv *formDataService) FilterList(page, limit uint32, formId uint64, filter map[string]interface{}) (error, uint32, FormDataList) {
	offset := (page - 1) * limit
	filter["form_id"] = formId
	return getList(srv.db, offset, limit, filter)
}

func (srv *formDataService) UpdateItem(id string, data map[string]interface{}) (error, int64) {
	if data == nil {
		data = make(map[string]interface{})
	}
	objId, _ := primitive.ObjectIDFromHex(id)
	data["_id"] = objId
	return updateItem(srv.db, data)
}

func (srv *formDataService) DeleteItem(id string) error {
	data := make(map[string]interface{})
	objId, _ := primitive.ObjectIDFromHex(id)
	data["_id"] = objId
	data["is_delete"] = true
	err, _ := updateItem(srv.db, data)
	return err
}
