// Package source
// Created by GoLand
// @User: lenora
// @Date: 2021/7/22
// @Time: 14:27

package source

import "github.com/jinzhu/gorm"

type DatasourceColumn struct {
	ID                     uint64 `gorm:"primaryKey;column:id;type:int(11) unsigned;not null" json:"id"`
	AppID                  uint64 `gorm:"index:app_id_idx;column:app_id;type:int(11) unsigned;not null;default:0" json:"app_id"`    // 应用id
	SchemataID             uint64 `gorm:"column:schemata_id;type:int(11) unsigned;not null;default:0" json:"schemata_id"`           // 数据源id
	TableSchema            string `gorm:"column:table_schema;type:varchar(64);not null;default:''" json:"table_schema"`             // 数据库名称
	TableName              string `gorm:"column:table_name;type:varchar(64);not null;default:''" json:"table_name"`                 // 数据表名称
	DataType               string `gorm:"column:data_type;type:varchar(64);not null;default:''" json:"data_type"`                   // 字段数据类型
	CharacterMaximumLength uint64 `gorm:"column:character_maximum_length;type:bigint(21) unsigned" json:"character_maximum_length"` // 字段长度
	Name                   string `gorm:"column:column_name;type:varchar(64);not null;default:''" json:"name"`                      // 字段名称
	Comment                string `gorm:"column:column_comment;type:varchar(1024);not null;default:''" json:"comment"`              // 字段备注
	ShowType               string `gorm:"column:show_type;type:varchar(20);not null" json:"show_type"`                              // 字段显示类型(业务类型) text:文本，num: 数字 date: 日期 radio:单选 checkbox:多选 user:成员 org:组织 relation: 关联 file: 文件
	FieldType              string `gorm:"column:field_type;type:varchar(20);not null" json:"field_type"`                            // 控件类型
	IsSystemField          *bool  `gorm:"column:is_system_field;type:tinyint(1);not null;default:0" json:"is_system_field"`         // 是否为系统字段0 不是系统字段 1是系统字段
}

// TableName get sql table name.获取数据库表名
func (m *DatasourceColumn) tableName() string {
	return "datasource_column"
}

var columnTablename = (&DatasourceColumn{}).tableName()

func getColumnList(db *gorm.DB, appId uint64, tbName []string, ids []uint64, tp []string) (error, []DatasourceColumn) {
	//var count uint32
	list := make([]DatasourceColumn, 0)
	query := db.Table(columnTablename)
	if appId != 0 {
		query = query.Where(DatasourceColumnColumns.AppID+" = ?", appId)
	}
	if tbName != nil {
		query = query.Where(DatasourceColumnColumns.TableName+" in (?)", tbName)
	}

	if ids != nil {
		query = query.Where(DatasourceColumnColumns.ID+" in (?)", ids)
	}
	if tp != nil {
		query = query.Where(DatasourceColumnColumns.ShowType+" in (?)", tp)
	}
	err := query.
		Find(&list).
		Error
	return err, list
}
