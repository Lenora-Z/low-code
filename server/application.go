// Package server
//Created by Goland
//@User: lenora
//@Date: 2021/3/8
//@Time: 11:23 上午
package server

import (
	"github.com/Lenora-Z/low-code/service/application"
	"github.com/Lenora-Z/low-code/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AppArg struct {
	Name        string `json:"name" validate:"required"`        //应用名称
	Icon        string `json:"icon" validate:"required"`        //应用icon
	Description string `json:"description" validate:"required"` //应用描述
}

// @Summary 新增应用
// @Tags Application
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param AppArg body AppArg true "请求体"
// @Success 200 {object} ApiResponse{result=object}
// @Router /app/new [post]
func (ds *defaultServer) NewApplication(ctx *gin.Context) {
	var arg AppArg
	if err := ctx.BindJSON(&arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}
	if err := validatorInstance().Struct(arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	if utils.Strlen(arg.Name) > 10 {
		logrus.Error("application name is too long")
		ds.InvalidParametersError(ctx)
		return
	}

	srv := application.NewApplicationService(ds.db)
	//名称查重
	if status, _ := srv.GetApplicationByName(arg.Name); !status {
		ds.ResponseError(ctx, 3002)
		return
	}
	err, item := srv.CreateApplication(arg.Name, arg.Icon, arg.Description)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	ds.ResponseSuccess(ctx, item)
	return
}

type AppListVO struct {
	List  []application.Application `json:"list"`
	Count uint32                    `json:"count"`
}

// @Summary 应用列表
// @Tags Application
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @Success 200 {object} ApiResponse{result=AppListVO}
// @Router /app/list [get]
func (ds *defaultServer) ApplicationList(ctx *gin.Context) {
	srv := application.NewApplicationService(ds.db)
	err, count, list := srv.ApplicationFullList()
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	ds.ResponseSuccess(ctx, AppListVO{
		List:  list,
		Count: count,
	})
	return
}

type AppUpdateArg struct {
	IdStruct
	AppArg
}

// @Summary 编辑应用
// @Tags Application
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param AppUpdateArg body AppUpdateArg true "请求体"
// @Success 200 {object} ApiResponse{result=object}
// @Router /app/edit [post]
func (ds *defaultServer) EditApplication(ctx *gin.Context) {
	var arg AppUpdateArg
	if err := ctx.BindJSON(&arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}
	if err := validatorInstance().Struct(arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	if utils.Strlen(arg.Name) > 10 {
		logrus.Error("application name is too long")
		ds.InvalidParametersError(ctx)
		return
	}

	srv := application.NewApplicationService(ds.db)
	//名称查重
	if status, item := srv.GetApplicationByName(arg.Name); !status && item.ID != arg.Id {
		ds.ResponseError(ctx, 3002)
		return
	}

	err, item := srv.UpdateApplication(arg.Id, arg.Name, arg.Icon, arg.Description)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	ds.ResponseSuccess(ctx, item)
	return
}

// @Summary 应用详情
// @Tags Application
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @Success 200 {object} ApiResponse{result=application.Application}
// @Router /app/detail [get]
func (ds *defaultServer) ApplicationDetail(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*CustomClaims)
	srv := application.NewApplicationService(ds.db)
	status, item := srv.GetApplication(claims.AppId)
	if status {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}
	ds.ResponseSuccess(ctx, item)
	return
}

type ApplicationByHashArgs struct {
	Hash string `form:"hash" json:"hash" validate:"required"` //应用hash
}

// ApplicationByHash @Summary 编辑应用
// @Tags Application
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param ApplicationByHashArgs body ApplicationByHashArgs true "请求体"
// @Success 200 {object} ApiResponse{result=application.Application}
// @Router /api/hash [post]
func (ds *defaultServer) ApplicationByHash(c *gin.Context) {
	args := new(ApplicationByHashArgs)
	if err := c.BindJSON(args); err != nil {
		ds.InvalidParametersError(c)
		return
	}
	if err := validatorInstance().Struct(args); err != nil {
		ds.InvalidParametersError(c)
		return
	}
	as := application.NewApplicationService(ds.db)
	res, app := as.GetApplicationByHash(args.Hash)
	if res {
		ds.ResponseError(c, NOT_FOUND)
		return
	}
	ds.ResponseSuccess(c, app)
	return
}
