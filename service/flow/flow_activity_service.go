// Package flow
// Created by GoLand
// @User: lenora
// @Date: 2021/7/31
// @Time: 18:29

package flow

import "github.com/jinzhu/gorm"

type FlowActivityService struct {
	ID                      uint64 `gorm:"primaryKey;column:id;type:int(11) unsigned;not null" json:"id"`
	AppID                   uint64 `gorm:"column:app_id;type:int(11) unsigned;not null;default:0" json:"app_id"`                                           // 应用id
	FlowID                  uint64 `gorm:"column:flow_id;type:int(11) unsigned;not null;default:0" json:"flow_id"`                                         // 流程id
	FlowActivityID          uint64 `gorm:"column:flow_activity_id;type:int(11) unsigned;not null;default:0" json:"flow_activity_id"`                       // 当前流程节点id
	ParamID                 uint64 `gorm:"column:param_id;type:int(11) unsigned;not null;default:0" json:"param_id"`                                       // 参数id
	InputNodeID             uint64 `gorm:"column:input_node_id;type:bigint(20) unsigned;not null;default:0" json:"input_node_id"`                          // 入参-前端生成的节点id
	InputFieldID            uint64 `gorm:"column:input_field_id;type:int(11) unsigned;not null;default:0" json:"input_field_id"`                           // 入参来源 控件id
	InputParamID            uint64 `gorm:"column:input_param_id;type:int(11) unsigned;not null;default:0" json:"input_param_id"`                           // 入参来源 服务出参id
	InputLowcodeUserID      uint64 `gorm:"column:input_lowcode_user_id;type:int(11) unsigned;not null;default:0" json:"input_lowcode_user_id"`             // 入参来源 低代码用户id
	InputText               string `gorm:"column:input_text;type:text;" json:"input_text"`                                                                 // 入参来源 输入文本
	InputFieldTableColumnID uint64 `gorm:"column:input_field_table_column_id;type:int(11) unsigned;not null;default:0" json:"input_field_table_column_id"` // 入参来源 列表控件显示字段表id
}

// TableName get sql table name.获取数据库表名
func (m *FlowActivityService) TableName() string {
	return "flow_activity_service"
}

var activityServiceTableName = (&FlowActivityService{}).TableName()

func createActivityService(db *gorm.DB, m FlowActivityService) (error, *FlowActivityService) {
	err := db.Table(activityServiceTableName).Create(&m).Error
	return err, &m
}

func delActivityService(db *gorm.DB, flowId uint64) error {
	err := db.Table(activityServiceTableName).Delete(&FlowActivityService{}, "flow_id = ?", flowId).Error
	return err
}
