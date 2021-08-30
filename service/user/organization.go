//Package user
//Created by Goland
//@User: lenora
//@Date: 2021/3/9
//@Time: 2:16 下午
package user

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Organization struct {
	ID        uint64    `gorm:"primary_key;column:id;type:int(11) unsigned;not null" json:"id"`
	AppID     uint64    `gorm:"index:app;column:app_id;type:int(11) unsigned;not null" json:"app_id"` // 应用id
	Name      string    `gorm:"column:name;type:varchar(255);not null" json:"name"`                   // 组织名称
	ParentID  uint64    `gorm:"column:parent_id;type:int(11) unsigned;not null" json:"parent_id"`     // 副组织id
	IsDelete  *bool     `gorm:"column:is_delete;type:tinyint(1) unsigned;not null" json:"is_delete"`  // 是否删除
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime" json:"updated_at"`
}

// TableName get sql table name.获取数据库表名
func (m *Organization) TableName() string {
	return "organization"
}

var orgTablename = (&Organization{}).TableName()

func createOrganization(db *gorm.DB, m *Organization) (error, *Organization) {
	err := db.Table(orgTablename).Create(m).Error
	return err, m
}

func organizationDetail(db *gorm.DB, id uint64) (bool, *Organization) {
	var item Organization
	status := db.Table(orgTablename).Where("id = ?", id).First(&item).RecordNotFound()
	return status, &item
}

func updateOrganization(db *gorm.DB, m *Organization) (error, *Organization) {
	m.UpdatedAt = time.Now()
	err := db.Table(orgTablename).Where("id = ?", m.ID).Update(m).Error
	return err, m
}

type OrganizationList []Organization

func getOrganizationList(db *gorm.DB, offset, limit uint32, appId uint64, parentId ...uint64) (error, uint32, OrganizationList) {
	var count uint32
	list := make([]Organization, 0)
	query := db.Table(orgTablename).Where("is_delete = 0")
	if len(parentId) > 0 {
		query = query.Where("parent_id = ?", parentId[0])
	}
	query = query.Where("app_id = ?", appId)
	err := query.Count(&count).
		Offset(offset).
		Limit(limit).
		Find(&list).
		Error
	return err, count, list
}

func getOrganizationListByIdGroup(db *gorm.DB, group []uint64) (error, OrganizationList) {
	var list OrganizationList
	err := db.Table(orgTablename).
		Where("id in (?)", group).
		Find(&list).
		Error
	return err, list
}

func (uList OrganizationList) Item(id uint64) Organization {
	o := Organization{}
	for _, v := range uList {
		if v.ID == id {
			o = v
			break
		}
	}
	return o
}

func (uList OrganizationList) Names() []string {
	userIds := make([]string, 0, len(uList))
	for _, v := range uList {
		userIds = append(userIds, v.Name)
	}
	return userIds
}
