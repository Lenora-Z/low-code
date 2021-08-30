//Created by Goland
//@User: lenora
//@Date: 2021/3/9
//@Time: 7:26 下午
package user

import (
	"github.com/Lenora-Z/low-code/utils"
)

func (srv *userService) CreateOrganization(appId, parentId uint64, name string) (error, *Organization) {
	return createOrganization(srv.db, &Organization{
		AppID:    appId,
		Name:     name,
		ParentID: parentId,
		IsDelete: &FALSE,
	})
}

func (srv *userService) OrganizationByParent(appId uint64, parentId ...uint64) (error, uint32, OrganizationList) {
	return getOrganizationList(srv.db, 0, utils.MAX_LIMIT, appId, parentId...)
}

func (srv *userService) UpdateOrganization(id uint64, name string) (error, *Organization) {
	return updateOrganization(srv.db, &Organization{
		ID:   id,
		Name: name,
	})
}

func (srv *userService) DeleteOrganization(id uint64) error {
	if err, _ := updateOrganization(srv.db, &Organization{
		ID:       id,
		IsDelete: &TRUE,
	}); err != nil {
		return err
	}
	return deleteOrgRole(srv.db, id, 0)
}

func (srv *userService) GetOrganizationByIdGroup(group []uint64) (error, OrganizationList) {
	return getOrganizationListByIdGroup(srv.db, group)
}

func (srv *userService) GetOrganizationRecursively(id, appId uint64, group []uint64) []uint64 {
	group = append(group, id)
	err, count, list := getOrganizationList(srv.db, 0, utils.MAX_LIMIT, appId, id)
	if err != nil {
		return nil
	}
	if count > 0 {
		for _, x := range list {
			group = srv.GetOrganizationRecursively(x.ID, appId, group)
		}
	}
	return group
}

func (srv *userService) OrganizationDetail(id uint64) (bool, *Organization) {
	return organizationDetail(srv.db, id)
}
