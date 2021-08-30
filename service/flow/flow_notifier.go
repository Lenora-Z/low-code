package flow

import "github.com/jinzhu/gorm"

type FlowNotifier struct {
	UserID   uint64 `gorm:"column:user_id;type:int(11) unsigned;not null;default:0" json:"user_id"` // 用户id
	Activity string `gorm:"column:activity;type:varchar(255);not null;default:''" json:"activity"`  // 任务id
	FlowID   uint64 `gorm:"column:flow_id;type:int(11) unsigned;not null;default:0" json:"flow_id"` // 流程id
}

// TableName get sql table name.获取数据库表名
func (m *FlowNotifier) TableName() string {
	return "flow_notifier"
}

var ccTableName = (&FlowNotifier{}).TableName()

func createNotifier(db *gorm.DB, m FlowNotifier) (error, *FlowNotifier) {
	err := db.Table(ccTableName).Create(&m).Error
	return err, &m
}

func deleteNotifier(db *gorm.DB, flowId uint64) error {
	err := db.Table(ccTableName).Delete(&FlowNotifier{}, FlowNotifierColumns.FlowID+" = ?", flowId).Error
	return err
}

func getNotifierList(db *gorm.DB, flowId, userId uint64) (error, uint32, NotifierList) {
	var count uint32
	list := make(NotifierList, 0)
	query := db.Table(ccTableName)
	if flowId > 0 {
		query = query.Where(FlowNotifierColumns.FlowID+" = ?", flowId)
	}
	if userId > 0 {
		query = query.Where(FlowNotifierColumns.UserID+" = ?", userId)
	}
	err := query.Count(&count).
		Find(&list).
		Error
	return err, count, list
}

func (nList NotifierList) Flows() []uint64 {
	flowIds := make([]uint64, 0, len(nList))
	for _, v := range nList {
		flowIds = append(flowIds, v.FlowID)
	}
	return flowIds
}

func (nList NotifierList) Acts() []string {
	notifierActs := make([]string, 0, len(nList))
	for _, v := range nList {
		notifierActs = append(notifierActs, v.Activity)
	}
	return notifierActs
}
