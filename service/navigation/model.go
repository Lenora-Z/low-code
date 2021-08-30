package navigation

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Navigation struct {
	ID        uint64    `gorm:"primaryKey;column:id;type:int(11) unsigned;not null" json:"id"`
	Name      string    `gorm:"column:name;type:varchar(255);not null" json:"name"`                  // 导航名称
	AppID     uint64    `gorm:"column:app_id;type:int(11) unsigned;not null" json:"app_id"`          // 应用id
	Number    string    `gorm:"column:number;type:varchar(100);not null" json:"number"`              // 编号
	Desc      string    `gorm:"column:desc;type:varchar(255)" json:"desc"`                           // 导航描述
	Content   []byte    `gorm:"column:content;type:text" json:"content"`                             // 导航内容
	Status    *bool     `gorm:"column:status;type:tinyint(1);not null" json:"status"`                // 表单状态 0-未生效 1-已生效
	IsOnline  *bool     `gorm:"column:is_online;type:tinyint(1) unsigned;not null" json:"is_online"` // 是否是上线表单 1-是
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime" json:"updated_at"`
}

// TableName get sql table name.获取数据库表名
func (m *Navigation) TableName() string {
	return "navigation"
}

var tablename = (&Navigation{}).TableName()

func createNavigation(db *gorm.DB, m Navigation) (error, *Navigation) {
	err := db.Table(tablename).Create(&m).Error
	return err, &m
}

func navigationDetail(db *gorm.DB, id uint64) (bool, *Navigation) {
	var item Navigation
	status := db.Table(tablename).Where("id = ?", id).First(&item).RecordNotFound()
	return status, &item
}

func navigationDetailByColumn(db *gorm.DB, group map[string]interface{}) (bool, *Navigation) {
	var item Navigation
	query := db.Table(tablename)
	for i, x := range group {
		query = query.Where(i+" = ?", x)
	}
	status := query.First(&item).RecordNotFound()
	return status, &item
}

func updateNavigation(db *gorm.DB, m Navigation) (error, *Navigation) {
	m.UpdatedAt = time.Now()
	err := db.Table(tablename).Where("id = ?", m.ID).Update(&m).Error
	return err, &m
}

func batchUpdateNavigation(db *gorm.DB, where, change map[string]interface{}) error {
	query := db.Table(tablename)
	for i, x := range where {
		query = query.Where(i+" = ?", x)
	}
	err := query.Updates(change).Error
	return err
}

func getNavigationList(db *gorm.DB, offset, limit uint32, appId uint64) (error, uint32, []Navigation) {
	var count uint32
	list := make([]Navigation, 0)
	query := db.Table(tablename).Where("app_id = ?", appId)
	err := query.Count(&count).
		Offset(offset).
		Limit(limit).
		Order("updated_at desc").
		Find(&list).
		Error
	return err, count, list
}
