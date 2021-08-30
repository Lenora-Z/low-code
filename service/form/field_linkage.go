package form

import "github.com/jinzhu/gorm"

type FieldLinkage struct {
	ID      uint64 `gorm:"primaryKey;column:id;type:int(11) unsigned;not null;default:0" json:"id"`
	FormID  uint64 `gorm:"index:form_id;column:form_id;type:int(11) unsigned;not null;default:0" json:"form_id"` // 表单id
	FieldID uint64 `gorm:"column:field_id;type:int(11) unsigned;not null;default:0" json:"field_id"`             // 级联控件id
	Content string `gorm:"column:content;type:text" json:"content"`                                              // 级联关系,格式:表单id#控件key,以&连接
}

// TableName get sql table name.获取数据库表名
func (m *FieldLinkage) TableName() string {
	return "field_linkage"
}

var linkageTablename = (&FieldLinkage{}).TableName()

func createLinkage(db *gorm.DB, m FieldLinkage) (error, *FieldLinkage) {
	err := db.Table(linkageTablename).Create(&m).Error
	return err, &m
}

func linkageDetailByColumn(db *gorm.DB, group map[string]interface{}) (bool, *FieldLinkage) {
	var item FieldLinkage
	query := db.Table(linkageTablename)
	for i, x := range group {
		query = query.Where(i+" = ?", x)
	}
	status := query.First(&item).RecordNotFound()
	return status, &item
}

func deleteLinkage(db *gorm.DB, formId uint64) error {
	err := db.Table(linkageTablename).Delete(&Field{}, "form_id = ?", formId).Error
	return err
}
