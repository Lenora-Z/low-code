// Package form
//Created by Goland
//@User: lenora
//@Date: 2021/3/12
//@Time: 10:21 上午
package form

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Form struct {
	ID                uint64    `gorm:"primary_key;column:id;type:int(11) unsigned;not null" json:"id"`
	Name              string    `gorm:"column:name;type:varchar(255);not null" json:"name"`                                             // 表单名称
	AppID             uint64    `gorm:"column:app_id;type:int(11) unsigned;not null" json:"app_id"`                                     // 应用id
	DatasourceTableID uint64    `gorm:"column:datasource_table_id;type:int(11) unsigned;not null;default:0" json:"datasource_table_id"` // 数据表id
	From              uint64    `gorm:"column:from;type:int(11);not null;default:0" json:"from"`                                        // 原表
	Type              int8      `gorm:"column:type;type:tinyint(1) unsigned;not null;default:1" json:"type"`                            // 表单类型 1-标准表单
	PageType          uint8     `gorm:"column:page_type;type:tinyint(1);not null;default:1" json:"page_type"`                           // 页面类型 1-表单 2-展示页
	Number            string    `gorm:"column:number;type:varchar(100);not null" json:"number"`                                         // 表单编号
	Desc              string    `gorm:"column:desc;type:varchar(255);not null" json:"desc"`                                             // 表单描述
	Content           []byte    `gorm:"column:content;type:text" json:"content"`                                                        // 表单控件json
	Footer            []byte    `gorm:"column:footer;type:text" json:"footer"`                                                          // 底部数据
	Status            *bool     `gorm:"column:status;type:tinyint(1);not null" json:"status"`                                           // 表单状态 0-未生效 1-已生效
	IsOnline          *bool     `gorm:"column:is_online;type:tinyint(1) unsigned;not null" json:"is_online"`                            // 是否是上线表单 1-是
	IsDelete          *bool     `gorm:"column:is_delete;type:tinyint(1) unsigned;not null;default:0" json:"is_delete"`                  // 删除状态 0-未删除 1-已删除
	CreatedAt         time.Time `gorm:"column:created_at;type:datetime;not null" json:"created_at"`
	UpdatedAt         time.Time `gorm:"column:updated_at;type:datetime" json:"updated_at"`
}

// TableName get sql table name.获取数据库表名
func (m *Form) TableName() string {
	return "form"
}

var formTablename = (&Form{}).TableName()

func createForm(db *gorm.DB, m Form) (error, *Form) {
	err := db.Table(formTablename).Create(&m).Error
	return err, &m
}

func formDetail(db *gorm.DB, id uint64) (bool, *Form) {
	var item Form
	status := db.Table(formTablename).Where("id = ?", id).First(&item).RecordNotFound()
	return status, &item
}

func formDetailByColumn(db *gorm.DB, group map[string]interface{}) (bool, *Form) {
	var item Form
	query := db.Table(formTablename)
	for i, x := range group {
		query = query.Where(i+" = ?", x)
	}
	status := query.First(&item).RecordNotFound()
	return status, &item
}

func updateForm(db *gorm.DB, m Form) (error, *Form) {
	m.UpdatedAt = time.Now()
	err := db.Table(formTablename).Where("id = ?", m.ID).Update(&m).Error
	return err, &m
}

func batchUpdateForm(db *gorm.DB, where, change map[string]interface{}) error {
	query := db.Table(formTablename)
	for i, x := range where {
		query = query.Where(i+" = ?", x)
	}
	err := query.Updates(change).Error
	return err
}

type FormList []Form

// 列表筛选
// @params db *gorm.DB db服务
// @params offset uint32 偏移量
// @params limit uint32 数据容量
// @params appId uint64 应用id
// @params t int8 表单类型
// @params pt int8 页面类型
// @params status bool 是否仅获取生效表单
// @params forms []uint64 表单id集合
// @return
func getFormList(db *gorm.DB, offset, limit uint32, appId uint64, t, pt int8, name string, status bool, forms []uint64) (error, uint32, FormList) {
	var count uint32
	list := make(FormList, 0)
	query := db.Table(formTablename).
		Where("app_id = ?", appId).
		Where("is_delete = 0")
	if t != 0 {
		query = query.Where("type = ?", t)
	}
	if pt != 0 {
		query = query.Where("page_type = ?", pt)
	}
	if name != "" {
		query = query.Where("name like ?", "%"+name+"%")
	}
	if status {
		query = query.Where("status = ?", ENABLE)
	}
	if len(forms) > 0 {
		query = query.Where("id in (?)", forms)
	}
	err := query.Count(&count).
		Order("updated_at desc").
		Offset(offset).
		Limit(limit).
		Find(&list).
		Error
	return err, count, list
}

func (list FormList) Ids() []uint64 {
	rList := make([]uint64, 0, cap(list))
	for _, x := range list {
		rList = append(rList, x.ID)
	}
	return rList
}
