// Package flow
// Created by GoLand
// @User: lenora
// @Date: 2021/7/28
// @Time: 17:02

package flow

import (
	"github.com/jinzhu/gorm"
)

type FlowActivity struct {
	ID        uint64 `gorm:"primaryKey;column:id;type:int(11) unsigned;not null" json:"id"`
	FlowID    uint64 `gorm:"column:flow_id;type:int(11) unsigned;not null;default:0" json:"flow_id"`       // 流程id
	ServiceID uint64 `gorm:"column:service_id;type:int(11) unsigned;not null;default:0" json:"service_id"` // 服务id
	NodeID    uint64 `gorm:"column:node_id;type:bigint(20) unsigned;not null;default:0" json:"node_id"`    // 前端生成的节点id
	Name      string `gorm:"column:name;type:varchar(255);not null;default:''" json:"name"`                // 节点名称
	Desc      string `gorm:"column:desc;type:varchar(255);not null;default:''" json:"desc"`                // 节点描述
	Type      uint8  `gorm:"column:type;type:tinyint(3) unsigned;not null;default:1" json:"type"`          // 节点类型 1-开始节点-数据表字段时间类型;2-开始节点-固定时间类型;3-数据操作类型4-外部服务类型
}

// TableName get sql table name.获取数据库表名
func (m *FlowActivity) TableName() string {
	return "flow_activity"
}

var activityTableName = (&FlowActivity{}).TableName()

func createActivity(db *gorm.DB, m FlowActivity) (error, *FlowActivity) {
	err := db.Table(activityTableName).Create(&m).Error
	return err, &m
}

func updateActivity(db *gorm.DB, m FlowActivity) (error, *FlowActivity) {
	err := db.Table(activityTableName).Where("id = ?", m.ID).Update(m).Error
	return err, &m
}

func delActivity(db *gorm.DB, flowId uint64) error {
	err := db.Table(activityTableName).Delete(&FlowActivity{}, "flow_id = ?", flowId).Error
	return err
}
