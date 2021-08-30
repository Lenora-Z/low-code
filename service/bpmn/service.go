//Created by Goland
//@User: lenora
//@Date: 2021/3/16
//@Time: 5:14 下午
package bpmn

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Lenora-Z/low-code/utils"
	"github.com/sirupsen/logrus"
	"net/url"
)

type BpmnService interface {
	//申请处理
	CompleteTask(id string, data []map[string]interface{}) error
	//流程部署
	DeployProcess(id uint64, name, desc, file string) (error, string)
	//流程启动
	ExecuteProcess(key string, maps []map[string]string, data map[string]interface{}) (error, string)
	//申请处理列表
	FlowTask(page, limit, key, status string) (error, *TaskListVO)
	//申请信息参数
	FlowTaskVariables(id, status string) (error, []map[string]interface{})
	//日志详情
	LogDetail(id string) (error, interface{})
	//流程日志列表
	ProcessLog(page, limit, key string) (error, *ProcessLogVO)
	//触发流程
	RunProcess(processId string, data map[string]interface{}, userId uint64) error
}

type bpmnService struct {
	api string
}

func NewBpmnService(api string) BpmnService {
	u := new(bpmnService)
	u.api = "http://" + api
	return u
}

func (srv *bpmnService) DeployProcess(id uint64, name, desc, file string) (error, string) {
	st := struct {
		Id   uint64 `json:"id"`
		Name string `json:"name"`
		Desc string `json:"description"`
		Node string `json:"node"`
	}{
		Id:   id,
		Name: name,
		Desc: desc,
		Node: file,
	}

	body, err := json.Marshal(st)
	if err != nil {
		return err, ""
	}

	ret, err := utils.SendRequest(&utils.HeaderRequest{
		Method: "POST",
		Url:    srv.api + PROCESS_DEPLOY,
		Body:   string(body),
	})
	logrus.Info("流程部署===结果返回===", ret)

	if err != nil {
		logrus.Error("deploy process wrong", err.Error())
		return err, ""
	}

	type Response struct {
		Id           string `json:"id"`
		Key          string `json:"key"`
		Name         string `json:"name"`
		DeploymentId string `json:"deploymentId"`
	}

	type DeployReturn struct {
		CommonReturn
		Data Response `json:"data"`
	}

	var result DeployReturn
	err = json.Unmarshal([]byte(ret), &result)
	if err != nil {
		logrus.Error("deploy result parsing wrong")
		return err, ""
	}

	if result.Code != SUCCESS_CODE {
		logrus.Error("deploy process failed")
		return errors.New("deploy process failed:" + result.Msg), ""
	}

	return nil, result.Data.Key
}

func (srv *bpmnService) FlowTask(page, limit, key, status string) (error, *TaskListVO) {
	param := url.Values{}
	param.Add("pageIndex", page)
	param.Add("pageSize", limit)
	param.Add("processDefinitionKeys", key)
	param.Add("status", status) //status=0时则表示所有状态的流程
	reqUrl := srv.api + TASK + "?" + param.Encode()
	logrus.Info("申请处理列表===请求地址===", reqUrl)

	ret, err := utils.SendRequest(&utils.HeaderRequest{
		Method: "GET",
		Url:    reqUrl,
	})
	logrus.Info("申请处理列表===结果返回===", ret)

	if err != nil {
		logrus.Error("get task error", err.Error())
		return err, nil
	}

	type DeployReturn struct {
		CommonReturn
		Data TaskListVO `json:"data"`
	}

	var result DeployReturn
	err = json.Unmarshal([]byte(ret), &result)
	if err != nil {
		logrus.Error("task result parsing wrong")
		return err, nil
	}

	if result.Code != SUCCESS_CODE {
		logrus.Error("get task failed")
		return errors.New("get task failed:" + result.Msg), nil
	}
	return nil, &result.Data
}

func (srv *bpmnService) FlowTaskVariables(id, status string) (error, []map[string]interface{}) {
	param := url.Values{}
	param.Add("taskId", id)
	param.Add("status", status)
	reqUrl := srv.api + TASK_VARIABLES + "?" + param.Encode()
	logrus.Info("待审核申请参数列表===请求地址===", reqUrl)

	ret, err := utils.SendRequest(&utils.HeaderRequest{
		Method: "GET",
		Url:    reqUrl,
	})
	logrus.Info("待审核申请参数列表===结果返回===", ret)

	if err != nil {
		logrus.Error("get task params error", err.Error())
		return err, nil
	}

	type DeployReturn struct {
		CommonReturn
		Data []map[string]interface{} `json:"data"`
	}

	var result DeployReturn
	err = json.Unmarshal([]byte(ret), &result)
	if err != nil {
		logrus.Error("task params result parsing wrong")
		return err, nil
	}

	if result.Code != SUCCESS_CODE {
		logrus.Error("get task params failed")
		return errors.New("get task params failed:" + result.Msg), nil
	}
	return nil, result.Data
}

func (srv *bpmnService) CompleteTask(id string, data []map[string]interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	body := fmt.Sprintf(`{"taskId": "%s","variables":%s}`, id, string(bytes))
	logrus.Info("申请处理===请求参数===", body)
	ret, err := utils.SendRequest(&utils.HeaderRequest{
		Method: "POST",
		Url:    srv.api + TASK_COMPLETE,
		Body:   body,
	})
	logrus.Info("申请处理===结果返回===", ret)

	if err != nil {
		logrus.Error("deal apply wrong", err.Error())
		return err
	}

	type DeployReturn struct {
		CommonReturn
		Data interface{} `json:"data"`
	}

	var result DeployReturn
	err = json.Unmarshal([]byte(ret), &result)
	if err != nil {
		logrus.Error("deal apply result parsing wrong")
		return err
	}

	if result.Code != SUCCESS_CODE {
		logrus.Error("deal apply failed")
		return errors.New("deal apply failed:" + result.Msg)
	}

	return nil
}

func (srv *bpmnService) ProcessLog(page, limit, key string) (error, *ProcessLogVO) {
	param := url.Values{}
	param.Add("pageIndex", page)
	param.Add("pageSize", limit)
	param.Add("processDefinitionKeys", key)
	reqUrl := srv.api + PROCESS_INSTANCES + "?" + param.Encode()
	logrus.Info("流程消息列表===请求地址===", reqUrl)

	ret, err := utils.SendRequest(&utils.HeaderRequest{
		Method: "GET",
		Url:    reqUrl,
	})
	logrus.Info("流程消息列表===结果返回===", ret)

	if err != nil {
		logrus.Error("get process log error", err.Error())
		return err, nil
	}

	type DeployReturn struct {
		CommonReturn
		Data ProcessLogVO `json:"data"`
	}

	var result DeployReturn
	err = json.Unmarshal([]byte(ret), &result)
	if err != nil {
		logrus.Error("process log parsing wrong")
		return err, nil
	}

	if result.Code != SUCCESS_CODE {
		logrus.Error("get process log failed")
		return errors.New("get process log failed:" + result.Msg), nil
	}
	return nil, &result.Data
}

func (srv *bpmnService) LogDetail(id string) (error, interface{}) {
	param := url.Values{}
	param.Add("processInstanceId", id)
	reqUrl := srv.api + PROCESS_LOG + "?" + param.Encode()
	logrus.Info("流程日志===请求地址===", reqUrl)

	ret, err := utils.SendRequest(&utils.HeaderRequest{
		Method: "GET",
		Url:    reqUrl,
	})
	logrus.Info("流程日志===结果返回===", ret)

	if err != nil {
		logrus.Error("get log detail error", err.Error())
		return err, nil
	}

	type DeployReturn struct {
		CommonReturn
		Data interface{} `json:"data"`
	}

	var result DeployReturn
	err = json.Unmarshal([]byte(ret), &result)
	if err != nil {
		logrus.Error("log detail parsing wrong")
		return err, nil
	}

	if result.Code != SUCCESS_CODE {
		logrus.Error("get log detail failed")
		return errors.New("get log detail failed:" + result.Msg), nil
	}
	return nil, &result.Data
}

func (srv *bpmnService) ExecuteProcess(key string, maps []map[string]string, data map[string]interface{}) (error, string) {
	bytes, err := json.Marshal(maps)
	if err != nil {
		return err, ""
	}
	dByte, err := json.Marshal(data)
	if err != nil {
		return err, ""
	}
	var body = fmt.Sprintf(`{"processDefinitionKey":"%s","mapping":%s,"variables":%s}`, key, string(bytes), string(dByte))
	logrus.Info("requestBody:", body)

	ret, err := utils.SendRequest(&utils.HeaderRequest{
		Method: "POST",
		Url:    srv.api + PROCESS_EXECUTE,
		Body:   body,
	})
	logrus.Info("流程启动===结果返回===", ret)

	if err != nil {
		logrus.Error("execute process wrong", err.Error())
		return err, ""
	}

	type res struct {
		ProcessInstanceId    string                 `json:"processInstanceId"`
		Variables            map[string]interface{} `json:"variables"`
		ProcessDefinitionKey string                 `json:"processDefinitionKey"`
	}

	type DeployReturn struct {
		CommonReturn
		Data res `json:"data"`
	}

	var result DeployReturn
	err = json.Unmarshal([]byte(ret), &result)
	if err != nil {
		logrus.Error("execute process parsing wrong")
		return err, ""
	}

	if result.Code != SUCCESS_CODE {
		logrus.Error("execute process failed")
		return errors.New("execute process failed:" + result.Msg), ""
	}

	return nil, result.Data.ProcessInstanceId

}

//触发流程
func (srv *bpmnService) RunProcess(processId string, data map[string]interface{}, userId uint64) error {
	if data == nil {
		data = make(map[string]interface{}, 0)
	}
	bytes, _ := json.Marshal(data)
	var body = fmt.Sprintf(`{"data":%s}`, string(bytes))
	logrus.Info(body)
	ret, err := utils.SendRequest(&utils.HeaderRequest{
		Method: "POST",
		Url:    srv.api + PROCESS_RUN + "?processId=" + processId + "&userId=" + fmt.Sprintf(`%d`, userId),
		//Url:    "http://192.168.9.23:8000" + PROCESS_RUN + "?processId=" + processId + "&userId=" + fmt.Sprintf(`%d`, userId),
		Body:   body,
	})
	if err != nil {
		return errors.New("获取bpmn响应出错：" + err.Error())
	}

	var resp = make(map[string]interface{})
	err = json.Unmarshal([]byte(ret), &resp)
	if err != nil {
		return errors.New("解析responseBody出错：" + err.Error())
	}
	//fmt.Println(body, "11111111111", resp)
	if resp["code"].(float64) != 200 {
		return errors.New(resp["msg"].(string))
	}
	return nil
}
