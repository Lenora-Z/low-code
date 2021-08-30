// Package source
// Created by GoLand
// @User: lenora
// @Date: 2021/7/22
// @Time: 14:33

package source

import "github.com/jinzhu/gorm"

type DatasourceColumnRelation struct {
	ID             uint64 `gorm:"primaryKey;column:id;type:int(11) unsigned;not null" json:"id"`
	SourceTableID  uint64 `gorm:"column:source_table_id;type:int(11) unsigned;not null;default:0" json:"source_table_id"`   // 源数据表id
	SourceColumnID uint64 `gorm:"column:source_column_id;type:int(11) unsigned;not null;default:0" json:"source_column_id"` // 源表字段id
	TargetTableID  uint64 `gorm:"column:target_table_id;type:int(11) unsigned;not null;default:0" json:"target_table_id"`   // 目标数据表id
	TargetColumnID uint64 `gorm:"column:target_column_id;type:int(11) unsigned;not null;default:0" json:"target_column_id"` // 目标表字段id
	Type           uint8  `gorm:"column:type;type:tinyint(3) unsigned;not null;default:1" json:"type"`                      // 关联类型.11对1，2:1对多,3:多对1,4:多对多
}

// TableName get sql table name.获取数据库表名
func (m *DatasourceColumnRelation) TableName() string {
	return "datasource_column_relation"
}

var relationTablename = (&DatasourceColumnRelation{}).TableName()

type RelationList []DatasourceColumnRelation

func getRelationList(db *gorm.DB, sourceId uint64) (error, RelationList) {
	//var count uint32
	list := make(RelationList, 0)
	err := db.Table(relationTablename).
		Where(DatasourceColumnRelationColumns.SourceTableID+" = ?", sourceId).
		Find(&list).
		Error
	return err, list
}

func (list RelationList) SourceColumns() []uint64 {
	uList := make([]uint64, 0, cap(list))
	for _, x := range list {
		uList = append(uList, x.SourceColumnID)
	}
	return uList
}

func (list RelationList) TargetTables() []uint64 {
	uList := make([]uint64, 0, cap(list))
	for _, x := range list {
		uList = append(uList, x.TargetTableID)
	}
	return uList
}

func (list RelationList) TargetColumns() []uint64 {
	uList := make([]uint64, 0, cap(list))
	for _, x := range list {
		uList = append(uList, x.TargetColumnID)
	}
	return uList
}
