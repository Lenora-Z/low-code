//Created by Goland
//@User: lenora
//@Date: 2021/3/16
//@Time: 4:59 下午
package user

import "github.com/jinzhu/gorm"

type RolePermission struct {
	ID       uint64 `gorm:"primary_key;column:id;type:int(11) unsigned;not null" json:"id"`
	RoleID   uint64 `gorm:"column:role_id;type:int(11) unsigned;not null" json:"role_id"` // 角色id
	ParentID uint64 `gorm:"column:parent_id;type:int(11);not null" json:"parent_id"`      // 父层id,permission表id
	FormID   uint64 `gorm:"column:form_id;type:int(11);not null" json:"form_id"`          // 表单id
}

// TableName get sql table name.获取数据库表名
func (m *RolePermission) TableName() string {
	return "role_permission"
}

var rmTablename = (&RolePermission{}).TableName()

func rolePermissionDetailByColumn(db *gorm.DB, roleId []uint64, pId, formId uint64) (bool, *RolePermission) {
	var item RolePermission
	query := db.Table(rmTablename)
	if len(roleId) > 0 {
		query = query.Where("role_id in (?)", roleId)
	}
	if pId != 0 {
		query = query.Where("parent_id = ?", pId)
	}
	if formId != 0 {
		query = query.Where("form_id = ?", formId)
	}
	status := query.First(&item).RecordNotFound()
	return status, &item
}

func createRolePermission(db *gorm.DB, m RolePermission) (error, *RolePermission) {
	err := db.Table(rmTablename).Create(&m).Error
	return err, &m
}

func deleteRolePermission(db *gorm.DB, roleId uint64) error {
	err := db.Table(rmTablename).Delete(&RolePermission{}, "role_id = ?", roleId).Error
	return err
}

func rolePermissionList(db *gorm.DB, pId uint64, roles []uint64) (error, PermissionList) {
	var list PermissionList
	query := db.Table(rmTablename).
		Where("role_id in (?)", roles)
	if pId == 0 {
		query = query.Select("distinct parent_id as id")
	} else {
		query = query.Where("parent_id = ?", pId).
			Select("distinct form_id as id")
	}
	err := query.Scan(&list).Error
	return err, list
}

func (l *PermissionList) Ids() []uint64 {
	ids := make([]uint64, 0, len(*l))
	for _, v := range *l {
		ids = append(ids, v.Id)
	}
	return ids
}
