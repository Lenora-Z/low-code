// Package source
// Created by GoLand
// @User: lenora
// @Date: 2021/7/22
// @Time: 10:41

package source

import (
	"github.com/jinzhu/gorm"
	"time"
)

type DatasourceTable struct {
	ID           uint64    `gorm:"primaryKey;column:id;type:int(11) unsigned;not null" json:"id"`
	CreatedAt    time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"created_at"`  // 创建时间
	AppID        uint64    `gorm:"index:app_id_idx;column:app_id;type:int(11) unsigned;not null;default:0" json:"app_id"` // 应用id
	SchemataID   uint64    `gorm:"column:schemata_id;type:int(11) unsigned;not null;default:0" json:"schemata_id"`        // 数据源id
	TableSchema  string    `gorm:"column:table_schema;type:varchar(64);not null;default:''" json:"schema"`                // 数据库名称
	TableName    string    `gorm:"column:table_name;type:varchar(64);not null;default:''" json:"name"`                    // 数据表名称
	TableComment string    `gorm:"column:table_comment;type:varchar(2048);not null;default:''" json:"comment"`            // 数据库表备注
}

// TableName get sql table name.获取数据库表名
func (m *DatasourceTable) tableName() string {
	return "datasource_table"
}

var tablename = (&DatasourceTable{}).tableName()

func tableDetail(db *gorm.DB, id uint64) (bool, *DatasourceTable) {
	var item DatasourceTable
	status := db.Table(tablename).Where("id = ?", id).First(&item).RecordNotFound()
	return status, &item
}

func tableDetailByColumn(db *gorm.DB, group map[string]interface{}) (bool, *DatasourceTable) {
	var item DatasourceTable
	query := db.Table(tablename)
	for i, x := range group {
		query = query.Where(i+" = ?", x)
	}
	status := query.First(&item).RecordNotFound()
	return status, &item
}

type TableList []DatasourceTable

func getTableList(db *gorm.DB, offset, limit uint32, appId uint64, ids []uint64, name string) (error, uint32, TableList) {
	var count uint32
	list := make(TableList, 0)
	query := db.Table(tablename).Where(DatasourceTableColumns.AppID+" = ?", appId)
	if ids != nil {
		query = query.Where(DatasourceTableColumns.ID+" in (?)", ids)
	}
	if name != "" {
		query = query.Where(DatasourceTableColumns.TableName+" like ? or "+DatasourceTableColumns.TableComment+" like ?", "%"+name+"%", "%"+name+"%")
	}
	err := query.Count(&count).
		Order(DatasourceTableColumns.CreatedAt + " desc").
		Offset(offset).
		Limit(limit).
		Find(&list).
		Error
	return err, count, list
}

func (list TableList) Ids() []uint64 {
	uList := make([]uint64, 0, cap(list))
	for _, x := range list {
		uList = append(uList, x.ID)
	}
	return uList
}

func (list TableList) Names() []string {
	uList := make([]string, 0, cap(list))
	for _, x := range list {
		uList = append(uList, x.TableName)
	}
	return uList
}
