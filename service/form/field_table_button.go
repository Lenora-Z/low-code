// Package form
// Created by GoLand
// @User: lenora
// @Date: 2021/7/26
// @Time: 15:25

package form

import "github.com/jinzhu/gorm"

type FieldTableButton struct {
	ID      uint64 `gorm:"primaryKey;column:id;type:int(11) unsigned;not null" json:"id"`
	FormID  uint64 `gorm:"index:form_field_idx;column:form_id;type:int(11) unsigned;not null;default:0" json:"form_id"`   // 表单id
	FieldID uint64 `gorm:"index:form_field_idx;column:field_id;type:int(11) unsigned;not null;default:0" json:"field_id"` // 控件id
	FlowID  uint64 `gorm:"column:flow_id;type:int(11) unsigned;not null;default:0" json:"flow_id"`                        // 流程id
	Name    string `gorm:"column:name;type:varchar(255);not null;default:''" json:"name"`                                 // 按钮名称
	Event   uint8  `gorm:"column:event;type:tinyint(3);not null;default:1" json:"event"`                                  // 触发事件类型
}

// TableName get sql table name.获取数据库表名
func (m *FieldTableButton) TableName() string {
	return "field_table_button"
}

var tableButtonTableName = (&FieldTableButton{}).TableName()

func createTableButton(db *gorm.DB, m FieldTableButton) (error, *FieldTableButton) {
	err := db.Table(tableButtonTableName).Create(&m).Error
	return err, &m
}

func deleteTableButton(db *gorm.DB, formId uint64) error {
	err := db.Table(tableButtonTableName).Delete(&FieldTableButton{}, "form_id = ?", formId).Error
	return err
}

func getTableButtonListWithCondition(db *gorm.DB, cond map[string]interface{}) (error, FieldButtons) {
	list := make(FieldButtons, 0)
	query := db.Table(tableButtonTableName)
	for i, x := range cond {
		query = query.Where(i+"= ?", x)
	}
	err := query.Scan(&list).Error
	return err, list
}
