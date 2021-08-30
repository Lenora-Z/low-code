//Created by Goland
//@User: lenora
//@Date: 2021/3/9
//@Time: 2:14 下午
package user

import (
	"github.com/jinzhu/gorm"
)

type UserService interface {
	//新建组织
	CreateOrganization(appId, parentId uint64, name string) (error, *Organization)
	//新建组织角色关系
	CreateOrgRole(appId, orgId, roleId uint64) (error, *OrgRole)
	//新建角色
	CreateRole(appId uint64, name string, permission map[uint64][]uint64) (error, *Role)
	//新建用户
	CreateUser(account, nickName, trueName, mobile, mail string, appId, groupId uint64) (error, *User)
	//删除组织
	DeleteOrganization(id uint64) error
	//删除组织下的全部角色关系
	DeleteOrgRoleByOrgId(groupId uint64) error
	//删除角色
	DeleteRole(id uint64) error
	//删除用户
	DeleteUser(id uint64) error
	//全部角色
	FullRoleList(appId uint64) (error, uint32, RoleList)
	//模糊检索应用内的全部用户
	GetAllUser(appId uint64, name string, group ...uint64) (error, uint32, UserList)
	//批量获取权限
	GetMultiPermission(ids []uint64) (error, uint32, []Permission)
	//批量获取用户by用户id
	GetMultiUser(ids []uint64) (error, UserList)
	//获取组织列表by组织id集合
	GetOrganizationByIdGroup(group []uint64) (error, OrganizationList)
	//获取组织角色对应关系by组织id
	GetRelationByOrgId(id uint64) (error, uint32, OrgRoleList)
	//获取组织角色对应关系by角色id
	GetRelationByRoleId(id []uint64) (error, uint32, OrgRoleList)
	//获取角色列表by角色id集合
	GetRoleByIdGroup(group []uint64) (error, RoleList)
	//获取用户by账号id
	GetUser(id uint64) (bool, *User)
	//获取用户by用户账号
	GetUserByAccount(name string, appId uint64) (bool, *User)
	//获取用户by真实姓名
	GetUserByTrueName(name string, appId uint64) (bool, *User)
	//用户列表
	GetUserList(page, limit uint32, appId uint64, name string, group []uint64) (error, uint32, UserList)
	//获取用户列表by组织id
	GetUserListByGroupId(groupId, appId uint64) (error, uint32, UserList)
	//组织列表
	OrganizationByParent(appId uint64, parentId ...uint64) (error, uint32, OrganizationList)
	//组织详情
	OrganizationDetail(id uint64) (bool, *Organization)
	//组织对应的全部角色
	PermissionOfRoles(pId uint64, roles []uint64) (error, PermissionList)
	//重置密码
	ResetUserPwd(id uint64, pwd ...string) (error, *User)
	//修改密码状态
	UpdatePwdStatus(id uint64, status *bool) (error, *User)
	//角色详情
	RoleDetail(id uint64) (bool, *Role)
	//角色详情by角色名
	RoleDetailByName(name string, appId ...uint64) (bool, *Role)
	//角色列表
	RoleList(page, limit uint32, appId uint64, name string) (error, uint32, RoleList)
	//检索用户权限
	SearchPermission(roleId []uint64, pId, formId uint64) (bool, *RolePermission)
	//更新组织
	UpdateOrganization(id uint64, name string) (error, *Organization)
	//更新角色
	UpdateRole(id uint64, name string, permission map[uint64][]uint64) (error, *Role)
	//更新用户信息
	UpdateUser(id uint64, account, nickName, trueName, mobile, mail string, groupId uint64) (error, *User)
	//递归获取组织id
	GetOrganizationRecursively(id, appId uint64, group []uint64) []uint64
}

type userService struct {
	db *gorm.DB
}

var TRUE, FALSE = true, false

func NewUserService(db *gorm.DB) UserService {
	u := new(userService)
	u.db = db
	return u
}
