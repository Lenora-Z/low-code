//Created by Goland
//@User: lenora
//@Date: 2021/3/16
//@Time: 8:18 下午
package user

import "github.com/jinzhu/gorm"

type Permission struct {
	ID       uint64 `gorm:"primary_key;column:id;type:int(11) unsigned;not null" json:"id"`
	Name     string `gorm:"column:name;type:varchar(255);not null" json:"name"`         // 权限名称
	URL      string `gorm:"column:url;type:varchar(255);not null" json:"url"`           // 路由地址
	AppID    uint64 `gorm:"column:app_id;type:int(11) unsigned;not null" json:"app_id"` // 应用id
	ParentID uint64 `gorm:"column:parent_id;type:int(11);not null" json:"parent_id"`    // 上层id
}

// TableName get sql table name.获取数据库表名
func (m *Permission) TableName() string {
	return "permission"
}

var pmTablename = (&Permission{}).TableName()

func getPermissionList(db *gorm.DB, offset, limit uint32, ids []uint64) (error, uint32, []Permission) {
	var count uint32
	list := make([]Permission, 0)
	query := db.Table(pmTablename)
	if len(ids) > 0 {
		query = query.Where("id in (?)", ids)
	}
	err := query.Count(&count).
		Offset(offset).
		Limit(limit).
		Find(&list).
		Error
	return err, count, list
}
