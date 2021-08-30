//Created by Goland
//@User: lenora
//@Date: 2021/3/9
//@Time: 7:27 下午
package user

import "github.com/Lenora-Z/low-code/utils"

func (srv *userService) CreateOrgRole(appId, orgId, roleId uint64) (error, *OrgRole) {
	return createOrgRole(srv.db, &OrgRole{
		AppID:   appId,
		GroupID: orgId,
		RoleID:  roleId,
	})
}

func (srv *userService) DeleteOrgRoleByOrgId(groupId uint64) error {
	return deleteOrgRole(srv.db, groupId, 0)
}

func (srv *userService) GetRelationByRoleId(id []uint64) (error, uint32, OrgRoleList) {
	return getOrgRoleList(srv.db, 0, utils.MAX_LIMIT, 0, id)
}

func (srv *userService) GetRelationByOrgId(id uint64) (error, uint32, OrgRoleList) {
	return getOrgRoleList(srv.db, 0, utils.MAX_LIMIT, id, nil)
}
