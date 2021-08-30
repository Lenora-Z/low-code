//Created by Goland
//@User: lenora
//@Date: 2021/3/9
//@Time: 4:17 下午
package user

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Role struct {
	ID          uint64    `gorm:"primary_key;column:id;type:int(11) unsigned;not null" json:"id"`
	AppID       uint64    `gorm:"column:app_id;type:int(11) unsigned;not null" json:"app_id"`       // 应用id
	Name        string    `gorm:"column:name;type:varchar(255);not null" json:"name"`               // 角色名
	Permissions string    `gorm:"column:permissions;type:varchar(255);not null" json:"permissions"` // 角色权限
	IsDelete    *bool     `gorm:"column:is_delete;type:tinyint(1);not null" json:"is_delete"`       // 是否被删除
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime;not null" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:datetime" json:"updated_at"`
}

// TableName get sql table name.获取数据库表名
func (m *Role) TableName() string {
	return "role"
}

var roleTablename = (&Role{}).TableName()

func createRole(db *gorm.DB, m *Role) (error, *Role) {
	err := db.Table(roleTablename).Create(m).Error
	return err, m
}

func roleDetail(db *gorm.DB, id uint64) (bool, *Role) {
	var item Role
	status := db.Table(roleTablename).Where("id = ?", id).First(&item).RecordNotFound()
	return status, &item
}

func roleDetailByColumn(db *gorm.DB, group map[string]interface{}) (bool, *Role) {
	var item Role
	query := db.Table(roleTablename)
	for i, x := range group {
		query = query.Where(i+" = ?", x)
	}
	status := query.First(&item).RecordNotFound()
	return status, &item
}

func updateRole(db *gorm.DB, m *Role) (error, *Role) {
	m.UpdatedAt = time.Now()
	err := db.Table(roleTablename).Where("id = ?", m.ID).Update(m).Error
	return err, m
}

func getRoleList(db *gorm.DB, offset, limit uint32, appId uint64, name string) (error, uint32, RoleList) {
	var count uint32
	list := make([]Role, 0)
	query := db.Table(roleTablename).Where("is_delete = 0")
	query = query.Where("app_id = ?", appId)
	if name != "" {
		query = query.Where("name like ?", "%"+name+"%")
	}
	err := query.Count(&count).
		Offset(offset).
		Limit(limit).
		Find(&list).
		Error
	return err, count, list
}

type RoleList []Role

func getRoleListByIdGroup(db *gorm.DB, group []uint64) (error, RoleList) {
	var list RoleList
	err := db.Table(roleTablename).
		Where("is_delete = 0").
		Where("id in (?)", group).
		Find(&list).
		Error
	return err, list
}

func (uList RoleList) Names() []string {
	userIds := make([]string, 0, len(uList))
	for _, v := range uList {
		userIds = append(userIds, v.Name)
	}
	return userIds
}
