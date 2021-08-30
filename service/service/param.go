//Created by Goland
//@User: lenora
//@Date: 2021/3/11
//@Time: 2:59 下午
package service

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Param struct {
	ID           uint64    `gorm:"primary_key;column:id;type:int(11) unsigned;not null" json:"id"`
	Name         string    `gorm:"column:name;type:varchar(255);not null" json:"name"`                     // 参数名称
	ServiceID    uint64    `gorm:"column:service_id;type:int(11) unsigned;not null" json:"service_id"`     // 服务id
	Check        *bool     `gorm:"column:check;type:tinyint(1) unsigned;not null;default:1" json:"check"`  // 是否必填
	Mode         uint8     `gorm:"column:mode;type:tinyint(1) unsigned;not null;default:1" json:"mode"`    // 入/出参 1-入 2-出
	Type         uint8     `gorm:"column:type;type:tinyint(3) unsigned;not null;default:1" json:"type"`    // 参数类型 1-string 2-uint 3-int 4-float 5-double 6-list
	DefaultValue string    `gorm:"column:default_value;type:varchar(255);default:''" json:"default_value"` // 默认值
	FixedValue   string    `gorm:"column:fixed_value;type:varchar(255);default:''" json:"fixed_value"`     // 固定值
	IsVisual     *bool     `gorm:"column:is_visual;type:tinyint(1);default:1" json:"is_visual"`            // 配置端是否可见 0-不可见 1-可见
	CreatedAt    time.Time `gorm:"column:created_at;type:datetime;not null" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;type:datetime" json:"updated_at"`
}

// TableName get sql table name.获取数据库表名
func (m *Param) TableName() string {
	return "param"
}

var paramTablename = (&Param{}).TableName()

//func create(db *gorm.DB, m *Sample) (error, *Sample) {
//	err := db.Table(tablename).Create(m).Error
//	return err, m
//}
//
func paramDetail(db *gorm.DB, id uint64) (bool, *Param) {
	var item Param
	status := db.Table(paramTablename).Where("id = ?", id).First(&item).RecordNotFound()
	return status, &item
}

//
//func detailByColumn(db *gorm.DB, group map[string]interface{}) (bool, *Sample) {
//	var item Sample
//	query := db.Table(tablename)
//	for i, x := range group {
//		query = query.Where(i+" = ?", x)
//	}
//	status := query.First(&item).RecordNotFound()
//	return status, &item
//}
//
func getParamList(db *gorm.DB, srvId []uint64, mode uint8) (error, []Param) {
	//var count uint32
	list := make([]Param, 0)
	query := db.Table(paramTablename)
	query = query.Where("is_visual = 1")
	if len(srvId) > 0 && srvId[0] != 0 {
		query = query.Where("service_id in (?)", srvId)
	}
	if mode > 0 {
		query = query.Where("mode = ?", mode)
	}
	err := query.Find(&list).
		Error
	return err, list
}
