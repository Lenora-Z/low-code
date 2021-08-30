//Created by Goland
//@User: lenora
//@Date: 2021/3/15
//@Time: 3:21 下午
package flow

const (
	DISABLE = iota
	ENABLE
)

const (
	PENDING int8 = iota
	VALID
	INVALID
)

const (
	ACT_TYPE uint8 = iota + 1
	COND_TYPE
)

const (
	IS uint8 = iota + 1
	IS_NOT
)

var TRUE, FALSE = true, false

var FlowAssigneeColumns = struct {
	UserID   string
	Activity string
	FlowID   string
}{
	UserID:   "user_id",
	Activity: "activity",
	FlowID:   "flow_id",
}

type ActAssignee struct {
	ActId    string   `json:"act_id"`   //事件标识
	Assignee []uint64 `json:"assignee"` //审批人id
}

var FlowNotifierColumns = struct {
	UserID string
	FlowID string
}{
	UserID: "user_id",
	FlowID: "flow_id",
}

type AssigneeList []FlowAssignee
type NotifierList []FlowNotifier

type ActNotifier struct {
	ActId    string   `json:"act_id"`   //事件标识
	Notifier []uint64 `json:"notifier"` //抄送人id
}
