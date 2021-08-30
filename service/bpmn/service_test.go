//Created by Goland
//@User: lenora
//@Date: 2021/3/16
//@Time: 7:01 下午
package bpmn

import (
	"github.com/sirupsen/logrus"
	"testing"
)

var api = "192.168.3.47:31800"

func TestBpmnService_DeployProcess(t *testing.T) {
	var str = "[{\"key\":5894158697,\"name\":\"开始节点\",\"description\":\"节点描述\",\"type\":0,\"props\":{\"name\":\"ccccc\",\"description\":\"ccccc\",\"type\":2,\"trigger_time\":123,\"trigger_interval\":33}},{\"key\":4738053273,\"name\":\"新增数据\",\"description\":\"\",\"type\":6,\"props\":{\"columns\":[{\"16\":{\"type\":3,\"value\":\"11111\"}}],\"name\":\"1231321321\",\"description\":\"31221312\",\"target_id\":2}},{\"key\":3251668235,\"name\":\"结束\",\"description\":\"节点描述\",\"type\":9}]"

	srv := NewBpmnService("192.168.3.47:31800")
	err, name := srv.DeployProcess(29, "示例流程111", "示例流程描述111", str)
	logrus.Info(name, err)
}

func TestBpmnService_FlowTask(t *testing.T) {
	srv := NewBpmnService(api)
	err, lists := srv.FlowTask("1", "10", "Process_1eu7jdf", "1")
	logrus.Info(err, lists)
}

func TestBpmnService_ProcessLog(t *testing.T) {
	srv := NewBpmnService(api)
	err, lists := srv.ProcessLog("1", "10", "Process_1eu7jdf")
	logrus.Info(err, lists)
}

func TestBpmnService_LogDetail(t *testing.T) {
	srv := NewBpmnService(api)
	err, item := srv.LogDetail("b3509539-86fe-11eb-aae1-1a788ae1a83b")
	logrus.Info(err, item)
}

func TestBpmnService_ExecuteProcess(t *testing.T) {
	srv := NewBpmnService(api)
	var m []map[string]string
	m = append(m, map[string]string{
		"userId": "message",
	})
	d := make(map[string]interface{})
	d["amount"] = 555
	d["message"] = "shdjksa"
	err, str := srv.ExecuteProcess("Process_1eu7jdf", m, d)
	logrus.Info(err, str)

}
