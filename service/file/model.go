package file

import (
	"github.com/jinzhu/gorm"
	"time"
)

type File struct {
	ID        uint32    `gorm:"primaryKey;column:id;type:int(11) unsigned;not null" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null" json:"created_at"`
	Path      string    `gorm:"column:path;type:varchar(255);not null" json:"path"`                            // 文件地址
	Name      string    `gorm:"column:name;type:varchar(255);not null;default:''" json:"name"`                 // 文件名称
	IsDelete  *bool     `gorm:"column:is_delete;type:tinyint(1) unsigned;not null;default:0" json:"is_delete"` // 是否删除
	AppID     uint64    `gorm:"column:app_id;type:int(11) unsigned;not null" json:"app_id"`                    // 应用id
	FileHash  string    `gorm:"column:file_hash;type:varchar(255);not null" json:"file_hash"`                  // 文件hash值
	Hash      string    `gorm:"unique;column:hash;type:varchar(50);not null" json:"hash"`                      // 数据hash
}

// TableName get sql table name.获取数据库表名
func (m *File) TableName() string {
	return "file"
}

var tablename = (&File{}).TableName()

func createFile(db *gorm.DB, m File) (error, *File) {
	err := db.Table(tablename).Create(&m).Error
	return err, &m
}

func fileDetailByColumn(db *gorm.DB, group map[string]interface{}) (bool, *File) {
	var item File
	query := db.Table(tablename)
	for i, x := range group {
		query = query.Where(i+" = ?", x)
	}
	status := query.First(&item).RecordNotFound()
	return status, &item
}
