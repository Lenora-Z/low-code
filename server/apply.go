//Created by Goland
//@User: lenora
//@Date: 2021/3/17
//@Time: 8:38 下午
package server

import (
	"fmt"
	"github.com/Lenora-Z/low-code/service/apply"
	"github.com/Lenora-Z/low-code/service/bpmn"
	"github.com/Lenora-Z/low-code/service/flow"
	"github.com/Lenora-Z/low-code/service/form"
	"github.com/Lenora-Z/low-code/service/form_data"
	"github.com/Lenora-Z/low-code/service/user"
	"github.com/Lenora-Z/low-code/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type TaskItem struct {
	Id                  string      `json:"id"`                    //申请id
	ProcessDefinitionId string      `json:"process_definition_id"` //实例id
	InstanceId          string      `json:"instance_id"`           //实例id
	CreateTime          string      `json:"create_time"`           //发起时间
	FlowName            string      `json:"flow_name"`             //流程名称
	FormName            string      `json:"form_name"`             //表单名称
	FormId              uint64      `json:"form_id"`               //表单id
	UserId              interface{} `json:"user_id"`               //用户id
	Account             string      `json:"account"`               //用户名
	Status              uint8       `json:"status"`                //处理状态 1待处理；2 已处理
}

type TaskListVO struct {
	List  []TaskItem `json:"list"`
	Count uint32     `json:"count"`
}

// @Summary 申请列表
// @Tags apply
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param Authorization header string true "Bearer +token"
// @param page query string true "页码"
// @param limit query string true "页容量"
// @param status query string true "申请状态"
// @Success 200 {object} ApiResponse{result=TaskListVO}
// @Router /api/apply/list [get]
func (ds *defaultServer) ApplyList(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*UserClaims)
	pageStr := ctx.Query("page")
	limitStr := ctx.Query("limit")
	statusStr := ctx.Query("status")
	if pageStr == "" || limitStr == "" {
		ds.InvalidParametersError(ctx)
		return
	}

	//if status := ds.validateClient(DEAL_APPLY, id.Uint64(), claims.UserId); !status {
	//	ds.ResponseError(ctx, ACCESS_DENY)
	//	return
	//}

	//获取用户可处理的全部流程
	srv := flow.NewFlowService(ds.db)
	err, flows, keys := srv.GetFlowListByAssignee(claims.UserId, claims.AppId)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	list := make([]TaskItem, 0, 100)
	if len(keys) <= 0 {
		ds.ResponseSuccess(ctx, TaskListVO{
			List:  list,
			Count: 0,
		})
		return
	}
	keyStr := strings.Join(keys, ",")

	bpmnSrv := bpmn.NewBpmnService(ds.conf.Engine.Api)
	err, ret := bpmnSrv.FlowTask(pageStr, limitStr, keyStr, statusStr)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	statusObj := utils.NewStr(statusStr)
	mongoSrv := form_data.NewFormDataService(ds.mongoDb)
	userSrv := user.NewUserService(ds.db)
	formSrv := form.NewFormService(ds.db)

	//根据实例id获取用户信息
	for _, x := range ret.TaskList {
		//mongo获取用户id
		found, item := mongoSrv.GetByInstanceId(x.ProcessInstanceId)
		if !found {
			ds.ResponseError(ctx, NOT_FOUND)
			return
		}

		//用户名称
		var account = ""
		userStr := fmt.Sprintf("%v", item["user_id"])
		userId := utils.NewStr(userStr)
		if userId.Uint64() != 0 {
			_, userItem := userSrv.GetUser(userId.Uint64())
			account = userItem.Account
		}

		//表单名称
		var formName = ""
		formId, ok := item["form_id"].(float64)
		if ok {
			_, formItem := formSrv.GetFormDetail(uint64(formId))
			formName = formItem.Name
		}

		//流程名称
		var name = ""
		for _, f := range flows {
			if f.Key == x.ProcessDefinitionId[0:strings.Index(x.ProcessDefinitionId, ":")] {
				name = f.Name
				break
			}
		}

		list = append(list, TaskItem{
			Id:                  x.Id,
			ProcessDefinitionId: x.ProcessDefinitionId,
			CreateTime:          x.CreateTime,
			InstanceId:          x.ProcessInstanceId,
			FlowName:            name,
			FormId:              uint64(formId),
			FormName:            formName,
			UserId:              item["user_id"],
			Account:             account,
			Status:              statusObj.Uint8(),
		})
	}

	ds.ResponseSuccess(ctx, TaskListVO{
		List:  list,
		Count: ret.Count,
	})
	return
}

type ApplyDetail struct {
	Variables []map[string]interface{} `json:"variables"`  //处理参数
	CreatedAt time.Time                `json:"created_at"` //填写时间
	Account   string                   `json:"account"`    //填写用户
	Assignee  string                   `json:"assignee"`   //处理人
	UpdatedAt time.Time                `json:"updated_at"` //处理时间
	Status    int8                     `json:"status"`     //处理状态
	Data      map[string]interface{}   `json:"data"`       //表单提交数据
}

// @Summary 申请要处理的参数
// @Tags apply
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param Authorization header string true "Bearer +token"
// @param id query string true "申请id"
// @param instance_id query string true "实例id"
// @param status query string true "状态值(同列表)"
// @Success 200 {object} ApiResponse{result=ApplyDetail}
// @Router /api/apply/param [get]
func (ds *defaultServer) ApplyParams(ctx *gin.Context) {
	//claims := ctx.MustGet(CLAIMS).(*UserClaims)
	taskIdStr := ctx.Query("id")
	instanceId := ctx.Query("instance_id")
	status := ctx.Query("status")
	if taskIdStr == "" || instanceId == "" || status == "" {
		ds.InvalidParametersError(ctx)
		return
	}

	//if status := ds.validateClient(DEAL_APPLY, id.Uint64(), claims.UserId); !status {
	//	ds.ResponseError(ctx, ACCESS_DENY)
	//	return
	//}
	bpmnSrv := bpmn.NewBpmnService(ds.conf.Engine.Api)
	err, ret := bpmnSrv.FlowTaskVariables(taskIdStr, status)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	result := ApplyDetail{
		Variables: ret,
	}

	userSrv := user.NewUserService(ds.db)
	code, msg, dataItem := ds.formDataItemDetail(2, instanceId, "show")
	if code != SUCCESS {
		ds.ResponseError(ctx, code, msg)
		return
	}
	result.Data = dataItem
	result.Account = fmt.Sprintf("%v", dataItem["user_name"])

	if status == "2" {
		srv := apply.NewService(ds.mongoDb)
		found, dealItem := srv.GetLogByProcessId(taskIdStr)
		if found {
			result.UpdatedAt = dealItem.Updated_At
			f, assignee := userSrv.GetUser(dealItem.Assignee)
			if !f {
				result.Assignee = assignee.TrueName
			}
			result.Status = dealItem.Status
		}
	}

	ds.ResponseSuccess(ctx, result)
	return
}

type applyArg struct {
	Id         string `json:"id"`
	InstanceId string `json:"instance_id"` //实例id
	Status     int8   `json:"status"`      //通过状态 0-不通过 1-通过
}

// @Summary 处理申请
// @Tags apply
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param Authorization header string true "Bearer +token"
// @param applyArg body applyArg true "请求体"
// @Success 200 {object} ApiResponse{result=object}
// @Router /api/apply/deal [post]
func (ds *defaultServer) DealApply(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*UserClaims)
	var arg applyArg
	if err := ctx.BindJSON(&arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}
	if err := validatorInstance().Struct(arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	//if status := ds.validateClient(DEAL_APPLY, arg.FormId, claims.UserId); !status {
	//	ds.ResponseError(ctx, ACCESS_DENY)
	//	return
	//}

	bpmnSrv := bpmn.NewBpmnService(ds.conf.Engine.Api)
	if err := bpmnSrv.CompleteTask(arg.Id, []map[string]interface{}{}); err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	mongoSrv := apply.NewService(ds.mongoDb)
	//instanceId := arg.InstanceId[0:strings.Index(arg.InstanceId, ":")]
	err, deal := mongoSrv.NewDealLog(arg.InstanceId, arg.Id, claims.UserId, arg.Status)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	logrus.Info("deal info:", deal) //新增记录的id

	ds.ResponseSuccess(ctx, nil)
	return
}

// @Summary 是否含未处理任务
// @description true-含有未处理  false-无未处理
// @Tags apply
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param Authorization header string true "Bearer +token"
// @Success 200 {object} ApiResponse{result=boolean}
// @Router /api/apply/status [get]
func (ds *defaultServer) ApplyStatus(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*UserClaims)

	srv := flow.NewFlowService(ds.db)
	err, _, keys := srv.GetFlowListByAssignee(claims.UserId, claims.AppId)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	keyStr := strings.Join(keys, ",")

	bpmnSrv := bpmn.NewBpmnService(ds.conf.Engine.Api)
	err, ret := bpmnSrv.FlowTask("1", fmt.Sprintf("%d", utils.MAX_LIMIT), keyStr, "1")
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	status := false

	if ret.Count > 0 {
		status = true
	}
	ds.ResponseSuccess(ctx, status)
	return
}

type NotifierTaskItem struct {
	Id                  string      `json:"id"`                    //申请id
	ProcessDefinitionId string      `json:"process_definition_id"` //实例id  节点名称
	CreateTime          string      `json:"create_time"`           //发起时间
	InstanceId          string      `json:"instance_id"`           //实例id
	FlowName            string      `json:"flow_name"`             //流程名称
	FormName            string      `json:"form_name"`             //表单名称
	FormId              uint64      `json:"form_id"`               //表单id
	UserId              interface{} `json:"user_id"`               //用户id
	Account             string      `json:"account"`               //用户名
	Status              uint8       `json:"status"`                //处理状态 1待处理；2已处理
}

type NotifierTaskListVO struct {
	List  []NotifierTaskItem `json:"list"`
	Count uint32             `json:"count"`
}

// @Summary 抄送列表
// @Tags apply
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param Authorization header string true "Bearer +token"
// @param page query string true "页码"
// @param limit query string true "页容量"
// @Success 200 {object} ApiResponse{result=NotifierTaskListVO}
// @Router /api/apply/list/notifier [get]
func (ds *defaultServer) NotifierApplyList(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*UserClaims)
	pageStr := ctx.Query("page")
	limitStr := ctx.Query("limit")
	if pageStr == "" || limitStr == "" {
		ds.InvalidParametersError(ctx)
		return
	}

	//mysql获取抄送该用户的全部流程及任务ID
	srv := flow.NewFlowService(ds.db)
	err, flows, keys := srv.GetFlowListByNotifier(claims.UserId, claims.AppId)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	list := make([]NotifierTaskItem, 0, 100)
	if len(keys) <= 0 {
		ds.ResponseSuccess(ctx, NotifierTaskListVO{
			List:  list,
			Count: 0,
		})
		return
	}
	keyStr := strings.Join(keys, ",")

	//bpmn获取任务节点信息
	statusStr := "0" //status=0时则表示所有状态的任务
	bpmnSrv := bpmn.NewBpmnService(ds.conf.Engine.Api)
	err, ret := bpmnSrv.FlowTask(pageStr, limitStr, keyStr, statusStr)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	mongoSrv := form_data.NewFormDataService(ds.mongoDb)
	userSrv := user.NewUserService(ds.db)
	formSrv := form.NewFormService(ds.db)

	//根据实例id获取用户信息
	for _, x := range ret.TaskList {
		//判断流程的状态
		var status uint8
		if len(x.DeleteReason) == 0 {
			status = 1 //未处理
		} else {
			status = 2 //已处理
		}

		//mongo获取用户id
		found, item := mongoSrv.GetByInstanceId(x.ProcessInstanceId)
		if !found {
			ds.ResponseError(ctx, NOT_FOUND)
			return
		}

		//用户名称
		var account = ""
		userStr := fmt.Sprintf("%v", item["user_id"])
		userId := utils.NewStr(userStr)
		if userId.Uint64() != 0 {
			_, userItem := userSrv.GetUser(userId.Uint64())
			account = userItem.Account
		}

		//表单名称
		var formName = ""
		formId, ok := item["form_id"].(float64)
		if ok {
			_, formItem := formSrv.GetFormDetail(uint64(formId))
			formName = formItem.Name
		}

		//流程名称
		var name = ""
		for _, f := range flows {
			if f.Key == x.ProcessDefinitionId[0:strings.Index(x.ProcessDefinitionId, ":")] {
				name = f.Name
				break
			}
		}

		list = append(list, NotifierTaskItem{
			Id:                  x.Id,
			ProcessDefinitionId: x.ProcessDefinitionId,
			CreateTime:          x.CreateTime,
			InstanceId:          x.ProcessInstanceId,
			FlowName:            name,
			UserId:              item["user_id"],
			FormName:            formName,
			FormId:              uint64(formId),
			Account:             account,
			Status:              status,
		})
	}

	ds.ResponseSuccess(ctx, NotifierTaskListVO{
		List:  list,
		Count: ret.Count,
	})
	return
}

//type NotifierApplyDetail struct {
//	Variables []map[string]interface{} `json:"variables"`  //处理参数
//	CreatedAt time.Time                `json:"created_at"` //填写时间
//	Account   string                   `json:"account"`    //填写用户
//	Assignee  string                   `json:"assignee"`   //处理人
//	UpdatedAt time.Time                `json:"updated_at"` //处理时间
//	Status    uint8                    `json:"status"`     //流程状态
//}

// @Summary 抄送流程的详情
// @Tags apply
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param Authorization header string true "Bearer +token"
// @param id query string true "流程id"
// @param instance_id query string true "实例id"
// @param status query string true "流程状态(同列表)"
// @Success 200 {object} ApiResponse{result=ApplyDetail}
// @Router /api/apply/notifier/param [get]
func (ds *defaultServer) NotifierApplyParam(ctx *gin.Context) {
	// 复用 /api/apply/param [get]

	//taskIdStr := ctx.Query("id")
	//statusStr := ctx.Query("status")
	//if taskIdStr == "" || statusStr == "" {
	//	ds.InvalidParametersError(ctx)
	//	return
	//}
	//
	//bpmnSrv := bpmn.NewBpmnService(ds.conf.Engine.Api)
	//err, ret := bpmnSrv.FlowTaskVariables(taskIdStr, statusStr)
	//if err != nil {
	//	ds.InternalServiceError(ctx, err.Error())
	//	return
	//}
	//
	//result := ApplyDetail{
	//	Variables: ret,
	//}
	//
	//if statusStr == "2" {
	//	srv := apply.NewService(ds.mongoDb)
	//	found, dealItem := srv.GetLogByProcessId(taskIdStr)
	//	if found {
	//		userSrv := user.NewUserService(ds.db)
	//		result.UpdatedAt = dealItem.Updated_At
	//		f, assignee := userSrv.GetUser(dealItem.Assignee)
	//		if f {
	//			result.Assignee = assignee.TrueName
	//		}
	//		logSrv := form_data.NewFormDataService(ds.mongoDb)
	//		logFound, log := logSrv.GetByInstanceId(dealItem.InstanceId)
	//		if logFound {
	//			t, _ := log["created_at"].(primitive.DateTime)
	//			result.CreatedAt = t.Time()
	//			userStr := fmt.Sprintf("%v", log["user_id"])
	//			userId := utils.NewStr(userStr)
	//			userFound, users := userSrv.GetUser(userId.Uint64())
	//			if userFound {
	//				result.Account = users.TrueName
	//			}
	//		}
	//	}
	//}
	//
	//ds.ResponseSuccess(ctx, result)
	//return
}
