package form_data

import "time"

const collectionName = "form_data"

type FormDataItem = map[string]interface{}

type FormDataItemStruct struct {
	ObjId      string    `json:"_id"`
	InstanceId string    `json:"instance_id"`
	UserId     uint64    `json:"user_id"`
	FormId     uint64    `json:"form_id"`
	FlowId     uint64    `json:"flow_id"`
	CreatedAt  time.Time `json:"created_at"`
	FormDataItem
}
