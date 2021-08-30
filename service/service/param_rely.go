//Created by Goland
//@User: lenora
//@Date: 2021/3/15
//@Time: 8:43 下午
package service

import (
	"github.com/jinzhu/gorm"
	"time"
)

type ParamRely struct {
	ID        uint64    `gorm:"primary_key;column:id;type:int(11) unsigned;not null" json:"id"`
	AppID     uint64    `gorm:"column:app_id;type:int(11) unsigned;not null" json:"app_id"`         // 应用id
	FlowID    uint64    `gorm:"column:flow_id;type:int(11) unsigned;not null" json:"flow_id"`       // 流程id
	ServiceID uint64    `gorm:"column:service_id;type:int(11) unsigned;not null" json:"service_id"` // 服务id
	ParamID   uint64    `gorm:"column:param_id;type:int(11) unsigned;not null" json:"param_id"`     // 参数id
	FormID    uint64    `gorm:"column:form_id;type:int(11) unsigned;not null" json:"form_id"`       // 表单id
	FieldID   uint64    `gorm:"column:field_id;type:int(11) unsigned;not null" json:"field_id"`     // 字段控件id
	Random    string    `gorm:"column:random;type:varchar(50);not null" json:"random"`              // 版本hash
	Status    uint8     `gorm:"column:status;type:tinyint(2);not null" json:"status"`               // 版本状态
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null" json:"created_at"`
}

// TableName get sql table name.获取数据库表名
func (m *ParamRely) TableName() string {
	return "param_rely"
}

var prTablename = (&ParamRely{}).TableName()

func createParamRely(db *gorm.DB, m ParamRely) (error, *ParamRely) {
	err := db.Table(prTablename).Create(&m).Error
	return err, &m
}

func batchUpdateRelies(db *gorm.DB, where, change map[string]interface{}) error {
	query := db.Table(prTablename)
	for i, x := range where {
		query = query.Where(i+" = ?", x)
	}
	err := query.Updates(change).Error
	return err
}

func getParamRelyList(db *gorm.DB, offset, limit uint32, srvId, formId, flowId []uint64) (error, uint32, []ParamRely) {
	var count uint32
	list := make([]ParamRely, 0)
	query := db.Table(prTablename).Where("status = ?", VALID)
	if len(srvId) > 0 {
		query = query.Where("service_id in (?)", srvId)
	}
	if len(formId) > 0 {
		query = query.Where("form_id in (?)", formId)
	}
	if len(flowId) > 0 {
		query = query.Where("flow_id in (?)", flowId)
	}
	err := query.Count(&count).
		Offset(offset).
		Limit(limit).
		Find(&list).
		Error
	return err, count, list
}
