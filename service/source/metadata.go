// Package source
// Created by GoLand
// @User: lenora
// @Date: 2021/8/3
// @Time: 14:37

package source

import "github.com/jinzhu/gorm"

type DatasourceMetadata struct {
	ID                 uint64 `gorm:"primaryKey;column:id;type:int(11) unsigned;not null" json:"id"`
	AppID              uint64 `gorm:"index:app_id_idx;column:app_id;type:int(11) unsigned;not null;default:0" json:"app_id"`            // 应用id
	DatasourceColumnID uint64 `gorm:"column:datasource_column_id;type:int(11) unsigned;not null;default:0" json:"datasource_column_id"` // 表字段id
	Group              string `gorm:"column:group;type:varchar(64);not null;default:''" json:"group"`                                   // 分组组名
	Key                string `gorm:"column:key;type:varchar(64);not null;default:''" json:"key"`                                       // key键
	Value              string `gorm:"column:value;type:varchar(64);not null;default:''" json:"value"`                                   // value值
}

// TableName get sql table name.获取数据库表名
func (m *DatasourceMetadata) TableName() string {
	return "datasource_metadata"
}

var metaTablename = (&DatasourceMetadata{}).TableName()

func getMetaList(db *gorm.DB, columnId uint64) (error, []DatasourceMetadata) {
	list := make([]DatasourceMetadata, 0)
	err := db.Table(metaTablename).
		Where(DatasourceMetadataColumns.DatasourceColumnID+" = ?", columnId).
		Find(&list).
		Error
	return err, list
}
