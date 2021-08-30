//Created by Goland
//@User: lenora
//@Date: 2021/3/9
//@Time: 7:26 下午
package user

import "github.com/Lenora-Z/low-code/utils"

func (srv *userService) CreateRole(appId uint64, name string, permission map[uint64][]uint64) (error, *Role) {
	err, item := createRole(srv.db, &Role{
		AppID:    appId,
		Name:     name,
		IsDelete: &FALSE,
	})
	if err != nil {
		return err, nil
	}
	for i, x := range permission {
		for _, v := range x {
			if err, _ := createRolePermission(srv.db, RolePermission{
				RoleID:   item.ID,
				ParentID: i,
				FormID:   v,
			}); err != nil {
				return err, nil
			}
		}
	}
	return nil, item
}

func (srv *userService) UpdateRole(id uint64, name string, permission map[uint64][]uint64) (error, *Role) {
	err, item := updateRole(srv.db, &Role{
		ID:   id,
		Name: name,
	})
	if err != nil {
		return err, nil
	}

	if err := deleteRolePermission(srv.db, id); err != nil {
		return err, nil
	}

	for i, x := range permission {
		for _, v := range x {
			if err, _ := createRolePermission(srv.db, RolePermission{
				RoleID:   id,
				ParentID: i,
				FormID:   v,
			}); err != nil {
				return err, nil
			}
		}
	}
	return nil, item

}

func (srv *userService) RoleList(page, limit uint32, appId uint64, name string) (error, uint32, RoleList) {
	offset := (page - 1) * limit
	return getRoleList(srv.db, offset, limit, appId, name)
}

func (srv *userService) RoleDetail(id uint64) (bool, *Role) {
	return roleDetail(srv.db, id)
}

func (srv *userService) RoleDetailByName(name string, appId ...uint64) (bool, *Role) {
	column := make(map[string]interface{})
	column["name"] = name
	column["is_delete"] = 0
	if len(appId) > 0 {
		column["app_id"] = appId[0]
	}
	return roleDetailByColumn(srv.db, column)
}

func (srv *userService) FullRoleList(appId uint64) (error, uint32, RoleList) {
	return getRoleList(srv.db, 0, utils.MAX_LIMIT, appId, "")
}

func (srv *userService) GetRoleByIdGroup(group []uint64) (error, RoleList) {
	return getRoleListByIdGroup(srv.db, group)
}

func (srv *userService) DeleteRole(id uint64) error {
	err, _ := updateRole(srv.db, &Role{
		ID:       id,
		IsDelete: &TRUE,
	})
	if err != nil {
		return err
	}
	return deleteOrgRole(srv.db, 0, id)
}

func (srv *userService) PermissionOfRoles(pId uint64, roles []uint64) (error, PermissionList) {
	return rolePermissionList(srv.db, pId, roles)
}

func (srv *userService) GetMultiPermission(ids []uint64) (error, uint32, []Permission) {
	return getPermissionList(srv.db, 0, utils.MAX_LIMIT, ids)
}

func (srv *userService) SearchPermission(roleId []uint64, pId, formId uint64) (bool, *RolePermission) {
	return rolePermissionDetailByColumn(srv.db, roleId, pId, formId)
}
