//Created by Goland
//@User: lenora
//@Date: 2021/3/15
//@Time: 2:55 下午
package server

import (
	"fmt"
	"github.com/Lenora-Z/low-code/service/bpmn"
	"github.com/Lenora-Z/low-code/service/flow"
	"github.com/Lenora-Z/low-code/service/service"
	"github.com/Lenora-Z/low-code/service/tritium"
	"github.com/Lenora-Z/low-code/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type FlowArg struct {
	Name string `json:"name" validate:"required"` //工作流名称
	Desc string `json:"desc"`                     //工作流描述
}

// NewFlow
// @Summary 新增流程[update]
// @Tags flow
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @param FlowArg body FlowArg true "请求体"
// @Success 200 {object} ApiResponse{result=object}
// @Router /flow/new [post]
func (ds *defaultServer) NewFlow(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*CustomClaims)
	var arg FlowArg
	if err := ctx.Bind(&arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}
	if err := validatorInstance().Struct(arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	querySrv := flow.NewFlowService(ds.db)
	if status, _ := querySrv.GetFlowDetailByName(arg.Name, claims.AppId); !status {
		ds.ResponseError(ctx, 3002)
		return
	}

	//新增流程
	srv := flow.NewFlowService(ds.db)
	err, item := srv.CreateFlow(claims.AppId, claims.UserId, arg.Name, fmt.Sprintf(`Process_%s`, utils.GetRandomStringSec(7)))
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	ds.ResponseSuccess(ctx, item)
	return
}

// NewFlowEasy
// @Summary 快速新增流程[new]
// @Tags flow
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @Success 200 {object} ApiResponse{result=object}
// @Router /flow/new/easy [post]
//func (ds *defaultServer) NewFlowEasy(ctx *gin.Context) {
//	claims := ctx.MustGet(CLAIMS).(*CustomClaims)
//
//	srv := flow.NewFlowService(ds.db)
//	err, item := srv.CreateFlow(
//		claims.AppId,
//		fmt.Sprintf(`未命名流程-%d`, utils.RandInt(100, 999)),
//		fmt.Sprintf(`Process_%s`, utils.GetRandomStringSec(7)),
//		"",
//		nil,
//		[]flow.ActAssignee{},
//		false,
//		[]flow.ActNotifier{},
//	)
//	if err != nil {
//		ds.InternalServiceError(ctx, err.Error())
//		return
//	}
//	ds.ResponseSuccess(ctx, item)
//	return
//}

type FlowItem struct {
	Id         uint64    `json:"id"`
	No         string    `json:"no"`     //流程编号
	Name       string    `json:"name"`   //流程名称
	Status     *bool     `json:"status"` //状态
	User       string    `json:"user"`   //创建人
	UpdateTime time.Time `json:"update_time"`
}

type FlowListVO struct {
	List  []FlowItem `json:"list"`
	Count uint32     `json:"count"`
}

// FlowList
// @Summary 流程列表[update]
// @Tags flow
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @param page query string false "页码"
// @param limit query string false "页容量"
// @param name query string false "流程名称"
// @Success 200 {object} ApiResponse{result=FlowListVO}
// @Router /flow/list [get]
func (ds *defaultServer) FlowList(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*CustomClaims)
	pageStr := ctx.Query("page")
	limitStr := ctx.Query("limit")
	name := ctx.Query("name")

	srv := flow.NewFlowService(ds.db)
	var err error
	var count uint32
	var list flow.FlowList
	users := make([]tritium.UserDetail, 0, 100)
	if pageStr == "" || limitStr == "" {
		err, count, list = srv.GetAllFlowList(claims.AppId)
	} else {
		err, count, list = srv.GetFlowList(utils.NewStr(pageStr).Uint32(), utils.NewStr(limitStr).Uint32(), claims.AppId, name)
	}

	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	if pageStr != "" && limitStr != "" {
		triSrv := tritium.NewTritiumService(ds.conf.TritiumConfig.Api)
		err, users = triSrv.BatchGetTriUserDetail(list.Users())
		if err != nil {
			ds.InternalServiceError(ctx, err.Error())
			return
		}
	}

	flowList := make([]FlowItem, 0, cap(list))
	for _, x := range list {
		var un string
		for _, u := range users {
			if x.UserID == u.Id {
				un = u.TrueName
				break
			}
		}
		flowList = append(flowList, FlowItem{
			Id:         x.ID,
			No:         x.Number,
			Name:       x.Name,
			User:       un,
			Status:     x.Status,
			UpdateTime: x.UpdatedAt,
		})
	}
	ds.ResponseSuccess(ctx, FlowListVO{
		List:  flowList,
		Count: count,
	})
	return
}

type FlowEditArg struct {
	IdStruct
	FlowArg
	Status bool   `json:"status"`
	Json   string `json:"json"`
}

// EditFlow
// @Summary 编辑流程[update]
// @Tags flow
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @param FlowEditArg body FlowEditArg true "请求体"
// @Success 200 {object} ApiResponse{result=object}
// @Router /flow/edit [post]
func (ds *defaultServer) EditFlow(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*CustomClaims)
	var arg FlowEditArg
	if err := ctx.BindJSON(&arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}
	if err := validatorInstance().Struct(arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	//判断流程状态
	querySrv := flow.NewFlowService(ds.db)
	status, f := querySrv.GetFlowDetail(arg.Id)
	if status {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}
	if *(f.IsOnline) == true {
		ds.ResponseError(ctx, 4002, "flow can't be changed")
		return
	}

	if status, f := querySrv.GetFlowDetailByName(arg.Name, claims.AppId); (!status) && (f.ID != arg.Id) {
		ds.ResponseError(ctx, 3002)
		return
	}

	var chainUp flow.ChainUpParams
	serviceSrv := service.NewOutsideService(ds.db)
	err, chainSrv := serviceSrv.GetServiceByType(service.CHAIN_UP)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	if len(chainSrv) <= 0 {
		logrus.Error("chain up service not found")
		ds.ResponseError(ctx, NOT_FOUND, "service not found")
		return
	}
	chainUp.ServiceId = chainSrv[0].ID
	err, chainParam := serviceSrv.GetParamList(chainUp.ServiceId, service.IN)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	if len(chainParam) <= 0 {
		logrus.Error("chain up param not found")
		ds.ResponseError(ctx, NOT_FOUND, "param not found")
		return
	}
	chainUp.ParamId = chainParam[0].ID

	var email flow.EmailParams
	err, emailSrv := serviceSrv.GetServiceByType(service.EMAIL)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	if len(emailSrv) <= 0 {
		logrus.Error("email service not found")
		ds.ResponseError(ctx, NOT_FOUND, "service not found")
		return
	}
	email.ServiceId = emailSrv[0].ID
	err, emailParam := serviceSrv.GetParamList(email.ServiceId, service.IN)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	if len(emailParam) < 4 {
		logrus.Error("get email param error")
		ds.InternalServiceError(ctx, "get param err")
		return
	}
	email.SenderId = emailParam[0].ID
	email.ReviverId = emailParam[1].ID
	email.Subject = emailParam[2].ID
	email.Content = emailParam[3].ID

	db := ds.db.Begin()
	srv := flow.NewFlowService(db)
	err, item := srv.UpdateFlow(
		arg.Id,
		arg.Name,
		arg.Json,
		arg.Status,
		chainUp,
		email,
	)
	if err != nil {
		db.Rollback()
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	if arg.Status {
		//部署至引擎
		bpmnSrv := bpmn.NewBpmnService(ds.conf.Engine.Api)
		err, key := bpmnSrv.DeployProcess(item.ID, item.Name, item.Desc, arg.Json)
		if err != nil {
			db.Rollback()
			ds.InternalServiceError(ctx, err.Error())
			return
		}

		if err, _ := srv.UpdateFlowKey(arg.Id, key); err != nil {
			db.Rollback()
			ds.InternalServiceError(ctx, err.Error())
			return
		}
	}

	db.Commit()
	ds.ResponseSuccess(ctx, nil)
	return
}

// @Summary 流程详情
// @Tags flow
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @param id query string true "流程id"
// @Success 200 {object} ApiResponse{result=flow.Flow}
// @Router /flow/detail [get]
func (ds *defaultServer) FlowDetail(ctx *gin.Context) {
	idStr := ctx.Query("id")
	if idStr == "" {
		ds.InvalidParametersError(ctx)
		return
	}

	id := utils.NewStr(idStr)
	srv := flow.NewFlowService(ds.db)
	status, item := srv.GetFlowDetail(id.Uint64())
	if status {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}
	ds.ResponseSuccess(ctx, item)
	return
}

type RelationItem struct {
	FormId uint64 `json:"form_id" validate:"required"` //表单id
	FlowId uint64 `json:"flow_id" validate:"required"` //流程id
}

type RelationArg struct {
	Relation []RelationItem `json:"relation" validate:"required"`
	Random   string         `json:"random" validate:"required"`
}

// @Summary 绑定表单流程映射
// @Tags maps
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @param RelationArg body RelationArg true "请求体"
// @Success 200 {object} ApiResponse{result=object}
// @Router /map/bind [post]
func (ds *defaultServer) FlowBound(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*CustomClaims)
	var arg RelationArg
	if err := ctx.BindJSON(&arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	db := ds.db.Begin()
	srv := flow.NewFlowService(db)
	for _, x := range arg.Relation {
		if err, _ := srv.CreateFlowMapping(x.FormId, x.FlowId, claims.AppId, arg.Random); err != nil {
			db.Rollback()
			ds.InternalServiceError(ctx, err.Error())
			return
		}
	}

	db.Commit()
	ds.ResponseSuccess(ctx, nil)
	return
}

// @Summary 已上线的映射关系
// @Tags maps
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @Success 200 {object} ApiResponse{result=[]flow.FlowMapping}
// @Router /map/list [get]
func (ds *defaultServer) OnlineMaps(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*CustomClaims)

	srv := flow.NewFlowService(ds.db)

	err, _, list := srv.GetOnlineMapping(claims.AppId)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	ds.ResponseSuccess(ctx, list)
	return
}

type LogItem struct {
	ProcessInstanceId    string `json:"process_instance_id"`    //实例id/实例编号
	ProcessDefinitionKey string `json:"process_definition_key"` //流程key
	StartTime            string `json:"start_time"`             //开始时间
	EndTime              string `json:"end_time"`               //结束时间
	State                string `json:"state"`                  //状态
	FlowNumber           string `json:"flow_number"`            //流程编号
	FlowName             string `json:"flow_name"`              //流程名称
}

type LogListVO struct {
	List  []LogItem `json:"list"`
	Count uint32    `json:"count"`
}

// @Summary 流程消息列表
// @Tags flow
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @param page query string true "页码"
// @param limit query string true "页容量"
// @Success 200 {object} ApiResponse{result=LogListVO}
// @Router /flow/record [get]
func (ds *defaultServer) FlowLogs(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*CustomClaims)
	pageStr := ctx.Query("page")
	limitStr := ctx.Query("limit")
	if pageStr == "" || limitStr == "" {
		ds.InvalidParametersError(ctx)
		return
	}

	//
	srv := flow.NewFlowService(ds.db)
	err, _, flows := srv.GetAllFlowList(claims.AppId)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	keyGroup := (&flows).Keys()
	keys := strings.Join(keyGroup, ",")

	bpmnSrv := bpmn.NewBpmnService(ds.conf.Engine.Api)
	err, ret := bpmnSrv.ProcessLog(pageStr, limitStr, keys)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	list := make([]LogItem, 0, cap(ret.ProcessInstances))
	for _, x := range ret.ProcessInstances {
		status, item := srv.GetFlowDetailByKey(x.ProcessDefinitionKey)
		if status {
			ds.ResponseError(ctx, NOT_FOUND)
			return
		}
		list = append(list, LogItem{
			ProcessInstanceId:    x.ProcessInstanceId,
			ProcessDefinitionKey: x.ProcessDefinitionKey,
			StartTime:            x.StartTime,
			EndTime:              x.EndTime,
			State:                x.State,
			FlowNumber:           item.Number,
			FlowName:             item.Name,
		})
	}
	ds.ResponseSuccess(ctx, LogListVO{
		List:  list,
		Count: ret.Count,
	})
	return

}

// @Summary 流程消息日志
// @Tags flow
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @param id query string true "实例id"
// @Success 200 {object} ApiResponse{result=interface{}}
// @Router /flow/record/log [get]
func (ds *defaultServer) FlowLogsDetail(ctx *gin.Context) {
	idStr := ctx.Query("id")

	srv := bpmn.NewBpmnService(ds.conf.Engine.Api)
	err, item := srv.LogDetail(idStr)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	ds.ResponseSuccess(ctx, item)
	return
}

// DeleteFlow
// @Summary 流程删除[new]
// @Tags flow
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @param IdStruct body IdStruct true "id"
// @Success 200 {object} ApiResponse{result=object}
// @Router /flow/delete [post]
func (ds *defaultServer) DeleteFlow(ctx *gin.Context) {
	var arg IdStruct
	if err := ctx.BindJSON(&arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}
	if err := validatorInstance().Struct(arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	srv := flow.NewFlowService(ds.db)
	notFound, item := srv.GetFlowDetail(arg.Id)
	if notFound || *item.IsDelete == true {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}
	if *item.IsOnline == true {
		ds.ResponseError(ctx, 3003, "form can't be deleted")
		return
	}

	if err := srv.DelFlow(arg.Id); err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	ds.ResponseSuccess(ctx, nil)
	return

}
