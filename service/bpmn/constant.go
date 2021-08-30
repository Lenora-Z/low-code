//Created by Goland
//@User: lenora
//@Date: 2021/3/16
//@Time: 5:14 下午
package bpmn

const (
	PROCESS_DEPLOY    = "/process/json/deploy"
	PROCESS_INSTANCES = "/process/instances"
	PROCESS_LOG       = "/process/log"
	PROCESS_EXECUTE   = "/process/start"
	TASK              = "/task"
	TASK_VARIABLES    = "/task/variables"
	TASK_COMPLETE     = "/task/complete"
	PROCESS_RUN       = "/process/engine"
)

const SUCCESS_CODE uint64 = 200

type CommonReturn struct {
	Code    uint64 `json:"code"`
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
}

type TaskItem struct {
	Id                  string `json:"id"`
	ProcessDefinitionId string `json:"processDefinitionId"`
	Name                string `json:"name"`
	ProcessInstanceId   string `json:"processInstanceId"`
	TaskDefinitionKey   string `json:"taskDefinitionKey"`
	CreateTime          string `json:"createTime"`
	DeleteReason        string `json:"deleteReason"`
}

type LogItem struct {
	ProcessInstanceId     string `json:"processInstanceId"`
	ProcessDefinitionName string `json:"processDefinitionName"`
	ProcessDefinitionKey  string `json:"processDefinitionKey"`
	StartTime             string `json:"startTime"`
	EndTime               string `json:"endTime"`
	StartActivityId       string `json:"startActivityId"`
	EndActivityId         string `json:"endActivityId"`
	State                 string `json:"state"`
}

type TaskListVO struct {
	Count    uint32     `json:"count"`
	TaskList []TaskItem `json:"taskList"`
}

type ProcessLogVO struct {
	Count            uint32    `json:"count"`
	ProcessInstances []LogItem `json:"processInstances"`
}
