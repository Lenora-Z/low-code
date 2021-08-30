package formdata

import "github.com/jinzhu/gorm"

// FieldRecords 关联控件表
type FieldRecords struct {
	ID                         uint32 `gorm:"primaryKey;column:id;type:int(11) unsigned;not null" json:"id"`
	FormID                     uint32 `gorm:"column:form_id;type:int(11) unsigned;not null;default:0" json:"formId"`                                           // 表单id
	FieldID                    uint32 `gorm:"column:field_id;type:int(11) unsigned;not null;default:0" json:"fieldId"`                                         // 控件id
	DatasourceColumnRelationID uint32 `gorm:"column:datasource_column_relation_id;type:int(11) unsigned;not null;default:0" json:"datasourceColumnRelationId"` // 表字段关联关系id
	DatasourceColumnIDs        string `gorm:"column:datasource_column_ids;type:varchar(255);not null;default:''" json:"datasourceColumnIds"`                   // 关联字段id列表,逗号分分割
	Mode                       int8   `gorm:"column:mode;type:tinyint(3);not null;default:1" json:"mode"`                                                      // 呈现方式 1-卡片 2-回填
	CountType                  int8   `gorm:"column:count_type;type:tinyint(3);not null;default:1" json:"countType"`                                           // 关联记录数量 1-单条 2-多条
	DetailStatus               bool   `gorm:"column:detail_status;type:tinyint(1) unsigned;not null;default:0" json:"detailStatus"`                            // 详情是否展示 0-关闭 1-开启
	Columns                    string `gorm:"column:columns;type:varchar(255);not null;default:''" json:"columns"`                                             // 显示字段,关联字段id列表,逗号分分割
}

// TableName get sql table name.获取数据库表名
func (f *FieldRecords) TableName() string {
	return "field_records"
}

var fieldRecordsTableName = (&FieldRecords{}).TableName()

func getFieldRecordsByFieldId(db *gorm.DB, fieldId uint32) (bool, *FieldRecords) {
	var item FieldRecords
	isNotFound := db.Table(fieldRecordsTableName).Where("field_id = ?", fieldId).First(&item).RecordNotFound()
	return isNotFound, &item
}
