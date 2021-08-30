// Package flow
// Created by GoLand
// @User: lenora
// @Date: 2021/8/2
// @Time: 21:33

package flow

import "github.com/jinzhu/gorm"

type FlowActivityGateway struct {
	ID                           uint64 `gorm:"primaryKey;column:id;type:int(11) unsigned;not null" json:"id"`
	FlowID                       uint64 `gorm:"column:flow_id;type:int(11) unsigned;not null;default:0" json:"flow_id"`                                                     // 流程id
	FlowActivityID               uint64 `gorm:"column:flow_activity_id;type:int(11) unsigned;not null;default:0" json:"flow_activity_id"`                                   // 当前流程节点id
	NodeID                       uint64 `gorm:"column:node_id;type:bigint(20) unsigned;not null;default:0" json:"node_id"`                                                  // 前端生成的节点id
	LeftInputFieldID             uint64 `gorm:"column:left_input_field_id;type:int(11) unsigned;not null;default:0" json:"left_input_field_id"`                             // 左值入参来源 控件id
	LeftInputParamID             uint64 `gorm:"column:left_input_param_id;type:int(11) unsigned;not null;default:0" json:"left_input_param_id"`                             // 左值入参来源 服务出参id
	LeftInputFieldTableColumnID  uint64 `gorm:"column:left_input_field_table_column_id;type:int(11) unsigned;not null;default:0" json:"left_input_field_table_column_id"`   // 左值入参来源 列表控件显示字段表id
	Expression                   uint8  `gorm:"column:expression;type:tinyint(3);not null;default:1" json:"expression"`                                                     // 操作符 1:是/等于 2:不是/不等于
	Op                           uint8  `gorm:"column:op;type:tinyint(3);not null;default:1" json:"op"`                                                                     // 关联符 0：无 1:且 2:或
	OpID                         uint64 `gorm:"column:op_id;type:int(11) unsigned;not null;default:0" json:"op_id"`                                                         // 关联本表主键id 0：无
	GroupID                      uint64 `gorm:"column:group_id;type:int(11) unsigned;not null;default:0" json:"group_id"`                                                   // 组编号
	RightInputFieldID            uint64 `gorm:"column:right_input_field_id;type:int(11) unsigned;not null;default:0" json:"right_input_field_id"`                           // 右值入参来源 控件id
	RightInputParamID            uint64 `gorm:"column:right_input_param_id;type:int(11) unsigned;not null;default:0" json:"right_input_param_id"`                           // 右值入参来源 服务出参id
	RightInputText               string `gorm:"column:right_input_text;type:varchar(255);not null;default:''" json:"right_input_text"`                                      // 右值入参来源 手动输入文本
	RightInputFieldTableColumnID uint64 `gorm:"column:right_input_field_table_column_id;type:int(11) unsigned;not null;default:0" json:"right_input_field_table_column_id"` // 右值入参来源 列表控件显示字段表id
}

// TableName get sql table name.获取数据库表名
func (m *FlowActivityGateway) TableName() string {
	return "flow_activity_gateway"
}

var activityGatewayTableName = (&FlowActivityGateway{}).TableName()

func createActivityGateway(db *gorm.DB, m FlowActivityGateway) (error, *FlowActivityGateway) {
	err := db.Table(activityGatewayTableName).Create(&m).Error
	return err, &m
}

func delActivityGateway(db *gorm.DB, flowId uint64) error {
	err := db.Table(activityGatewayTableName).Delete(&FlowActivityGateway{}, "flow_id = ?", flowId).Error
	return err
}
