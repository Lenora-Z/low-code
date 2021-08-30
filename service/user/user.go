//Created by Goland
//@User: lenora
//@Date: 2021/3/9
//@Time: 7:57 下午
package user

import (
	"github.com/jinzhu/gorm"
	"time"
)

type User struct {
	ID        uint64    `gorm:"primary_key;column:id;type:int(11) unsigned;not null" json:"id"`
	AppID     uint64    `gorm:"column:app_id;type:int(11) unsigned;not null" json:"app_id"`              // 应用id
	Account   string    `gorm:"column:account;type:varchar(255);not null" json:"account"`                // 账号
	Password  string    `gorm:"column:password;type:varchar(255);not null" json:"password"`              // 密码
	Nickname  string    `gorm:"column:nickname;type:varchar(255)" json:"nickname"`                       // 昵称
	TrueName  string    `gorm:"column:true_name;type:varchar(255);not null;default:''" json:"true_name"` // 真实姓名
	Mobile    string    `gorm:"column:mobile;type:varchar(20)" json:"mobile"`                            // 手机号
	Mail      string    `gorm:"column:mail;type:varchar(255)" json:"mail"`                               // 邮箱
	GroupID   uint64    `gorm:"column:group_id;type:int(11);not null" json:"group_id"`                   // 组织id
	PwdStatus *bool     `gorm:"column:pwd_status;type:tinyint(2);not null" json:"pwd_status"`            // 1已重置
	IsDelete  *bool     `gorm:"column:is_delete;type:tinyint(1) unsigned;not null" json:"is_delete"`     // 是否被删除 1-已删除
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime" json:"updated_at"`
}

// TableName get sql table name.获取数据库表名
func (m *User) TableName() string {
	return "user"
}

var userTablename = (&User{}).TableName()

func createUser(db *gorm.DB, m *User) (error, *User) {
	err := db.Table(userTablename).Create(m).Error
	return err, m
}

func userDetail(db *gorm.DB, id uint64) (bool, *User) {
	var item User
	status := db.Table(userTablename).Where("id = ?", id).First(&item).RecordNotFound()
	return status, &item
}

func userDetailByColumn(db *gorm.DB, group map[string]interface{}) (bool, *User) {
	var item User
	query := db.Table(userTablename)
	for i, x := range group {
		query = query.Where(i+" = ?", x)
	}
	status := query.First(&item).RecordNotFound()
	return status, &item
}

func updateUser(db *gorm.DB, m *User) (error, *User) {
	m.UpdatedAt = time.Now()
	err := db.Table(userTablename).Where("id = ?", m.ID).Update(m).Error
	return err, m
}

type UserList []User

func getUserList(db *gorm.DB, offset, limit uint32, appId uint64, name string, groupId []uint64, trueName string) (error, uint32, UserList) {
	var count uint32
	list := make(UserList, 0)
	query := db.Table(userTablename).Where("is_delete = 0").Where("app_id = ?", appId)
	if name != "" {
		query = query.Where("account = ? or nickname = ? or true_name = ?", name, name, name)
	}
	if trueName != "" {
		query = query.Where("true_name like ?", "%"+trueName+"%")
	}
	if groupId != nil {
		query = query.Where("group_id in (?)", groupId)
	}
	err := query.Count(&count).
		Offset(offset).
		Limit(limit).
		Find(&list).
		Error
	return err, count, list
}

func getMultiUserDetail(db *gorm.DB, ids []uint64) (error, UserList) {
	list := make(UserList, 0)
	err := db.Table(userTablename).Where("id in (?)", ids).Find(&list).Error
	return err, list
}

func (uList UserList) GetUser(id uint64) User {
	for _, v := range uList {
		if v.ID == id {
			return v
		}
	}
	return User{}
}

func (uList UserList) TrueNames() []string {
	arr := make([]string, 0, cap(uList))
	for _, v := range uList {
		arr = append(arr, v.TrueName)
	}
	return arr
}

func (uList UserList) UserIds() []uint64 {
	arr := make([]uint64, 0, cap(uList))
	for _, v := range uList {
		arr = append(arr, v.ID)
	}
	return arr
}

func (uList UserList) Groups() []uint64 {
	arr := make([]uint64, 0, cap(uList))
	for _, v := range uList {
		arr = append(arr, v.GroupID)
	}
	return arr
}
