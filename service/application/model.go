//Created by Goland
//@User: lenora
//@Date: 2021/3/8
//@Time: 1:57 下午
package application

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Application struct {
	ID        uint64    `gorm:"primary_key;column:id;type:int(11) unsigned;not null" json:"id"`
	Name      string    `gorm:"column:name;type:varchar(255);not null" json:"name"`                // 应用名称
	Icon      string    `gorm:"column:icon;type:varchar(50);not null" json:"icon"`                 // icon
	Desc      string    `gorm:"column:desc;type:varchar(255)" json:"desc"`                         // 应用描述
	Status    int8      `gorm:"column:status;type:tinyint(1) unsigned;not null" json:"status"`     // 应用状态 0-未发布 1-上线 2-下线
	VersionID uint64    `gorm:"column:version_id;type:int(5) unsigned;not null" json:"version_id"` // 当前版本id
	AppHash   string    `gorm:"column:app_hash;type:varchar(255);not null" json:"app_hash"`        // 应用hash
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime" json:"updated_at"`
}

// TableName get sql table name.获取数据库表名
func (m *Application) TableName() string {
	return "application"
}

var tablename = (&Application{}).TableName()

func createApp(db *gorm.DB, m *Application) (error, *Application) {
	err := db.Table(tablename).Create(m).Error
	return err, m
}

func appDetail(db *gorm.DB, id uint64) (bool, *Application) {
	var item Application
	status := db.Table(tablename).Where("id = ?", id).First(&item).RecordNotFound()
	return status, &item
}

func appDetailByColumn(db *gorm.DB, group map[string]interface{}) (bool, *Application) {
	var item Application
	query := db.Table(tablename)
	for i, x := range group {
		query = query.Where(i+" = ?", x)
	}
	status := query.First(&item).RecordNotFound()
	return status, &item
}

func getAppList(db *gorm.DB, offset, limit uint32) (error, uint32, []Application) {
	var count uint32
	list := make([]Application, 0)
	query := db.Table(tablename)
	err := query.Count(&count).
		Offset(offset).
		Limit(limit).
		Find(&list).
		Error
	return err, count, list
}

func updateApp(db *gorm.DB, m *Application) (error, *Application) {
	m.UpdatedAt = time.Now()
	err := db.Table(tablename).Where("id = ?", m.ID).Updates(m).Error
	return err, m
}
