// Package flow
// Created by GoLand
// @User: lenora
// @Date: 2021/7/31
// @Time: 18:04

package flow

import (
	"github.com/jinzhu/gorm"
	"time"
)

type FlowActivityStartTime struct {
	ID              uint64    `gorm:"primaryKey;column:id;type:int(11) unsigned;not null" json:"id"`
	FlowID          uint64    `gorm:"column:flow_id;type:int(11) unsigned;not null;default:0" json:"flow_id"`                   // 流程id
	FlowActivityID  uint64    `gorm:"column:flow_activity_id;type:int(11) unsigned;not null;default:0" json:"flow_activity_id"` // 节点id
	TriggerTime     time.Time `gorm:"column:trigger_time;type:datetime" json:"trigger_time"`                                    // 触发时间
	TriggerInterval int64     `gorm:"column:trigger_interval;type:int(5);not null;default:0" json:"trigger_interval"`           // 触发间隔 0：当前时间 -n:提前n小时 n:延后n小时
}

// TableName get sql table name.获取数据库表名
func (m *FlowActivityStartTime) TableName() string {
	return "flow_activity_start_time"
}

var activityTimeTableName = (&FlowActivityStartTime{}).TableName()

func createActivityTime(db *gorm.DB, m FlowActivityStartTime) (error, *FlowActivityStartTime) {
	err := db.Table(activityTimeTableName).Create(&m).Error
	return err, &m
}

func delActivityTime(db *gorm.DB, flowId uint64) error {
	err := db.Table(activityTimeTableName).Delete(&FlowActivityStartTime{}, "flow_id = ?", flowId).Error
	return err
}
