package navigation

import (
	"fmt"
	"github.com/Lenora-Z/low-code/utils"
	"github.com/jinzhu/gorm"
	"time"
)

type NavigationService interface {
	//新增导航栏
	CreateNavigation(appId uint64, name, desc string, content []byte, status bool) (error, *Navigation)
	//获取导航栏详情
	GetNavDetail(id uint64) (bool, *Navigation)
	//获取导航栏by名称
	GetNavDetailByName(name string, appId ...uint64) (bool, *Navigation)
	//获取生效表单
	GetValidNavigation(appId uint64) (bool, *Navigation)
	//获取导航栏列表
	NavigationList(page, limit uint32, appId uint64) (error, uint32, []Navigation)
	//更新导航栏
	UpdateNavigation(id uint64, name, desc string, content []byte, status bool) (error, *Navigation)
	//导航栏上下线
	UpdateAppNavigation(appId uint64, status bool) error
	//更新导航栏内容
	UpdateNavigationContent(id uint64, content []byte) (error, *Navigation)
}

type navigationService struct {
	db *gorm.DB
}

func NewNavigationService(db *gorm.DB) NavigationService {
	u := new(navigationService)
	u.db = db
	return u
}

func (srv *navigationService) CreateNavigation(appId uint64, name, desc string, content []byte, status bool) (error, *Navigation) {
	number := fmt.Sprintf("NV%s%d", time.Now().Format(utils.DateAttrFormatStr), utils.RandInt(1000, 9999))
	return createNavigation(srv.db, Navigation{
		Name:     name,
		AppID:    appId,
		Number:   number,
		Desc:     desc,
		Content:  content,
		Status:   &status,
		IsOnline: &FALSE,
	})
}

func (srv *navigationService) GetNavDetail(id uint64) (bool, *Navigation) {
	return navigationDetail(srv.db, id)
}

func (srv *navigationService) GetNavDetailByName(name string, appId ...uint64) (bool, *Navigation) {
	column := make(map[string]interface{})
	column["name"] = name
	if len(appId) > 0 {
		column["app_id"] = appId
	}
	return navigationDetailByColumn(srv.db, column)
}

func (srv *navigationService) GetValidNavigation(appId uint64) (bool, *Navigation) {
	column := make(map[string]interface{})
	column["app_id"] = appId
	column["status"] = &TRUE
	//column["is_online"] = &TRUE
	return navigationDetailByColumn(srv.db, column)
}

func (srv *navigationService) NavigationList(page, limit uint32, appId uint64) (error, uint32, []Navigation) {
	offset := (page - 1) * limit
	return getNavigationList(srv.db, offset, limit, appId)
}

func (srv *navigationService) UpdateNavigation(id uint64, name, desc string, content []byte, status bool) (error, *Navigation) {
	return updateNavigation(srv.db, Navigation{
		ID:      id,
		Name:    name,
		Desc:    desc,
		Content: content,
		Status:  &status,
	})
}

func (srv *navigationService) UpdateNavigationContent(id uint64, content []byte) (error, *Navigation) {
	return updateNavigation(srv.db, Navigation{
		ID:      id,
		Content: content,
	})
}

func (srv *navigationService) UpdateAppNavigation(appId uint64, status bool) error {
	if err := batchUpdateNavigation(
		srv.db,
		map[string]interface{}{
			"app_id": appId,
			"status": &TRUE,
		},
		map[string]interface{}{
			"is_online": &FALSE,
		},
	); err != nil {
		return err
	}
	if status {
		return batchUpdateNavigation(
			srv.db,
			map[string]interface{}{
				"app_id": appId,
				"status": &TRUE,
			},
			map[string]interface{}{
				"is_online": &TRUE,
			},
		)
	}
	return nil
}
