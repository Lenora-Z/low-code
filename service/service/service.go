//Created by Goland
//@User: lenora
//@Date: 2021/3/11
//@Time: 2:58 下午
package service

import (
	"github.com/Lenora-Z/low-code/utils"
	"github.com/jinzhu/gorm"
)

type OutsideService interface {
	//新增参数依赖
	CreateParamRely(appId, flowId, serviceId, paramId, formId, FieldId uint64, rand string) (error, *ParamRely)
	//参数详情
	GetParamDetail(id uint64) (bool, *Param)
	//获取参数依赖
	GetParamRely(srvId, flowId, formId []uint64) (error, uint32, []ParamRely)
	//参数列表
	GetParamList(srvId uint64, mode uint8) (error, []Param)
	// GetParamListBySrvGroup 获取多个服务的全部参数
	GetParamListBySrvGroup(srvId []uint64) (error, []Param)
	// GetServiceByType 获取外部服务by类型
	GetServiceByType(t ...int8) (error, ServiceList)
	// GetAllServiceList 服务列表
	GetAllServiceList() (error, []Service)
	//参数依赖生效
	OnlineRelies(appId uint64, rand string) error
	//参数依赖失效
	OfflineRelies(appId uint64) error
}

type outsideService struct {
	db *gorm.DB
}

func NewOutsideService(db *gorm.DB) OutsideService {
	u := new(outsideService)
	u.db = db
	return u
}

func (srv *outsideService) CreateParamRely(appId, flowId, serviceId, paramId, formId, FieldId uint64, rand string) (error, *ParamRely) {
	return createParamRely(srv.db, ParamRely{
		AppID:     appId,
		FlowID:    flowId,
		ServiceID: serviceId,
		ParamID:   paramId,
		FormID:    formId,
		FieldID:   FieldId,
		Random:    rand,
		Status:    PENDING,
	})
}

func (srv *outsideService) GetParamRely(srvId, flowId, formId []uint64) (error, uint32, []ParamRely) {
	return getParamRelyList(srv.db, 0, utils.MAX_LIMIT, srvId, formId, flowId)
}

func (srv *outsideService) GetParamList(srvId uint64, mode uint8) (error, []Param) {
	return getParamList(srv.db, []uint64{srvId}, mode)
}

func (srv *outsideService) GetParamListBySrvGroup(srvId []uint64) (error, []Param) {
	return getParamList(srv.db, srvId, IN)
}

func (srv *outsideService) GetParamDetail(id uint64) (bool, *Param) {
	return paramDetail(srv.db, id)
}

func (srv *outsideService) GetAllServiceList() (error, []Service) {
	return getServiceList(srv.db, nil)
}

func (srv *outsideService) GetServiceByType(t ...int8) (error, ServiceList) {
	return getServiceList(srv.db, t)
}

func (srv *outsideService) OnlineRelies(appId uint64, rand string) error {
	//旧版本失效
	if err := batchUpdateRelies(
		srv.db,
		map[string]interface{}{
			"app_id": appId,
			"status": VALID,
		},
		map[string]interface{}{
			"status": INVALID,
		},
	); err != nil {
		return err
	}
	//新版本生效
	return batchUpdateRelies(
		srv.db,
		map[string]interface{}{
			"app_id": appId,
			"status": PENDING,
			"random": rand,
		},
		map[string]interface{}{
			"status": VALID,
		},
	)
}

func (srv *outsideService) OfflineRelies(appId uint64) error {
	return batchUpdateRelies(
		srv.db,
		map[string]interface{}{
			"app_id": appId,
			"status": VALID,
		},
		map[string]interface{}{
			"status": INVALID,
		},
	)
}
