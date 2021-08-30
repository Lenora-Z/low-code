package formdata

import "github.com/jinzhu/gorm"

type FieldTableColumn struct {
	ID                         uint32 `gorm:"primaryKey;column:id;type:int(11) unsigned;not null" json:"id"`
	FormID                     uint32 `gorm:"column:form_id;type:int(11) unsigned;not null;default:0" json:"formId"`                                           // 表单id
	FieldID                    uint32 `gorm:"column:field_id;type:int(11) unsigned;not null;default:0" json:"fieldId"`                                         // 控件id
	DatasourceColumnRelationID uint32 `gorm:"column:datasource_column_relation_id;type:int(11) unsigned;not null;default:0" json:"datasourceColumnRelationId"` // 表字段关联表id,0表示本数据表字段
	DatasourceColumnID         uint32 `gorm:"column:datasource_column_id;type:int(11) unsigned;not null;default:0" json:"datasourceColumnId"`                  // 表字段id
	ShowName                   string `gorm:"column:show_name;type:varchar(255);not null;default:''" json:"showName"`                                          // 表头显示名称
	IsCondition                bool   `gorm:"column:is_condition;type:tinyint(1) unsigned;not null;default:0" json:"isCondition"`                              // 是否作为筛选条件 0-关闭 1-开启
}

// TableName get sql table name.获取数据库表名
func (f *FieldTableColumn) TableName() string {
	return "field_table_column"
}

var fieldTableColumnTableName = (&FieldTableColumn{}).TableName()

func getFieldTableColumnByFieldId(db *gorm.DB, fieldId uint32) (bool, FieldTableColumnList) {
	var item FieldTableColumnList
	isNotFound := db.Table(fieldTableColumnTableName).Where("field_id = ?", fieldId).Order("datasource_column_id asc").Find(&item).RecordNotFound()
	return isNotFound, item
}
