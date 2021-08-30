package formdata

import "github.com/jinzhu/gorm"

// FieldLinkage 级联控件属性表
type FieldLinkage struct {
	ID      uint32 `gorm:"primaryKey;column:id;type:int(11) unsigned;not null;default:0" json:"id"`
	FormID  uint32 `gorm:"column:form_id;type:int(11) unsigned;not null;default:0" json:"formId"`   // 表单id
	FieldID uint32 `gorm:"column:field_id;type:int(11) unsigned;not null;default:0" json:"fieldId"` // 级联控件id
	Content string `gorm:"column:content;type:text" json:"content"`                                 // 级联关系,格式:表字段关联关系id#表字段id#显示文案,以@连接
}

// TableName get sql table name.获取数据库表名
func (f *FieldLinkage) TableName() string {
	return "field_linkage"
}

var fieldLinkageTableName = (&FieldLinkage{}).TableName()

func getFieldLinkageByFieldId(db *gorm.DB, id uint32) (bool, *FieldLinkage) {
	var item FieldLinkage
	isNotFound := db.Table(fieldLinkageTableName).Where("field_id = ?", id).First(&item).RecordNotFound()
	return isNotFound, &item
}
