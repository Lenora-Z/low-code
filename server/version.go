//Created by Goland
//@User: lenora
//@Date: 2021/3/11
//@Time: 3:34 下午
package server

import (
	"github.com/Lenora-Z/low-code/service/application"
	"github.com/Lenora-Z/low-code/service/flow"
	"github.com/Lenora-Z/low-code/service/form"
	"github.com/Lenora-Z/low-code/service/navigation"
	"github.com/Lenora-Z/low-code/service/service"
	"github.com/Lenora-Z/low-code/service/version"
	"github.com/Lenora-Z/low-code/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type OnlineVersionVO struct {
	version.Version
	AppHash string `json:"app_hash"`
}

// @Summary 上线版本
// @Tags version
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @Success 200 {object} ApiResponse{result=OnlineVersionVO}
// @Router /version/online [get]
func (ds *defaultServer) OnlineVersion(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*CustomClaims)
	srv := version.NewVersionService(ds.db)
	_, item := srv.OnlineVersion(claims.AppId)
	if item.ID == 0 {
		ds.ResponseSuccess(ctx, nil)
	} else {
		ret := OnlineVersionVO{
			Version: item,
			AppHash: claims.UUid,
		}
		ds.ResponseSuccess(ctx, ret)
	}
	return
}

type VersionArg struct {
	Domain      string `json:"domain"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Random      string `json:"random"`
}

type VersionVO struct {
	Domain string `json:"domain"`
	Hash   string `json:"hash"`
}

// @Summary 发布新版本
// @Tags version
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @param VersionArg body VersionArg true "请求体"
// @Success 200 {object} ApiResponse{result=VersionVO}
// @Router /version/publish [post]
func (ds *defaultServer) PublishApplication(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*CustomClaims)
	var arg VersionArg
	if err := ctx.BindJSON(&arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}
	if err := validatorInstance().Struct(arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	db := ds.db.Begin()
	//新建版本
	srv := version.NewVersionService(db)
	err, item := srv.CreateVersion(claims.AppId, arg.Domain, arg.Version, arg.Description)
	if err != nil {
		db.Rollback()
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	appSrv := application.NewApplicationService(db)
	if status, _ := application.NewApplicationService(ds.db).GetApplication(claims.AppId); status {
		db.Rollback()
		ds.ResponseError(ctx, NOT_FOUND, "app not found")
		return
	}

	//表单上线
	formSrv := form.NewFormService(db)
	if err := formSrv.OnlineForm(claims.AppId); err != nil {
		db.Rollback()
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	//流程上线
	flowSrv := flow.NewFlowService(db)
	if err := flowSrv.OnlineFlow(claims.AppId); err != nil {
		db.Rollback()
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	//版本映射关系更新
	if err := flowSrv.OnlineMapping(claims.AppId, item.ID, arg.Random); err != nil {
		db.Rollback()
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	//参数依赖更新
	outSrv := service.NewOutsideService(db)
	if err := outSrv.OnlineRelies(claims.AppId, arg.Random); err != nil {
		db.Rollback()
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	//导航栏上线
	navSrv := navigation.NewNavigationService(db)
	if err := navSrv.UpdateAppNavigation(claims.AppId, true); err != nil {
		db.Rollback()
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	//应用上线
	if err, _ := appSrv.UpdateApplicationVersion(claims.AppId, item.ID, application.ONLINE); err != nil {
		db.Rollback()
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	db.Commit()
	ds.ResponseSuccess(ctx, VersionVO{
		Domain: item.Domain,
		Hash:   claims.UUid,
	})
	return
}

type VersionListVO struct {
	List  []version.Version `json:"list" validate:"required"`
	Count uint32            `json:"count" validate:"required"`
}

// @Summary 版本列表
// @Tags version
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @param page query string true "页码"
// @param limit query string true "页容量"
// @Success 200 {object} ApiResponse{result=VersionListVO}
// @Router /version/list [get]
func (ds *defaultServer) VersionHistory(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*CustomClaims)
	pageStr := ctx.Query("page")
	limitStr := ctx.Query("limit")
	if pageStr == "" || limitStr == "" {
		ds.InvalidParametersError(ctx)
		return
	}

	srv := version.NewVersionService(ds.db)
	err, count, list := srv.VersionHistory(utils.NewStr(pageStr).Uint32(), utils.NewStr(limitStr).Uint32(), claims.AppId)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	ds.ResponseSuccess(ctx, VersionListVO{
		List:  list,
		Count: count,
	})
	return
}

// @Summary 版本下线
// @Tags version
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @param IdStruct body IdStruct true "请求体"
// @Success 200 {object} ApiResponse{result=object}
// @Router /version/offline [post]
func (ds *defaultServer) OfflineApplication(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*CustomClaims)
	var arg IdStruct
	if err := ctx.BindJSON(&arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}
	if err := validatorInstance().Struct(arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	db := ds.db.Begin()
	srv := version.NewVersionService(db)
	//版本状态修改
	err, item := srv.UpdateVersionStatus(arg.Id, version.OFFLINE)
	if err != nil {
		db.Rollback()
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	appSrv := application.NewApplicationService(db)
	//应用状态修改
	if err, _ := appSrv.UpdateApplicationVersion(claims.AppId, item.ID, application.OFFLINE); err != nil {
		db.Rollback()
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	//表单下线
	formSrv := form.NewFormService(db)
	if err := formSrv.OfflineForm(claims.AppId); err != nil {
		db.Rollback()
		logrus.Error("form offline failed:", err)
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	//流程下线
	flowSrv := flow.NewFlowService(db)
	if err := flowSrv.OfflineFlow(claims.AppId); err != nil {
		db.Rollback()
		logrus.Error("flow offline failed:", err)
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	//映射关系下线
	if err := flowSrv.OfflineMapping(claims.AppId); err != nil {
		db.Rollback()
		logrus.Error("mapping offline failed:", err)
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	//参数依赖失效
	outSrv := service.NewOutsideService(db)
	if err := outSrv.OfflineRelies(claims.AppId); err != nil {
		db.Rollback()
		logrus.Error("relies offline failed:", err)
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	//导航栏下线
	navSrv := navigation.NewNavigationService(db)
	if err := navSrv.UpdateAppNavigation(claims.AppId, false); err != nil {
		db.Rollback()
		logrus.Error("navigation offline failed:", err)
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	db.Commit()
	ds.ResponseSuccess(ctx, nil)
	return
}
