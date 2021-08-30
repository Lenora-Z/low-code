package server

import (
	"github.com/Lenora-Z/low-code/service/user"
	"github.com/Lenora-Z/low-code/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type OrganizationArg struct {
	Name     string `json:"name" validate:"required"`
	ParentId uint64 `json:"parent_id"`
}

// @Summary 新增组织
// @Tags organization
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @param OrganizationArg body OrganizationArg true "请求体"
// @Success 200 {object} ApiResponse{result=object}
// @Router /organization/new [post]
func (ds *defaultServer) newOrganization(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*CustomClaims)
	var arg OrganizationArg
	if err := ctx.BindJSON(&arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}
	if err := validatorInstance().Struct(arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	if utils.Strlen(arg.Name) > 20 {
		logrus.Error("organization name is too long")
		ds.InvalidParametersError(ctx)
		return
	}

	srv := user.NewUserService(ds.db)
	err, item := srv.CreateOrganization(claims.AppId, arg.ParentId, arg.Name)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	ds.ResponseSuccess(ctx, item)
	return
}

type OrgItem struct {
	user.Organization
	Child []OrgItem `json:"child"`
}

// @Summary 组织列表
// @Tags organization
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @Success 200 {object} ApiResponse{result=[]OrgItem}
// @Router /organization/list [get]
func (ds *defaultServer) OrganizationList(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*CustomClaims)

	srv := user.NewUserService(ds.db)
	err, count, items := srv.OrganizationByParent(claims.AppId)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	//生成组织树
	list := make([]OrgItem, 0, 100)
	if count > 0 {
		for _, x := range items {
			if x.ParentID == 0 {
				item := OrgItem{x, []OrgItem{}}
				item.getChildren(items)
				list = append(list, item)
			}
		}
	}
	ds.ResponseSuccess(ctx, list)
	return
}

type OrganizationEditArg struct {
	Name string `json:"name" validate:"required"`
	Id   uint64 `json:"id"`
}

// @Summary 编辑组织
// @Tags organization
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @param OrganizationEditArg body OrganizationEditArg true "请求体"
// @Success 200 {object} ApiResponse{result=object}
// @Router /organization/edit [post]
func (ds *defaultServer) EditOrganization(ctx *gin.Context) {
	var arg OrganizationEditArg
	if err := ctx.BindJSON(&arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}
	if err := validatorInstance().Struct(arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	if utils.Strlen(arg.Name) > 20 {
		logrus.Error("organization name is too long")
		ds.InvalidParametersError(ctx)
		return
	}

	srv := user.NewUserService(ds.db)
	err, item := srv.UpdateOrganization(arg.Id, arg.Name)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	ds.ResponseSuccess(ctx, item)
	return
}

// @Summary 删除组织
// @Tags organization
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @param IdStruct body IdStruct true "请求体"
// @Success 200 {object} ApiResponse{result=object}
// @Router /organization/delete [post]
func (ds *defaultServer) DeleteOrganization(ctx *gin.Context) {
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

	srv := user.NewUserService(ds.db)
	//用户检测
	err, count, _ := srv.GetUserListByGroupId(arg.Id, claims.AppId)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	if count > 0 {
		ds.InternalServiceError(ctx, "该组织下有用户，不允许删除")
		return
	}

	//组织检测
	err, count, _ = srv.OrganizationByParent(claims.AppId, arg.Id)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	if count > 0 {
		ds.InternalServiceError(ctx, "该组织下有组织，不允许删除")
		return
	}

	if err := srv.DeleteOrganization(arg.Id); err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	ds.ResponseSuccess(ctx, nil)
	return
}

type OrgDetail struct {
	Id    uint64   `json:"id"`
	Name  string   `json:"name"`
	Roles []uint64 `json:"roles"` //角色id集合
}

// @Summary 组织详情
// @Tags organization
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @param id query string true "组织id"
// @Success 200 {object} ApiResponse{result=OrgDetail}
// @Router /organization/detail [get]
func (ds *defaultServer) OrganizationDetail(ctx *gin.Context) {
	idStr := ctx.Query("id")
	if idStr == "" {
		ds.InvalidParametersError(ctx)
		return
	}

	srv := user.NewUserService(ds.db)
	status, item := srv.OrganizationDetail(utils.NewStr(idStr).Uint64())
	if status {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}

	err, _, roles := srv.GetRelationByOrgId(item.ID)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	ds.ResponseSuccess(ctx, OrgDetail{
		Id:    item.ID,
		Name:  item.Name,
		Roles: roles.RoleIds(),
	})
	return
}

type OrgRoleArg struct {
	Id    uint64   `json:"id" validate:"required"`
	Roles []uint64 `json:"roles" validate:"required"`
}

// @Summary 组织角色配置
// @Tags organization
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @param OrgRoleArg body OrgRoleArg true "请求体"
// @Success 200 {object} ApiResponse{result=object}
// @Router /organization/config/role [post]
func (ds *defaultServer) ConfigOrgRole(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*CustomClaims)
	var arg OrgRoleArg
	if err := ctx.BindJSON(&arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}
	if err := validatorInstance().Struct(arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}
	db := ds.db.Begin()
	srv := user.NewUserService(db)
	if err := srv.DeleteOrgRoleByOrgId(arg.Id); err != nil {
		db.Rollback()
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	for _, x := range arg.Roles {
		if err, _ := srv.CreateOrgRole(claims.AppId, arg.Id, x); err != nil {
			db.Rollback()
			ds.InternalServiceError(ctx, err.Error())
			return
		}
	}
	db.Commit()
	ds.ResponseSuccess(ctx, nil)
	return
}

// @Summary 组织列表
// @Tags api/organization
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param Authorization header string true "Bearer +token"
// @Success 200 {object} ApiResponse{result=[]OrgItem}
// @Router /api/organization/list [get]
func (ds *defaultServer) ClientOrganizationList(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*UserClaims)

	srv := user.NewUserService(ds.db)
	err, count, items := srv.OrganizationByParent(claims.AppId)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	list := make([]OrgItem, 0, 100)
	if count > 0 {
		for _, x := range items {
			if x.ParentID == 0 {
				item := OrgItem{x, []OrgItem{}}
				item.getChildren(items)
				list = append(list, item)
			}
		}
	}
	ds.ResponseSuccess(ctx, list)
	return
}

func (options *OrgItem) getChildren(group []user.Organization) {
	for _, x := range group {
		if x.ParentID == options.ID {
			options.Child = append(options.Child, OrgItem{x, []OrgItem{}})
		}
	}
	if len(options.Child) > 0 {
		for i, _ := range options.Child {
			options.Child[i].getChildren(group)
		}
	}
	return
}
