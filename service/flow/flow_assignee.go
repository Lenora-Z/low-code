package flow

import "github.com/jinzhu/gorm"

type FlowAssignee struct {
	UserID   uint64 `gorm:"column:user_id;type:int(11) unsigned;not null;default:0" json:"user_id"` // 用户id
	Activity string `gorm:"column:activity;type:varchar(255);not null;default:''" json:"activity"`  // 任务id
	FlowID   uint64 `gorm:"column:flow_id;type:int(11) unsigned;not null;default:0" json:"flow_id"` // 流程id
}

// TableName get sql table name.获取数据库表名
func (m *FlowAssignee) TableName() string {
	return "flow_assignee"
}

var assigneeTableName = (&FlowAssignee{}).TableName()

func createAssignee(db *gorm.DB, m FlowAssignee) (error, *FlowAssignee) {
	err := db.Table(assigneeTableName).Create(&m).Error
	return err, &m
}

func deleteAssignee(db *gorm.DB, flowId uint64) error {
	err := db.Table(assigneeTableName).Delete(&FlowAssignee{}, FlowAssigneeColumns.FlowID+" = ?", flowId).Error
	return err
}

func getAssigneeList(db *gorm.DB, flowId, userId uint64) (error, uint32, AssigneeList) {
	var count uint32
	list := make(AssigneeList, 0)
	query := db.Table(assigneeTableName)
	if flowId > 0 {
		query = query.Where(FlowAssigneeColumns.FlowID+" = ?", flowId)
	}
	if userId > 0 {
		query = query.Where(FlowAssigneeColumns.UserID+" = ?", userId)
	}
	err := query.Count(&count).
		Find(&list).
		Error
	return err, count, list
}

func (uList AssigneeList) Flows() []uint64 {
	flowIds := make([]uint64, 0, len(uList))
	for _, v := range uList {
		flowIds = append(flowIds, v.FlowID)
	}
	return flowIds
}

func (uList AssigneeList) Users() []uint64 {
	userIds := make([]uint64, 0, len(uList))
	for _, v := range uList {
		userIds = append(userIds, v.UserID)
	}
	return userIds
}

func (uList AssigneeList) Acts() []string {
	assigneeActs := make([]string, 0, len(uList))
	for _, v := range uList {
		assigneeActs = append(assigneeActs, v.Activity)
	}
	return assigneeActs
}
