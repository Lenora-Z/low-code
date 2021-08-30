//Created by Goland
//@User: lenora
//@Date: 2021/3/11
//@Time: 2:58 下午
package service

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Service struct {
	ID        uint64    `gorm:"primary_key;column:id;type:int(11) unsigned;not null" json:"id"`
	Name      string    `gorm:"column:name;type:varchar(255);not null" json:"name"`                  // 服务名称
	Type      int8      `gorm:"column:type;type:tinyint(3) unsigned;not null;default:1" json:"type"` // 服务类型 1-上链服务 2-发送邮件 3-代码包
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime" json:"updated_at"`
}

// TableName get sql table name.获取数据库表名
func (m *Service) TableName() string {
	return "service"
}

var srvTablename = (&Service{}).TableName()

//func create(db *gorm.DB, m *Sample) (error, *Sample) {
//	err := db.Table(tablename).Create(m).Error
//	return err, m
//}
//
//func detail(db *gorm.DB, id uint64) (bool, *Sample) {
//	var item Sample
//	status := db.Table(tablename).Where("id = ?", id).First(&item).RecordNotFound()
//	return status, &item
//}
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
//func update(db *gorm.DB, m *Sample) (error, *Sample) {
//	m.UpdatedAt = time.Now()
//	err := db.Table(tablename).Where("id = ?", m.ID).Update(m).Error
//	return err, m
//}
//

type ServiceList []Service

func getServiceList(db *gorm.DB, t []int8) (error, ServiceList) {
	list := make(ServiceList, 0)
	query := db.Table(srvTablename)
	if len(t) > 0 {
		query = query.Where("type in (?)", t)
	}
	err := query.Find(&list).
		Error
	return err, list
}

func (list ServiceList) Ids() []uint64 {
	sList := make([]uint64, 0, cap(list))
	for _, x := range list {
		sList = append(sList, x.ID)
	}
	return sList
}
