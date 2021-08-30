//Created by Goland
//@User: lenora
//@Date: 2021/3/12
//@Time: 4:08 下午
package flow

import (
	"fmt"
	"github.com/Lenora-Z/low-code/utils"
	"github.com/jinzhu/gorm"
	"time"
)

type FlowService interface {
	//新增流程
	CreateFlow(appId, userId uint64, name, key string) (error, *Flow)
	//新增映射
	CreateFlowMapping(formId, flowId, appId uint64, rand string) (error, *FlowMapping)
	// DelFlow 流程删除
	DelFlow(id uint64) error
	// DelMappingByFormId 删除某表单下的全部映射关系
	DelMappingByFormId(id uint64) error
	//已发布的全部流程
	GetAllFlowList(appId uint64) (error, uint32, FlowList)
	//流程详情
	GetFlowDetail(id uint64) (bool, *Flow)
	//获取流程by  key
	GetFlowDetailByKey(key string) (bool, *Flow)
	//获取流程by名称
	GetFlowDetailByName(name string, appId ...uint64) (bool, *Flow)
	//流程列表
	GetFlowList(page, limit uint32, appId uint64, name string) (error, uint32, FlowList)
	//获取某用户可处理的全部流程
	GetFlowListByAssignee(userId, appId uint64) (error, FlowList, []string)
	//获取表单的对应映射
	GetMappingByFormId(appId uint64, formId []uint64) (error, uint32, FlowMappingList)
	//获取生效映射
	GetOnlineMapping(appId uint64) (error, uint32, FlowMappingList)
	//流程下线
	OfflineFlow(appId uint64) error
	//映射失效
	OfflineMapping(appId uint64) error
	//流程上线
	OnlineFlow(appId uint64) error
	//映射生效
	OnlineMapping(appId, versionId uint64, rand string) error
	//更新流程
	UpdateFlow(id uint64, name, jsonContent string, status bool, chainUp ChainUpParams, email EmailParams) (error, *Flow)
	//更新流程key
	UpdateFlowKey(id uint64, key string) (error, *Flow)
	//获取抄送某用户的全部流程
	GetFlowListByNotifier(userId, appId uint64) (error, FlowList, []string)
}

type flowService struct {
	db *gorm.DB
}

func NewFlowService(db *gorm.DB) FlowService {
	u := new(flowService)
	u.db = db
	return u
}

func (srv *flowService) CreateFlow(appId, userId uint64, name, key string) (error, *Flow) {
	number := fmt.Sprintf("AL%s%d", time.Now().Format(utils.DateAttrFormatStr), utils.RandInt(1000, 9999))
	err, item := createFlow(srv.db, Flow{
		Number:   number,
		Name:     name,
		AppID:    appId,
		UserID:   userId,
		Key:      key,
		Status:   &FALSE,
		IsOnline: &FALSE,
		IsDelete: &FALSE,
	})
	return err, item
}

func (srv *flowService) CreateFlowMapping(formId, flowId, appId uint64, rand string) (error, *FlowMapping) {
	return createFlowMapping(srv.db, FlowMapping{
		FormID: formId,
		FlowID: flowId,
		AppID:  appId,
		Random: rand,
		Status: PENDING,
	})
}

func (srv *flowService) GetMappingByFormId(appId uint64, formId []uint64) (error, uint32, FlowMappingList) {
	return getFMList(srv.db, 0, utils.MAX_LIMIT, appId, formId, nil, "")
}

func (srv *flowService) GetMappingByFlowId(appId, flowId uint64) (error, uint32, FlowMappingList) {
	return getFMList(srv.db, 0, utils.MAX_LIMIT, appId, nil, []uint64{flowId}, "")
}

func (srv *flowService) GetOnlineMapping(appId uint64) (error, uint32, FlowMappingList) {
	return getFMList(srv.db, 0, utils.MAX_LIMIT, appId, nil, nil, "")
}

func (srv *flowService) GetAllFlowList(appId uint64) (error, uint32, FlowList) {
	return getFlowList(srv.db, 0, utils.MAX_LIMIT, appId, "", true)
}

func (srv *flowService) GetFlowDetail(id uint64) (bool, *Flow) {
	return flowDetail(srv.db, id)
}

func (srv *flowService) GetFlowDetailByKey(key string) (bool, *Flow) {
	column := make(map[string]interface{})
	column["key"] = key
	return flowDetailByColumn(srv.db, column)
}

func (srv *flowService) GetFlowDetailByName(name string, appId ...uint64) (bool, *Flow) {
	column := make(map[string]interface{})
	column["name"] = name
	if len(appId) > 0 {
		column["app_id"] = appId
	}
	return flowDetailByColumn(srv.db, column)
}

func (srv *flowService) GetFlowList(page, limit uint32, appId uint64, name string) (error, uint32, FlowList) {
	offset := (page - 1) * limit
	return getFlowList(srv.db, offset, limit, appId, name, false)
}

func (srv *flowService) GetFlowListByAssignee(userId, appId uint64) (error, FlowList, []string) {
	err, _, list := getAssigneeList(srv.db, 0, userId)
	if err != nil {
		return err, nil, nil
	}
	ids := list.Flows()
	err, _, flows := getFlowList(srv.db, 0, utils.MAX_LIMIT, appId, "", true, ids...)
	return err, flows, list.Acts()
}

func (srv *flowService) GetFlowListByNotifier(userId, appId uint64) (error, FlowList, []string) {
	err, _, list := getNotifierList(srv.db, 0, userId)
	if err != nil {
		return err, nil, nil
	}
	ids := list.Flows()
	err, _, flows := getFlowList(srv.db, 0, utils.MAX_LIMIT, appId, "", true, ids...)
	return err, flows, list.Acts()
}

func (srv *flowService) UpdateFlow(id uint64, name, jsonContent string, status bool, chainUp ChainUpParams, email EmailParams) (error, *Flow) {
	if err := srv.delFlowDetail(id); err != nil {
		return err, nil
	}
	err, item := updateFlow(srv.db, Flow{
		ID:     id,
		Name:   name,
		JSON:   jsonContent,
		Status: &status,
	})
	if err != nil {
		return err, nil
	}

	if status {
		err = srv.flowDecode(item.ID, item.AppID, []byte(jsonContent), chainUp, email)
	}

	return err, item
}

func (srv *flowService) UpdateFlowKey(id uint64, key string) (error, *Flow) {
	return updateFlow(srv.db, Flow{
		ID:  id,
		Key: key,
	})
}

func (srv *flowService) OnlineFlow(appId uint64) error {
	where := make(map[string]interface{})
	change := make(map[string]interface{})
	where["app_id"] = appId
	where["status"] = &TRUE
	change["is_online"] = &TRUE
	return batchUpdateFlow(srv.db, where, change)
}

func (srv *flowService) OfflineFlow(appId uint64) error {
	where := make(map[string]interface{})
	change := make(map[string]interface{})
	where["app_id"] = appId
	where["status"] = &TRUE
	change["is_online"] = &FALSE
	return batchUpdateFlow(srv.db, where, change)
}

func (srv *flowService) OnlineMapping(appId, versionId uint64, rand string) error {
	//旧版本失效
	if err := batchUpdateMapping(
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
	return batchUpdateMapping(
		srv.db,
		map[string]interface{}{
			"app_id": appId,
			"status": PENDING,
			"random": rand,
		},
		map[string]interface{}{
			"status":     VALID,
			"version_id": versionId,
		},
	)
}

func (srv *flowService) OfflineMapping(appId uint64) error {
	return batchUpdateMapping(
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

func (srv *flowService) DelMappingByFormId(id uint64) error {
	return batchDeleteMapping(srv.db, map[string]interface{}{"form_id": id})
}

func (srv *flowService) DelFlow(id uint64) error {
	err, _ := updateFlow(srv.db, Flow{
		ID:       id,
		IsDelete: &TRUE,
	})
	return err
}

func (srv *flowService) delFlowDetail(id uint64) error {
	if err := delActivity(srv.db, id); err != nil {
		return err
	}
	if err := delActivityData(srv.db, id); err != nil {
		return err
	}
	if err := delActivityGateway(srv.db, id); err != nil {
		return err
	}
	if err := delActivityService(srv.db, id); err != nil {
		return err
	}
	if err := delActivityTable(srv.db, id); err != nil {
		return err
	}
	if err := delActivityTime(srv.db, id); err != nil {
		return err
	}
	return nil
}

func (srv *flowService) createAssignee(flowId uint64, assignee []ActAssignee) error {
	for _, x := range assignee {
		for _, v := range x.Assignee {
			if err, _ := createAssignee(srv.db, FlowAssignee{
				UserID:   v,
				Activity: x.ActId,
				FlowID:   flowId,
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

func (srv *flowService) createNotifier(flowId uint64, notifier []ActNotifier) error {
	for _, x := range notifier {
		for _, v := range x.Notifier {
			if err, _ := createNotifier(srv.db, FlowNotifier{
				UserID:   v,
				Activity: x.ActId,
				FlowID:   flowId,
			}); err != nil {
				return err
			}
		}
	}
	return nil
}
