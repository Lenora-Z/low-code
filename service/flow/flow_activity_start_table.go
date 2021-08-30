// Package flow
// Created by GoLand
// @User: lenora
// @Date: 2021/7/31
// @Time: 18:11

package flow

import "github.com/jinzhu/gorm"

type FlowActivityStartTable struct {
	ID                 uint64 `gorm:"primaryKey;column:id;type:int(11) unsigned;not null" json:"id"`
	FlowID             uint64 `gorm:"column:flow_id;type:int(11) unsigned;not null;default:0" json:"flow_id"`                           // 流程id
	FlowActivityID     uint64 `gorm:"column:flow_activity_id;type:int(11) unsigned;not null;default:0" json:"flow_activity_id"`         // 节点id
	DatasourceTableID  uint64 `gorm:"column:datasource_table_id;type:int(11) unsigned;not null;default:0" json:"datasource_table_id"`   // 表id
	DatasourceColumnID uint64 `gorm:"column:datasource_column_id;type:int(11) unsigned;not null;default:0" json:"datasource_column_id"` // 表字段id
	TriggerInterval    int64  `gorm:"column:trigger_interval;type:int(5);not null;default:0" json:"trigger_interval"`                   // 触发间隔 0：当前时间 -n:提前n小时 n:延后n小时
}

// TableName get sql table name.获取数据库表名
func (m *FlowActivityStartTable) TableName() string {
	return "flow_activity_start_table"
}

var activityTableTableName = (&FlowActivityStartTable{}).TableName()

func createActivityTable(db *gorm.DB, m FlowActivityStartTable) (error, *FlowActivityStartTable) {
	err := db.Table(activityTableTableName).Create(&m).Error
	return err, &m
}

func delActivityTable(db *gorm.DB, flowId uint64) error {
	err := db.Table(activityTableTableName).Delete(&FlowActivityStartTable{}, "flow_id = ?", flowId).Error
	return err
}
