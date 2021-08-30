// Package form
// Created by GoLand
// @User: lenora
// @Date: 2021/7/26
// @Time: 15:25

package form

import "github.com/jinzhu/gorm"

type FieldButton struct {
	ID      uint64 `gorm:"primaryKey;column:id;type:int(11) unsigned;not null" json:"id"`
	FormID  uint64 `gorm:"index:form_field_idx;column:form_id;type:int(11) unsigned;not null;default:0" json:"form_id"`   // 表单id
	FieldID uint64 `gorm:"index:form_field_idx;column:field_id;type:int(11) unsigned;not null;default:0" json:"field_id"` // 控件id
	FlowID  uint64 `gorm:"column:flow_id;type:int(11) unsigned;not null;default:0" json:"flow_id"`                        // 流程id
	Name    string `gorm:"column:name;type:varchar(255);not null;default:''" json:"name"`                                 // 按钮名称
	Event   uint8  `gorm:"column:event;type:tinyint(3);not null;default:1" json:"event"`                                  // 触发事件类型
}

// TableName get sql table name.获取数据库表名
func (m *FieldButton) TableName() string {
	return "field_button"
}

var buttonTableName = (&FieldButton{}).TableName()

func createButton(db *gorm.DB, m FieldButton) (error, *FieldButton) {
	err := db.Table(buttonTableName).Create(&m).Error
	return err, &m
}

func deleteButton(db *gorm.DB, formId uint64) error {
	err := db.Table(buttonTableName).Delete(&FieldButton{}, "form_id = ?", formId).Error
	return err
}

type FieldButtons []FieldButton

func getButtonListWithCondition(db *gorm.DB, cond map[string]interface{}) (error, FieldButtons) {
	list := make(FieldButtons, 0)
	query := db.Table(buttonTableName)
	for i, x := range cond {
		query = query.Where(i+"= ?", x)
	}
	err := query.Find(&list).Error
	return err, list
}

func (list FieldButtons) Forms() []uint64 {
	uList := make([]uint64, 0, cap(list))
	for _, x := range list {
		uList = append(uList, x.FormID)
	}
	return uList
}

func (list FieldButtons) Fields() []uint64 {
	uList := make([]uint64, 0, cap(list))
	for _, x := range list {
		uList = append(uList, x.FieldID)
	}
	return uList
}
