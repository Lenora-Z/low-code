//Created by Goland
//@User: lenora
//@Date: 2021/3/12
//@Time: 4:09 下午
package flow

import (
	"github.com/jinzhu/gorm"
	"time"
)

type FlowMapping struct {
	ID        uint64    `gorm:"primary_key;column:id;type:int(11) unsigned;not null" json:"id"`
	FormID    uint64    `gorm:"column:form_id;type:int(11) unsigned;not null" json:"form_id"` // 表单id
	FlowID    uint64    `gorm:"column:flow_id;type:int(11);not null" json:"flow_id"`          // 流程id
	VersionID uint64    `gorm:"column:version_id;type:int(11);not null" json:"version_id"`    // 版本id
	AppID     uint64    `gorm:"column:app_id;type:int(11);not null" json:"app_id"`            // 应用id
	Relation  string    `gorm:"column:relation;type:text" json:"relation"`                    // 映射关系
	Status    int8      `gorm:"column:status;type:tinyint(1);not null" json:"status"`         // 状态值 0-未生效 1-已生效 2-已失效
	Random    string    `gorm:"column:random;type:varchar(20);not null" json:"random"`        // 随机数
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null" json:"created_at"`
}

// TableName get sql table name.获取数据库表名
func (m *FlowMapping) TableName() string {
	return "flow_mapping"
}

var fmTablename = (&FlowMapping{}).TableName()

func createFlowMapping(db *gorm.DB, m FlowMapping) (error, *FlowMapping) {
	err := db.Table(fmTablename).Create(&m).Error
	return err, &m
}

func batchDeleteMapping(db *gorm.DB, where map[string]interface{}) error {
	query := db.Table(fmTablename)
	for i, x := range where {
		query = query.Where(i+" = ?", x)
	}
	return query.Delete(&FlowMapping{}).Error
}

func batchUpdateMapping(db *gorm.DB, where, change map[string]interface{}) error {
	query := db.Table(fmTablename)
	for i, x := range where {
		query = query.Where(i+" = ?", x)
	}
	err := query.Updates(change).Error
	return err
}

type FlowMappingList []FlowMapping

func getFMList(db *gorm.DB, offset, limit uint32, appId uint64, forms, flows []uint64, rand string) (error, uint32, FlowMappingList) {
	var count uint32
	list := make([]FlowMapping, 0)
	query := db.Table(fmTablename).Where("app_id  = ?", appId)
	if len(forms) > 0 {
		query = query.Where("form_id in (?)", forms).Where("flow_id <> 0")
	}
	if len(flows) > 0 {
		query = query.Where("flow_id in (?)", flows).Where("form_id <> 0")
	}
	if rand != "" {
		query = query.Where("random = ?", rand)
	}
	err := query.Count(&count).
		Offset(offset).
		Limit(limit).
		Find(&list).
		Error
	return err, count, list
}

func (uList FlowMappingList) Flows() []uint64 {
	userIds := make([]uint64, 0, len(uList))
	for _, v := range uList {
		userIds = append(userIds, v.FlowID)
	}
	return userIds
}

func (uList FlowMappingList) Forms() []uint64 {
	userIds := make([]uint64, 0, len(uList))
	for _, v := range uList {
		userIds = append(userIds, v.FormID)
	}
	return userIds
}

func (uList FlowMappingList) ItemByForm(id uint64) FlowMappingList {
	list := make(FlowMappingList, 0, cap(uList))
	for _, x := range uList {
		if x.FormID == id {
			list = append(list, x)
		}
	}
	return list
}
