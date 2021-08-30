// Package flow
// Created by GoLand
// @User: lenora
// @Date: 2021/8/2
// @Time: 19:47

package flow

import "github.com/jinzhu/gorm"

type FlowActivityData struct {
	ID                      uint64 `gorm:"primaryKey;column:id;type:int(11) unsigned;not null" json:"id"`
	AppID                   uint64 `gorm:"column:app_id;type:int(11) unsigned;not null;default:0" json:"app_id"`                                           // 应用id
	FlowID                  uint64 `gorm:"column:flow_id;type:int(11) unsigned;not null;default:0" json:"flow_id"`                                         // 流程id
	FlowActivityID          uint64 `gorm:"column:flow_activity_id;type:int(11) unsigned;not null;default:0" json:"flow_activity_id"`                       // 当前流程节点id
	DatasourceTableID       uint64 `gorm:"column:datasource_table_id;type:int(11) unsigned;not null;default:0" json:"datasource_table_id"`                 // 数据表id
	DatasourceColumnID      uint64 `gorm:"column:datasource_column_id;type:int(11) unsigned;not null;default:0" json:"datasource_column_id"`               // 数据表字段id
	Type                    uint8  `gorm:"column:type;type:tinyint(3);not null;default:1" json:"type"`                                                     // 参数字段类型 1:操作字段 2:条件字段
	Expression              uint8  `gorm:"column:expression;type:tinyint(3);not null;default:1" json:"expression"`                                         // 操作符 1:是/等于 2:不是/不等于
	Op                      uint8  `gorm:"column:op;type:tinyint(3);not null;default:1" json:"op"`                                                         // 关联符 0：无 1:且 2:或
	OpID                    uint64 `gorm:"column:op_id;type:int(11) unsigned;not null;default:0" json:"op_id"`                                             // 关联本表主键id
	GroupID                 uint64 `gorm:"column:group_id;type:int(11) unsigned;not null;default:0" json:"group_id"`                                       // 组编号
	InputNodeID             uint64 `gorm:"column:input_node_id;type:bigint(20) unsigned;not null;default:0" json:"input_node_id"`                          // 入参-前端生成的节点id
	InputFieldID            uint64 `gorm:"column:input_field_id;type:int(11) unsigned;not null;default:0" json:"input_field_id"`                           // 入参来源 控件id
	InputParamID            uint64 `gorm:"column:input_param_id;type:int(11) unsigned;not null;default:0" json:"input_param_id"`                           // 入参来源 服务出参id
	InputText               string `gorm:"column:input_text;type:varchar(255);not null;default:''" json:"input_text"`                                      // 入参来源 输入文本
	InputFieldTableColumnID uint64 `gorm:"column:input_field_table_column_id;type:int(11) unsigned;not null;default:0" json:"input_field_table_column_id"` // 入参来源 列表控件显示字段表id
}

// TableName get sql table name.获取数据库表名
func (m *FlowActivityData) TableName() string {
	return "flow_activity_data"
}

var activityDataTableName = (&FlowActivityData{}).TableName()

func createActivityData(db *gorm.DB, m FlowActivityData) (error, *FlowActivityData) {
	err := db.Table(activityDataTableName).Create(&m).Error
	return err, &m
}

func delActivityData(db *gorm.DB, flowId uint64) error {
	err := db.Table(activityDataTableName).Delete(&FlowActivityData{}, "flow_id = ?", flowId).Error
	return err
}
