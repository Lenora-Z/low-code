package formdata

import "github.com/jinzhu/gorm"

type FieldTable struct {
	ID                uint32 `gorm:"primaryKey;column:id;type:int(11) unsigned;not null" json:"id"`
	FormID            uint32 `gorm:"column:form_id;type:int(11) unsigned;not null;default:0" json:"formId"`                        // 表单id
	FieldID           uint32 `gorm:"column:field_id;type:int(11) unsigned;not null;default:0" json:"fieldId"`                      // 控件id
	DatasourceTableID uint32 `gorm:"column:datasource_table_id;type:int(11) unsigned;not null;default:0" json:"datasourceTableId"` // 数据表id
	IsExport          bool   `gorm:"column:is_export;type:tinyint(1) unsigned;not null;default:0" json:"isExport"`                 // 导出配置 0: 关闭导出 1:允许导出
	IsFilter          bool   `gorm:"column:is_filter;type:tinyint(1) unsigned;not null;default:0" json:"isFilter"`                 // 数据过滤配置  0:关闭过滤  1:允许过滤
}

// TableName get sql table name.获取数据库表名
func (f *FieldTable) TableName() string {
	return "field_table"
}

var fieldTableTableName = (&FieldTable{}).TableName()

func getFieldTableByFieldId(db *gorm.DB, fieldId uint32) (bool, *FieldTable) {
	var item FieldTable
	isNotFound := db.Table(fieldTableTableName).Where("field_id = ?", fieldId).First(&item).RecordNotFound()
	return isNotFound, &item
}
