//Created by Goland
//@User: lenora
//@Date: 2021/3/12
//@Time: 4:09 下午
package flow

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Flow struct {
	ID           uint64    `gorm:"primary_key;column:id;type:int(11) unsigned;not null" json:"id"`
	Number       string    `gorm:"column:number;type:varchar(100);not null" json:"number"`                        // 流程号
	Name         string    `gorm:"column:name;type:varchar(255);not null" json:"name"`                            // 流程名称
	AppID        uint64    `gorm:"column:app_id;type:int(11);not null" json:"app_id"`                             // 应用id
	Key          string    `gorm:"unique;column:key;type:varchar(255);not null" json:"key"`                       // bpmn中的流程key
	Desc         string    `gorm:"column:desc;type:text" json:"desc"`                                             // 表单描述
	UserID       uint64    `gorm:"column:user_id;type:int(11) unsigned;not null;default:0" json:"user_id"`        // 用户id
	ServiceGroup string    `gorm:"column:service_group;type:varchar(255);not null" json:"service_group"`          // 使用到的服务集合
	Assignee     []byte    `gorm:"column:assignee;type:varchar(255);not null;default:''" json:"assignee"`         // 处理用户
	Notifier     []byte    `gorm:"column:notifier;type:varchar(255);not null;default:''" json:"notifier"`         // 抄送人
	XML          string    `gorm:"column:xml;type:text;not null" json:"xml"`                                      // xml文件内容
	JSON         string    `gorm:"column:json;type:text" json:"json"`                                             // 工作流json
	Status       *bool     `gorm:"column:status;type:tinyint(1);not null" json:"status"`                          // 流程状态
	IsOnline     *bool     `gorm:"column:is_online;type:tinyint(1) unsigned;not null" json:"is_online"`           // 是否已上线
	IsDelete     *bool     `gorm:"column:is_delete;type:tinyint(1) unsigned;not null;default:0" json:"is_delete"` // 删除状态
	CreatedAt    time.Time `gorm:"column:created_at;type:datetime;not null" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;type:datetime" json:"updated_at"`
}

// TableName get sql table name.获取数据库表名
func (m *Flow) TableName() string {
	return "flow"
}

var flowTablename = (&Flow{}).TableName()

func createFlow(db *gorm.DB, m Flow) (error, *Flow) {
	err := db.Table(flowTablename).Create(&m).Error
	return err, &m
}

func flowDetail(db *gorm.DB, id uint64) (bool, *Flow) {
	var item Flow
	status := db.Table(flowTablename).Where("id = ?", id).First(&item).RecordNotFound()
	return status, &item
}

func flowDetailByColumn(db *gorm.DB, group map[string]interface{}) (bool, *Flow) {
	var item Flow
	query := db.Table(flowTablename)
	for i, x := range group {
		query = query.Where("`"+i+"` = ?", x)
	}
	status := query.First(&item).RecordNotFound()
	return status, &item
}

func updateFlow(db *gorm.DB, m Flow) (error, *Flow) {
	m.UpdatedAt = time.Now()
	err := db.Table(flowTablename).Where("id = ?", m.ID).Update(&m).Error
	return err, &m
}

func batchUpdateFlow(db *gorm.DB, where, change map[string]interface{}) error {
	query := db.Table(flowTablename)
	for i, x := range where {
		query = query.Where(i+" = ?", x)
	}
	err := query.Updates(change).Error
	return err
}

//
//func delete(db *gorm.DB, id uint64) error {
//	err := db.Table(tablename).Delete(&Sample{},"id = ?",id).Error
//	return err
//}
//

type FlowList []Flow

func getFlowList(db *gorm.DB, offset, limit uint32, appId uint64, name string, status bool, ids ...uint64) (error, uint32, FlowList) {
	var count uint32
	var list FlowList
	query := db.Table(flowTablename).
		Where("app_id = ?", appId).
		Where("is_delete = 0")
	if name != "" {
		query = query.Where("name like ?", "%"+name+"%")
	}
	if status {
		query = query.Where("status = ?", ENABLE)
	}
	if len(ids) > 0 {
		query = query.Where("id in (?)", ids)
	}
	err := query.Count(&count).
		Order("created_at desc").
		Offset(offset).
		Limit(limit).
		Find(&list).
		Error
	return err, count, list
}

func (uList FlowList) Keys() []string {
	keys := make([]string, 0, len(uList))
	for _, v := range uList {
		keys = append(keys, v.Key)
	}
	return keys
}

func (uList FlowList) Users() []uint64 {
	keys := make([]uint64, 0, len(uList))
	for _, v := range uList {
		keys = append(keys, v.UserID)
	}
	return keys
}
