// Package form
// Created by GoLand
// @User: lenora
// @Date: 2021/7/26
// @Time: 15:25

package form

import "github.com/jinzhu/gorm"

type FieldRecords struct {
	ID                         uint64 `gorm:"primaryKey;column:id;type:int(11) unsigned;not null" json:"id"`
	FormID                     uint64 `gorm:"index:form_field_idx;column:form_id;type:int(11) unsigned;not null;default:0" json:"form_id"`                        // 表单id
	FieldID                    uint64 `gorm:"index:form_field_idx;column:field_id;type:int(11) unsigned;not null;default:0" json:"field_id"`                      // 控件id
	DatasourceColumnRelationID uint64 `gorm:"column:datasource_column_relation_id;type:int(11) unsigned;not null;default:0" json:"datasource_column_relation_id"` // 表字段关联关系id
	DatasourceColumnIDs        string `gorm:"column:datasource_column_ids;type:varchar(255);not null;default:''" json:"datasource_column_ids"`                    // 关联字段id列表,逗号分分割
	Mode                       uint8  `gorm:"column:mode;type:tinyint(3);not null;default:1" json:"mode"`                                                         // 呈现方式 1-卡片 2-回填
	CountType                  uint8  `gorm:"column:count_type;type:tinyint(3);not null;default:1" json:"count_type"`                                             // 关联记录数量 1-单条 2-多条
	DetailStatus               bool   `gorm:"column:detail_status;type:tinyint(1) unsigned;not null;default:0" json:"detail_status"`                              // 详情是否展示 0-关闭 1-开启
	Columns                    string `gorm:"column:columns;type:varchar(255);not null;default:''" json:"columns"`                                                // 显示字段,关联字段id列表,逗号分分割
}

// TableName get sql table name.获取数据库表名
func (m *FieldRecords) TableName() string {
	return "field_records"
}

var recordsTableName = (&FieldRecords{}).TableName()

func createRecords(db *gorm.DB, m FieldRecords) (error, *FieldRecords) {
	err := db.Table(recordsTableName).Create(&m).Error
	return err, &m
}

func deleteRecords(db *gorm.DB, formId uint64) error {
	err := db.Table(recordsTableName).Delete(&FieldRecords{}, "form_id = ?", formId).Error
	return err
}
