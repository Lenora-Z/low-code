package formdata

import "github.com/jinzhu/gorm"

// FieldButton 按钮控件表
type FieldButton struct {
	ID      uint32 `gorm:"primaryKey;column:id;type:int(11) unsigned;not null" json:"id"`
	FormID  uint32 `gorm:"column:form_id;type:int(11) unsigned;not null;default:0" json:"formId"`   // 表单id
	FieldID uint32 `gorm:"column:field_id;type:int(11) unsigned;not null;default:0" json:"fieldId"` // 控件id
	FlowID  uint32 `gorm:"column:flow_id;type:int(11) unsigned;not null;default:0" json:"flowId"`   // 流程id
	Name    string `gorm:"column:name;type:varchar(255);not null;default:''" json:"name"`           // 按钮名称
	Event   int8   `gorm:"column:event;type:tinyint(3);not null;default:1" json:"event"`            // 触发事件类型
}

// TableName get sql table name.获取数据库表名
func (f *FieldButton) TableName() string {
	return "field_button"
}

var fieldButtonTableName = (&FieldButton{}).TableName()

func getFieldButtonByFieldId(db *gorm.DB, id uint32) (bool, *FieldButton) {
	var item FieldButton
	isNotFound := db.Table(fieldButtonTableName).Where("field_id = ?", id).First(&item).RecordNotFound()
	return isNotFound, &item
}
