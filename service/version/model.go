//Created by Goland
//@User: lenora
//@Date: 2021/3/11
//@Time: 3:35 下午
package version

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Version struct {
	ID        uint64    `gorm:"primary_key;column:id;type:int(10) unsigned;not null" json:"id"`
	AppID     uint64    `gorm:"index:app;column:app_id;type:int(11) unsigned;not null" json:"app_id"` // 应用id
	Domain    string    `gorm:"column:domain;type:varchar(255)" json:"domain"`                        // 应用域名
	Version   string    `gorm:"column:version;type:varchar(255);not null" json:"version"`             // 版本号
	Status    int8      `gorm:"column:status;type:tinyint(255);not null" json:"status"`               // 状态值
	Note      string    `gorm:"column:note;type:text" json:"note"`                                    // 版本介绍
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime" json:"updated_at"`
}

// TableName get sql table name.获取数据库表名
func (m *Version) TableName() string {
	return "version"
}

var tablename = (&Version{}).TableName()

func createVersion(db *gorm.DB, m Version) (error, *Version) {
	err := db.Table(tablename).Create(&m).Error
	return err, &m
}

//
//func detail(db *gorm.DB, id uint64) (bool, *Sample) {
//	var item Sample
//	status := db.Table(tablename).Where("id = ?", id).First(&item).RecordNotFound()
//	return status, &item
//}
//
func versionDetailByColumn(db *gorm.DB, group map[string]interface{}) (bool, Version) {
	var item Version
	query := db.Table(tablename)
	for i, x := range group {
		query = query.Where(i+" = ?", x)
	}
	status := query.First(&item).RecordNotFound()
	return status, item
}

//
func updateVersion(db *gorm.DB, m Version) (error, Version) {
	m.UpdatedAt = time.Now()
	err := db.Table(tablename).Where("id = ?", m.ID).Update(&m).First(&m).Error
	return err, m
}

func batchUpdateVersionByColumn(db *gorm.DB, items, where map[string]interface{}) error {
	query := db.Table(tablename)
	for i, x := range where {
		query = query.Where(i+" = ?", x)
	}
	err := query.Updates(items).Error
	return err
}

func getVersionList(db *gorm.DB, offset, limit uint32, appId uint64, status ...int8) (error, uint32, []Version) {
	var count uint32
	list := make([]Version, 0)
	query := db.Table(tablename).Where("app_id = ?", appId)
	if len(status) > 0 {
		query = query.Where("status in (?)", status)
	}
	err := query.Count(&count).
		Offset(offset).
		Limit(limit).
		Find(&list).
		Error
	return err, count, list
}
