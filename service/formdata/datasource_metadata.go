package formdata

import (
	"github.com/jinzhu/gorm"
)

type DatasourceMetadata struct {
	ID                 uint32 `gorm:"primaryKey;column:id;type:int(11) unsigned;not null" json:"id"`
	AppID              uint32 `gorm:"column:app_id;type:int(11) unsigned;not null;default:0" json:"appId"`                            // 应用id
	DatasourceColumnID uint32 `gorm:"column:datasource_column_id;type:int(11) unsigned;not null;default:0" json:"datasourceColumnId"` // 表字段id
	Group              string `gorm:"column:group;type:varchar(64);not null;default:''" json:"group"`                                 // 分组组名
	Key                string `gorm:"column:key;type:varchar(64);not null;default:''" json:"key"`                                     // key键
	Value              string `gorm:"column:value;type:varchar(64);not null;default:''" json:"value"`                                 // value值
}

// TableName get sql table name.获取数据库表名
func (f *DatasourceMetadata) TableName() string {
	return "datasource_metadata"
}

var datasourceMetadataTableName = (&DatasourceMetadata{}).TableName()

func getDatasourceMetadataById(db *gorm.DB, id uint32) (bool, []DatasourceMetadata) {
	var items []DatasourceMetadata
	isNotFound := db.Table(datasourceMetadataTableName).Where("datasource_column_id = ?", id).Find(&items).RecordNotFound()
	return isNotFound, items
}

func getDatasourceMetadataByKey(db *gorm.DB, id uint32, key string) (bool, string) {
	var items DatasourceMetadata
	isNotFound := db.Table(datasourceMetadataTableName).
		Where("datasource_column_id = ?", id).Where("`key` = ?", key).First(&items).RecordNotFound()
	return isNotFound, items.Value
}

func getDatasourceMetadataByValue(db *gorm.DB, id uint32, value string) (bool, string) {
	var items DatasourceMetadata
	isNotFound := db.Table(datasourceMetadataTableName).
		Where("datasource_column_id = ?", id).Where("`value` = ?", value).First(&items).RecordNotFound()
	return isNotFound, items.Key
}
