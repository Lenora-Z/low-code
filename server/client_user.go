//Created by Goland
//@User: lenora
//@Date: 2021/3/16
//@Time: 7:20 下午
package server

import (
	"github.com/Lenora-Z/low-code/service/application"
	"github.com/Lenora-Z/low-code/service/form"
	"github.com/Lenora-Z/low-code/service/user"
	"github.com/Lenora-Z/low-code/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"strings"
)

type UserLoginArg struct {
	Account  string `json:"account" validate:"required"`
	Password string `json:"password" validate:"required"`
	AppHash  string `json:"app_hash" validate:"required"`
}

// @Summary 用户登录
// @Tags client/user
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app_hash header string true "应用hash"
// @param UserLoginArg body UserLoginArg true "请求体"
// @Success 200 {object} ApiResponse{result=string}
// @Router /api/user/login [post]
func (ds *defaultServer) UserLogin(ctx *gin.Context) {
	var err error
	var arg UserLoginArg
	if err = ctx.BindJSON(&arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}
	if err = validatorInstance().Struct(arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	//获取用户对应的应用并判断应用状态
	appSrv := application.NewApplicationService(ds.db)
	status, appItem := appSrv.GetApplicationByHash(arg.AppHash)
	if status {
		logrus.Error("app is not exist")
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}

	if appItem.Status != application.ONLINE {
		logrus.Error("app is not online")
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}

	//用户校验
	srv := user.NewUserService(ds.db)
	status, item := srv.GetUserByAccount(arg.Account, appItem.ID)
	if status {
		logrus.Error("user not found")
		ds.ResponseError(ctx, 2002)
		return
	}
	if status := utils.CheckPassword(item.Password, arg.Password); !status {
		logrus.Error("password incorrect")
		ds.ResponseError(ctx, 2002)
		return
	}
	if *item.PwdStatus {
		item.PwdStatus = &FALSE
		err, item = srv.UpdatePwdStatus(item.ID, item.PwdStatus)
		if err != nil {
			logrus.Error("password incorrect")
			ds.InternalServiceError(ctx, err.Error())
			return
		}
	}

	//生成登录token
	token, err := ds.CreateToken(item.ID, appItem.ID, item.Mobile, item.Account)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	ds.ResponseSuccess(ctx, token)
	return
}

type UserPwdArg struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}

// @Summary 用户修改密码
// @Tags client/user
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param Authorization header string true "Bearer +token"
// @param UserPwdArg body UserPwdArg true "请求体"
// @Success 200 {object} ApiResponse{result=object}
// @Router /api/user/password/reset [post]
func (ds *defaultServer) UserChangePwd(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*UserClaims)
	var arg UserPwdArg
	if err := ctx.BindJSON(&arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}
	if err := validatorInstance().Struct(arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	srv := user.NewUserService(ds.db)
	status, item := srv.GetUser(claims.UserId)
	if status {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}
	if status := utils.CheckPassword(item.Password, arg.OldPassword); !status {
		logrus.Error("password incorrect")
		ds.ResponseError(ctx, 2002)
		return
	}

	if err, _ := srv.ResetUserPwd(claims.UserId, arg.NewPassword); err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	ds.ResponseSuccess(ctx, nil)
	return
}

type Permission struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
}

type UserPermission struct {
	Permission
	Child []Permission `json:"child"`
}

type UserPermissionVO []UserPermission

// @Summary 用户权限
// @Tags client/user
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param Authorization header string true "Bearer +token"
// @Success 200 {object} ApiResponse{result=UserPermissionVO}
// @Router /api/user/permission [get]
func (ds *defaultServer) UserPermission(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*UserClaims)

	//获取用户组织
	srv := user.NewUserService(ds.db)
	status, item := srv.GetUser(claims.UserId)
	if status {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}

	//根据组织信息获取到相应全部角色的权限(一级权限)
	err, _, roles := srv.GetRelationByOrgId(item.GroupID)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	err, permissions := srv.PermissionOfRoles(0, roles.RoleIds())
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	ids := permissions.Ids()

	ret := make([]UserPermission, 0, 50)
	if len(ids) > 0 {
		//权限-路由信息获取
		err, _, list := srv.GetMultiPermission(ids)
		if err != nil {
			ds.InternalServiceError(ctx, err.Error())
			return
		}

		formSrv := form.NewFormService(ds.db)
		//遍历获取子集路由并获取对应路由信息
		for _, v := range list {
			if v.ID == DEAL_APPLY {
				err, permis := srv.PermissionOfRoles(v.ID, roles.RoleIds())
				if err != nil {
					ds.InternalServiceError(ctx, err.Error())
					return
				}
				err, _, forms := formSrv.MultiFormList(claims.AppId, permis.Ids())
				if err != nil {
					ds.InternalServiceError(ctx, err.Error())
					return
				}
				child := make([]Permission, 0, cap(forms))
				for _, t := range forms {
					child = append(child, Permission{
						Id:   t.ID,
						Name: t.Name,
					})
				}
				ret = append(ret, UserPermission{
					Permission: Permission{
						Id:   v.ID,
						Name: v.Name,
					},
					Child: child,
				})
			}
		}
	}
	ds.ResponseSuccess(ctx, ret)
	return
}

type UserInfoVO struct {
	Id               uint64 `json:"id"`                //用户id
	TrueName         string `json:"true_name"`         //真实姓名
	Nickname         string `json:"nickname"`          //昵称
	Mobile           string `json:"mobile"`            //手机号
	AppId            uint64 `json:"app_id"`            //应用id
	AppName          string `json:"app_name"`          //应用名称
	OrganizationName string `json:"organization_name"` //组织名称
}

// @Summary 用户信息
// @Tags client/user
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param Authorization header string true "Bearer +token"
// @Success 200 {object} ApiResponse{result=UserInfoVO}
// @Router /api/user/info [get]
func (ds *defaultServer) UserInfo(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*UserClaims)
	srv := user.NewUserService(ds.db)
	status, item := srv.GetUser(claims.UserId)
	if status || (*item.IsDelete) == true {
		logrus.Error("user not exist")
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}

	_, org := srv.OrganizationDetail(item.GroupID)

	ds.ResponseSuccess(ctx, UserInfoVO{
		Id:               claims.UserId,
		TrueName:         item.TrueName,
		Nickname:         item.Nickname,
		Mobile:           item.Nickname,
		AppId:            claims.AppId,
		AppName:          claims.AppName,
		OrganizationName: org.Name,
	})
	return

}

type ClientUserItem struct {
	Id       uint64 `json:"id"`        //用户id
	TrueName string `json:"true_name"` //真实姓名
}

type ClientUserListVO struct {
	List  []ClientUserItem `json:"list"`
	Count uint32           `json:"count"`
}

// @Summary 用户列表[update]
// @Tags client/user
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param Authorization header string true "Bearer +token"
// @param name query string false "用户名称"
// @param org_id query string false "组织id,多个用','分隔"
// @Success 200 {object} ApiResponse{result=ClientUserListVO}
// @Router /api/user/list [get]
func (ds *defaultServer) ClientUserList(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*UserClaims)
	name := ctx.Query("name")
	group := ctx.Query("org_id")
	var err error
	var count uint32
	var users user.UserList
	srv := user.NewUserService(ds.db)
	if group == "" {
		err, count, users = srv.GetAllUser(claims.AppId, name)
	} else {
		var groupIds = make([]uint64, len(strings.Split(group, ",")))
		for _, groupId := range strings.Split(group, ",") {
			groupIds = append(groupIds, utils.NewStr(groupId).Uint64())
		}
		err, count, users = srv.GetAllUser(claims.AppId, name, groupIds...)
	}
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	userList := make([]ClientUserItem, 0, cap(users))
	for _, x := range users {
		userList = append(userList, ClientUserItem{
			Id:       x.ID,
			TrueName: x.TrueName,
		})
	}
	ds.ResponseSuccess(ctx, ClientUserListVO{
		List:  userList,
		Count: count,
	})
	return
}
