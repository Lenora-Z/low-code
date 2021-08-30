package formdata

import (
	"github.com/jinzhu/gorm"
	"time"
)

type DatasourceTable struct {
	ID           uint32    `gorm:"primaryKey;column:id;type:int(11) unsigned;not null" json:"id"`
	CreatedAt    time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"createdAt"` // 创建时间
	AppID        uint32    `gorm:"column:app_id;type:int(11) unsigned;not null;default:0" json:"appId"`                 // 应用id
	SchemataID   uint32    `gorm:"column:schemata_id;type:int(11) unsigned;not null;default:0" json:"schemataId"`       // 数据源id
	TableSchema  string    `gorm:"column:table_schema;type:varchar(64);not null;default:''" json:"tableSchema"`         // 数据库名称
	TableName    string    `gorm:"column:table_name;type:varchar(64);not null;default:''" json:"tableName"`             // 数据表名称
	TableComment string    `gorm:"column:table_comment;type:varchar(2048);not null;default:''" json:"tableComment"`     // 数据库表备注
}

// TableName get sql table name.获取数据库表名
func (f *DatasourceTable) GetTableName() string {
	return "datasource_table"
}

var datasourceTableTableName = (&DatasourceTable{}).GetTableName()

func getDatasourceTableById(db *gorm.DB, id uint32) (bool, *DatasourceTable) {
	var item DatasourceTable
	isNotFound := db.Table(datasourceTableTableName).Where("id = ?", id).First(&item).RecordNotFound()
	return isNotFound, &item
}
