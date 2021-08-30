//Created by Goland
//@User: lenora
//@Date: 2021/3/9
//@Time: 11:48 上午
package server

import (
	"github.com/Lenora-Z/low-code/service/user"
	"github.com/Lenora-Z/low-code/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
)

type UserArg struct {
	GroupId  uint64 `json:"group_id" validate:"required"`
	Account  string `json:"account" validate:"required"`
	NickName string `json:"nickname"`
	TrueName string `json:"true_name" validate:"required"`
	Mobile   string `json:"mobile"`
	Mail     string `json:"mail"`
}

// @Summary 新增用户
// @Tags user
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @param UserArg body UserArg true "请求体"
// @Success 200 {object} ApiResponse{result=object}
// @Router /user/new [post]
func (ds *defaultServer) NewUser(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*CustomClaims)
	var arg UserArg
	if err := ctx.BindJSON(&arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}
	if err := validatorInstance().Struct(arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	if utils.Strlen(arg.Account) > 25 {
		logrus.Error("user account is too long")
		ds.InvalidParametersError(ctx)
		return
	}

	srv := user.NewUserService(ds.db)
	//用户账号查重
	if status, _ := srv.GetUserByAccount(arg.Account, claims.AppId); !status {
		logrus.Error("the user account is exist:", arg.Account)
		ds.ResponseError(ctx, 3002)
		return
	}
	//真实姓名查重
	if status, _ := srv.GetUserByTrueName(arg.TrueName, claims.AppId); !status {
		logrus.Error("the user true name is exist:", arg.TrueName)
		ds.ResponseError(ctx, 3002)
		return
	}
	err, _ := srv.CreateUser(
		arg.Account,
		arg.NickName,
		arg.TrueName,
		arg.Mobile,
		arg.Mail,
		claims.AppId,
		arg.GroupId,
	)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	ds.ResponseSuccess(ctx, nil)
	return
}

type UserListArg struct {
	Page    uint32   `json:"page" validate:"required"`  //页码
	Limit   uint32   `json:"limit" validate:"required"` //页容量
	Name    string   `json:"name"`                      //用户名称
	GroupId uint64   `json:"group_id"`                  //组织id
	Roles   []uint64 `json:"roles"`                     //角色集合
}

type UserListItem struct {
	Id        uint64    `json:"id"`
	Account   string    `json:"account"`    //账号/用户名
	NickName  string    `json:"nickname"`   //昵称
	TrueName  string    `json:"true_name"`  //真实姓名
	CreatedAt time.Time `json:"created_at"` //创建时间
	UpdatedAt time.Time `json:"updated_at"` //更新时间
	GroupName string    `json:"group_name"` //组织
	RoleName  []string  `json:"role_name"`  //角色
}

type UserListVO struct {
	List  []UserListItem `json:"list"`
	Count uint32         `json:"count"`
}

// @Summary 用户列表
// @Tags user
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @param UserListArg body UserListArg true "请求体"
// @Success 200 {object} ApiResponse{result=UserListVO}
// @Router /user/list [post]
func (ds *defaultServer) UserList(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*CustomClaims)
	var arg UserListArg
	if err := ctx.BindJSON(&arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}
	if err := validatorInstance().Struct(arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	srv := user.NewUserService(ds.db)
	//根据角色获取到的组织id
	var RolesGroupId, OrgGroupId []uint64
	group := make([]uint64, 0)
	if len(arg.Roles) > 0 {
		err, _, rList := srv.GetRelationByRoleId(arg.Roles)
		if err != nil {
			ds.InternalServiceError(ctx, err.Error())
			return
		}
		RolesGroupId = (&rList).OrgIds()
	}
	if arg.GroupId != 0 {
		OrgGroupId = srv.GetOrganizationRecursively(arg.GroupId, claims.AppId, []uint64{})
	}
	if (len(arg.Roles) > 0) && (arg.GroupId != 0) {
		temp := utils.IntersectionInt(RolesGroupId, OrgGroupId)
		if temp != nil {
			group = temp
		}
	} else if arg.GroupId == 0 {
		group = RolesGroupId
	} else {
		group = OrgGroupId
	}

	err, count, users := srv.GetUserList(arg.Page, arg.Limit, claims.AppId, arg.Name, group)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	err, groupList := srv.GetOrganizationByIdGroup(users.Groups())
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	list := make([]UserListItem, 0, cap(users))
	for _, x := range users {
		err, _, rList := srv.GetRelationByOrgId(x.GroupID)
		if err != nil {
			ds.InternalServiceError(ctx, err.Error())
			return
		}

		err, roleList := srv.GetRoleByIdGroup((&rList).RoleIds())
		if err != nil {
			ds.InternalServiceError(ctx, err.Error())
			return
		}

		list = append(list, UserListItem{
			Id:        x.ID,
			Account:   x.Account,
			NickName:  x.Nickname,
			TrueName:  x.TrueName,
			CreatedAt: x.CreatedAt,
			UpdatedAt: x.UpdatedAt,
			GroupName: groupList.Item(x.GroupID).Name,
			RoleName:  (&roleList).Names(),
		})
	}
	ds.ResponseSuccess(ctx, UserListVO{
		List:  list,
		Count: count,
	})
	return
}

type UserAttr struct {
	Id       uint64 `json:"id"`
	Account  string `json:"account"`   //账号
	TrueName string `json:"true_name"` //真实姓名
	Mail     string `json:"mail"`      //邮箱
}

// AllUserList
// @Summary 全部用户[update]
// @Tags user
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @param name query string false "用户真实姓名"
// @Success 200 {object} ApiResponse{result=[]UserAttr}
// @Router /user/all [get]
func (ds *defaultServer) AllUserList(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*CustomClaims)
	name := ctx.Query("name")

	srv := user.NewUserService(ds.db)
	err, _, users := srv.GetAllUser(claims.AppId, name)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	list := make([]UserAttr, 0, cap(users))
	for _, x := range users {
		list = append(list, UserAttr{
			Id:       x.ID,
			Account:  x.Account,
			TrueName: x.TrueName,
			Mail:     x.Mail,
		})
	}
	ds.ResponseSuccess(ctx, list)
	return
}

type UserEditArg struct {
	IdStruct
	UserArg
}

// @Summary 编辑用户
// @Tags user
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @param UserEditArg body UserEditArg true "请求体"
// @Success 200 {object} ApiResponse{result=object}
// @Router /user/edit [post]
func (ds *defaultServer) EditUser(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*CustomClaims)
	var arg UserEditArg
	if err := ctx.BindJSON(&arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}
	if err := validatorInstance().Struct(arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	if utils.Strlen(arg.Account) > 25 {
		logrus.Error("user account is too long")
		ds.InvalidParametersError(ctx)
		return
	}

	srv := user.NewUserService(ds.db)
	if status, item := srv.GetUserByAccount(arg.Account, claims.AppId); (!status) && (item.ID != arg.Id) {
		logrus.Error("the user account is exist:", arg.Account)
		ds.ResponseError(ctx, 3002)
		return
	}
	if status, item := srv.GetUserByTrueName(arg.TrueName, claims.AppId); (!status) && (item.ID != arg.Id) {
		logrus.Error("the user true name is exist:", arg.TrueName)
		ds.ResponseError(ctx, 3002)
		return
	}
	err, _ := srv.UpdateUser(
		arg.Id,
		arg.Account,
		arg.NickName,
		arg.TrueName,
		arg.Mobile,
		arg.Mail,
		arg.GroupId,
	)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	ds.ResponseSuccess(ctx, nil)
	return
}

// @Summary 用户密码重置
// @Tags user
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @param IdStruct body IdStruct true "请求体"
// @Success 200 {object} ApiResponse{result=object}
// @Router /user/password/reset [post]
func (ds *defaultServer) PasswordReset(ctx *gin.Context) {
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
	err, _ := srv.ResetUserPwd(arg.Id)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	ds.ResponseSuccess(ctx, nil)
	return
}

// @Summary 删除用户
// @Tags user
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @param IdStruct body IdStruct true "请求体"
// @Success 200 {object} ApiResponse{result=object}
// @Router /user/delete [post]
func (ds *defaultServer) DeleteUser(ctx *gin.Context) {
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
	if err := srv.DeleteUser(arg.Id); err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	ds.ResponseSuccess(ctx, nil)
	return
}

// @Summary 用户详情
// @Tags user
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @param id query string true "用户id"
// @Success 200 {object} ApiResponse{result=user.User}
// @Router /user/detail [get]
func (ds *defaultServer) UserDetail(ctx *gin.Context) {
	idStr := ctx.Query("id")
	if idStr == "" {
		ds.InvalidParametersError(ctx)
		return
	}

	srv := user.NewUserService(ds.db)
	status, item := srv.GetUser(utils.NewStr(idStr).Uint64())
	if status {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}
	item.Password = ""
	ds.ResponseSuccess(ctx, item)
	return
}
