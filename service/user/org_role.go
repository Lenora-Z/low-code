//Created by Goland
//@User: lenora
//@Date: 2021/3/9
//@Time: 5:42 下午
package user

import (
	"github.com/jinzhu/gorm"
)

type OrgRole struct {
	ID      uint64 `gorm:"primary_key;column:id;type:int(11) unsigned;not null" json:"id"`
	GroupID uint64 `gorm:"column:group_id;type:int(11) unsigned;not null" json:"group_id"` // 组织id
	RoleID  uint64 `gorm:"column:role_id;type:int(11) unsigned;not null" json:"role_id"`   // 角色id
	AppID   uint64 `gorm:"column:app_id;type:int(11) unsigned;not null" json:"app_id"`     // 应用id
}

// TableName get sql table name.获取数据库表名
func (m *OrgRole) TableName() string {
	return "org_role"
}

var ogTablename = (&OrgRole{}).TableName()

func createOrgRole(db *gorm.DB, m *OrgRole) (error, *OrgRole) {
	err := db.Table(ogTablename).Create(m).Error
	return err, m
}

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

type OrgRoleList []OrgRole

func getOrgRoleList(db *gorm.DB, offset, limit uint32, groupId uint64, roleId []uint64) (error, uint32, OrgRoleList) {
	var count uint32
	list := make([]OrgRole, 0, limit)
	query := db.Table(ogTablename)
	if groupId != 0 {
		query = query.Where("group_id = ?", groupId)
	}
	if len(roleId) > 0 {
		query = query.Where("role_id in (?)", roleId)
	}
	err := query.
		Offset(offset).
		Limit(limit).
		Find(&list).
		Error
	return err, count, list
}

func deleteOrgRole(db *gorm.DB, groupId, roleId uint64) error {
	query := db.Table(ogTablename)
	if groupId != 0 {
		query = query.Where("group_id = ?", groupId)
	}
	if roleId != 0 {
		query = query.Where("role_id = ?", roleId)
	}
	err := query.Delete(&OrgRole{}).Error
	return err
}

func (uList OrgRoleList) OrgIds() []uint64 {
	userIds := make([]uint64, 0, len(uList))
	for _, v := range uList {
		userIds = append(userIds, v.GroupID)
	}
	return userIds
}

func (uList OrgRoleList) RoleIds() []uint64 {
	roleIds := make([]uint64, 0, len(uList))
	for _, v := range uList {
		roleIds = append(roleIds, v.RoleID)
	}
	return roleIds
}
