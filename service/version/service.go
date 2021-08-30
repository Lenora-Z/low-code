//Created by Goland
//@User: lenora
//@Date: 2021/3/11
//@Time: 3:35 下午
package version

import "github.com/jinzhu/gorm"

type VersionService interface {
	//发布版本
	CreateVersion(appId uint64, domain, version, note string) (error, *Version)
	//上线版本信息
	OnlineVersion(appId uint64) (bool, Version)
	//历史版本
	VersionHistory(page, limit uint32, appId uint64) (error, uint32, []Version)
	//更新版本状态
	UpdateVersionStatus(id uint64, status int8) (error, Version)
}

type versionService struct {
	db *gorm.DB
}

func NewVersionService(db *gorm.DB) VersionService {
	u := new(versionService)
	u.db = db
	return u
}

func (srv *versionService) CreateVersion(appId uint64, domain, version, note string) (error, *Version) {
	//覆盖已上线版本
	where := make(map[string]interface{})
	where["status"] = ONLINE
	where["app_id"] = appId
	update := make(map[string]interface{})
	update["status"] = COVERED
	if err := batchUpdateVersionByColumn(srv.db, update, where); err != nil {
		return err, nil
	}
	//发布新版本
	return createVersion(srv.db, Version{
		AppID:   appId,
		Domain:  domain,
		Version: version,
		Status:  ONLINE,
		Note:    note,
	})
}

func (srv *versionService) OnlineVersion(appId uint64) (bool, Version) {
	column := make(map[string]interface{})
	column["status"] = ONLINE
	column["app_id"] = appId
	return versionDetailByColumn(srv.db, column)
}

func (srv *versionService) VersionHistory(page, limit uint32, appId uint64) (error, uint32, []Version) {
	offset := (page - 1) * limit
	return getVersionList(srv.db, offset, limit, appId, OFFLINE, COVERED)
}

func (srv *versionService) UpdateVersionStatus(id uint64, status int8) (error, Version) {
	return updateVersion(srv.db, Version{
		ID:     id,
		Status: status,
	})
}
