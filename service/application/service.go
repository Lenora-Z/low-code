//Created by Goland
//@User: lenora
//@Date: 2021/3/8
//@Time: 1:57 下午
package application

import (
	"github.com/Lenora-Z/low-code/utils"
	"github.com/jinzhu/gorm"
)

type ApplicationService interface {
	//应用列表
	ApplicationList(page, limit uint32) (error, uint32, []Application)
	//应用列表
	ApplicationFullList() (error, uint32, []Application)
	//新增应用
	CreateApplication(name, icon, desc string) (error, *Application)
	//应用详情
	GetApplication(id uint64) (bool, *Application)
	// GetApplicationByHash 获取应用by hash
	GetApplicationByHash(hash string) (bool, *Application)
	//获取应用by名称
	GetApplicationByName(name string) (bool, *Application)
	//更新应用
	UpdateApplication(id uint64, name, icon, desc string) (error, *Application)
	//更新版本号
	UpdateApplicationVersion(id, version uint64, status int8) (error, *Application)
}

type applicationService struct {
	db *gorm.DB
}

func NewApplicationService(db *gorm.DB) ApplicationService {
	u := new(applicationService)
	u.db = db
	return u
}

func (srv *applicationService) CreateApplication(name, icon, desc string) (error, *Application) {
	hash, err := utils.GetUUid()
	if err != nil {
		return err, nil
	}
	return createApp(srv.db, &Application{
		Name:    name,
		Icon:    icon,
		Desc:    desc,
		AppHash: hash,
		Status:  PENDING,
	})
}

func (srv *applicationService) ApplicationList(page, limit uint32) (error, uint32, []Application) {
	offset := (page - 1) * limit
	return getAppList(srv.db, offset, limit)
}

func (srv *applicationService) ApplicationFullList() (error, uint32, []Application) {
	return getAppList(srv.db, 0, utils.MAX_LIMIT)
}

func (srv *applicationService) UpdateApplication(id uint64, name, icon, desc string) (error, *Application) {
	return updateApp(srv.db, &Application{
		ID:   id,
		Name: name,
		Icon: icon,
		Desc: desc,
	})
}

func (srv *applicationService) UpdateApplicationVersion(id, version uint64, status int8) (error, *Application) {
	return updateApp(srv.db, &Application{
		ID:        id,
		VersionID: version,
		Status:    status,
	})
}

func (srv *applicationService) GetApplication(id uint64) (bool, *Application) {
	return appDetail(srv.db, id)
}

func (srv *applicationService) GetApplicationByName(name string) (bool, *Application) {
	column := make(map[string]interface{})
	column["name"] = name
	return appDetailByColumn(srv.db, column)
}

func (srv *applicationService) GetApplicationByHash(hash string) (bool, *Application) {
	column := make(map[string]interface{})
	column["app_hash"] = hash
	return appDetailByColumn(srv.db, column)
}
