package server

import (
	"github.com/Lenora-Z/low-code/service/application"
	"github.com/Lenora-Z/low-code/service/navigation"
	"github.com/Lenora-Z/low-code/service/routing"
	"github.com/Lenora-Z/low-code/service/user"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// EditNavigation
// @Summary 编辑导航栏[update]
// @Tags navigation
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @param routing.RouterItemList body routing.RouterItemList true "请求体"
// @Success 200 {object} ApiResponse{result=object}
// @Router /nav/edit [post]
func (ds *defaultServer) EditNavigation(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*CustomClaims)
	var arg routing.RouterItemList
	if err := ctx.BindJSON(&arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}
	//if err := validatorInstance().Struct(arg); err != nil {
	//	ds.InvalidParametersError(ctx)
	//	return
	//}

	//已上线的应用导航的不可修改
	if claims.AppStatus == application.ONLINE {
		logrus.Error("app is online")
		ds.ResponseError(ctx, 3003, "navigation can't be changed")
		return
	}

	db := ds.db.Begin()
	//路由信息更新
	routingSrv := routing.NewRoutingService(db)
	if err := routingSrv.UpdateRouting(claims.AppId, arg); err != nil {
		db.Rollback()
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	db.Commit()
	ds.ResponseSuccess(ctx, nil)
	return

}

type NavigationVO struct {
	*navigation.Navigation
	Content routing.RouterItemList `json:"content"`
}

// @Summary 获取应用导航栏
// @Tags api/navigation
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param Authorization header string true "Bearer +token"
// @Success 200 {object} ApiResponse{result=[]RouterItem}
// @Router /api/nav [get]
func (ds *defaultServer) AppNavigation(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*UserClaims)

	userSrv := user.NewUserService(ds.db)
	err, _, relation := userSrv.GetRelationByOrgId(claims.GroupId)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	err, Permissions := userSrv.PermissionOfRoles(4, relation.RoleIds())
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	srv := routing.NewRoutingService(ds.db)
	err, routes := srv.GetRouteGroup(Permissions.Ids())
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	//生成路由树
	list := make([]RouterItem, 0, cap(routes))
	if len(routes) > 0 {
		for _, x := range routes {
			if x.ParentID == 0 {
				item := RouterItem{x, []RouterItem{}}
				item.getChildren(routes)
				list = append(list, item)
			}
		}
	}
	ds.ResponseSuccess(ctx, list)
	return
}
