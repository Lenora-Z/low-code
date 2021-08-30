// Package form
// Created by GoLand
// @User: lenora
// @Date: 2021/7/26
// @Time: 15:25

package form

import "github.com/jinzhu/gorm"

type FieldTableFilter struct {
	ID                         uint64 `gorm:"primaryKey;column:id;type:int(11) unsigned;not null" json:"id"`
	FormID                     uint64 `gorm:"index:form_field_idx;column:form_id;type:int(11) unsigned;not null;default:0" json:"form_id"`                        // 表单id
	FieldID                    uint64 `gorm:"index:form_field_idx;column:field_id;type:int(11) unsigned;not null;default:0" json:"field_id"`                      // 控件id
	DatasourceColumnRelationID uint64 `gorm:"column:datasource_column_relation_id;type:int(11) unsigned;not null;default:0" json:"datasource_column_relation_id"` // 表字段关联表id,0表示本数据表字段
	DatasourceColumnID         uint64 `gorm:"column:datasource_column_id;type:int(11) unsigned;not null;default:0" json:"datasource_column_id"`                   // 表字段id
	FieldType                  string `gorm:"column:field_type;type:varchar(20);not null" json:"field_type"`                                                      // 控件类型
	FieldTypeCondition         uint8  `gorm:"column:field_type_condition;type:tinyint(3);not null;default:1" json:"field_type_condition"`                         // 控件筛选选项条件
	FieldTypeConditionValue    string `gorm:"column:field_type_condition_value;type:varchar(255);not null;default:''" json:"field_type_condition_value"`          // 控件筛选选项值
}

// TableName get sql table name.获取数据库表名
func (m *FieldTableFilter) TableName() string {
	return "field_table_filter"
}

var tableFilterTableName = (&FieldTableFilter{}).TableName()

func createTableFilter(db *gorm.DB, m FieldTableFilter) (error, *FieldTableFilter) {
	err := db.Table(tableFilterTableName).Create(&m).Error
	return err, &m
}

func deleteTableFilter(db *gorm.DB, formId uint64) error {
	err := db.Table(tableFilterTableName).Delete(&FieldTableFilter{}, "form_id = ?", formId).Error
	return err
}
