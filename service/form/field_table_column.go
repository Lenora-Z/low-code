// Package form
// Created by GoLand
// @User: lenora
// @Date: 2021/7/26
// @Time: 15:25

package form

import "github.com/jinzhu/gorm"

type FieldTableColumn struct {
	ID                         uint64 `gorm:"primaryKey;column:id;type:int(11) unsigned;not null" json:"id"`
	FormID                     uint64 `gorm:"index:form_field_idx;column:form_id;type:int(11) unsigned;not null;default:0" json:"form_id"`                        // 表单id
	FieldID                    uint64 `gorm:"index:form_field_idx;column:field_id;type:int(11) unsigned;not null;default:0" json:"field_id"`                      // 控件id
	DatasourceColumnRelationID uint64 `gorm:"column:datasource_column_relation_id;type:int(11) unsigned;not null;default:0" json:"datasource_column_relation_id"` // 表字段关联表id,0表示本数据表字段
	DatasourceColumnID         uint64 `gorm:"column:datasource_column_id;type:int(11) unsigned;not null;default:0" json:"datasource_column_id"`                   // 表字段id
	ShowName                   string `gorm:"column:show_name;type:varchar(255);not null;default:''" json:"show_name"`                                            // 表头显示名称
	IsCondition                bool   `gorm:"column:is_condition;type:tinyint(1) unsigned;not null;default:0" json:"is_condition"`                                // 是否作为筛选条件 0-关闭 1-开启
}

// TableName get sql table name.获取数据库表名
func (m *FieldTableColumn) TableName() string {
	return "field_table_column"
}

var tableColumnTableName = (&FieldTableColumn{}).TableName()

func createTableColumn(db *gorm.DB, m FieldTableColumn) (error, *FieldTableColumn) {
	err := db.Table(tableColumnTableName).Create(&m).Error
	return err, &m
}

func deleteTableColumn(db *gorm.DB, formId uint64) error {
	err := db.Table(tableColumnTableName).Delete(&FieldTableColumn{}, "form_id = ?", formId).Error
	return err
}

func getTableColumnList(db *gorm.DB, fields uint64) (error, []FieldTableColumn) {
	list := make([]FieldTableColumn, 0)
	query := db.Table(tableColumnTableName)
	if fields != 0 {
		query = query.Where("field_id = ?", fields)
	}
	err := query.
		Find(&list).
		Error
	return err, list
}
