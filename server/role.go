package server

import (
	"github.com/Lenora-Z/low-code/service/user"
	"github.com/Lenora-Z/low-code/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type RoleArg struct {
	Name        string              `json:"name" validate:"required"`
	Permissions map[string][]uint64 `json:"permissions" validate:"required"`
}

// @Summary 新增角色
// @Tags role
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @param RoleArg body RoleArg true "请求体"
// @Success 200 {object} ApiResponse{result=object}
// @Router /role/new [post]
func (ds *defaultServer) NewRole(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*CustomClaims)
	var arg RoleArg
	if err := ctx.BindJSON(&arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}
	if err := validatorInstance().Struct(arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	if utils.Strlen(arg.Name) > 10 {
		logrus.Error("role name is too long")
		ds.InvalidParametersError(ctx)
		return
	}

	db := ds.db.Begin()
	srv := user.NewUserService(db)
	permission := make(map[uint64][]uint64, 0)
	for i, x := range arg.Permissions {
		key := utils.NewStr(i).Uint64()
		permission[key] = x
	}
	//唯一性检查
	res, _ := srv.RoleDetailByName(arg.Name, claims.AppId)
	if !res {
		db.Rollback()
		ds.ResponseError(ctx, 3002)
		return
	}
	err, item := srv.CreateRole(claims.AppId, arg.Name, permission)
	if err != nil {
		db.Rollback()
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	db.Commit()
	ds.ResponseSuccess(ctx, item)
	return
}

type RoleEditArg struct {
	IdStruct
	RoleArg
}

// @Summary 编辑角色
// @Tags role
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @param RoleEditArg body RoleEditArg true "请求体"
// @Success 200 {object} ApiResponse{result=object}
// @Router /role/edit [post]
func (ds *defaultServer) EditRole(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*CustomClaims)
	var arg RoleEditArg
	if err := ctx.BindJSON(&arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}
	if err := validatorInstance().Struct(arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	if utils.Strlen(arg.Name) > 10 {
		logrus.Error("role name is too long")
		ds.InvalidParametersError(ctx)
		return
	}

	db := ds.db.Begin()
	srv := user.NewUserService(db)
	permission := make(map[uint64][]uint64, 0)
	for i, x := range arg.Permissions {
		key := utils.NewStr(i).Uint64()
		permission[key] = x
	}
	//唯一性检查
	if res, roles := user.NewUserService(ds.db).RoleDetailByName(arg.Name, claims.AppId); !res && roles.ID != arg.Id {
		ds.ResponseError(ctx, 3002)
		return
	}
	err, item := srv.UpdateRole(arg.Id, arg.Name, permission)
	if err != nil {
		db.Rollback()
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	db.Commit()
	ds.ResponseSuccess(ctx, item)
	return
}

type RoleListItem struct {
	user.Role
	OrganizationGroup []string `json:"organization_group"`
}

type RoleListVO struct {
	List  []RoleListItem `json:"list"`
	Count uint32         `json:"count"`
}

// @Summary 角色列表
// @Tags role
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @param page query string false "页码"
// @param limit query string false "页容量"
// @param name query string false "角色名"
// @Success 200 {object} ApiResponse{result=RoleListVO}
// @Router /role/list [get]
func (ds *defaultServer) RoleList(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*CustomClaims)
	pageStr := ctx.Query("page")
	limitStr := ctx.Query("limit")
	nameStr := ctx.Query("name")

	srv := user.NewUserService(ds.db)
	var err error
	var count uint32
	var list user.RoleList
	if pageStr == "" || limitStr == "" {
		err, count, list = srv.FullRoleList(claims.AppId)
	} else {
		err, count, list = srv.RoleList(utils.NewStr(pageStr).Uint32(), utils.NewStr(limitStr).Uint32(), claims.AppId, nameStr)
	}
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	retList := make([]RoleListItem, 0, cap(list))
	for _, x := range list {
		oList := make(user.OrganizationList, 0)
		if pageStr != "" && limitStr != "" {
			err, _, rList := srv.GetRelationByRoleId([]uint64{x.ID})
			if err != nil {
				ds.InternalServiceError(ctx, err.Error())
				return
			}
			if len(rList) > 0 {
				err, oList = srv.GetOrganizationByIdGroup((&rList).OrgIds())
				if err != nil {
					ds.InternalServiceError(ctx, err.Error())
					return
				}
			}
		}
		retList = append(retList, RoleListItem{
			Role:              x,
			OrganizationGroup: (oList).Names(),
		})
	}

	ds.ResponseSuccess(ctx, RoleListVO{
		List:  retList,
		Count: count,
	})
	return
}

// @Summary 删除角色
// @Tags role
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @param IdStruct body IdStruct true "请求体"
// @Success 200 {object} ApiResponse{result=object}
// @Router /role/delete [post]
func (ds *defaultServer) DeleteRole(ctx *gin.Context) {
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
	if err := srv.DeleteRole(arg.Id); err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	ds.ResponseSuccess(ctx, nil)
	return

}

type RoleDetail struct {
	Id          uint64              `json:"id"`
	Name        string              `json:"name"`
	Permissions map[uint64][]uint64 `json:"permissions"` //权限id集合
}

// @Summary 角色详情
// @Tags role
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @param id query string true "角色id"
// @Success 200 {object} ApiResponse{result=RoleDetail}
// @Router /role/detail [get]
func (ds *defaultServer) RoleDetail(ctx *gin.Context) {
	idStr := ctx.Query("id")
	if idStr == "" {
		ds.InvalidParametersError(ctx)
		return
	}

	srv := user.NewUserService(ds.db)
	status, item := srv.RoleDetail(utils.NewStr(idStr).Uint64())
	if status {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}

	list := make(map[uint64][]uint64, 0)
	for i := 1; i <= 4; i++ {
		_, p := srv.PermissionOfRoles(uint64(i), []uint64{item.ID})
		list[uint64(i)] = p.Ids()
	}

	ds.ResponseSuccess(ctx, RoleDetail{
		Id:          item.ID,
		Name:        item.Name,
		Permissions: list,
	})
	return
}
