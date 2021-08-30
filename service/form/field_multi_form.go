package form

import "github.com/jinzhu/gorm"

type FieldMultiForm struct {
	FormID  uint64 `gorm:"primaryKey;column:form_id;type:int(11) unsigned;not null;default:0" json:"form_id"`               // 父表单id
	FieldID uint64 `gorm:"primaryKey;index:field;column:field_id;type:int(11) unsigned;not null;default:0" json:"field_id"` // 控件id
	ChildID uint64 `gorm:"column:child_id;type:int(11) unsigned;not null;default:0" json:"child_id"`
	Method  uint8  `gorm:"column:method;type:tinyint(1) unsigned;not null;default:1" json:"method"` // 表单引入方式 1-关联 2-原表复制
	Mode    uint8  `gorm:"column:mode;type:tinyint(1) unsigned;not null;default:1" json:"mode"`     // 呈现样式 1-单条记录 2-多条记录 3-表格
}

// TableName get sql table name.获取数据库表名
func (m *FieldMultiForm) TableName() string {
	return "field_multi_form"
}

var multiFormTablename = (&FieldMultiForm{}).TableName()

func createMultiForm(db *gorm.DB, m FieldMultiForm) (error, *FieldMultiForm) {
	err := db.Table(multiFormTablename).Create(&m).Error
	return err, &m
}

func deleteMultiForm(db *gorm.DB, formId uint64) error {
	err := db.Table(multiFormTablename).Delete(&FieldMultiForm{}, "form_id = ?", formId).Error
	return err
}

type MultiFormList []FieldMultiForm

func getRelationList(db *gorm.DB, id uint64) (error, MultiFormList) {
	list := make(MultiFormList, 0)
	query := db.Table(multiFormTablename).Where("form_id = ?", id)
	err := query.Find(&list).Error
	return err, list
}

func (list MultiFormList) Children() []uint64 {
	ret := make([]uint64, 0, cap(list))
	for _, x := range list {
		ret = append(ret, x.ChildID)
	}
	return ret
}

func (list MultiFormList) Fields() []uint64 {
	ret := make([]uint64, 0, cap(list))
	for _, x := range list {
		ret = append(ret, x.FieldID)
	}
	return ret
}
