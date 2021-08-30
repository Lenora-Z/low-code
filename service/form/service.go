// Package form Created by Goland
//@User: lenora
//@Date: 2021/3/12
//@Time: 10:21 上午
package form

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Lenora-Z/low-code/utils"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"time"
)

type FormService interface {
	// ButtonByFlowId 按钮列表by流程id
	ButtonByFlowId(flowId uint64) (error, FieldButtons)
	// CopyForm 表单复制
	CopyForm(id uint64, content []FieldFormat, foot FieldFormat) (error, *Form)
	// CreateForm 新建表单
	CreateForm(appId, sourceId uint64, t int8, pageType uint8, name, desc string, content []FieldFormat, foot []byte, status bool) (error, *Form, []uint64)
	// DelForm 删除表单
	DelForm(id uint64) error
	// FieldButtonByFlowId 按钮控件列表by流程id
	FieldButtonByFlowId(flowId uint64) (error, FieldButtons)
	// FormList 表单列表
	FormList(page, limit uint32, appId uint64, t int8, name string) (error, uint32, FormList)
	// FormListByPageType 表单列表by页面展示类型
	FormListByPageType(appId uint64, pt int8, status bool) (error, uint32, FormList)
	//已发布的全部表单
	FullFormList(appId uint64, t int8) (error, uint32, FormList)
	//字段列表
	FieldList(formId uint64, fields ...string) (error, []Field)
	//字段列表by多个表单id
	FieldListByMultiFormIds(ids []uint64) (error, []Field)
	// FieldListByMultiIds 字段列表by多个id
	FieldListByMultiIds(ids []uint64) (error, []Field)
	// GetAllRelationForm 递归获取全部关联表单id
	GetAllRelationForm(id, appId uint64, group FormTreeList) FormTreeList
	//字段详情
	GetFieldDetail(id uint64) (bool, *Field)
	//字段详情by key
	GetFieldDetailByKey(key string) (bool, *Field)
	//表单详情
	GetFormDetail(id uint64) (bool, *Form)
	//表单详情by名称
	GetFormDetailByName(name string, appId ...uint64) (bool, *Form)
	//获取级联控件配置
	GetLinkageField(id uint64) (bool, []SingleLinkage)
	// GetRelationForm 获取关联表
	GetRelationForm(id uint64) (error, MultiFormList)
	// GetFormDetailByNo 表单详情by编号
	GetFormDetailByNo(number string, appId ...uint64) (bool, *Form)
	// MultiFormList 批量获取表单
	MultiFormList(appId uint64, ids []uint64) (error, uint32, FormList)
	//表单下线
	OfflineForm(appId uint64) error
	//表单上线
	OnlineForm(appId uint64) error
	// TableButtonByFlowId 列表控件按钮列表by流程id
	TableButtonByFlowId(flowId uint64) (error, FieldButtons)
	// TableColumns 列表控件字段列表by列表控件id
	TableColumns(id uint64) (error, []FieldTableColumn)
	// UpdateForm 编辑表单
	UpdateForm(id uint64, pageType uint8, name, desc string, content []FieldFormat, foot []byte, status bool) (error, *Form, []uint64)
}

type formService struct {
	db *gorm.DB
}

func NewFormService(db *gorm.DB) FormService {
	u := new(formService)
	u.db = db
	return u
}

func (srv *formService) CopyForm(id uint64, content []FieldFormat, foot FieldFormat) (error, *Form) {
	notFound, origin := formDetail(srv.db, id)
	if notFound || *origin.IsDelete == true {
		return errors.New("form not found"), nil
	}
	contentByte, _ := json.Marshal(content)
	footerByte, _ := json.Marshal(foot)
	number := fmt.Sprintf("AF%s%d", time.Now().Format(utils.DateAttrFormatStr), utils.RandInt(1000, 9999))
	err, item := createForm(srv.db, Form{
		Name:     origin.Name + "-副本",
		AppID:    origin.AppID,
		From:     id,
		Type:     origin.Type,
		PageType: origin.PageType,
		Number:   number,
		Desc:     origin.Desc,
		Content:  contentByte,
		Footer:   footerByte,
		Status:   origin.Status,
		IsOnline: &FALSE,
	})
	if err != nil {
		return err, nil
	}
	btnStatus := &FALSE
	if *item.Status == true {
		for _, x := range content {
			err, _ = srv.createField(x.Layers, item.ID, item.PageType, btnStatus)
			if err != nil {
				return err, nil
			}
		}
	}
	return err, item
}

func (srv *formService) CreateForm(appId, sourceId uint64, t int8, pageType uint8, name, desc string, content []FieldFormat, foot []byte, status bool) (error, *Form, []uint64) {
	number := fmt.Sprintf("AF%s%d", time.Now().Format(utils.DateAttrFormatStr), utils.RandInt(1000, 9999))
	contentByte, _ := json.Marshal(content)

	//新增表单
	err, item := createForm(srv.db, Form{
		Name:              name,
		AppID:             appId,
		DatasourceTableID: sourceId,
		Type:              t,
		PageType:          pageType,
		Number:            number,
		Desc:              desc,
		Content:           contentByte,
		Footer:            foot,
		Status:            &status,
		IsOnline:          &FALSE,
		IsDelete:          &FALSE,
	})
	if err != nil {
		return err, nil, nil
	}
	flowIds := make([]uint64, 0, cap(content))
	bs := false
	btnStatus := &bs
	if status {
		for _, x := range content {
			err, flows := srv.createField(x.Layers, item.ID, item.PageType, btnStatus)
			if err != nil {
				return err, nil, nil
			}
			if len(flows) > 0 {
				flowIds = append(flowIds, flows...)
			}
		}

		if pageType == FORM_PAGE && *btnStatus == false {
			logrus.Error("button is required")
			return errors.New("param error"), nil, nil
		}
	}
	return err, item, flowIds
}

func (srv *formService) FieldList(formId uint64, fields ...string) (error, []Field) {
	return getFieldList(srv.db, []uint64{formId}, nil, fields...)
}

func (srv *formService) FieldListByMultiFormIds(ids []uint64) (error, []Field) {
	return getFieldList(srv.db, ids, nil)
}

func (srv *formService) FieldListByMultiIds(ids []uint64) (error, []Field) {
	return getFieldList(srv.db, nil, ids)
}

func (srv *formService) FormList(page, limit uint32, appId uint64, t int8, name string) (error, uint32, FormList) {
	offset := (page - 1) * limit
	return getFormList(srv.db, offset, limit, appId, t, 0, name, false, nil)
}

func (srv *formService) FormListByPageType(appId uint64, pt int8, status bool) (error, uint32, FormList) {
	return getFormList(srv.db, 0, utils.MAX_LIMIT, appId, 0, pt, "", status, nil)
}

func (srv *formService) FullFormList(appId uint64, t int8) (error, uint32, FormList) {
	return getFormList(srv.db, 0, utils.MAX_LIMIT, appId, t, 0, "", true, nil)
}

func (srv *formService) GetFieldDetail(id uint64) (bool, *Field) {
	return fieldDetail(srv.db, id)
}

func (srv *formService) GetFieldDetailByKey(key string) (bool, *Field) {
	column := make(map[string]interface{})
	column["`key`"] = key
	return fieldDetailByColumn(srv.db, column)
}

func (srv *formService) GetFormDetail(id uint64) (bool, *Form) {
	return formDetail(srv.db, id)
}

func (srv *formService) GetFormDetailByName(name string, appId ...uint64) (bool, *Form) {
	column := make(map[string]interface{})
	column["name"] = name
	column["is_delete"] = 0
	if len(appId) > 0 {
		column["app_id"] = appId
	}
	return formDetailByColumn(srv.db, column)
}

func (srv *formService) GetLinkageField(id uint64) (bool, []SingleLinkage) {
	column := make(map[string]interface{})
	column["field_id"] = id

	notFound, item := linkageDetailByColumn(srv.db, column)
	if notFound {
		return notFound, nil
	}
	return false, linkageDecode(item.Content)
}

func (srv *formService) GetRelationForm(id uint64) (error, MultiFormList) {
	return getRelationList(srv.db, id)
}

func (srv *formService) GetAllRelationForm(id, appId uint64, group FormTreeList) FormTreeList {
	err, list := getRelationList(srv.db, id)
	if err != nil {
		return nil
	}
	if len(list) <= 0 {
		return FormTreeList{}
	}
	err, _, forms := srv.MultiFormList(appId, list.Children())
	if err != nil {
		return nil
	}
	for _, x := range forms {
		var searchId uint64
		if x.From == 0 {
			searchId = x.ID
		} else {
			searchId = x.From
		}
		child := srv.GetAllRelationForm(searchId, appId, group)
		group = append(group, FormTree{
			Id:       x.ID,
			Children: child,
		})
	}
	return group
}

func (srv *formService) GetFormDetailByNo(number string, appId ...uint64) (bool, *Form) {
	column := make(map[string]interface{})
	column["number"] = number
	column["is_delete"] = 0
	if len(appId) > 0 {
		column["app_id"] = appId
	}
	return formDetailByColumn(srv.db, column)
}

func (srv *formService) MultiFormList(appId uint64, ids []uint64) (error, uint32, FormList) {
	return getFormList(srv.db, 0, utils.MAX_LIMIT, appId, 0, 0, "", true, ids)
}

func (srv *formService) UpdateForm(id uint64, pageType uint8, name, desc string, content []FieldFormat, foot []byte, status bool) (error, *Form, []uint64) {
	contentByte, _ := json.Marshal(content)
	//表单更新
	err, item := updateForm(srv.db, Form{
		ID:      id,
		Name:    name,
		Desc:    desc,
		Content: contentByte,
		Footer:  foot,
		Status:  &status,
	})
	if err != nil {
		return err, nil, nil
	}

	//删除原有的控件
	if err := srv.deleteField(id); err != nil {
		return err, nil, nil
	}

	flowIds := make([]uint64, 0, cap(content))
	bs := false
	btnStatus := &bs
	if status {
		for _, x := range content {
			err, flows := srv.createField(x.Layers, item.ID, pageType, btnStatus)
			if err != nil {
				return err, nil, nil
			}
			if len(flows) > 0 {
				flowIds = append(flowIds, flows...)
			}
		}

		if pageType == FORM_PAGE && *btnStatus == false {
			logrus.Error("button is required")
			return errors.New("param error"), nil, nil
		}
	}
	return err, item, flowIds

}

func (srv *formService) OnlineForm(appId uint64) error {
	where := make(map[string]interface{})
	change := make(map[string]interface{})
	where["app_id"] = appId
	where["status"] = &TRUE
	change["is_online"] = &TRUE
	return batchUpdateForm(srv.db, where, change)
}

func (srv *formService) OfflineForm(appId uint64) error {
	where := make(map[string]interface{})
	change := make(map[string]interface{})
	where["app_id"] = appId
	where["status"] = &TRUE
	change["is_online"] = &FALSE
	return batchUpdateForm(srv.db, where, change)
}

func (srv *formService) DelForm(id uint64) error {
	if err := srv.deleteField(id); err != nil {
		return err
	}
	err, _ := updateForm(srv.db, Form{
		ID:       id,
		IsDelete: &TRUE,
	})

	return err
}

func (srv *formService) ButtonByFlowId(flowId uint64) (error, FieldButtons) {
	err, btn := srv.FieldButtonByFlowId(flowId)
	if err != nil {
		return err, nil
	}

	err, tbBtn := srv.TableButtonByFlowId(flowId)
	if err != nil {
		return err, nil
	}
	btn = append(btn, tbBtn...)
	return nil, btn
}

func (srv *formService) FieldButtonByFlowId(flowId uint64) (error, FieldButtons) {
	return getButtonListWithCondition(srv.db, map[string]interface{}{
		"flow_id": flowId,
		"event":   PROCESS,
	})
}

func (srv *formService) TableButtonByFlowId(flowId uint64) (error, FieldButtons) {
	return getTableButtonListWithCondition(srv.db, map[string]interface{}{
		"flow_id": flowId,
		"event":   PROCESS,
	})
}

func (srv *formService) TableColumns(id uint64) (error, []FieldTableColumn) {
	return getTableColumnList(srv.db, id)
}

func (srv *formService) createField(item []map[string]interface{}, formId uint64, pageType uint8, bool2 *bool) (error, []uint64) {
	flows := make([]uint64, 0, cap(item))
	for _, x := range item {
		filedContent, _ := json.Marshal(x)
		//控件类型
		t, ok := x["type"].(string)
		if !ok {
			return errors.New("param error"), nil
		}
		if (pageType == FORM_PAGE && FIELD_USE_TYPE[t] == D) || (pageType == SHOW_PAGE && FIELD_USE_TYPE[t] == C) {
			return errors.New("field type does not match page type"), nil
		}
		if t == BUTTON {
			*bool2 = TRUE
		}
		//控件关键字
		key, ok := x["name"].(string)
		if !ok {
			return errors.New("param error"), nil
		}
		if t == TABS {
			logrus.Info("layout field:", key, ">>>>")
			tabByte, err := json.Marshal(x["tabs"])
			if err != nil {
				logrus.Error("tabs marshal error")
				return errors.New("param error"), nil
			}
			tabs := make([]ChildTabs, 0)
			if err := json.Unmarshal(tabByte, &tabs); err != nil {
				logrus.Error("tabs format error")
				return errors.New("param error"), nil
			}
			for _, tab := range tabs {
				for _, layer := range tab.Children {
					//TODO
					err, ret := srv.createField(layer.Layers, formId, pageType, bool2)
					if err != nil {
						return err, nil
					}
					if len(ret) > 0 {
						flows = append(flows, ret...)
					}
				}
			}
			continue
		}

		var title string
		var necessary, repeat bool
		//控件标题
		title, ok = x["title"].(string)
		if !ok || title == "" || utils.Strlen(title) > 40 {
			return errors.New("param error"), nil
		}
		//是否必填
		necessary, ok = x["require"].(bool)
		if !ok {
			return errors.New("param error"), nil
		}
		//是否做唯一检验
		repeat, ok = x["repeat"].(bool)
		if !ok {
			return errors.New("param error"), nil
		}
		sourceId, ok := x["sourceFieldId"].(float64)
		if (!ok || sourceId <= 0) && FIELD_USE_TYPE[t] == C {
			return errors.New("param error"), nil
		}

		err, f := createField(srv.db, Field{
			Title:              title,
			Key:                key,
			Type:               t,
			FormID:             formId,
			IsOnly:             &repeat,
			IsNecessary:        &necessary,
			DatasourceColumnID: uint64(sourceId),
			Content:            filedContent,
		})
		if err != nil {
			return err, nil
		}
		err, flowId := srv.formatField(formId, f.ID, t, filedContent)
		if err != nil {
			return err, nil
		}
		if len(flowId) > 0 {
			flows = append(flows, flowId...)
		}
	}
	return nil, flows
}

func (srv *formService) deleteField(id uint64) error {
	if err := deleteMultiForm(srv.db, id); err != nil {
		return err
	}
	if err := deleteLinkage(srv.db, id); err != nil {
		return err
	}
	if err := deleteRecords(srv.db, id); err != nil {
		logrus.Error("delete field-records failed:", err)
		return err
	}
	if err := deleteButton(srv.db, id); err != nil {
		logrus.Error("delete field-button failed:", err)
		return err
	}
	if err := deleteTable(srv.db, id); err != nil {
		logrus.Error("delete field-table failed:", err)
		return err
	}
	if err := deleteTableColumn(srv.db, id); err != nil {
		logrus.Error("delete field-table-column failed:", err)
		return err
	}
	if err := deleteTableButton(srv.db, id); err != nil {
		logrus.Error("delete field-table-button failed:", err)
		return err
	}
	if err := deleteTableFilter(srv.db, id); err != nil {
		logrus.Error("delete field-table-button failed:", err)
		return err
	}

	return deleteField(srv.db, id)
}

func (list FormTreeList) Ids() []uint64 {
	return list.ids([]uint64{})
}

func (list FormTreeList) ids(group []uint64) []uint64 {
	for _, x := range list {
		group = append(group, x.Id)
		if len(x.Children) > 0 {
			group = x.Children.ids(group)
		}
	}
	return group
}
