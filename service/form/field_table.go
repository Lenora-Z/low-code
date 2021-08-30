// Package form
// Created by GoLand
// @User: lenora
// @Date: 2021/7/26
// @Time: 15:25

package form

import "github.com/jinzhu/gorm"

type FieldTable struct {
	ID                uint64 `gorm:"primaryKey;column:id;type:int(11) unsigned;not null" json:"id"`
	FormID            uint64 `gorm:"index:form_field_idx;column:form_id;type:int(11) unsigned;not null;default:0" json:"form_id"`    // 表单id
	FieldID           uint64 `gorm:"index:form_field_idx;column:field_id;type:int(11) unsigned;not null;default:0" json:"field_id"`  // 控件id
	DatasourceTableID uint64 `gorm:"column:datasource_table_id;type:int(11) unsigned;not null;default:0" json:"datasource_table_id"` // 数据表id
	IsExport          bool   `gorm:"column:is_export;type:tinyint(1) unsigned;not null;default:0" json:"is_export"`                  // 导出配置 0: 关闭导出 1:允许导出
	IsFilter          bool   `gorm:"column:is_filter;type:tinyint(1) unsigned;not null;default:0" json:"is_filter"`                  // 数据过滤配置  0:关闭过滤  1:允许过滤
}

// TableName get sql table name.获取数据库表名
func (m *FieldTable) TableName() string {
	return "field_table"
}

var tableTableName = (&FieldTable{}).TableName()

func createTable(db *gorm.DB, m FieldTable) (error, *FieldTable) {
	err := db.Table(tableTableName).Create(&m).Error
	return err, &m
}

func deleteTable(db *gorm.DB, formId uint64) error {
	err := db.Table(tableTableName).Delete(&FieldTable{}, "form_id = ?", formId).Error
	return err
}
