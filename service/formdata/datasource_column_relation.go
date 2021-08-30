package formdata

import "github.com/jinzhu/gorm"

type DatasourceColumnRelation struct {
	ID             uint32 `gorm:"primaryKey;column:id;type:int(11) unsigned;not null" json:"id"`
	SourceTableID  uint32 `gorm:"column:source_table_id;type:int(11) unsigned;not null;default:0" json:"sourceTableId"`   // 源数据表id
	SourceColumnID uint32 `gorm:"column:source_column_id;type:int(11) unsigned;not null;default:0" json:"sourceColumnId"` // 源表字段id
	TargetTableID  uint32 `gorm:"column:target_table_id;type:int(11) unsigned;not null;default:0" json:"targetTableId"`   // 目标数据表id
	TargetColumnID uint32 `gorm:"column:target_column_id;type:int(11) unsigned;not null;default:0" json:"targetColumnId"` // 目标表字段id
	Type           uint8  `gorm:"column:type;type:tinyint(3) unsigned;not null;default:1" json:"type"`                    // 关联类型.1:1对1，2:1对多,3:多对1,4:多对多
}

// TableName get sql table name.获取数据库表名
func (f *DatasourceColumnRelation) TableName() string {
	return "datasource_column_relation"
}

var datasourceColumnRelationTableName = (&DatasourceColumnRelation{}).TableName()

func getDatasourceColumnRelationById(db *gorm.DB, id uint32) (bool, *DatasourceColumnRelation) {
	var item DatasourceColumnRelation
	isNotFound := db.Table(datasourceColumnRelationTableName).Where("id = ?", id).First(&item).RecordNotFound()
	return isNotFound, &item
}
