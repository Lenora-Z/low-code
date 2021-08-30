// Package server Created by Goland
//@User: lenora
//@Date: 2021/3/12
//@Time: 10:20 上午

package server

import (
	"encoding/json"
	"github.com/Lenora-Z/low-code/service/flow"
	"github.com/Lenora-Z/low-code/service/form"
	"github.com/Lenora-Z/low-code/utils"
	"github.com/gin-gonic/gin"
	"time"
)

type FormArg struct {
	Name          string             `json:"name" validate:"required"`      //表单名称
	Type          int8               `json:"type" validate:"required"`      //表单类型 1-标准表单 2-弹窗
	PageType      uint8              `json:"page_type" validate:"required"` //页面类型 1-表单 2-展示页
	SourceTableId uint64             `json:"source_table_id"`               //数据表id
	Description   string             `json:"description"`                   //表单描述
	Content       []form.FieldFormat `json:"plugins"`                       //控件信息
	Status        bool               `json:"status"`                        //表单状态
	Footer        form.FieldFormat   `json:"footer"`                        //底部配置
}

type FormListItem struct {
	Id         uint64    `json:"id"`
	No         string    `json:"no"`        //表单编号
	Type       int8      `json:"type"`      //表单类型
	PageType   uint8     `json:"page_type"` //页面类型
	Name       string    `json:"name"`      //表单名称
	IsBound    *bool     `json:"is_bound"`  //是否有流程依赖
	Status     *bool     `json:"status"`    //表单状态
	UpdateTime time.Time `json:"update_time"`
}

type FormListVO struct {
	List  []FormListItem `json:"list"`
	Count uint32         `json:"count"`
}

type FormEditArg struct {
	IdStruct
	FormArg
}

type FormCopyArg struct {
	IdStruct
	Content []form.FieldFormat `json:"plugins"` //控件信息
	Footer  form.FieldFormat   `json:"footer"`  //底部配置
}

type FormCopyVO struct {
	Id       uint64 `json:"id"`        //表单id
	Name     string `json:"name"`      //表单名称
	OriginId uint64 `json:"origin_id"` //源表id
}

type FormItemVO struct {
	*form.Form
	Plugins []form.FieldFormat `json:"plugins"`
	Footer  form.FieldFormat   `json:"footer"`
}

type FieldAttrItem struct {
	Id     uint64                 `json:"id"`    //控件id
	Title  string                 `json:"title"` //控件名称
	Key    string                 `json:"key"`   //控件key
	Type   string                 `json:"type"`  //控件类型
	Config map[string]interface{} `json:"config"`
}

type ServiceFields struct {
	FormId    uint64 `json:"form_id"`
	FormName  string `json:"form_name"`
	FieldId   uint64 `json:"field_id"`
	FieldName string `json:"field_name"`
}

type ServiceFieldsVO []ServiceFields

type RelatedForm struct {
	FormId    uint64          `json:"form_id"`    //表单id
	Name      string          `json:"name"`       //表单名称
	FieldName string          `json:"field_name"` //控件名称
	Fields    []FieldAttrItem `json:"fields"`     //字段集合
}

type FlowButtons struct {
	FormId     uint64 `json:"form_id"`     //表单id
	FormName   string `json:"form_name"`   //表单名称
	ButtonName string `json:"button_name"` //按钮名称
}

type FlowMappingVO struct {
	Id      uint64 `json:"id"`       //表单id
	Name    string `json:"name"`     //表单名称
	TableId uint64 `json:"table_id"` //数据表id
}

type FlowFormMap struct {
	FormList  []FlowMappingVO `json:"form_list"`  //页面列表
	TableList []FlowMappingVO `json:"table_list"` //列表控件列表
}

// @Summary 新增表单[update]
// @Tags form
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @param FormArg body FormArg true "请求体"
// @Success 200 {object} ApiResponse{result=object}
// @Router /form/new [post]
func (ds *defaultServer) NewForm(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*CustomClaims)
	var arg FormArg
	if err := ctx.BindJSON(&arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}
	if err := validatorInstance().Struct(arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	if arg.PageType == form.FORM_PAGE && arg.SourceTableId == 0 {
		ds.InvalidParametersError(ctx)
		return

	}

	//名字不可重复
	if status, _ := (form.NewFormService(ds.db)).GetFormDetailByName(arg.Name, claims.AppId); !status {
		ds.ResponseError(ctx, 3002)
		return
	}

	if arg.Status {
		if len(arg.Content) == 0 {
			ds.InvalidParametersError(ctx)
			return
		}
	}

	db := ds.db.Begin()
	srv := form.NewFormService(db)
	footer, _ := json.Marshal(arg.Footer)
	err, item, flows := srv.CreateForm(
		claims.AppId,
		arg.SourceTableId,
		arg.Type,
		arg.PageType,
		arg.Name,
		arg.Description,
		arg.Content,
		footer,
		arg.Status,
	)
	if err != nil {
		db.Rollback()
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	if arg.Status {
		flowIds := make([]uint64, 0, cap(flows))
		flowSrv := flow.NewFlowService(db)
		for _, x := range flows {
			if x != 0 && !utils.IsContainUInt64(flowIds, x) {
				if err, _ := flowSrv.CreateFlowMapping(item.ID, x, claims.AppId, ""); err != nil {
					db.Rollback()
					ds.InternalServiceError(ctx, err.Error())
					return
				}
				flowIds = append(flowIds, x)
			}
		}
	}

	db.Commit()
	ds.ResponseSuccess(ctx, nil)
	return
}

// GetFormList
// @Summary 表单列表[update]
// @Description 页码不传获取的为生效表单 传了页码则可分页获取全部表单
// @Tags form
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @param page query string false "页码"
// @param limit query string false "页容量"
// @param type query string false "表单分类 1-标准表单 2-弹窗"
// @param page_type query string false "表单分类 1-表单 2-显示"
// @param name query string false "表单名称"
// @Success 200 {object} ApiResponse{result=FormListVO}
// @Router /form/list [get]
func (ds *defaultServer) GetFormList(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*CustomClaims)
	pageStr := ctx.Query("page")
	limitStr := ctx.Query("limit")
	name := ctx.Query("name")

	typeStr := ctx.Query("type")
	if typeStr == "" {
		typeStr = "0"
	}

	page := utils.NewStr(pageStr).Uint32()
	limit := utils.NewStr(limitStr).Uint32()
	types := utils.NewStr(typeStr).Int8()

	srv := form.NewFormService(ds.db)
	var err error
	var count uint32
	var list form.FormList

	formList := make([]FormListItem, 0, 100)
	if page == 0 || limit == 0 {
		err, count, list = srv.FullFormList(claims.AppId, types)
		if err != nil {
			ds.InternalServiceError(ctx, err.Error())
			return
		}
		for _, x := range list {
			formList = append(formList, FormListItem{
				Id:         x.ID,
				No:         x.Number,
				Type:       x.Type,
				PageType:   x.PageType,
				Name:       x.Name,
				IsBound:    &FALSE,
				Status:     x.Status,
				UpdateTime: x.UpdatedAt,
			})
		}
	} else {
		err, count, list = srv.FormList(page, limit, claims.AppId, types, name)
		if err != nil {
			ds.InternalServiceError(ctx, err.Error())
			return
		}
		flowSrv := flow.NewFlowService(ds.db)
		err, _, mapping := flowSrv.GetMappingByFormId(claims.AppId, list.Ids())
		if err != nil {
			ds.InternalServiceError(ctx, err.Error())
			return
		}
		for _, x := range list {
			status := false
			c := len(mapping.ItemByForm(x.ID))
			if c > 0 {
				status = true
			}
			formList = append(formList, FormListItem{
				Id:         x.ID,
				No:         x.Number,
				Type:       x.Type,
				PageType:   x.PageType,
				Name:       x.Name,
				IsBound:    &status,
				Status:     x.Status,
				UpdateTime: x.UpdatedAt,
			})
		}
	}
	ds.ResponseSuccess(ctx, FormListVO{
		List:  formList,
		Count: count,
	})
	return
}

// @Summary 编辑表单
// @Tags form
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @param FormEditArg body FormEditArg true "请求体"
// @Success 200 {object} ApiResponse{result=object}
// @Router /form/edit [post]
func (ds *defaultServer) EditForm(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*CustomClaims)
	var arg FormEditArg
	if err := ctx.BindJSON(&arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}
	if err := validatorInstance().Struct(arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	preSrv := form.NewFormService(ds.db)

	if status, f := preSrv.GetFormDetailByName(arg.Name, claims.AppId); (!status) && (f.ID != arg.Id) {
		ds.ResponseError(ctx, 3002)
		return
	}

	if arg.Status {
		if len(arg.Content) == 0 {
			ds.InvalidParametersError(ctx)
			return
		}
	}

	//判断表单状态t
	status, f := preSrv.GetFormDetail(arg.Id)
	if status || *f.IsDelete {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}
	//发布版本中的表单不可修改
	if *(f.IsOnline) == true {
		ds.ResponseError(ctx, 3003, "form can't be changed")
		return
	}

	db := ds.db.Begin()
	srv := form.NewFormService(db)
	footer, _ := json.Marshal(arg.Footer)
	err, item, flows := srv.UpdateForm(arg.Id, f.PageType, arg.Name, arg.Description, arg.Content, footer, arg.Status)
	if err != nil {
		db.Rollback()
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	flowSrv := flow.NewFlowService(db)
	if err := flowSrv.DelMappingByFormId(arg.Id); err != nil {
		db.Rollback()
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	if arg.Status {
		flowIds := make([]uint64, 0, cap(flows))
		for _, x := range flows {
			if x != 0 && !utils.IsContainUInt64(flowIds, x) {
				if err, _ := flowSrv.CreateFlowMapping(item.ID, x, claims.AppId, ""); err != nil {
					db.Rollback()
					ds.InternalServiceError(ctx, err.Error())
					return
				}
				flowIds = append(flowIds, x)
			}
		}
	}

	db.Commit()
	ds.ResponseSuccess(ctx, nil)
	return
}

// CopyForm
//复制表单
func (ds *defaultServer) CopyForm(ctx *gin.Context) {
	var arg FormCopyArg
	if err := ctx.BindJSON(&arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}
	if err := validatorInstance().Struct(arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	db := ds.db.Begin()
	srv := form.NewFormService(db)
	err, item := srv.CopyForm(arg.Id, arg.Content, arg.Footer)
	if err != nil {
		db.Rollback()
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	db.Commit()
	ds.ResponseSuccess(ctx, FormCopyVO{item.ID, item.Name, item.From})
	return
}

// @Summary 表单详情
// @Tags form
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @param id query string true "表单id"
// @Success 200 {object} ApiResponse{result=FormItemVO}
// @Router /form/detail [get]
func (ds *defaultServer) FormDetail(ctx *gin.Context) {
	idStr := ctx.Query("id")
	if idStr == "" {
		ds.InvalidParametersError(ctx)
		return
	}

	id := utils.NewStr(idStr)
	srv := form.NewFormService(ds.db)
	status, item := srv.GetFormDetail(id.Uint64())
	if status || *item.IsDelete == true {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}

	//控件内容及底部信息以json形式返回前端
	ret := make([]form.FieldFormat, 0, cap(item.Content))
	if err := json.Unmarshal(item.Content, &ret); err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	foot := form.FieldFormat{}
	if item.Footer != nil {
		if err := json.Unmarshal(item.Footer, &foot); err != nil {
			ds.InternalServiceError(ctx, err.Error())
			return
		}
	}
	item.Content = []byte{}
	ds.ResponseSuccess(ctx, FormItemVO{
		Form:    item,
		Plugins: ret,
		Footer:  foot,
	})
	return
}

// @Summary 表单字段列表
// @Tags form
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @param id query string true "表单id"
// @Success 200 {object} ApiResponse{result=[]FieldAttrItem}
// @Router /form/fields [get]
func (ds *defaultServer) FieldList(ctx *gin.Context) {
	idStr := ctx.Query("id")
	if idStr == "" {
		ds.InvalidParametersError(ctx)
		return
	}

	srv := form.NewFormService(ds.db)
	err, list := srv.FieldList(utils.NewStr(idStr).Uint64())
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	ll := make([]FieldAttrItem, 0, cap(list))

	//仅返回收集数据型的控件
	for _, x := range list {
		if form.FIELD_USE_TYPE[x.Type] != form.C {
			continue
		}
		cf := make(map[string]interface{})
		_ = json.Unmarshal(x.Content, &cf)
		ll = append(ll, FieldAttrItem{
			Id:     x.ID,
			Title:  x.Title,
			Key:    x.Key,
			Type:   x.Type,
			Config: cf,
		})
	}

	ds.ResponseSuccess(ctx, ll)
	return
}

// RelatedForm
// @Summary 获取关联列表
// @Tags field
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @param id query string false "表单id"
// @Success 200 {object} ApiResponse{result=[]RelatedForm}
// @Router /form/relation [get]
func (ds *defaultServer) RelatedForm(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*CustomClaims)
	idStr := ctx.Query("id")
	if idStr == "" {
		ds.InvalidParametersError(ctx)
		return
	}
	id := utils.NewStr(idStr)
	srv := form.NewFormService(ds.db)
	//查询表单信息
	notFound, item := srv.GetFormDetail(id.Uint64())
	if notFound || *item.IsDelete == true {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}

	if *item.Status == false {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}
	//if(复制来的表单){查询关系使用源表单的id}
	formId := id.Uint64()
	if item.From != 0 {
		formId = item.From
	}
	err, list := srv.GetRelationForm(formId)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	retList := make([]RelatedForm, 0, cap(list))

	if len(list) <= 0 {
		ds.ResponseSuccess(ctx, retList)
		return
	}

	//子表单详情
	err, _, forms := srv.MultiFormList(claims.AppId, list.Children())
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	if len(forms) <= 0 {
		ds.ResponseSuccess(ctx, retList)
		return
	}

	err, formField := srv.FieldList(id.Uint64())
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	formFields := make([]form.Field, 0, cap(formField))
	for _, ff := range formField {
		if ff.Type == form.MULTIFORM {
			formFields = append(formFields, ff)
		}
	}

	//获取全部控件属性
	err, fields := srv.FieldListByMultiFormIds(list.Children())
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	for _, x := range list {
		for _, fm := range forms {
			if x.ChildID == fm.ID {
				for _, mf := range formFields {
					if mf.ID == x.FieldID {
						fs := make([]FieldAttrItem, 0, cap(fields))
						for _, flds := range fields {
							if form.FIELD_USE_TYPE[flds.Type] != form.C || flds.Type == form.FILE || flds.Type == form.AUTOGRAPH || flds.Type == form.LINKAGE {
								continue
							}
							if flds.FormID == x.ChildID {
								fs = append(fs, FieldAttrItem{
									Id:    flds.ID,
									Title: flds.Title,
									Key:   flds.Key,
									Type:  flds.Type,
								})
							}
						}
						retList = append(retList, RelatedForm{
							FormId:    x.ChildID,
							Name:      fm.Name,
							FieldName: mf.Title,
							Fields:    fs,
						})
					}
				}
			}
		}
	}

	ds.ResponseSuccess(ctx, retList)
	return

}

// CollectFormList
// @Summary 获取表单页表单
// @Tags form
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @param status query string false "是否获取生效表单,否时不传或传空值,是时传递有效字符串即可"
// @Success 200 {object} ApiResponse{result=[]FormItemVO}
// @Router /form/list/collect [get]
func (ds *defaultServer) CollectFormList(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*CustomClaims)
	statusStr := ctx.Query("status")
	var status bool
	if statusStr != "" {
		status = true
	}
	srv := form.NewFormService(ds.db)
	err, _, forms := srv.FormListByPageType(claims.AppId, form.FORM_PAGE, status)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	list := make([]FormItemVO, 0, cap(forms))
	for i, x := range forms {
		ret := make([]form.FieldFormat, 0, cap(x.Content))
		if x.Content != nil {
			if err := json.Unmarshal(x.Content, &ret); err != nil {
				ds.InternalServiceError(ctx, err.Error())
				return
			}
		}
		foot := form.FieldFormat{}
		if x.Footer != nil {
			if err := json.Unmarshal(x.Footer, &foot); err != nil {
				ds.InternalServiceError(ctx, err.Error())
				return
			}
		}
		forms[i].Content = []byte{}
		aa := FormItemVO{
			Form:    &forms[i],
			Plugins: ret,
			Footer:  foot,
		}
		list = append(list, aa)
	}
	ds.ResponseSuccess(ctx, list)
	return

}

// DeleteForm
// @Summary 表单删除[new]
// @Tags form
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @param IdStruct body IdStruct true "id"
// @Success 200 {object} ApiResponse{result=object}
// @Router /form/delete [post]
func (ds *defaultServer) DeleteForm(ctx *gin.Context) {
	var arg IdStruct
	if err := ctx.BindJSON(&arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}
	if err := validatorInstance().Struct(arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}
	notFound, item := form.NewFormService(ds.db).GetFormDetail(arg.Id)
	if notFound || *item.IsDelete == true {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}
	if *item.IsOnline == true {
		ds.ResponseError(ctx, 3003, "form can't be deleted")
		return
	}

	db := ds.db.Begin()
	srv := form.NewFormService(db)
	if err := srv.DelForm(arg.Id); err != nil {
		db.Rollback()
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	flowSrv := flow.NewFlowService(db)
	if err := flowSrv.DelMappingByFormId(arg.Id); err != nil {
		db.Rollback()
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	db.Commit()
	ds.ResponseSuccess(ctx, nil)
	return

}

// FlowButton
// @Summary 绑定流程的按钮列表[new]
// @Tags form
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @param id query string true "流程id"
// @Success 200 {object} ApiResponse{result=[]FlowButtons}
// @Router /form/flow/button [get]
func (ds *defaultServer) FlowButton(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*CustomClaims)
	idStr := ctx.Query("id")
	if idStr == "" {
		ds.InvalidParametersError(ctx)
		return
	}

	flowId := utils.NewStr(idStr).Uint64()
	srv := form.NewFormService(ds.db)
	err, buttons := srv.ButtonByFlowId(flowId)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	err, _, forms := srv.MultiFormList(claims.AppId, buttons.Forms())
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	list := make([]FlowButtons, 0, cap(buttons))
	for _, b := range buttons {
		var formId uint64
		var name string
		for _, f := range forms {
			if b.FormID == f.ID {
				formId = f.ID
				name = f.Name
				break
			}
		}
		if formId == 0 {
			continue
		}
		list = append(list, FlowButtons{
			FormId:     formId,
			FormName:   name,
			ButtonName: b.Name,
		})
	}
	ds.ResponseSuccess(ctx, list)
	return

}

// FlowMapping
// @Summary 流程绑定的所有表单[new]
// @Tags form
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @param id query string true "流程id"
// @Success 200 {object} ApiResponse{result=FlowFormMap}
// @Router /form/mapping [get]
func (ds *defaultServer) FlowMapping(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*CustomClaims)
	idStr := ctx.Query("id")
	if idStr == "" {
		ds.InvalidParametersError(ctx)
		return
	}

	flowId := utils.NewStr(idStr).Uint64()

	srv := form.NewFormService(ds.db)
	err, btnMap := srv.FieldButtonByFlowId(flowId)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	formMap := make([]FlowMappingVO, 0, cap(btnMap))
	if len(btnMap) > 0 {
		err, _, forms := srv.MultiFormList(claims.AppId, btnMap.Forms())
		if err != nil {
			ds.InternalServiceError(ctx, err.Error())
			return
		}

		for _, x := range forms {
			formMap = append(formMap, FlowMappingVO{
				Id:      x.ID,
				Name:    x.Name,
				TableId: x.DatasourceTableID,
			})
		}
	}

	err, tableMap := srv.TableButtonByFlowId(flowId)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	fieldMap := make([]FlowMappingVO, 0, cap(tableMap))
	if len(tableMap) > 0 {
		err, tables := srv.FieldListByMultiIds(tableMap.Fields())
		if err != nil {
			ds.InternalServiceError(ctx, err.Error())
			return
		}

		for _, x := range tables {
			fieldMap = append(fieldMap, FlowMappingVO{
				Id:      x.ID,
				Name:    x.Title,
				TableId: x.DatasourceColumnID,
			})
		}
	}

	ds.ResponseSuccess(ctx, FlowFormMap{
		FormList:  formMap,
		TableList: fieldMap,
	})
	return

}

// TablaFieldColumns
// @Summary 列表控件字段列表[new]
// @Tags form
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @param id query string true "列表控件id"
// @Success 200 {object} ApiResponse{result=[]form.FieldTableColumn}
// @Router /form/table/columns [get]
func (ds *defaultServer) TablaFieldColumns(ctx *gin.Context) {
	idStr := ctx.Query("id")
	if idStr == "" {
		ds.InvalidParametersError(ctx)
		return
	}

	id := utils.NewStr(idStr).Uint64()
	srv := form.NewFormService(ds.db)
	err, list := srv.TableColumns(id)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	ds.ResponseSuccess(ctx, list)
	return

}
