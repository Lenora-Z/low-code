//Created by Goland
//@User: lenora
//@Date: 2021/3/11
//@Time: 3:24 下午
package server

import (
	"github.com/Lenora-Z/low-code/service/service"
	"github.com/Lenora-Z/low-code/utils"
	"github.com/gin-gonic/gin"
)

// @Summary 服务列表(全部)
// @Tags service
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @Success 200 {object} ApiResponse{result=[]service.Service}
// @Router /service/list/all [get]
func (ds *defaultServer) ServiceList(ctx *gin.Context) {
	srv := service.NewOutsideService(ds.db)
	err, list := srv.GetAllServiceList()
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	ds.ResponseSuccess(ctx, list)
	return
}

type Params struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
}

type AppServiceItem struct {
	FlowId      uint64   `json:"flow_id"`
	FlowName    string   `json:"flow_name"`
	ServiceId   uint64   `json:"service_id"`
	ServiceName string   `json:"service_name"`
	Params      []Params `json:"params"`
}

type AppServiceListVO []AppServiceItem

type FormFieldArg struct {
	FormId  uint64 `json:"form_id"`
	FieldId uint64 `json:"field_id"`
}

type Relies struct {
	FlowId    uint64         `json:"flow_id"`
	ServiceId uint64         `json:"service_id"`
	ParamId   uint64         `json:"param_id"`
	Relies    []FormFieldArg `json:"relies"`
}

type NewRelyArg struct {
	Random   string
	Relation []Relies `json:"relation"`
}

// @Summary 服务参数依赖绑定
// @Tags service
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @param NewRelyArg body NewRelyArg true "请求体"
// @Success 200 {object} ApiResponse{result=object}
// @Router /service/param/bind [post]
func (ds *defaultServer) NewParamRelies(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*CustomClaims)
	var arg NewRelyArg
	if err := ctx.BindJSON(&arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}
	//if err := validatorInstance().Struct(&arg); err != nil {
	//	logrus.Error(err)
	//	ds.InvalidParametersError(ctx)
	//	return
	//}

	db := ds.db.Begin()
	srv := service.NewOutsideService(db)
	for _, x := range arg.Relation {
		for _, v := range x.Relies {
			if err, _ := srv.CreateParamRely(claims.AppId, x.FlowId, x.ServiceId, x.ParamId, v.FormId, v.FieldId, arg.Random); err != nil {
				db.Rollback()
				ds.InternalServiceError(ctx, err.Error())
				return
			}
		}
	}

	db.Commit()
	ds.ResponseSuccess(ctx, nil)
	return
}

// ServiceParams
// @Summary 服务参数列表
// @Tags service
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @param id query string true "服务id,传0时获取全部服务"
// @param type query string true "入/出参   in/out"
// @Success 200 {object} ApiResponse{result=[]service.Param}
// @Router /service/params [get]
func (ds *defaultServer) ServiceParams(ctx *gin.Context) {
	idStr := ctx.Query("id")
	t := ctx.Query("type")
	if t == "" || idStr == "" {
		ds.InvalidParametersError(ctx)
		return
	}

	id := utils.NewStr(idStr).Uint64()
	types := service.OUT
	if t == "in" {
		types = service.IN
	}

	srv := service.NewOutsideService(ds.db)
	err, list := srv.GetParamList(id, types)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	ds.ResponseSuccess(ctx, list)
	return
}

type PackageItem struct {
	Id     uint64          `json:"id"`     //代码包id
	Name   string          `json:"name"`   //代码包名称
	Params []service.Param `json:"params"` //代码包出参
}

// PackageService
// @Summary 代码包列表
// @Tags service
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @Success 200 {object} ApiResponse{result=[]PackageItem}
// @Router /service/package [get]
func (ds *defaultServer) PackageService(ctx *gin.Context) {
	srv := service.NewOutsideService(ds.db)
	err, srvs := srv.GetServiceByType(service.PACKAGE)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	err, params := srv.GetParamListBySrvGroup(srvs.Ids())
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	list := make([]PackageItem, 0, cap(srvs))
	for _, s := range srvs {
		param := make([]service.Param, 0, cap(params))
		for _, p := range params {
			if p.ServiceID == s.ID {
				param = append(param, p)
			}
		}
		list = append(list, PackageItem{
			Id:     s.ID,
			Name:   s.Name,
			Params: param,
		})
	}
	ds.ResponseSuccess(ctx, list)
	return

}
