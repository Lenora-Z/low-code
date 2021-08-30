//Created by Goland
//@User: lenora
//@Date: 2021/3/12
//@Time: 10:21 上午
package form

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Field struct {
	ID                 uint64    `gorm:"primary_key;column:id;type:int(11) unsigned;not null" json:"id"`
	Key                string    `gorm:"uniqueIndex:form_field;index:keys;column:key;type:varchar(100);not null" json:"key"`               // 控件id
	Title              string    `gorm:"column:title;type:varchar(255);not null" json:"title"`                                             // 控件名称
	Type               string    `gorm:"column:type;type:varchar(20);not null" json:"type"`                                                // 控件类型
	FormID             uint64    `gorm:"unique_index:form_field;column:form_id;type:int(11);not null" json:"form_id"`                      // 表单id                       // 表单id
	DatasourceColumnID uint64    `gorm:"column:datasource_column_id;type:int(11) unsigned;not null;default:0" json:"datasource_column_id"` // 数据表字段id
	IsOnly             *bool     `gorm:"column:is_only;type:tinyint(1);not null" json:"is_only"`                                           // 是否唯一
	IsNecessary        *bool     `gorm:"column:is_necessary;type:tinyint(1);not null" json:"is_necessary"`                                 // 是否必填
	Content            []byte    `gorm:"column:content;type:text" json:"content"`                                                          // 控件内容json
	CreatedAt          time.Time `gorm:"column:created_at;type:datetime;not null" json:"created_at"`
	UpdatedAt          time.Time `gorm:"column:updated_at;type:datetime" json:"updated_at"`
}

// fieldTablename get sql table name.获取数据库表名
func (m *Field) fieldTablename() string {
	return "field"
}

var fieldTablename = (&Field{}).fieldTablename()

func createField(db *gorm.DB, m Field) (error, *Field) {
	err := db.Table(fieldTablename).Create(&m).Error
	return err, &m
}

func deleteField(db *gorm.DB, formId uint64) error {
	err := db.Table(fieldTablename).Delete(&Field{}, "form_id = ?", formId).Error
	return err
}

func fieldDetail(db *gorm.DB, id uint64) (bool, *Field) {
	var item Field
	status := db.Table(fieldTablename).Where("id = ?", id).First(&item).RecordNotFound()
	return status, &item
}

func fieldDetailByColumn(db *gorm.DB, group map[string]interface{}) (bool, *Field) {
	var item Field
	query := db.Table(fieldTablename)
	for i, x := range group {
		query = query.Where(i+" = ?", x)
	}
	status := query.First(&item).RecordNotFound()
	return status, &item
}

func getFieldList(db *gorm.DB, formId, ids []uint64, keys ...string) (error, []Field) {
	list := make([]Field, 0)
	query := db.Table(fieldTablename)
	if len(formId) > 0 {
		query = query.Where("form_id in (?)", formId)
	}
	if len(ids) > 0 {
		query = query.Where("id in (?)", ids)
	}
	if len(keys) > 0 {
		query = query.Where("`key` in (?)", keys)
	}
	err := query.Find(&list).Error
	return err, list
}
