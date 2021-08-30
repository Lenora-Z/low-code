//Created by Goland
//@User: lenora
//@Date: 2021/3/16
//@Time: 9:01 下午
package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Lenora-Z/low-code/service/bpmn"
	"github.com/Lenora-Z/low-code/service/data_log"
	"github.com/Lenora-Z/low-code/service/flow"
	"github.com/Lenora-Z/low-code/service/form"
	"github.com/Lenora-Z/low-code/service/form_data"
	"github.com/Lenora-Z/low-code/service/formdata"
	"github.com/Lenora-Z/low-code/service/mongodb"
	"github.com/Lenora-Z/low-code/service/user"
	"github.com/Lenora-Z/low-code/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	DATA_LIST uint64 = iota + 1
	DEAL_APPLY
	FORM_ITEM
)

type FormDetailVO struct {
	FormItemVO
	Child map[uint64]FormItemVO `json:"child"`
}

type FormDataArg struct {
	Id    uint64                   `json:"id" form:"id" validate:"required"`
	Data  []map[string]interface{} `json:"data" form:"data"`   //要提交的数据(日期与日期范围使用时间戳格式,日期范围传递时间戳格式的数组)
	Child []FormDataArg            `json:"child" form:"child"` //关联的多表单的数据
}

/*// AddFormData
// @Summary 提交表单数据[update]
// @Tags client
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param Authorization header string true "Bearer +token"
// @param FormDataArg body FormDataArg true "请求体"
// @Success 200 {object} ApiResponse{result=object}
// @Router /api/form/submit [post]*/
func (ds *defaultServer) AddFormData(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*UserClaims)
	var arg FormDataArg
	if err := ctx.Bind(&arg); err != nil {
		logrus.Error(err)
		ds.InvalidParametersError(ctx)
		return
	}
	if err := validatorInstance().Struct(arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	formSrv := form.NewFormService(ds.db)
	flowSrv := flow.NewFlowService(ds.db)
	formData := form_data.NewFormDataService(ds.mongoDb)
	logSrv := data_log.NewService(ds.mongoDb)

	//获取所有要共同提交的表单id
	//tree := formSrv.GetAllRelationForm(arg.Id, claims.AppId, []form.FormTree{})
	//for _, id := range tree.Ids() {
	//	_temp := make([]map[string]interface{}, 0)
	//	for _, c := range arg.Child {
	//		if c.Id == id {
	//			_temp = c.Data
	//			break
	//		}
	//	}
	//	if len(_temp) == 0 {
	//		ds.InvalidParametersError(ctx)
	//		return
	//	}
	//}

	type ChildData struct {
		form.Form
		Data []map[string]interface{}
	}
	type Children struct {
		ChildData
		Mode uint8
	}

	err, child := formSrv.GetRelationForm(arg.Id)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	childData := make([]ChildData, 0, cap(child))
	children := make([]Children, 0, cap(childData))

	if len(child) > 0 {
		err, _, childForm := formSrv.MultiFormList(claims.AppId, child.Children())
		if err != nil {
			ds.InternalServiceError(ctx, err.Error())
			return
		}
		for _, childItem := range childForm {
			_temp := false
			for _, c := range arg.Child {
				if c.Id == childItem.ID {
					_temp = true
					childData = append(childData, ChildData{
						Form: childItem,
						Data: c.Data,
					})
					break
				}
			}
			if !_temp {
				ds.InvalidParametersError(ctx)
				return
			}
		}

		for _, chi := range childData {
			for _, c := range child {
				if chi.ID == c.ChildID {
					children = append(children, Children{
						ChildData: chi,
						Mode:      c.Mode,
					})
				}
			}
		}

	}

	//mongo存储
	session, err := ds.mongoDb.Client().StartSession()
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	if err = session.StartTransaction(); err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	if err = mongo.WithSession(context.Background(), session, func(sessionContext mongo.SessionContext) error {
		//存储本表单数据
		err, self := ds.submitFormData(
			ctx,
			Srv{formSrv, flowSrv, formData, logSrv},
			arg.Id,
			claims.UserId,
			[]map[string]interface{}{arg.Data[0]},
			nil,
		)
		if err != nil {
			return err
		}
		//存储子表单数据
		for _, v := range children {
			if v.Mode == form.SINGLE {
				v.Data = []map[string]interface{}{v.Data[0]}
			}
			if err, _ := ds.submitFormData(
				ctx,
				Srv{formSrv, flowSrv, formData, logSrv},
				v.ID,
				claims.UserId,
				v.Data, self[0],
			); err != nil {
				return err
			}

		}

		if err := session.CommitTransaction(context.Background()); err != nil {
			return err
		}
		return nil
	}); err != nil {
		if errs := session.AbortTransaction(context.Background()); errs != nil {
			ds.InternalServiceError(ctx, errs.Error())
			return
		}
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	session.EndSession(context.Background())

	ds.ResponseSuccess(ctx, nil)
	return

}

type FormDataEditArg struct {
	ObjId string `json:"obj_id" form:"obj_id" validate:"required"` //数据id
	FormDataArg
}

/*// @Summary 修改单条表单数据
// @Tags client
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param Authorization header string true "Bearer +token"
// @param route-id header string true "路由id"
// @param FormDataEditArg body FormDataEditArg true "请求体"
// @Success 200 {object} ApiResponse{result=object}
// @Router /api/form/edit [post]*/
func (ds *defaultServer) EditFormData(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*UserClaims)
	var arg FormDataEditArg
	if err := ctx.Bind(&arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}
	if err := validatorInstance().Struct(arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	dataSrv := form_data.NewFormDataService(ds.mongoDb)
	status, item := dataSrv.GetByObjId(arg.ObjId)
	if !status {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}
	formStr := fmt.Sprintf("%v", item["form_id"])
	formId := utils.NewStr(formStr).Uint64()
	if formId == 0 {
		logrus.Error("form_id is not uint64")
		ds.InternalServiceError(ctx)
		return
	}

	formSrv := form.NewFormService(ds.db)
	err, fields := formSrv.FieldList(formId)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	//控件字段值收集
	err, data := ds.validateItem(ctx, arg.Data[0], fields, claims.AppId, 0, data_log.EDIT, arg.ObjId)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	//数据更新
	if err, _ := dataSrv.UpdateItem(arg.ObjId, data); err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	//日志记录
	logSrv := data_log.NewService(ds.mongoDb)
	if err, _ := logSrv.NewLog(claims.UserId, formId, data_log.EDIT, data); err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	ds.ResponseSuccess(ctx, nil)
	return
}

/*// FormDataDetail
// @Summary 表单数据详情
// @Tags client
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param Authorization header string true "Bearer +token"
// @param route-id header string true "路由id"
// @param id query string true "记录id"
// @param type query string false "查看方式(编辑用-edit   查看用-show)"
// @Success 200 {object} ApiResponse{result=object}
// @Router /api/form/detail [get]*/
func (ds *defaultServer) FormDataDetail(ctx *gin.Context) {
	objId := ctx.Query("id")

	types := ctx.Query("type")
	code, msg, item := ds.formDataItemDetail(1, objId, types)
	if code != SUCCESS {
		ds.ResponseError(ctx, code, msg)
		return
	}

	dataSrv := form_data.NewFormDataService(ds.mongoDb)
	err, _, child := dataSrv.GetListByParentId(item["_id"])
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	children := make(map[uint64][]interface{}, 0)
	for _, v := range child {
		idStr := fmt.Sprintf("%v", item["form_id"])
		id := utils.NewStr(idStr).Uint64()
		if id == 0 {
			continue
		}
		children[id] = append(children[id], v)
	}
	item["children"] = children

	ds.ResponseSuccess(ctx, item)
	return
}

type FormDataListArg struct {
	Id     string                 `json:"id" validate:"required"`    //控件id(table控件)
	Page   uint32                 `json:"page" validate:"required"`  //页码
	Limit  uint32                 `json:"limit" validate:"required"` //页容量
	Search map[string]interface{} `json:"search"`                    //查询条件(k-v格式传值,k是控件key,v是筛选条件,时间区间以时间戳(s)数组传递)
}

type FormDataItemVO struct {
	form_data.FormDataItem
	UserName string `json:"user_name"`
}

type FormDataListVO struct {
	List  []map[string]interface{} `json:"list"`
	Count uint32                   `json:"count"`
}

type Column struct {
	FormId uint64 `json:"formId"`
	Title  string `json:"title"`
	Key    string `json:"key"`
	Name   string `json:"name"`
	Type   string `json:"type"`
}

type ContentStruct struct {
	FormId        uint64   `json:"formId"`
	Columns       []Column `json:"columns"`
	HasCustom     bool     `json:"hasCustomColumn"`
	DataFilter    bool     `json:"dataFilter"`
	DataFilterKey string   `json:"dataFilterKey"`
	Condition     []string `json:"condition"`
}

/*// FormDataList
// @Summary 表单数据列表(控件id)
// @Tags client
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param Authorization header string true "Bearer +token"
// @param route-id header string true "路由id"
// @param FormDataListArg body FormDataListArg true "请求体"
// @Success 200 {object} ApiResponse{result=FormDataListVO}
// @Router /api/form/data [post]*/
func (ds *defaultServer) FormDataList(ctx *gin.Context) {
	var arg FormDataListArg
	if err := ctx.BindJSON(&arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}
	if err := validatorInstance().Struct(arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	//获取控件信息并校验类型
	srv := form.NewFormService(ds.db)
	status, field := srv.GetFieldDetailByKey(arg.Id)
	if status || field.Type != form.TABLES {
		logrus.Error("field not found")
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}
	kk := string(field.Content)
	logrus.Info("field info:", kk)

	//格式化表单字段配置
	var content ContentStruct
	err := json.Unmarshal(field.Content, &content)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	//格式化筛选条件
	filter := make(map[string]interface{})
	mongoSrv := form_data.NewFormDataService(ds.mongoDb)
	if len(arg.Search) > 0 {
		appId := ctx.MustGet(CLAIMS).(*UserClaims).AppId
		//字表查询字段筛选
		search := searchData(content, arg.Search)
		if len(search) > 1 {
			ids := make([]primitive.ObjectID, 0, 100)
			for k, v := range search {
				if k == 0 {
					continue
				}
				err, childFilter := ds.getFilter(appId, k, v, ds.getLanguage(ctx))
				if err != nil {
					ds.InternalServiceError(ctx, err.Error())
					return
				}
				err, count, childList := mongoSrv.FilterList(1, utils.MAX_LIMIT, k, childFilter)
				if count <= 0 {
					ds.ResponseSuccess(ctx, FormDataListVO{List: form_data.FormDataList{}, Count: 0})
					return
				}
				ids = append(ids, childList.ObjIds()...)
			}
			err, filter = ds.getFilter(appId, content.FormId, search[0], ds.getLanguage(ctx))
			if err != nil {
				ds.InternalServiceError(ctx, err.Error())
				return
			}
			filter["_id"] = mongodb.In(ids)

		} else {
			err, filter = ds.getFilter(appId, content.FormId, arg.Search, ds.getLanguage(ctx))
			if err != nil {
				ds.InternalServiceError(ctx, err.Error())
				return
			}
		}
	}
	if content.DataFilter && len(content.Condition) > 0 {
		notFound, filterKey := srv.GetFieldDetailByKey(content.DataFilterKey)
		if !notFound {
			switch filterKey.Type {
			case form.SINGLE_CHOICE:
				if len(content.Condition) == 1 {
					filter[content.DataFilterKey] = content.Condition[0]
				} else {
					filter[content.DataFilterKey] = mongodb.In(content.Condition)
				}
			case form.MULTIPLE_CHOICE:
				sort.Strings(content.Condition)
				filter[content.DataFilterKey] = strings.Join(content.Condition, "、")
			}
		}
	}

	//mongodb中数据
	err, count, list := mongoSrv.FilterList(arg.Page, arg.Limit, content.FormId, filter)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	userSrv := user.NewUserService(ds.db)
	err, users := userSrv.GetMultiUser(list.UserIds())
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	logrus.Info(users)

	//特殊控件的处理
	for i, x := range list {
		for _, c := range content.Columns {
			if content.FormId == c.FormId {
				v, ok := x[c.Key]
				if !ok {
					continue
				}
				switch c.Type {
				//成员控件获取用户的真实姓名
				case form.MEMBER:
					e, names := ds.getMembers(v, userSrv)
					if e != nil && e.Error() != "" {
						ds.InternalServiceError(ctx, e.Error())
						return
					}
					list[i][c.Key] = strings.Join(names, ",")
				//组织控件获取组成名称
				case form.ORGANIZATION:
					e, names := ds.getOrganizations(v, userSrv)
					if e != nil && e.Error() != "" {
						ds.InternalServiceError(ctx, e.Error())
						return
					}
					list[i][c.Key] = strings.Join(names, ",")
				case "created_at":
					t, ok := v.(primitive.DateTime)
					if ok {
						list[i][c.Key] = t.Time().Format(utils.TimeFormatStr)
					} else {
						ti, ok := v.(float64)
						if !ok {
							continue
						}
						list[i][c.Key] = utils.Time(ti).Format(utils.TimeFormatStr)
					}
				default:
					continue
				}
			} else {
				parent, _ := x["_id"].(primitive.ObjectID)
				err, _, child := mongoSrv.FilterList(1, utils.MAX_LIMIT, c.FormId, map[string]interface{}{
					"parent_id": parent.Hex(),
				})
				if err != nil {
					ds.InternalServiceError(ctx, err.Error())
					return
				}
				con := make([]string, 0, cap(child))
				for _, chi := range child {
					con = append(con, fmt.Sprintf("%v", chi[c.Key]))
				}
				list[i][c.Key] = strings.Join(con, "、")
			}
		}
		//是否展示用户名称
		if content.HasCustom {
			x["user_name"] = ""
			idStr := fmt.Sprintf("%v", x["user_id"])
			id := utils.NewStr(idStr).Uint64()
			if id == 0 {
				continue
			}
			list[i]["user_name"] = users.GetUser(id).TrueName
		}
	}

	ds.ResponseSuccess(ctx, FormDataListVO{
		List:  list,
		Count: count,
	})
	return
}

/*//RelatedFormData @Summary 获取联级数据
//@Tags client
//@Accept  json
//@Produce  json
//@param lang header string false "语言(zh/en),默认en"
//@param app-id header string true "应用id"
//@param form_id query string true "表单id"
//@param id query string false "关联数据id"
//@Success 200 {object} ApiResponse{result=[]map[string]interface{}}
//@Router /api/form/related [get]*/
func (ds *defaultServer) RelatedFormData(ctx *gin.Context) {
	formStr := ctx.Query("form_id")
	id := ctx.Query("id")
	if formStr == "" {
		ds.InvalidParametersError(ctx)
		return
	}
	formId := utils.NewStr(formStr).Uint64()
	filter := make(map[string]interface{})
	if id != "" {
		filter["parent_id"] = id
	}
	dataSrv := form_data.NewFormDataService(ds.mongoDb)
	err, _, list := dataSrv.FilterList(1, utils.MAX_LIMIT, formId, filter)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	for i, x := range list {
		id, ok := x["_id"].(primitive.ObjectID)
		if !ok {
			continue
		}
		code, msg, item := ds.formDataItemDetail(1, id.Hex(), "show")
		if code != SUCCESS {
			ds.ResponseError(ctx, code, msg)
			return
		}
		list[i] = item
	}

	ds.ResponseSuccess(ctx, list)
	return
}

type FormDataLogArg struct {
	Id    uint64 `json:"id" validate:"required"`    //控件id
	Page  uint32 `json:"page" validate:"required"`  //页码
	Limit uint32 `json:"limit" validate:"required"` //页容量
}

/*// FormDataLog @Summary 表单数据列表(表单id)
// @Tags client
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param Authorization header string true "Bearer +token"
// @param route-id header string true "路由id"
// @param FormDataLogArg body FormDataLogArg true "请求体"
// @Success 200 {object} ApiResponse{result=FormDataListVO}
// @Router /api/form/log [post]*/
func (ds *defaultServer) FormDataLog(ctx *gin.Context) {
	var arg FormDataLogArg
	if err := ctx.BindJSON(&arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}
	if err := validatorInstance().Struct(arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	mongoSrv := form_data.NewFormDataService(ds.mongoDb)
	err, count, list := mongoSrv.GetListByFormId(arg.Page, arg.Limit, arg.Id)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	userSrv := user.NewUserService(ds.db)

	err, columns := form.NewFormService(ds.db).FieldList(arg.Id)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	for i, x := range list {
		for _, c := range columns {
			v, ok := x[c.Key]
			if !ok {
				continue
			}
			switch c.Type {
			case form.MEMBER:
				userIdInterface, ok := v.(bson.A)
				if !ok {
					continue
				}
				userIds := make([]uint64, 0, cap(userIdInterface))
				for _, ii := range userIdInterface {
					g, s := ii.(float64)
					if !s {
						continue
					}
					userIds = append(userIds, uint64(g))

				}
				err, u := userSrv.GetMultiUser(userIds)
				if err != nil {
					ds.InternalServiceError(ctx, err.Error())
					return
				}
				list[i][c.Key] = strings.Join(u.TrueNames(), ",")
			case form.ORGANIZATION:
				orgs, ok := v.(bson.A)
				if !ok {
					continue
				}
				org := make([]uint64, 0, cap(orgs))
				for _, o := range orgs {
					g, s := o.(float64)
					if !s {
						continue
					}
					org = append(org, uint64(g))
				}
				err, organization := userSrv.GetOrganizationByIdGroup(org)
				if err != nil {
					ds.InternalServiceError(ctx, err.Error())
					return
				}
				list[i][c.Key] = strings.Join(organization.Names(), ",")
			default:
				continue
			}
		}
	}

	ds.ResponseSuccess(ctx, FormDataListVO{
		List:  list,
		Count: count,
	})
	return
}

/*// FormDataListExport
// @Summary 表单数据列表导出
// @Tags client
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param Authorization header string true "Bearer +token"
// @param route-id header string true "路由id"
// @param id query string true "控件id"
// @Success 200 {object} ApiResponse{result=object}
// @Router /api/form/download [get]*/
func (ds *defaultServer) FormDataListExport(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		ds.InvalidParametersError(ctx)
		return
	}

	srv := form.NewFormService(ds.db)
	status, field := srv.GetFieldDetailByKey(id)
	if status || field.Type != form.TABLES {
		logrus.Error("field not found")
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}

	//格式化导出字段配置
	type contentStruct struct {
		FormId        uint64   `json:"formId"`
		NeedExport    bool     `json:"needExport"`
		ExportColumns []Column `json:"exportColumns"`
		HasCustom     bool     `json:"hasExportColumn"`
		DataFilter    bool     `json:"dataFilter"`
		DataFilterKey string   `json:"dataFilterKey"`
		Condition     []string `json:"condition"`
	}
	var content contentStruct
	_ = json.Unmarshal(field.Content, &content)
	if !content.NeedExport {
		ds.ResponseError(ctx, 2004)
		return
	}

	filter := make(map[string]interface{})
	if content.DataFilter && len(content.Condition) > 0 {
		notFound, filterKey := srv.GetFieldDetailByKey(content.DataFilterKey)
		if !notFound {
			switch filterKey.Type {
			case form.SINGLE_CHOICE:
				if len(content.Condition) == 1 {
					filter[content.DataFilterKey] = content.Condition[0]
				} else {
					filter[content.DataFilterKey] = mongodb.In(content.Condition)
				}
			case form.MULTIPLE_CHOICE:
				sort.Strings(content.Condition)
				filter[content.DataFilterKey] = strings.Join(content.Condition, "、")
			}
		}
	}

	//获取数据
	dataSrv := form_data.NewFormDataService(ds.mongoDb)
	err, _, list := dataSrv.FilterList(1, utils.MAX_LIMIT, content.FormId, filter)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	userSrv := user.NewUserService(ds.db)
	err, users := userSrv.GetMultiUser(list.UserIds())
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	datas := make([]map[string]interface{}, 0, cap(list))

	//字段值处理
	for i, x := range list {
		item := make(map[string]interface{})
		for _, c := range content.ExportColumns {
			if content.FormId == c.FormId {
				item[c.Name] = ""
				if c.Type == "index" {
					item[c.Name] = i + 1
				}
				v, ok := x[c.Key]
				if !ok {
					continue
				}
				item[c.Name] = ds.getItemData(c.Type, userSrv, v)
			} else {
				parent, _ := x["_id"].(primitive.ObjectID)
				err, _, child := dataSrv.FilterList(1, utils.MAX_LIMIT, c.FormId, map[string]interface{}{
					"parent_id": parent.Hex(),
				})
				if err != nil {
					ds.InternalServiceError(ctx, err.Error())
					return
				}
				con := make([]string, 0, cap(child))
				for _, chi := range child {
					if chi[c.Key] == nil {
						continue
					}
					con = append(con, fmt.Sprintf("%v", ds.getItemData(c.Type, userSrv, chi[c.Key])))
				}
				item[c.Name] = strings.Join(con, "、")
			}
		}
		//是否展示用户名称
		if content.HasCustom {
			idStr := fmt.Sprintf("%v", x["user_id"])
			id := utils.NewStr(idStr).Uint64()
			if id == 0 {
				continue
			}
			item["用户名称"] = users.GetUser(id).TrueName
		}
		datas = append(datas, item)
	}

	//生成导出数据表头
	top := make([]string, 0, cap(content.ExportColumns))
	for _, c := range content.ExportColumns {
		top = append(top, c.Name)
	}

	//创建excel文件
	err, src, contentType, header := utils.CreateExcel("form_data", top, datas)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	ctx.DataFromReader(http.StatusOK, int64(len(src)), contentType, bytes.NewReader(src), header)
	return
}

type DataIdStruct struct {
	Id string `json:"id"`
}

/*// DeleteFormData @Summary 删除单条表单数据
// @Tags client
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param Authorization header string true "Bearer +token"
// @param route-id header string true "路由id"
// @param DataIdStruct body DataIdStruct true "请求体"
// @Success 200 {object} ApiResponse{result=object}
// @Router /api/form/delete [post]*/
func (ds *defaultServer) DeleteFormData(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*UserClaims)
	var arg DataIdStruct
	if err := ctx.Bind(&arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}
	if err := validatorInstance().Struct(arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	dataSrv := form_data.NewFormDataService(ds.mongoDb)
	status, item := dataSrv.GetByObjId(arg.Id)
	if !status {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}

	formStr := fmt.Sprintf("%v", item["form_id"])
	formId := utils.NewStr(formStr).Uint64()
	if formId == 0 {
		logrus.Error("form_id is not uint64")
		ds.InternalServiceError(ctx)
		return
	}

	deleteStatus, ok := item["is_delete"].(bool)
	if !ok {
		logrus.Error("form_id is not boolean")
		ds.InternalServiceError(ctx)
		return
	}
	if deleteStatus {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}

	if err := dataSrv.DeleteItem(arg.Id); err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	data := map[string]interface{}{
		"data_id": arg.Id,
	}

	//日志记录
	logSrv := data_log.NewService(ds.mongoDb)
	if err, _ := logSrv.NewLog(claims.UserId, formId, data_log.DELETE, data); err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	ds.ResponseSuccess(ctx, nil)
	return
}

// GetFormDetail
// @Summary 使用端表单结构详情
// @Tags client
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param Authorization header string true "Bearer +token"
// @param route-id header string true "路由id"
// @param id query string true "表单id"
// @Success 200 {object} ApiResponse{result=FormDetailVO}
// @Router /api/form [get]
func (ds *defaultServer) GetFormDetail(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*UserClaims)
	idStr := ctx.Query("id")
	if idStr == "" {
		ds.InvalidParametersError(ctx)
		return
	}
	id := utils.NewStr(idStr)

	srv := form.NewFormService(ds.db)

	//tree := srv.GetAllRelationForm(id.Uint64(), claims.AppId, []form.FormTree{})
	list := make(map[uint64]FormItemVO)
	err, code, item := formDetail(srv, id.Uint64(), ds.getLanguage(ctx))
	if err != nil {
		ds.ResponseError(ctx, code, err.Error())
		return
	}
	//list = append(list, *item)
	//list.allFormDetail(srv, tree, ds.getLanguage(ctx))
	var searchId = item.ID
	if item.From != 0 {
		searchId = item.From
	}
	err, childForm := srv.GetRelationForm(searchId)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	if len(childForm) > 0 {
		err, _, children := srv.MultiFormList(claims.AppId, childForm.Children())
		if err != nil {
			ds.ResponseError(ctx, code, err.Error())
			return
		}
		for _, x := range children {
			err, code, v := formDetail(srv, x.ID, ds.getLanguage(ctx))
			if err != nil {
				ds.ResponseError(ctx, code, err.Error())
				return
			}
			list[x.ID] = *v
		}
	}
	ds.ResponseSuccess(ctx, FormDetailVO{
		FormItemVO: *item,
		Child:      list,
	})
	return
}

type TableDataArg struct {
	FieldName string                 `json:"fieldName" validate:"required"` //列表控件key
	Filter    map[uint32]interface{} `json:"filter"`                        //筛选字段及值
	Limit     uint32                 `json:"limit"`                         //每页条数，默认为10
	Page      uint32                 `json:"page"`                          //展示第几页，默认为1
}

// GetTableDataByFilter
// @Summary 筛选查看列表数据[new]
// @Tags client
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param Authorization header string true "Bearer +token"
// @param TableDataArg body TableDataArg true "请求体"
// @Success 200 {object} ApiResponse{result=object}
// @Router /api/form/table/list [post]
func (ds *defaultServer) GetTableDataByFilter(ctx *gin.Context) {
	var arg TableDataArg
	if err := ctx.Bind(&arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	if err := validatorInstance().Struct(arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	var page = arg.Page
	var limit = arg.Limit
	if arg.Page == 0 {
		page = 1 //默认值
	} else if arg.Page >= 100 {
		ds.ResponseError(ctx, WRONG_PARAM, "页码设置过大")
		return
	}
	if arg.Limit == 0 {
		limit = 10 //默认值
	} else if arg.Limit >= 200 {
		ds.ResponseError(ctx, WRONG_PARAM, "条数设置过多")
		return
	}

	formdataSrv := formdata.NewFormdataService(ds.db)
	formSrv := form.NewFormService(ds.db)
	userSrv := user.NewUserService(ds.db)
	//获取控件id
	isNotFound, field := formSrv.GetFieldDetailByKey(arg.FieldName)
	if isNotFound {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}

	//查询本业务数据表表名
	isNotFound, fieldTable := formdataSrv.GetFieldTableByFieldId(uint32(field.ID))
	if isNotFound {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}
	isNotFound, datasourceTable := formdataSrv.GetDatasourceTableById(fieldTable.DatasourceTableID)
	if isNotFound {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}

	//查询要显示的业务数据表字段
	var mainTables = make(map[string][]string, 0)      //本表信息     k=本表表名    v=[本表的显示字段a，本表的显示字段b...]
	var relationTables = make(map[string][]string, 0)  //关联表信息   k=关联表表名  v=[本表关联字段a，关联表的关联字段b，关联表的显示字段ids，关联表的显示字段c，关联表的显示字段d...]
	var filterCondition = make(map[string][]string, 0) //过滤条件     k=表名        v=["字段a = 值1","字段b like '%值2%'",...]
	var relatedFieldTableColumnList = make(formdata.FieldTableColumnList, 0)
	var mainFieldTableColumnList = make(formdata.FieldTableColumnList, 0)
	isNotFound, fieldTableColumnList := formdataSrv.GetFieldTableColumnByFieldId(uint32(field.ID))
	if isNotFound {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}
	for _, v := range fieldTableColumnList {
		if v.DatasourceColumnRelationID == 0 {
			//本数据表的字段
			mainFieldTableColumnList = append(mainFieldTableColumnList, v)
		} else {
			//关联数据表的字段
			relatedFieldTableColumnList = append(relatedFieldTableColumnList, v)
		}
	}

	var reflectIds = make([]uint32, 0) //需要查找映射关系的表字段id

	//本数据表的字段
	mainColIds := mainFieldTableColumnList.DatasourceColIds()
	err, mainColNames, _, mainReflectIds := formdataSrv.GetBusinessColNameByIds(mainColIds)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	mainTables[datasourceTable.TableName] = mainColNames
	mainColIds = append([]uint32{0}, mainColIds...)
	reflectIds = append(reflectIds, mainReflectIds...)

	//关联数据表的字段
	for relationId, colIds := range relatedFieldTableColumnList.RelationIdsAndColIds() {
		isNotFound, datasourceColumnRelation := formdataSrv.GetDatasourceColumnRelationById(relationId)
		if isNotFound {
			ds.ResponseError(ctx, NOT_FOUND)
			return
		}
		//查业务数据表表名
		isNotFound, datasourceTable := formdataSrv.GetDatasourceTableById(datasourceColumnRelation.TargetTableID)
		if isNotFound {
			ds.ResponseError(ctx, NOT_FOUND)
			return
		}

		//查业务数据表字段名
		//关联表的关联字段名
		err, relateColName, _, _ := formdataSrv.GetBusinessColNameByIds([]uint32{datasourceColumnRelation.TargetColumnID})
		if err != nil {
			ds.InternalServiceError(ctx, err.Error())
			return
		}
		//本表的关联字段名
		err, sourceColName, _, _ := formdataSrv.GetBusinessColNameByIds([]uint32{datasourceColumnRelation.SourceColumnID})
		if err != nil {
			ds.InternalServiceError(ctx, err.Error())
			return
		}
		//关联表要显示的字段名
		err, relateColNames, _, relateReflectIds := formdataSrv.GetBusinessColNameByIds(colIds)
		if err != nil {
			ds.InternalServiceError(ctx, err.Error())
			return
		}
		reflectIds = append(reflectIds, relateReflectIds...)

		var relateColIds string
		for _, colId := range colIds {
			relateColIds += strconv.FormatUint(uint64(colId), 10) + ","
		}
		relateColIds = strings.TrimRight(relateColIds, ",")

		relationTables[datasourceTable.TableName] = append([]string{sourceColName[0], relateColName[0], relateColIds}, relateColNames...)
	}

	//数据过滤信息
	if fieldTable.IsFilter == true {
		//需要数据过滤
		isNotFound, fieldTableFilterList := formdataSrv.GetFieldTableFilterByFieldId(uint32(field.ID))
		if isNotFound {
			ds.ResponseError(ctx, NOT_FOUND)
			return
		}

		for _, v := range fieldTableFilterList {
			err, whereSentences := formdataSrv.GetWhereSentence(&v)
			if err != nil {
				ds.ResponseError(ctx, NOT_FOUND)
				return
			}
			if v.DatasourceColumnRelationID == 0 {
				if _, ok := filterCondition[datasourceTable.TableName]; !ok {
					filterCondition[datasourceTable.TableName] = whereSentences
				} else {
					filterCondition[datasourceTable.TableName] = append(filterCondition[datasourceTable.TableName], whereSentences...)
				}
			} else {
				isNotFound, datasourceColumnRelation := formdataSrv.GetDatasourceColumnRelationById(v.DatasourceColumnRelationID)
				if isNotFound {
					ds.ResponseError(ctx, NOT_FOUND)
					return
				}
				//查业务数据表表名
				isNotFound, datasourceTable := formdataSrv.GetDatasourceTableById(datasourceColumnRelation.TargetTableID)
				if isNotFound {
					ds.ResponseError(ctx, NOT_FOUND)
					return
				}
				if _, ok := filterCondition[datasourceTable.TableName]; !ok {
					filterCondition[datasourceTable.TableName] = whereSentences
				} else {
					filterCondition[datasourceTable.TableName] = append(filterCondition[datasourceTable.TableName], whereSentences...)
				}
			}
		}
	}

	//处理筛选字段
	for colId, value := range arg.Filter {
		isNotFound, datasourceColumn := formdataSrv.GetDatasourceColumnById(colId)
		if isNotFound {
			ds.ResponseError(ctx, NOT_FOUND)
			return
		}

		var whereSentences []string
		if datasourceColumn.FieldType == "SingleChoice" {
			//单选控件传入选项id
			if v, ok := value.(string); !ok {
				ds.ResponseError(ctx, FAIL, "单选控件数据格式错误")
				return
			} else {
				whereSentences = []string{fmt.Sprintf("%s = %s", datasourceColumn.ColumnName, v)}
			}
		} else if datasourceColumn.FieldType == "Number" {
			//数值控件传入数字
			if v, ok := value.(float64); !ok {
				ds.ResponseError(ctx, FAIL, "数值控件数据格式错误")
				return
			} else {
				whereSentences = []string{fmt.Sprintf("%s = %d", datasourceColumn.ColumnName, int(v))}
			}
		} else if datasourceColumn.FieldType == "DateTime" {
			//日期控件传入时间戳
			switch value.(type) {
			case []interface{}:
				var times = make([]string, 0)
				for _, tmIntr := range value.([]interface{}) {
					if tm, ok := tmIntr.(float64); !ok {
						ds.ResponseError(ctx, FAIL, "日期控件数据格式错误")
						return
					} else {
						if len(times) == 1 {
							times = append(times, time.Unix(int64(tm)+3600*24, 0).Format("2006-01-02"))
						} else {
							times = append(times, time.Unix(int64(tm), 0).Format("2006-01-02"))
						}
					}
				}
				whereSentences = []string{fmt.Sprintf("%s > '%s'", datasourceColumn.ColumnName, times[0]), fmt.Sprintf("%s < '%s'", datasourceColumn.ColumnName, times[1])}
			case float64:
				whereSentences = []string{fmt.Sprintf("%s like '%s%%'", datasourceColumn.ColumnName, time.Unix(int64(value.(float64)), 0).Format("2006-01-02"))}
			default:
				ds.ResponseError(ctx, FAIL, "日期控件数据格式错误")
				return
			}
		} else if datasourceColumn.FieldType == "Member" {
			//成员控件传入成员id，默认只处理第一个   [12]
			if v, ok := value.([]interface{}); !ok {
				ds.ResponseError(ctx, FAIL, "成员控件数据格式错误")
				return
			} else {
				if vv, ok := v[0].(float64); !ok {
					ds.ResponseError(ctx, FAIL, "成员控件数据格式错误")
					return
				} else {
					whereSentences = []string{fmt.Sprintf("%s = %d", datasourceColumn.ColumnName, int(vv))}
				}
			}
		} else {
			if v, ok := value.(string); !ok {
				ds.ResponseError(ctx, FAIL, "控件数据格式错误")
				return
			} else {
				whereSentences = []string{fmt.Sprintf("%s like '%%%s%%'", datasourceColumn.ColumnName, v)}
			}
		}

		if _, ok := filterCondition[datasourceColumn.TableName]; !ok {
			filterCondition[datasourceColumn.TableName] = whereSentences
		} else {
			filterCondition[datasourceColumn.TableName] = append(filterCondition[datasourceColumn.TableName], whereSentences...)
		}
	}

	//查询业务数据
	err, total, data, relateColIds := formdataSrv.GetBusinessTableData(ds.businessDb, mainTables, relationTables, filterCondition, page, limit)
	if err != nil {
		logrus.Error("获取业务数据出错：", err.Error())
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	allColIds := append(mainColIds, relateColIds...)
	list := make([]map[string]interface{}, 0)
	for _, rowData := range data {
		var col = make(map[string]interface{}, 0)
		for k, v := range rowData {
			if k == 0 {
				col["_id"] = v
			} else {
				colId := strconv.FormatUint(uint64(allColIds[k]), 10)
				isCreatorId := formdataSrv.IsCreatorIdColumn(ds.db, allColIds[k])
				if !isCreatorId {
					userId, _ := strconv.ParseUint(v, 10, 64)
					isNotFound, userItem := userSrv.GetUser(userId)
					if isNotFound {
						col[colId] = ""
					} else {
						col[colId] = userItem.TrueName
					}
					continue
				}

				//令时间类型返回为空
				if v == "null" {
					col[colId] = nil
				} else {
					col[colId] = v
				}

				if len(v) == 0 {
					col[colId] = ""
					continue
				}

				for _, id := range reflectIds {
					if id == allColIds[k] {
						isNotFound, value := formdataSrv.GetDatasourceMetadataByKey(id, v)
						if isNotFound {
							col[colId] = ""
						} else {
							col[colId] = value
						}
						break
					}
				}
			}
		}
		list = append(list, col)
	}

	ds.ResponseSuccess(ctx, map[string]interface{}{
		"list":  list,
		"total": total,
	})
	return
}

type TableData struct {
	FieldArg
	Data map[uint32]interface{} `json:"data" validate:"required"` //更新的字段及值
}

type FieldArg struct {
	FieldName string `json:"fieldName" validate:"required"` //控件key
	RowId     uint32 `json:"rowId" validate:"required"`     //行数据id
}

// UpdateTableData
// @Summary 修改列表单条数据，仅可修改列表关联的本数据表[new]
// @Tags client
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param Authorization header string true "Bearer +token"
// @param TableData body TableData true "请求体"
// @Success 200 {object} ApiResponse{result=object}
// @Router /api/form/table/edit [post]
func (ds *defaultServer) UpdateTableData(ctx *gin.Context) {
	var arg TableData
	if err := ctx.Bind(&arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	if err := validatorInstance().Struct(arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	formdataSrv := formdata.NewFormdataService(ds.db)
	formSrv := form.NewFormService(ds.db)
	isNotFound, field := formSrv.GetFieldDetailByKey(arg.FieldName)
	if isNotFound {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}

	//查询本业务数据表表名
	isNotFound, fieldTable := formdataSrv.GetFieldTableByFieldId(uint32(field.ID))
	if isNotFound {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}
	isNotFound, datasourceTable := formdataSrv.GetDatasourceTableById(fieldTable.DatasourceTableID)
	if isNotFound {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}

	//处理要更新的字段
	var newVaules string
	for colId, valueInter := range arg.Data {
		var value string
		isNotFound, datasourceColumn := formdataSrv.GetDatasourceColumnById(colId)
		if isNotFound {
			ds.ResponseError(ctx, NOT_FOUND)
			return
		}

		switch datasourceColumn.DataType {
		case "int":
			if valueInter != nil {
				if datasourceColumn.FieldType == "Member" {
					if data, ok := valueInter.([]interface{}); !ok {
						ds.ResponseError(ctx, FAIL, "成员控件数据格式错误")
						return
					} else {
						var memberId = make([]float64, 0)
						for _, v := range data {
							if id, ok := v.(float64); !ok {
								ds.ResponseError(ctx, FAIL, "成员控件数据格式错误")
								return
							} else {
								memberId = append(memberId, id)
							}
						}
						value = fmt.Sprintf(datasourceColumn.ColumnName+" = %d", int(memberId[0]))
					}
				} else {
					if v, ok := valueInter.(float64); !ok {
						ds.ResponseError(ctx, FAIL, "数据类型错误")
						return
					} else {
						value = fmt.Sprintf(datasourceColumn.ColumnName+" = %d", int(v))
					}
				}
			}
		case "tinyint":
			if valueInter != nil {
				if datasourceColumn.FieldType == "SingleChoice" {
					if item, ok := valueInter.(string); !ok {
						ds.ResponseError(ctx, FAIL, "单选控件数据格式错误")
						return
					} else {
						isNotFound, datasourceMetadatas := formdataSrv.GetDatasourceMetadataById(datasourceColumn.ID)
						if isNotFound {
							ds.ResponseError(ctx, NOT_FOUND)
							return
						}
						var isExist bool
						for _, v := range datasourceMetadatas {
							if v.Value == item {
								code, _ := strconv.Atoi(v.Key)
								value = fmt.Sprintf(datasourceColumn.ColumnName+" = %d", code)

								isExist = true
								break
							}
						}

						if !isExist {
							ds.ResponseError(ctx, FAIL, "单选选项不存在")
							return
						}
					}
				} else {
					if v, ok := valueInter.(float64); !ok {
						ds.ResponseError(ctx, FAIL, "数据类型错误")
						return
					} else {
						value = fmt.Sprintf(datasourceColumn.ColumnName+" = %d", int(v))
					}
				}
			}
		case "char":
			if valueInter != nil {
				if datasourceColumn.FieldType == "Records" {
					if data, ok := valueInter.([]interface{}); !ok {
						ds.ResponseError(ctx, FAIL, "关联记录控件数据格式错误")
						return
					} else {
						var relationIds string
						for _, v := range data {
							if id, ok := v.(float64); !ok {
								ds.ResponseError(ctx, FAIL, "关联记录控件数据格式错误")
								return
							} else {
								idStr := strconv.Itoa(int(id))
								relationIds += idStr + ","
							}
						}
						value = fmt.Sprintf(datasourceColumn.ColumnName+" = '%s'", strings.TrimRight(relationIds, ","))
					}
				} else if datasourceColumn.FieldType == "MultipleChoice" {
					if data, ok := valueInter.([]interface{}); !ok {
						ds.ResponseError(ctx, FAIL, "多选控件数据格式错误")
						return
					} else {
						var keys string
						isNotFound, datasourceMetadatas := formdataSrv.GetDatasourceMetadataById(datasourceColumn.ID)
						if isNotFound {
							ds.ResponseError(ctx, NOT_FOUND)
							return
						}
						for _, v := range data {
							if item, ok := v.(string); !ok {
								ds.ResponseError(ctx, FAIL, "多选控件选项值数据格式错误")
								return
							} else {
								var isExist bool
								for _, v := range datasourceMetadatas {
									if v.Value == item {
										keys += v.Key + ","

										isExist = true
										break
									}
								}
								if !isExist {
									ds.ResponseError(ctx, FAIL, "多选控件选项值不存在")
									return
								}
							}
						}
						value = fmt.Sprintf(datasourceColumn.ColumnName+" = '%s'", strings.TrimRight(keys, ","))
					}
				} else {
					if v, ok := valueInter.(string); !ok {
						ds.ResponseError(ctx, FAIL, "数据类型错误")
						return
					} else {
						value = fmt.Sprintf(datasourceColumn.ColumnName+" = '%s'", v)
					}
				}
			}
		case "varchar":
			if v, ok := valueInter.(string); !ok {
				ds.ResponseError(ctx, FAIL, "数据类型错误")
				return
			} else {
				value = fmt.Sprintf(datasourceColumn.ColumnName+" = '%s'", v)
			}
		case "datetime":
			if valueInter == nil {
				value = fmt.Sprintf(datasourceColumn.ColumnName + " = null")
			} else {
				if v, ok := valueInter.(float64); !ok {
					ds.ResponseError(ctx, FAIL, "日期控件数据格式错误")
					return
				} else {
					value = fmt.Sprintf(datasourceColumn.ColumnName+" = '%s'", time.Unix(int64(v), 0).Format("2006-01-02 15:04:05"))
				}
			}
		default:
			if v, ok := valueInter.(string); !ok {
				ds.ResponseError(ctx, FAIL, "不支持的控件数据格式")
				return
			} else {
				value = fmt.Sprintf(datasourceColumn.ColumnName+" = '%s'", v)
			}
		}

		if value != "" {
			newVaules += value + ","
		}
	}
	newVaules = strings.TrimRight(newVaules, ",")

	err := formdataSrv.UpdateTableData(ds.businessDb, datasourceTable.TableName, arg.RowId, newVaules)
	if err != nil {
		logrus.Error("更新业务数据出错：", err.Error())
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	ds.ResponseSuccess(ctx, nil)
	return
}

// FormatTableData
// @Summary 填充列表单条数据至表单格式(为修改单条表单数据功能服务)[new]
// @Tags client
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param Authorization header string true "Bearer +token"
// @param FieldArg body FieldArg true "请求体"
// @Success 200 {object} ApiResponse{result=object}
// @Router /api/form/table/format [post]
func (ds *defaultServer) FormatTableData(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*UserClaims)
	var arg FieldArg
	if err := ctx.Bind(&arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	if err := validatorInstance().Struct(arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	formdataSrv := formdata.NewFormdataService(ds.db)
	formSrv := form.NewFormService(ds.db)
	isNotFound, field := formSrv.GetFieldDetailByKey(arg.FieldName)
	if isNotFound {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}

	//查询本业务数据表表名
	isNotFound, fieldTable := formdataSrv.GetFieldTableByFieldId(uint32(field.ID))
	if isNotFound {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}
	isNotFound, datasourceTable := formdataSrv.GetDatasourceTableById(fieldTable.DatasourceTableID)
	if isNotFound {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}

	//查询本业务数据表字段
	var mainTables = make(map[string][]string, 0) //本表信息     k=本表表名    v=[本表的显示字段a，本表的显示字段b...]
	err, mainColNames, mainColIds, mainColViewNames, mainColFieldTypes := formdataSrv.GetBusinessColNameByTableName(uint32(claims.AppId), datasourceTable.TableName)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	mainTables[datasourceTable.TableName] = mainColNames

	//查询业务数据
	err, data := formdataSrv.GetBusinessTableDataById(ds.businessDb, mainTables, nil, arg.RowId)
	if err != nil {
		logrus.Error("获取业务数据出错：", err.Error())
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	var result = make([]map[string]interface{}, 0)
	for k, colValue := range data {
		if k > 0 { //包含行id值
			var item = make(map[string]interface{}, 0)
			item["colId"] = mainColIds[k-1]
			item["colName"] = mainColViewNames[k-1]
			if colValue == "null" {
				item["value"] = nil
			} else {
				item["value"] = colValue
			}
			if len(colValue) == 0 {
				item["value"] = ""
			}
			item["ctrlType"] = mainColFieldTypes[k-1]
			item["options"] = []string{}

			if mainColFieldTypes[k-1] == "MultipleChoice" || mainColFieldTypes[k-1] == "SingleChoice" {
				isNotFound, datasourceMetadatas := formdataSrv.GetDatasourceMetadataById(mainColIds[k-1])
				if isNotFound {
					ds.ResponseError(ctx, NOT_FOUND)
					return
				}
				var keysAndValues = make(map[string]string, len(datasourceMetadatas))
				var values = make([]string, 0)
				for _, v := range datasourceMetadatas {
					keysAndValues[v.Key] = v.Value
					values = append(values, v.Value)
				}
				item["options"] = values
				if v, ok := keysAndValues[colValue]; ok {
					item["value"] = v
				} else {
					item["value"] = ""
				}
			}

			result = append(result, item)
		}
	}

	ds.ResponseSuccess(ctx, result)
	return
}

// DeleteTableData
// @Summary 删除列表单条数据，仅删除列表关联的本数据表[new]
// @Tags client
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param Authorization header string true "Bearer +token"
// @param FieldArg body FieldArg true "请求体"
// @Success 200 {object} ApiResponse{result=object}
// @Router /api/form/table/delete [delete]
func (ds *defaultServer) DeleteTableData(ctx *gin.Context) {
	var arg FieldArg
	if err := ctx.Bind(&arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	if err := validatorInstance().Struct(arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	formdataSrv := formdata.NewFormdataService(ds.db)
	formSrv := form.NewFormService(ds.db)
	isNotFound, field := formSrv.GetFieldDetailByKey(arg.FieldName)
	if isNotFound {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}

	//查询本业务数据表表名
	isNotFound, fieldTable := formdataSrv.GetFieldTableByFieldId(uint32(field.ID))
	if isNotFound {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}
	isNotFound, datasourceTable := formdataSrv.GetDatasourceTableById(fieldTable.DatasourceTableID)
	if isNotFound {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}

	err := formdataSrv.DeleteTableData(ds.businessDb, datasourceTable.TableName, arg.RowId)
	if err != nil {
		logrus.Error("删除业务数据出错：", err.Error())
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	ds.ResponseSuccess(ctx, nil)
	return
}

// TableDataDetail
// @Summary 查看列表单条数据详情，暂只显示列表关联的本数据表详情[new]
// @Tags client
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param Authorization header string true "Bearer +token"
// @param FieldArg body FieldArg true "请求体"
// @Success 200 {object} ApiResponse{result=object}
// @Router /api/form/table/detail [post]
func (ds *defaultServer) TableDataDetail(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*UserClaims)
	var arg FieldArg
	if err := ctx.Bind(&arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	if err := validatorInstance().Struct(arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	formdataSrv := formdata.NewFormdataService(ds.db)
	formSrv := form.NewFormService(ds.db)
	isNotFound, field := formSrv.GetFieldDetailByKey(arg.FieldName)
	if isNotFound {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}

	//查询本业务数据表表名
	isNotFound, fieldTable := formdataSrv.GetFieldTableByFieldId(uint32(field.ID))
	if isNotFound {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}
	isNotFound, datasourceTable := formdataSrv.GetDatasourceTableById(fieldTable.DatasourceTableID)
	if isNotFound {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}

	//查询要显示的业务数据表字段
	var mainTables = make(map[string][]string, 0)     //本表信息     k=本表表名    v=[本表的字段a，本表的字段b...]
	var relationTables = make(map[string][]string, 0) //关联表信息   k=关联表表名  v=[本表关联字段a，关联表的关联字段b，关联表的显示字段c，关联表的显示字段d...]
	isNotFound, fieldTableColumnList := formdataSrv.GetFieldTableColumnByFieldId(uint32(field.ID))
	if isNotFound {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}
	var relatedFieldTableColumnList = make(formdata.FieldTableColumnList, 0)
	for _, v := range fieldTableColumnList {
		if v.DatasourceColumnRelationID != 0 {
			relatedFieldTableColumnList = append(relatedFieldTableColumnList, v)
		}
	}

	var colIds = make([]uint32, 0)
	var colFieldTypes = make([]string, 0)
	var colNames = make([]string, 0)

	//本数据表的字段
	err, mainColNames, mainColIds, mainColViewNames, mainColFieldTypes := formdataSrv.GetBusinessColNameByTableName(uint32(claims.AppId), datasourceTable.TableName)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	mainTables[datasourceTable.TableName] = mainColNames
	colNames = append(colNames, mainColViewNames...)
	colFieldTypes = append(colFieldTypes, mainColFieldTypes...)
	colIds = append(colIds, mainColIds...)

	//关联数据表的字段
	for relationId, relateColIds := range relatedFieldTableColumnList.RelationIdsAndColIds() {
		isNotFound, datasourceColumnRelation := formdataSrv.GetDatasourceColumnRelationById(relationId)
		if isNotFound {
			ds.ResponseError(ctx, NOT_FOUND)
			return
		}
		//查业务数据表表名
		isNotFound, datasourceTable := formdataSrv.GetDatasourceTableById(datasourceColumnRelation.TargetTableID)
		if isNotFound {
			ds.ResponseError(ctx, NOT_FOUND)
			return
		}

		//查业务数据表字段名
		//关联表的关联字段名
		err, relateColName, _, _ := formdataSrv.GetBusinessColNameByIds([]uint32{datasourceColumnRelation.TargetColumnID})
		if err != nil {
			ds.InternalServiceError(ctx, err.Error())
			return
		}
		//本表的关联字段名
		err, sourceColName, _, _ := formdataSrv.GetBusinessColNameByIds([]uint32{datasourceColumnRelation.SourceColumnID})
		if err != nil {
			ds.InternalServiceError(ctx, err.Error())
			return
		}
		//关联表要显示的字段名
		err, relateColNames, relateShowColNames, _ := formdataSrv.GetBusinessColNameByIds(relateColIds)
		if err != nil {
			ds.InternalServiceError(ctx, err.Error())
			return
		}
		colNames = append(colNames, relateShowColNames...)
		colIds = append(colIds, relateColIds...)

		for _, colId := range relateColIds {
			isNotFound, datasourceColumn := formdataSrv.GetDatasourceColumnById(colId)
			if isNotFound {
				ds.ResponseError(ctx, NOT_FOUND)
				return
			}
			colFieldTypes = append(colFieldTypes, datasourceColumn.FieldType)
		}

		relationTables[datasourceTable.TableName] = append([]string{sourceColName[0], relateColName[0]}, relateColNames...)
	}

	//查询业务数据
	err, data := formdataSrv.GetBusinessTableDataById(ds.businessDb, mainTables, relationTables, arg.RowId)
	if err != nil {
		logrus.Error("获取业务数据出错：", err.Error())
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	var result = make([]map[string]interface{}, 0)
	for k, colValue := range data {
		if k > 0 { //包含行id值
			var item = make(map[string]interface{}, 0)
			//item["colId"] = colIds[k-1]
			item["colName"] = colNames[k-1]
			if colValue == "null" {
				item["value"] = nil
			} else {
				item["value"] = colValue
			}
			if len(colValue) == 0 {
				item["value"] = ""
			}
			item["ctrlType"] = colFieldTypes[k-1]

			if colFieldTypes[k-1] == "MultipleChoice" || colFieldTypes[k-1] == "SingleChoice" {
				if colValue != "" {
					isNotFound, value := formdataSrv.GetDatasourceMetadataByKey(colIds[k-1], colValue)
					if isNotFound {
						item["value"] = ""
					} else {
						item["value"] = value
					}
				}
			}

			result = append(result, item)
		}
	}

	ds.ResponseSuccess(ctx, result)
	return
}

type RelatedFieldArg struct {
	FieldName string                 `json:"fieldName" validate:"required"` //关联记录控件key
	Filter    map[string]interface{} `json:"filter"`                        //查询条件(表字段id:值)
}

// GetRelatedData
// @Summary 查询关联记录的数据[new]
// @Tags client
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param Authorization header string true "Bearer +token"
// @param RelatedFieldArg body RelatedFieldArg true "请求体"
// @Success 200 {object} ApiResponse{result=object}
// @Router /api/form/relate [post]
func (ds *defaultServer) GetRelatedData(ctx *gin.Context) {
	var arg RelatedFieldArg
	if err := ctx.Bind(&arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	if err := validatorInstance().Struct(arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	formdataSrv := formdata.NewFormdataService(ds.db)
	formSrv := form.NewFormService(ds.db)

	isNotFound, field := formSrv.GetFieldDetailByKey(arg.FieldName)
	if isNotFound {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}

	var reflectIds = make([]uint32, 0) //需要查找映射关系的表字段id

	//查询关联控件信息
	isNotFound, fieldRecords := formdataSrv.GetFieldRecordsByFieldId(uint32(field.ID))
	if isNotFound {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}
	//关联表的显示字段
	var colIds []uint32
	for _, colIdStr := range strings.Split(fieldRecords.Columns, ",") {
		colId, _ := strconv.Atoi(colIdStr)
		colIds = append(colIds, uint32(colId))
	}
	err, colNames, _, mainReflectIds := formdataSrv.GetBusinessColNameByIds(colIds)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	reflectIds = append(reflectIds, mainReflectIds...)

	//查询关联的业务数据表名
	isNotFound, datasourceColumnRelation := formdataSrv.GetDatasourceColumnRelationById(fieldRecords.DatasourceColumnRelationID)
	if isNotFound {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}
	isNotFound, datasourceTable := formdataSrv.GetDatasourceTableById(datasourceColumnRelation.TargetTableID)
	if isNotFound {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}

	var data [][]string
	var total int
	if len(arg.Filter) == 0 {
		//卡片选择 不带条件搜索
		err, total, data, _ = formdataSrv.GetBusinessTableData(ds.businessDb, map[string][]string{datasourceTable.TableName: colNames}, nil, nil, 1, utils.MAX_LIMIT)
		if err != nil {
			logrus.Error("获取业务数据出错：", err.Error())
			ds.InternalServiceError(ctx, err.Error())
			return
		}
	} else {
		//自动回填 带条件搜索
		var filterCondition = make(map[string][]string)
		filterCondition[datasourceTable.TableName] = make([]string, 0)
		for colIdStr, value := range arg.Filter {
			colId, _ := strconv.Atoi(colIdStr)
			isNotFound, datasourceColumn := formdataSrv.GetDatasourceColumnById(uint32(colId))
			if isNotFound {
				ds.ResponseError(ctx, NOT_FOUND)
				return
			}
			condition := datasourceColumn.ColumnName + " = '" + value.(string) + "'"
			filterCondition[datasourceTable.TableName] = append(filterCondition[datasourceTable.TableName], condition)
		}

		err, total, data, _ = formdataSrv.GetBusinessTableData(ds.businessDb, map[string][]string{datasourceTable.TableName: colNames}, nil, filterCondition, 1, utils.MAX_LIMIT)
		if err != nil {
			logrus.Error("获取业务数据出错：", err.Error())
			ds.InternalServiceError(ctx, err.Error())
			return
		}
	}

	colIds = append([]uint32{0}, colIds...)
	//升序排序
	sort.Slice(colIds, func(i, j int) bool {
		return colIds[i] < colIds[j]
	})
	list := make([]map[string]interface{}, 0)
	for _, rowData := range data {
		var col = make(map[string]interface{}, 0)
		for k, v := range rowData {
			if k == 0 {
				col["_id"] = v
			} else {
				if v == "null" {
					col[strconv.FormatUint(uint64(colIds[k]), 10)] = nil
				} else {
					col[strconv.FormatUint(uint64(colIds[k]), 10)] = v
				}
				if len(v) == 0 {
					col[strconv.FormatUint(uint64(colIds[k]), 10)] = ""
					continue
				}
				for _, id := range reflectIds {
					if id == colIds[k] {
						isNotFound, value := formdataSrv.GetDatasourceMetadataByKey(id, v)
						if isNotFound {
							col[strconv.FormatUint(uint64(colIds[k]), 10)] = ""
						} else {
							col[strconv.FormatUint(uint64(colIds[k]), 10)] = value
						}
						break
					}
				}
			}
		}
		list = append(list, col)
	}

	ds.ResponseSuccess(ctx, map[string]interface{}{
		"list":  list,
		"total": total,
	})
}

type Option struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// GetTableOptionData
// @Summary 查询列表控件上筛选字段的选项集[new]
// @Tags client
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param Authorization header string true "Bearer +token"
// @param columnId query string true "表字段id"
// @Success 200 {object} ApiResponse{result=[]Option}
// @Router /api/form/table/option/list [get]
func (ds *defaultServer) GetTableOptionData(ctx *gin.Context) {
	columnIdStr := ctx.Query("columnId")
	if columnIdStr == "" {
		ds.InvalidParametersError(ctx)
		return
	}

	columnId := utils.NewStr(columnIdStr).Uint32()
	formdataSrv := formdata.NewFormdataService(ds.db)

	isNotFound, datasourceMetadatas := formdataSrv.GetDatasourceMetadataById(columnId)
	if isNotFound {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}

	var data = make([]Option, 0)
	for _, v := range datasourceMetadatas {
		data = append(data, Option{
			Key:   v.Key,
			Value: v.Value,
		})
	}

	ds.ResponseSuccess(ctx, data)
	return
}

type CallBpmnArg struct {
	FieldName  string                 `json:"fieldName"  validate:"required"` //按钮控件key/列表控件key
	IsInTable  bool                   `json:"isInTable" validate:"required"`  //是否为列表内的按钮
	ButtonName string                 `json:"buttonName"`                     //按钮名称
	Data       map[string]interface{} `json:"data"`                           //流程需要的数据     列表内为字段id:值；表单内为控件key:值
}

// CallBpmn
// @Summary 触发流程[new]
// @Tags client
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param Authorization header string true "Bearer +token"
// @param CallBpmnArg body CallBpmnArg true "请求体"
// @Success 200 {object} ApiResponse{result=object}
// @Router /api/flow/run [post]
func (ds *defaultServer) CallBpmn(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*UserClaims)
	var arg CallBpmnArg
	if err := ctx.Bind(&arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	formdataSrv := formdata.NewFormdataService(ds.db)
	formSrv := form.NewFormService(ds.db)
	flowSrv := flow.NewFlowService(ds.db)
	isNotFound, field := formSrv.GetFieldDetailByKey(arg.FieldName)
	if isNotFound {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}

	var event int8
	var flowId uint32
	var data = make(map[string]interface{}, 0)
	if arg.IsInTable {
		//列表内的按钮，需要传递行数据
		isNotFound, fieldTableButton := formdataSrv.GetFieldTableButtonByFieldId(field.ID, arg.ButtonName)
		if isNotFound {
			ds.ResponseError(ctx, NOT_FOUND)
			return
		}

		for colIdStr, value := range arg.Data { //字段id:值
			if colIdStr == "index" || colIdStr == "_id" {
				continue
			}
			colId, _ := strconv.ParseUint(colIdStr, 10, 32)
			isNotFound, datasourceColumn := formdataSrv.GetDatasourceColumnById(uint32(colId))
			if isNotFound {
				ds.ResponseError(ctx, NOT_FOUND)
				return
			}

			if datasourceColumn.FieldType == "DateTime" {
				if value == nil {
					data[colIdStr] = ""
				} else {
					if v, ok := value.(string); !ok {
						ds.ResponseError(ctx, FAIL, "日期类型数据格式错误")
						return
					} else {
						data[colIdStr] = strings.Replace(strings.TrimRight(v, "+08:00"), "T", " ", -1) //key为数据表字段id
					}
				}
			} else if datasourceColumn.FieldType == "SingleChoice" {
				if value != nil {
					if v, ok := value.(string); !ok {
						ds.ResponseError(ctx, FAIL, "单选类型数据格式错误")
						return
					} else {
						if v != "" {
							isNotFound, itemKey := formdataSrv.GetDatasourceMetadataByValue(uint32(colId), v)
							if isNotFound {
								ds.ResponseError(ctx, FAIL, "单选选项不存在")
								return
							}
							data[colIdStr] = itemKey //key为数据表字段id
						}
					}
				}
			} else {
				data[colIdStr] = value //key为数据表字段id
			}
		}

		flowId = fieldTableButton.FlowID
		event = fieldTableButton.Event
	} else {
		//表单内的按钮，需要传递表单数据

		isNotFound, fieldButton := formdataSrv.GetFieldButtonByFieldId(uint32(field.ID))
		if isNotFound {
			ds.ResponseError(ctx, NOT_FOUND)
			return
		}

		flowId = fieldButton.FlowID
		event = fieldButton.Event

		for fieldName, value := range arg.Data { //控件key:值
			var hasValue bool
			isNotFound, field := formSrv.GetFieldDetailByKey(fieldName)
			if isNotFound {
				ds.ResponseError(ctx, NOT_FOUND)
				return
			}

			switch field.Type {
			case "File", "Autograph":
				if value != nil {
					if itemData, ok := value.([]interface{}); !ok {
						ds.ResponseError(ctx, FAIL, "附件/签名控件数据格式错误")
						return
					} else {
						var paths string
						for _, v := range itemData {
							if path, ok := v.(string); !ok {
								ds.ResponseError(ctx, FAIL, "附件/签名控件数据格式错误")
								return
							} else {
								if len(path) == 0 {
									ds.ResponseError(ctx, FAIL, "附件/签名控件数据格式错误")
									return
								}
								paths += path + ","
							}
						}
						data[fmt.Sprintf("%d", field.ID)] = strings.TrimRight(paths, ",") //key为控件id
					}
					hasValue = true
				}
			case "Member":
				if value != nil {
					if itemDatas, ok := value.([]interface{}); !ok {
						ds.ResponseError(ctx, FAIL, "成员控件数据格式错误")
						return
					} else {
						var memberId = make([]float64, 0)
						for _, v := range itemDatas {
							if id, ok := v.(float64); !ok {
								ds.ResponseError(ctx, FAIL, "成员控件数据格式错误")
								return
							} else {
								memberId = append(memberId, id)
							}
						}
						data[fmt.Sprintf("%d", field.ID)] = memberId[0] //key为控件id   只取了第一个
					}
					hasValue = true
				}
			case "CascadeControl":
				if value != nil {
					isNotFound, fieldLinkage := formdataSrv.GetFieldLinkageByFieldId(uint32(field.ID))
					if isNotFound {
						ds.ResponseError(ctx, NOT_FOUND)
						return
					}

					var grades = make([][]string, 0)
					for _, grade := range strings.Split(fieldLinkage.Content, "@") { //表字段关联关系id#表字段id#显示文案@表字段关联关系id#表字段id#显示文案
						grades = append(grades, strings.Split(grade, "#"))
					}

					var itemValues = make([]string, 0)
					if subValues, ok := value.([]interface{}); !ok {
						ds.ResponseError(ctx, FAIL, "级联控件数据格式错误")
						return
					} else {
						for _, v := range subValues {
							if data, ok := v.(string); !ok {
								ds.ResponseError(ctx, FAIL, "级联控件数据格式错误")
								return
							} else {
								if len(data) == 0 {
									ds.ResponseError(ctx, FAIL, "级联控件数据格式错误")
									return
								}
								itemValues = append(itemValues, data)
							}
						}
					}
					hasValue = true

					for k, grade := range grades {
						relationId, _ := strconv.ParseUint(grade[0], 10, 32)
						isNotFound, datasourceColumnRelation := formdataSrv.GetDatasourceColumnRelationById(uint32(relationId))
						if isNotFound {
							ds.ResponseError(ctx, NOT_FOUND)
							return
						}

						isNotFound, sourceDatasourceColumn := formdataSrv.GetDatasourceColumnById(datasourceColumnRelation.SourceColumnID)
						if isNotFound {
							ds.ResponseError(ctx, NOT_FOUND)
							return
						}

						isNotFound, targetDatasourceTable := formdataSrv.GetDatasourceTableById(datasourceColumnRelation.TargetTableID)
						if isNotFound {
							ds.ResponseError(ctx, NOT_FOUND)
							return
						}

						isNotFound, sourceDatasourceTable := formdataSrv.GetDatasourceTableById(datasourceColumnRelation.SourceTableID)
						if isNotFound {
							ds.ResponseError(ctx, NOT_FOUND)
							return
						}

						if relationId == 0 {
							//本表单的字段

							colId, _ := strconv.ParseUint(grade[1], 10, 32)
							isNotFound, datasourceColumn := formdataSrv.GetDatasourceColumnById(uint32(colId))
							if isNotFound {
								ds.ResponseError(ctx, NOT_FOUND)
								return
							}

							err, datas := formdataSrv.GetColValues(ds.businessDb, sourceDatasourceTable.TableName, datasourceColumn.ColumnName, fmt.Sprintf("id = %s", itemValues[k]))
							if err != nil {
								logrus.Error("获取业务数据出错：", err.Error())
								ds.InternalServiceError(ctx, err.Error())
								return
							}

							//和表单内其他控件传递的key不一样
							data[grade[1]] = datas[0][1] //key为数据表字段id
						} else {
							//关联表的字段

							if sourceDatasourceColumn.FieldType == "Records" {
								//本表的字段为关联字段，只需记录关联数据的id

								//和表单内其他控件传递的key不一样
								data[grade[1]] = itemValues[k] //key为数据表字段id
							} else {
								//本表的字段非关联类字段，则存储对应的值

								colId, _ := strconv.ParseUint(grade[1], 10, 64)
								isNotFound, targetDatasourceColumn := formdataSrv.GetDatasourceColumnById(uint32(colId))
								if isNotFound {
									ds.ResponseError(ctx, NOT_FOUND)
									return
								}

								err, datas := formdataSrv.GetColValues(ds.businessDb, targetDatasourceTable.TableName, targetDatasourceColumn.ColumnName, fmt.Sprintf("id = %s", itemValues[k]))
								if err != nil {
									logrus.Error("获取业务数据出错：", err.Error())
									ds.InternalServiceError(ctx, err.Error())
									return
								}

								if sourceDatasourceColumn.DataType == "datetime" {
									//和表单内其他控件传递的key不一样
									data[grade[1]] = strings.Replace(strings.TrimRight(datas[0][1], "+08:00"), "T", " ", -1) //key为数据表字段id
								} else {
									//和表单内其他控件传递的key不一样
									data[grade[1]] = datas[0][1] //key为数据表字段id
								}
							}
						}
					}
				}
			case "Records":
				if value != nil {
					if itemDatas, ok := value.([]interface{}); !ok {
						ds.ResponseError(ctx, FAIL, "关联记录控件数据格式错误")
						return
					} else {
						var relationIds string
						for _, v := range itemDatas {
							if id, ok := v.(float64); !ok {
								ds.ResponseError(ctx, FAIL, "关联记录控件数据格式错误")
								return
							} else {
								idStr := strconv.Itoa(int(id))
								relationIds += idStr + ","
							}
						}
						data[fmt.Sprintf("%d", field.ID)] = strings.TrimRight(relationIds, ",") //key为控件id
					}
					hasValue = true
				}
			case "SingleChoice":
				if v, ok := value.(string); !ok {
					ds.ResponseError(ctx, FAIL, "单选控件数据格式错误")
					return
				} else {
					if len(v) != 0 {
						isNotFound, itemKey := formdataSrv.GetDatasourceMetadataByValue(uint32(field.DatasourceColumnID), value.(string))
						if isNotFound {
							ds.ResponseError(ctx, FAIL, "单选选项不存在")
							return
						}
						data[fmt.Sprintf("%d", field.ID)] = itemKey //key为控件id
						hasValue = true
					}
				}
			case "DateTime":
				if v, ok := value.(float64); !ok {
					ds.ResponseError(ctx, FAIL, "日期控件数据格式错误")
					return
				} else {
					if v != 0 {
						hasValue = true
					}
					data[fmt.Sprintf("%d", field.ID)] = time.Unix(int64(v), 10).Format("2006-01-02 15:04:05")
				}
			case "DatetimeRange":
				if value != nil {
					if itemDatas, ok := value.([]interface{}); !ok {
						ds.ResponseError(ctx, FAIL, "日期范围控件数据格式错误")
						return
					} else {
						if _, ok := itemDatas[0].(float64); !ok {
							ds.ResponseError(ctx, FAIL, "日期范围控件数据格式错误")
							return
						} else if _, ok1 := itemDatas[1].(float64); !ok1 {
							ds.ResponseError(ctx, FAIL, "日期范围控件数据格式错误")
							return
						}
						startTime := time.Unix(int64(itemDatas[0].(float64)), 10).Format("2006-01-02 15:04:05")
						endTime := time.Unix(int64(itemDatas[1].(float64)), 10).Format("2006-01-02 15:04:05")
						data[fmt.Sprintf("%d", field.ID)] = startTime + "," + endTime //key为控件id
						hasValue = true
					}
				}
			default:
				data[fmt.Sprintf("%d", field.ID)] = value //key为控件id
				if value != nil {
					hasValue = true
				}
			}
			if *field.IsNecessary && !hasValue {
				ds.ResponseError(ctx, FAIL, "必填项没有值")
				return
			}
		}
	}

	if event != 1 {
		ds.ResponseError(ctx, FAIL, "不是触发流程的按钮")
		return
	}

	isNotFound, flowDetail := flowSrv.GetFlowDetail(uint64(flowId))
	if isNotFound {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}

	bpmnSrv := bpmn.NewBpmnService(ds.conf.Engine.Api)
	err := bpmnSrv.RunProcess(flowDetail.Key, data, claims.UserId)
	if err != nil {
		logrus.Error("触发流程出错：", err.Error())
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	ds.ResponseSuccess(ctx, "触发流程成功")
	return
}

type ExportArg struct {
	FieldName string                 `json:"fieldName" validate:"required"` //列表控件key
	Filter    map[uint32]interface{} `json:"filter"`                        //筛选字段及值
}

// ExportTableData
// @Summary 导出列表数据[new]
// @Tags client
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param Authorization header string true "Bearer +token"
// @param ExportArg body ExportArg true "请求体"
// @Success 200 {object} ApiResponse{result=object}
// @Router /api/form/table/download [post]
func (ds *defaultServer) ExportTableData(ctx *gin.Context) {
	var arg ExportArg
	if err := ctx.Bind(&arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	if err := validatorInstance().Struct(arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	formdataSrv := formdata.NewFormdataService(ds.db)
	formSrv := form.NewFormService(ds.db)
	//获取控件id
	isNotFound, field := formSrv.GetFieldDetailByKey(arg.FieldName)
	if isNotFound {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}

	//查询本业务数据表表名
	isNotFound, fieldTable := formdataSrv.GetFieldTableByFieldId(uint32(field.ID))
	if isNotFound {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}

	if fieldTable.IsExport == false {
		ds.ResponseError(ctx, 2004, "不允许导出列表数据")
		return
	}

	isNotFound, datasourceTable := formdataSrv.GetDatasourceTableById(fieldTable.DatasourceTableID)
	if isNotFound {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}

	//查询要显示的业务数据表字段
	var mainTables = make(map[string][]string, 0)      //本表信息     k=本表表名    v=[本表的显示字段a，本表的显示字段b...]
	var relationTables = make(map[string][]string, 0)  //关联表信息   k=关联表表名  v=[本表关联字段a，关联表的关联字段b，关联表的显示字段ids，关联表的显示字段c，关联表的显示字段d...]
	var filterCondition = make(map[string][]string, 0) //过滤条件     k=表名        v=["字段a = 值1","字段b like '%值2%'",...]
	isNotFound, fieldTableColumnList := formdataSrv.GetFieldTableColumnByFieldId(uint32(field.ID))
	if isNotFound {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}
	var relatedFieldTableColumnList = make(formdata.FieldTableColumnList, 0)
	var mainFieldTableColumnList = make(formdata.FieldTableColumnList, 0)
	for _, v := range fieldTableColumnList {
		if v.DatasourceColumnRelationID == 0 {
			//本数据表的字段
			mainFieldTableColumnList = append(mainFieldTableColumnList, v)
		} else {
			//关联数据表的字段
			relatedFieldTableColumnList = append(relatedFieldTableColumnList, v)
		}
	}

	var reflectIds = make([]uint32, 0) //需要查找映射关系的表字段id

	//本数据表的字段
	mainColIds := mainFieldTableColumnList.DatasourceColIds()
	err, _, mainColNames, mainReflectIds := formdataSrv.GetBusinessColNameByIds(mainColIds)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	mainTables[datasourceTable.TableName] = mainColNames
	mainColIds = append([]uint32{0}, mainColIds...)
	reflectIds = append(reflectIds, mainReflectIds...)

	//关联数据表的字段
	for relationId, colIds := range relatedFieldTableColumnList.RelationIdsAndColIds() {
		isNotFound, datasourceColumnRelation := formdataSrv.GetDatasourceColumnRelationById(relationId)
		if isNotFound {
			ds.ResponseError(ctx, NOT_FOUND)
			return
		}
		//查业务数据表表名
		isNotFound, datasourceTable := formdataSrv.GetDatasourceTableById(datasourceColumnRelation.TargetTableID)
		if isNotFound {
			ds.ResponseError(ctx, NOT_FOUND)
			return
		}

		//查业务数据表字段名
		//关联表的关联字段名
		err, relateColName, _, _ := formdataSrv.GetBusinessColNameByIds([]uint32{datasourceColumnRelation.TargetColumnID})
		if err != nil {
			ds.InternalServiceError(ctx, err.Error())
			return
		}
		//本表的关联字段名
		err, sourceColName, _, _ := formdataSrv.GetBusinessColNameByIds([]uint32{datasourceColumnRelation.SourceColumnID})
		if err != nil {
			ds.InternalServiceError(ctx, err.Error())
			return
		}
		//关联表要显示的字段名
		err, _, relateShowColNames, relateReflectIds := formdataSrv.GetBusinessColNameByIds(colIds)
		if err != nil {
			ds.InternalServiceError(ctx, err.Error())
			return
		}
		reflectIds = append(reflectIds, relateReflectIds...)

		var relateColIds string
		for _, colId := range colIds {
			relateColIds += strconv.FormatUint(uint64(colId), 10) + ","
		}
		relateColIds = strings.TrimRight(relateColIds, ",")

		relationTables[datasourceTable.TableName] = append([]string{sourceColName[0], relateColName[0], relateColIds}, relateShowColNames...)
	}

	//数据过滤信息
	if fieldTable.IsFilter {
		//需要数据过滤
		isNotFound, fieldTableFilterList := formdataSrv.GetFieldTableFilterByFieldId(uint32(field.ID))
		if isNotFound {
			ds.ResponseError(ctx, NOT_FOUND)
			return
		}

		for _, v := range fieldTableFilterList {
			err, whereSentences := formdataSrv.GetWhereSentence(&v)
			if err != nil {
				ds.ResponseError(ctx, NOT_FOUND)
				return
			}
			if v.DatasourceColumnRelationID == 0 {
				if _, ok := filterCondition[datasourceTable.TableName]; !ok {
					filterCondition[datasourceTable.TableName] = whereSentences
				} else {
					filterCondition[datasourceTable.TableName] = append(filterCondition[datasourceTable.TableName], whereSentences...)
				}
			} else {
				isNotFound, datasourceColumnRelation := formdataSrv.GetDatasourceColumnRelationById(v.DatasourceColumnRelationID)
				if isNotFound {
					ds.ResponseError(ctx, NOT_FOUND)
					return
				}
				//查业务数据表表名
				isNotFound, datasourceTable := formdataSrv.GetDatasourceTableById(datasourceColumnRelation.TargetTableID)
				if isNotFound {
					ds.ResponseError(ctx, NOT_FOUND)
					return
				}
				if _, ok := filterCondition[datasourceTable.TableName]; !ok {
					filterCondition[datasourceTable.TableName] = whereSentences
				} else {
					filterCondition[datasourceTable.TableName] = append(filterCondition[datasourceTable.TableName], whereSentences...)
				}
			}
		}
	}

	//处理筛选字段
	for colId, value := range arg.Filter {
		isNotFound, datasourceColumn := formdataSrv.GetDatasourceColumnById(colId)
		if isNotFound {
			ds.ResponseError(ctx, NOT_FOUND)
			return
		}

		//todo: 要区分不同的value类型
		var whereSentence string
		if datasourceColumn.FieldType == "SingleChoice" {
			//单选控件传入选项id
			if v, ok := value.(string); !ok {
				ds.ResponseError(ctx, FAIL, "单选控件数据格式错误")
				return
			} else {
				whereSentence = datasourceColumn.ColumnName + fmt.Sprintf(" = %s", v)
			}
		} else if datasourceColumn.FieldType == "Number" {
			//数值控件传入数字
			if v, ok := value.(float64); !ok {
				ds.ResponseError(ctx, FAIL, "数值控件数据格式错误")
				return
			} else {
				whereSentence = datasourceColumn.ColumnName + fmt.Sprintf(" = %d", int(v))
			}
		} else if datasourceColumn.FieldType == "DateTime" {
			//日期控件传入时间戳
			if v, ok := value.(float64); !ok {
				ds.ResponseError(ctx, FAIL, "日期控件数据格式错误")
				return
			} else {
				whereSentence = datasourceColumn.ColumnName + fmt.Sprintf(" = '%s'", time.Unix(int64(v), 10).Format("2006-01-02 15:04:05"))
			}
		} else if datasourceColumn.FieldType == "Member" {
			//成员控件传入成员id，默认只处理第一个   [12]
			if v, ok := value.([]interface{}); !ok {
				ds.ResponseError(ctx, FAIL, "成员控件数据格式错误")
				return
			} else {
				if vv, ok := v[0].(float64); !ok {
					ds.ResponseError(ctx, FAIL, "成员控件数据格式错误")
					return
				} else {
					whereSentence = datasourceColumn.ColumnName + fmt.Sprintf(" = %d", int(vv))
				}
			}
		} else {
			if v, ok := value.(string); !ok {
				ds.ResponseError(ctx, FAIL, "控件数据格式错误")
				return
			} else {
				whereSentence = datasourceColumn.ColumnName + fmt.Sprintf(" = '%s'", v)
			}
		}

		if _, ok := filterCondition[datasourceColumn.TableName]; !ok {
			filterCondition[datasourceColumn.TableName] = []string{whereSentence}
		} else {
			filterCondition[datasourceColumn.TableName] = append(filterCondition[datasourceColumn.TableName], whereSentence)
		}
	}

	//查询业务数据
	err, _, data, relateColIds := formdataSrv.GetBusinessTableData(ds.businessDb, mainTables, relationTables, filterCondition, 1, utils.MAX_LIMIT)
	if err != nil {
		logrus.Error("获取业务数据出错：", err.Error())
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	//todo:需注意字段显示名与值的对应关系
	//生成导出数据表头
	var relatedTableNames = make([]string, 0)
	for _, colNames := range relationTables {
		relatedTableNames = append(relatedTableNames, colNames[2:]...)
	}
	top := append(mainTables[datasourceTable.TableName], relatedTableNames...)

	allColIds := append(mainColIds, relateColIds...)
	list := make([]map[string]interface{}, 0)
	for _, rowData := range data {
		var col = make(map[string]interface{}, 0)
		for k, v := range rowData {
			if k == 0 {
				col["_id"] = v
			} else {
				col[strconv.FormatUint(uint64(allColIds[k]), 10)] = v
				if len(v) == 0 {
					col[strconv.FormatUint(uint64(allColIds[k]), 10)] = ""
					continue
				}
				for _, id := range reflectIds {
					if id == allColIds[k] {
						isNotFound, value := formdataSrv.GetDatasourceMetadataByKey(id, v)
						if isNotFound {
							col[strconv.FormatUint(uint64(allColIds[k]), 10)] = ""
						} else {
							col[strconv.FormatUint(uint64(allColIds[k]), 10)] = value
						}
						break
					}
				}
			}
		}
		list = append(list, col)
	}

	//创建excel文件
	err, src, contentType, header := utils.CreateExcel("table_data", top, list)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	ctx.DataFromReader(http.StatusOK, int64(len(src)), contentType, bytes.NewReader(src), header)
	return
}

type CascadeArg struct {
	FieldName string                 `json:"fieldName" validate:"required"` //级联控件key
	Data      map[string]interface{} `json:"data"`                          //上一级信息(字段id:值)
}

// GetCascadeData
// @Summary 查询级联数据[new]
// @Tags client
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param Authorization header string true "Bearer +token"
// @param CascadeArg body CascadeArg true "请求体"
// @Success 200 {object} ApiResponse{result=object}
// @Router /api/form/cascade [post]
func (ds *defaultServer) GetCascadeData(ctx *gin.Context) {
	var arg CascadeArg
	if err := ctx.Bind(&arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	if err := validatorInstance().Struct(arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	formdataSrv := formdata.NewFormdataService(ds.db)
	formSrv := form.NewFormService(ds.db)

	isNotFound, field := formSrv.GetFieldDetailByKey(arg.FieldName)
	if isNotFound {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}
	isNotFound, fieldLinkage := formdataSrv.GetFieldLinkageByFieldId(uint32(field.ID))
	if isNotFound {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}

	var grades = make([][]string, 0)
	for _, grade := range strings.Split(fieldLinkage.Content, "@") { //表字段关联关系id#表字段id#显示文案@表字段关联关系id#表字段id#显示文案
		grades = append(grades, strings.Split(grade, "#"))
	}

	var tableName, colName string
	var reflectId uint32
	var whereSentence string
	if len(arg.Data) == 0 {
		//查询第一级
		relationId, _ := strconv.ParseUint(grades[0][0], 10, 32)
		if relationId == 0 {
			//本表单的数据

			isNotFound, formItem := formSrv.GetFormDetail(uint64(fieldLinkage.FormID))
			if isNotFound || *formItem.IsDelete == true {
				ds.ResponseError(ctx, NOT_FOUND)
				return
			}

			isNotFound, datasourceTable := formdataSrv.GetDatasourceTableById(uint32(formItem.DatasourceTableID))
			if isNotFound {
				ds.ResponseError(ctx, NOT_FOUND)
				return
			}

			tableName = datasourceTable.TableName
		} else {
			//关联表的数据

			isNotFound, datasourceColumnRelation := formdataSrv.GetDatasourceColumnRelationById(uint32(relationId))
			if isNotFound {
				ds.ResponseError(ctx, NOT_FOUND)
				return
			}

			isNotFound, datasourceTable := formdataSrv.GetDatasourceTableById(datasourceColumnRelation.TargetTableID)
			if isNotFound {
				ds.ResponseError(ctx, NOT_FOUND)
				return
			}

			tableName = datasourceTable.TableName
		}

		colId, _ := strconv.ParseUint(grades[0][1], 10, 32)
		isNotFound, datasourceColumn := formdataSrv.GetDatasourceColumnById(uint32(colId))
		if isNotFound {
			ds.ResponseError(ctx, NOT_FOUND)
			return
		}

		if datasourceColumn.FieldType == "MultipleChoice" || datasourceColumn.FieldType == "SingleChoice" {
			reflectId = uint32(colId)
		}
		colName = datasourceColumn.ColumnName
	} else {
		//查询非第一级
		//两级若为同一表则需加上两级之间的关联关系

		for k, v := range grades {
			if chosenRowId, ok := arg.Data[v[1]]; ok {
				relationId, _ := strconv.ParseUint(grades[k+1][0], 10, 32)
				if relationId == 0 {
					//要查询的一级为本表单的数据

					isNotFound, formItem := formSrv.GetFormDetail(uint64(fieldLinkage.FormID))
					if isNotFound || *formItem.IsDelete == true {
						ds.ResponseError(ctx, NOT_FOUND)
						return
					}

					isNotFound, datasourceTable := formdataSrv.GetDatasourceTableById(uint32(formItem.DatasourceTableID))
					if isNotFound {
						ds.ResponseError(ctx, NOT_FOUND)
						return
					}

					if grades[k][0] == "0" {
						//上一级也为本表单的数据

						colId, _ := strconv.ParseUint(grades[k][1], 10, 32)
						isNotFound, datasourceColumn := formdataSrv.GetDatasourceColumnById(uint32(colId))
						if isNotFound {
							ds.ResponseError(ctx, NOT_FOUND)
							return
						}
						err, data := formdataSrv.GetColValues(ds.businessDb, datasourceTable.TableName, datasourceColumn.ColumnName, fmt.Sprintf("id = %s", chosenRowId))
						if err != nil {
							logrus.Error("获取业务数据出错：", err.Error())
							ds.InternalServiceError(ctx, err.Error())
							return
						}
						if datasourceColumn.DataType == "int" || datasourceColumn.DataType == "tinyint" {
							whereSentence = datasourceColumn.ColumnName + " = " + data[0][1]
						} else {
							whereSentence = datasourceColumn.ColumnName + " = '" + data[0][1] + "'"
						}
					}

					tableName = datasourceTable.TableName
				} else {
					//要查询的一级为关联表的数据
					var datasourceColumnRelationUp = &formdata.DatasourceColumnRelation{}
					relationUpId, _ := strconv.ParseUint(grades[k][0], 10, 32)
					if relationUpId != 0 {
						//上一级为关联表数据
						isNotFound, datasourceColumnRelationUp = formdataSrv.GetDatasourceColumnRelationById(uint32(relationUpId))
						if isNotFound {
							ds.ResponseError(ctx, NOT_FOUND)
							return
						}
					} else {
						//上一级为本表单数据
						datasourceColumnRelationUp.TargetTableID = 0
					}

					isNotFound, datasourceColumnRelation := formdataSrv.GetDatasourceColumnRelationById(uint32(relationId))
					if isNotFound {
						ds.ResponseError(ctx, NOT_FOUND)
						return
					}

					isNotFound, datasourceTable := formdataSrv.GetDatasourceTableById(datasourceColumnRelation.TargetTableID)
					if isNotFound {
						ds.ResponseError(ctx, NOT_FOUND)
						return
					}

					if datasourceColumnRelationUp.TargetTableID == datasourceColumnRelation.TargetTableID {
						//上一级与要查询的一级为同一张关联表的数据
						upColId, _ := strconv.ParseUint(grades[k][1], 10, 32)
						isNotFound, datasourceColumn := formdataSrv.GetDatasourceColumnById(uint32(upColId))
						if isNotFound {
							ds.ResponseError(ctx, NOT_FOUND)
							return
						}

						err, data := formdataSrv.GetColValues(ds.businessDb, datasourceTable.TableName, datasourceColumn.ColumnName, fmt.Sprintf("id = %s", chosenRowId))
						if err != nil {
							logrus.Error("获取业务数据出错：", err.Error())
							ds.InternalServiceError(ctx, err.Error())
							return
						}
						if datasourceColumn.DataType == "int" || datasourceColumn.DataType == "tinyint" {
							whereSentence = datasourceColumn.ColumnName + " = " + data[0][1]
						} else {
							whereSentence = datasourceColumn.ColumnName + " = '" + data[0][1] + "'"
						}
					}

					tableName = datasourceTable.TableName
				}

				colId, _ := strconv.ParseUint(grades[k+1][1], 10, 32)
				isNotFound, datasourceColumn := formdataSrv.GetDatasourceColumnById(uint32(colId))
				if isNotFound {
					ds.ResponseError(ctx, NOT_FOUND)
					return
				}
				if datasourceColumn.FieldType == "MultipleChoice" || datasourceColumn.FieldType == "SingleChoice" {
					reflectId = uint32(colId)
				}
				colName = datasourceColumn.ColumnName
				break
			}
		}
	}

	//查询业务数据
	err, data := formdataSrv.GetColValues(ds.businessDb, tableName, colName, whereSentence)
	if err != nil {
		logrus.Error("获取业务数据出错：", err.Error())
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	var result = make(map[string]interface{})
	for _, item := range data {
		if reflectId != 0 {
			//需要映射的字段

			isNotFound, value := formdataSrv.GetDatasourceMetadataByKey(reflectId, item[1])
			if isNotFound {
				result[item[0]] = ""
			} else {
				result[item[0]] = value
			}
		} else {
			result[item[0]] = item[1]
			if len(item[1]) == 0 {
				result[item[0]] = ""
			}
		}
	}

	ds.ResponseSuccess(ctx, result)
	return
}

// GetRelatedDataDetail
// @Summary 查询某条关联记录的详情数据[new]
// @Tags client
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param Authorization header string true "Bearer +token"
// @param FieldArg body FieldArg true "请求体"
// @Success 200 {object} ApiResponse{result=object}
// @Router /api/form/relate/detail [post]
func (ds *defaultServer) GetRelatedDataDetail(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*UserClaims)
	var arg FieldArg
	if err := ctx.Bind(&arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	if err := validatorInstance().Struct(arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	formdataSrv := formdata.NewFormdataService(ds.db)
	formSrv := form.NewFormService(ds.db)
	isNotFound, field := formSrv.GetFieldDetailByKey(arg.FieldName)
	if isNotFound {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}

	//查询关联控件信息
	isNotFound, fieldRecords := formdataSrv.GetFieldRecordsByFieldId(uint32(field.ID))
	if isNotFound {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}

	if fieldRecords.DetailStatus == false {
		ds.ResponseError(ctx, 2004, "不允许查看详情")
		return
	}

	//查询关联的业务数据表名
	isNotFound, datasourceColumnRelation := formdataSrv.GetDatasourceColumnRelationById(fieldRecords.DatasourceColumnRelationID)
	if isNotFound {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}
	isNotFound, datasourceTable := formdataSrv.GetDatasourceTableById(datasourceColumnRelation.TargetTableID)
	if isNotFound {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}

	//关联数据表的字段
	err, colNames, colIds, viewColNames, colFieldTypes := formdataSrv.GetBusinessColNameByTableName(uint32(claims.AppId), datasourceTable.TableName)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	mainTables := map[string][]string{datasourceTable.TableName: colNames}

	//查询业务数据
	err, data := formdataSrv.GetBusinessTableDataById(ds.businessDb, mainTables, nil, arg.RowId)
	if err != nil {
		logrus.Error("获取业务数据出错：", err.Error())
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	var result = make([]map[string]interface{}, 0)
	for k, colValue := range data {
		if k > 0 { //包含行id值
			var item = make(map[string]interface{}, 0)
			item["colName"] = viewColNames[k-1]
			if colValue == "null" {
				item["value"] = nil
			} else {
				item["value"] = colValue
			}
			if len(colValue) == 0 {
				item["value"] = ""
			}
			item["ctrlType"] = colFieldTypes[k-1]

			if colFieldTypes[k-1] == "MultipleChoice" || colFieldTypes[k-1] == "SingleChoice" {
				isNotFound, value := formdataSrv.GetDatasourceMetadataByKey(colIds[k-1], colValue)
				if isNotFound {
					item["value"] = ""
				} else {
					item["value"] = value
				}
			}

			result = append(result, item)
		}
	}

	ds.ResponseSuccess(ctx, result)
	return
}

type FormdataArg struct {
	FormId uint32                 `json:"id" validate:"required"`   //本表单id
	Data   map[string]interface{} `json:"data" validate:"required"` //要提交的数据 控件key:值
}

// SubmitFormData
// @Summary 提交表单数据[new]
// @Tags client
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param Authorization header string true "Bearer +token"
// @param FormdataArg body FormdataArg true "请求体"
// @Success 200 {object} ApiResponse{result=object}
// @Router /api/form/submit [post]
func (ds *defaultServer) SubmitFormData(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*UserClaims)
	var arg FormdataArg
	if err := ctx.Bind(&arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	if err := validatorInstance().Struct(arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	formdataSrv := formdata.NewFormdataService(ds.db)
	formSrv := form.NewFormService(ds.db)

	//查询本表单关联的数据表名
	isNotFound, formTable := formSrv.GetFormDetail(uint64(arg.FormId))
	if isNotFound || *formTable.IsDelete == true || *formTable.IsOnline == false || *formTable.Status == false {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}
	isNotFound, datasourceTable := formdataSrv.GetDatasourceTableById(uint32(formTable.DatasourceTableID))
	if isNotFound {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}

	//查询本表单所有的控件信息
	err, fields := formSrv.FieldList(uint64(arg.FormId))
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	var multiFormData = make(map[string]interface{}, 0) //多表单控件数据
	var recordsData = make(map[string]interface{}, 0)   //关联记录控件数据
	var cascadeData = make(map[string]interface{}, 0)   //级联选择控件数据
	var mainFormData = make(map[string]interface{}, 0)  //非特殊控件的控件数据
	for fieldName, value := range arg.Data {
		if strings.Contains(fieldName, "MultiForm") {
			multiFormData[fieldName] = value
		} else if strings.Contains(fieldName, "Records") {
			recordsData[fieldName] = value
		} else if strings.Contains(fieldName, "CascadeControl") {
			cascadeData[fieldName] = value
		} else {
			if strings.Contains(fieldName, "File") || strings.Contains(fieldName, "Autograph") {
				if data, ok := value.([]interface{}); !ok {
					ds.ResponseError(ctx, FAIL, "附件/签名控件数据格式错误")
					return
				} else {
					var paths string
					for _, v := range data {
						if path, ok := v.(string); !ok {
							ds.ResponseError(ctx, FAIL, "附件/签名控件数据格式错误")
							return
						} else {
							paths += path + ","
						}
					}
					mainFormData[fieldName] = strings.TrimRight(paths, ",")
				}
			} else if strings.Contains(fieldName, "DatetimeRange") {
				if itemDatas, ok := value.([]interface{}); !ok {
					ds.ResponseError(ctx, FAIL, "日期范围控件数据格式错误")
					return
				} else {
					if _, ok := itemDatas[0].(float64); !ok {
						ds.ResponseError(ctx, FAIL, "日期范围控件数据格式错误")
						return
					} else if _, ok2 := itemDatas[1].(float64); !ok2 {
						ds.ResponseError(ctx, FAIL, "日期范围控件数据格式错误")
						return
					}
					startTime := time.Unix(int64(itemDatas[0].(float64)), 0).Format("2006-01-02 15:04:05")
					endTime := time.Unix(int64(itemDatas[1].(float64)), 0).Format("2006-01-02 15:04:05")
					mainFormData[fieldName] = startTime + "," + endTime
				}
			} else {
				mainFormData[fieldName] = value
			}
		}
	}

	var colData = make(map[string][]interface{}, 0)
	colData["creator_id"] = []interface{}{"int", claims.UserId}
	for _, field := range fields {
		var isMatched, hasValue bool

		//非特殊控件的控件数据
		for fieldName, value := range mainFormData {
			if field.Key == fieldName {
				isNotFound, datasourceColumn := formdataSrv.GetDatasourceColumnById(uint32(field.DatasourceColumnID))
				if isNotFound {
					ds.ResponseError(ctx, NOT_FOUND)
					return
				}
				var model string
				switch datasourceColumn.DataType {
				case "int":
					if value != nil {
						model = "int"
						if datasourceColumn.FieldType == "Member" {
							if data, ok := value.([]interface{}); !ok {
								ds.ResponseError(ctx, FAIL, "成员控件数据格式错误")
								return
							} else {
								var memberId = make([]float64, 0)
								for _, v := range data {
									if id, ok := v.(float64); !ok {
										ds.ResponseError(ctx, FAIL, "成员控件数据格式错误")
										return
									} else {
										memberId = append(memberId, id)
									}
								}
								colData[datasourceColumn.ColumnName] = []interface{}{model, int(memberId[0])}
							}
						} else {
							if v, ok := value.(float64); !ok {
								ds.ResponseError(ctx, FAIL, "数据类型错误")
								return
							} else {
								//验证数值范围限制  临时解决办法
								min, max := utils.GetNumberRange(string(field.Content))
								if min == max && min == -1 {
									//无限制
									colData[datasourceColumn.ColumnName] = []interface{}{model, int(v)}
								} else {
									if int(v) >= min && int(v) <= max {
										colData[datasourceColumn.ColumnName] = []interface{}{model, int(v)}
									} else {
										ds.ResponseError(ctx, FAIL, "数值类型超过限制范围")
										return
									}
								}
							}
						}
						hasValue = true
					}
				case "tinyint":
					if value != nil {
						model = "int"
						if datasourceColumn.FieldType == "SingleChoice" {
							if item, ok := value.(string); !ok {
								ds.ResponseError(ctx, FAIL, "单选控件数据格式错误")
								return
							} else {
								isNotFound, datasourceMetadatas := formdataSrv.GetDatasourceMetadataById(datasourceColumn.ID)
								if isNotFound {
									ds.ResponseError(ctx, NOT_FOUND)
									return
								}
								var isExist bool
								for _, v := range datasourceMetadatas {
									if v.Value == item {
										code, _ := strconv.Atoi(v.Key)
										colData[datasourceColumn.ColumnName] = []interface{}{model, code}

										isExist = true
										break
									}
								}
								if !isExist {
									ds.ResponseError(ctx, FAIL, "单选选项不存在")
									return
								}
							}
						} else {
							if v, ok := value.(float64); !ok {
								ds.ResponseError(ctx, FAIL, "数据类型错误")
								return
							} else {
								colData[datasourceColumn.ColumnName] = []interface{}{model, int(v)}
							}
						}
						hasValue = true
					}
				case "char", "varchar":
					if datasourceColumn.FieldType == "MultipleChoice" {
						if value != nil {
							if itemDatas, ok := value.([]interface{}); !ok {
								ds.ResponseError(ctx, FAIL, "多选控件数据格式错误")
								return
							} else {
								var keys string
								isNotFound, datasourceMetadatas := formdataSrv.GetDatasourceMetadataById(datasourceColumn.ID)
								if isNotFound {
									ds.ResponseError(ctx, NOT_FOUND)
									return
								}
								for _, option := range itemDatas {
									if item, ok := option.(string); !ok {
										ds.ResponseError(ctx, FAIL, "多选控件选项值数据格式错误")
										return
									} else {
										var isExist bool
										for _, v := range datasourceMetadatas {
											if v.Value == item {
												keys += v.Key + ","

												isExist = true
												break
											}
										}
										if !isExist {
											ds.ResponseError(ctx, FAIL, "多选控件选项值不存在")
											return
										}
									}
								}
								model = "string"
								keys = strings.TrimRight(keys, ",")
								colData[datasourceColumn.ColumnName] = []interface{}{model, keys}
								if len(keys) != 0 {
									hasValue = true
								}
							}
						}
					} else {
						if v, ok := value.(string); !ok {
							ds.ResponseError(ctx, FAIL, "数据类型错误")
							return
						} else {
							if datasourceColumn.FieldType == "Input" {
								//验证输入框字数范围限制  临时解决办法
								min, max := utils.GetInputRange(string(field.Content))
								if min == max && min == -1 {
									//无限制
									model = "string"
									colData[datasourceColumn.ColumnName] = []interface{}{model, v}
									if len(v) != 0 {
										hasValue = true
									}
								} else {
									if utf8.RuneCountInString(v) >= min && utf8.RuneCountInString(v) <= max {
										model = "string"
										colData[datasourceColumn.ColumnName] = []interface{}{model, v}
										if len(v) != 0 {
											hasValue = true
										}
									} else {
										ds.ResponseError(ctx, FAIL, "输入框类型超过字数限制")
										return
									}
								}
							} else {
								model = "string"
								colData[datasourceColumn.ColumnName] = []interface{}{model, v}
								if len(v) != 0 {
									hasValue = true
								}
							}
						}
					}
				case "datetime":
					if v, ok := value.(float64); !ok {
						ds.ResponseError(ctx, FAIL, "数据类型错误")
						return
					} else {
						model = "datetime"
						colData[datasourceColumn.ColumnName] = []interface{}{model, int64(v)}
						if v != 0 {
							hasValue = true
						}
					}
				default:
					if v, ok := value.(string); !ok {
						ds.ResponseError(ctx, FAIL, "不支持的控件数据格式")
						return
					} else {
						model = "string"
						colData[datasourceColumn.ColumnName] = []interface{}{model, v}
						if len(v) != 0 {
							hasValue = true
						}
					}
				}
				isMatched = true
				break
			}
		}

		if isMatched {
			if *field.IsNecessary && !hasValue {
				ds.ResponseError(ctx, FAIL, "必填项没有值")
				return
			}
			continue
		}

		//处理关联记录控件数据
		for fieldName, value := range recordsData {
			if field.Key == fieldName {
				//查询关联记录控件信息
				isNotFound, fieldRecords := formdataSrv.GetFieldRecordsByFieldId(uint32(field.ID))
				if isNotFound {
					ds.ResponseError(ctx, NOT_FOUND)
					return
				}

				//查询关联的业务数据表名
				isNotFound, datasourceColumnRelation := formdataSrv.GetDatasourceColumnRelationById(fieldRecords.DatasourceColumnRelationID)
				if isNotFound {
					ds.ResponseError(ctx, NOT_FOUND)
					return
				}

				isNotFound, datasourceColumn := formdataSrv.GetDatasourceColumnById(datasourceColumnRelation.SourceColumnID)
				if isNotFound {
					ds.ResponseError(ctx, NOT_FOUND)
					return
				}

				var model string
				switch datasourceColumn.DataType {
				case "int":
					if value != nil {
						model = "int"
						if datasourceColumn.FieldType == "Member" {
							if data, ok := value.([]interface{}); !ok {
								ds.ResponseError(ctx, FAIL, "成员控件数据格式错误")
								return
							} else {
								var memberId = make([]float64, 0)
								for _, v := range data {
									if id, ok := v.(float64); !ok {
										ds.ResponseError(ctx, FAIL, "成员控件数据格式错误")
										return
									} else {
										memberId = append(memberId, id)
									}
								}
								colData[datasourceColumn.ColumnName] = []interface{}{model, int(memberId[0])}
							}
						} else {
							if v, ok := value.(float64); !ok {
								ds.ResponseError(ctx, FAIL, "数据类型错误")
								return
							} else {
								//验证数值范围限制  临时解决办法
								min, max := utils.GetNumberRange(string(field.Content))
								if min == max && min == -1 {
									//无限制
									colData[datasourceColumn.ColumnName] = []interface{}{model, int(v)}
								} else {
									if int(v) >= min && int(v) <= max {
										colData[datasourceColumn.ColumnName] = []interface{}{model, int(v)}
									} else {
										ds.ResponseError(ctx, FAIL, "数值类型超过限制范围")
										return
									}
								}
							}
						}
						hasValue = true
					}
				case "tinyint":
					if value != nil {
						model = "int"
						if datasourceColumn.FieldType == "SingleChoice" {
							if item, ok := value.(string); !ok {
								ds.ResponseError(ctx, FAIL, "单选控件数据格式错误")
								return
							} else {
								isNotFound, datasourceMetadatas := formdataSrv.GetDatasourceMetadataById(datasourceColumn.ID)
								if isNotFound {
									ds.ResponseError(ctx, NOT_FOUND)
									return
								}
								var isExist bool
								for _, v := range datasourceMetadatas {
									if v.Value == item {
										code, _ := strconv.Atoi(v.Key)
										colData[datasourceColumn.ColumnName] = []interface{}{model, code}

										isExist = true
										break
									}
								}
								if !isExist {
									ds.ResponseError(ctx, FAIL, "单选选项不存在")
									return
								}
							}
						} else {
							if v, ok := value.(float64); !ok {
								ds.ResponseError(ctx, FAIL, "数据类型错误")
								return
							} else {
								colData[datasourceColumn.ColumnName] = []interface{}{model, int(v)}
							}
						}
						hasValue = true
					}
				case "char":
					if value != nil {
						model = "string"
						if datasourceColumn.FieldType == "Records" {
							if data, ok := value.([]interface{}); !ok {
								ds.ResponseError(ctx, FAIL, "关联记录控件数据格式错误")
								return
							} else {
								var relationIds string
								for _, v := range data {
									if id, ok := v.(float64); !ok {
										ds.ResponseError(ctx, FAIL, "关联记录控件数据格式错误")
										return
									} else {
										idStr := strconv.Itoa(int(id))
										relationIds += idStr + ","
									}
								}
								colData[datasourceColumn.ColumnName] = []interface{}{model, strings.TrimRight(relationIds, ",")}
							}
						} else {
							if v, ok := value.(string); !ok {
								ds.ResponseError(ctx, FAIL, "数据类型错误")
								return
							} else {
								if datasourceColumn.FieldType == "Input" {
									//验证输入框字数范围限制  临时解决办法
									min, max := utils.GetInputRange(string(field.Content))
									if min == max && min == -1 {
										//无限制
										colData[datasourceColumn.ColumnName] = []interface{}{model, v}
										if len(v) != 0 {
											hasValue = true
										}
									} else {
										if utf8.RuneCountInString(v) >= min && utf8.RuneCountInString(v) <= max {
											colData[datasourceColumn.ColumnName] = []interface{}{model, v}
											if len(v) != 0 {
												hasValue = true
											}
										} else {
											ds.ResponseError(ctx, FAIL, "输入框类型超过字数限制")
											return
										}
									}
								} else {
									colData[datasourceColumn.ColumnName] = []interface{}{model, v}
									if len(v) != 0 {
										hasValue = true
									}
								}
							}
						}
					}
				case "varchar":
					if v, ok := value.(string); !ok {
						ds.ResponseError(ctx, FAIL, "数据类型错误")
						return
					} else {
						if datasourceColumn.FieldType == "Input" {
							//验证输入框字数范围限制  临时解决办法
							min, max := utils.GetInputRange(string(field.Content))
							if min == max && min == -1 {
								//无限制
								model = "string"
								colData[datasourceColumn.ColumnName] = []interface{}{model, v}
								if len(v) != 0 {
									hasValue = true
								}
							} else {
								if utf8.RuneCountInString(v) >= min && utf8.RuneCountInString(v) <= max {
									model = "string"
									colData[datasourceColumn.ColumnName] = []interface{}{model, v}
									if len(v) != 0 {
										hasValue = true
									}
								} else {
									ds.ResponseError(ctx, FAIL, "输入框类型超过字数限制")
									return
								}
							}
						} else {
							model = "string"
							colData[datasourceColumn.ColumnName] = []interface{}{model, v}
							if len(v) != 0 {
								hasValue = true
							}
						}
					}
				case "datetime":
					if v, ok := value.(float64); !ok {
						ds.ResponseError(ctx, FAIL, "数据类型错误")
						return
					} else {
						model = "datetime"
						colData[datasourceColumn.ColumnName] = []interface{}{model, int64(v)}
						if v != 0 {
							hasValue = true
						}
					}
				default:
					if v, ok := value.(string); !ok {
						ds.ResponseError(ctx, FAIL, "不支持的控件数据格式")
						return
					} else {
						model = "string"
						colData[datasourceColumn.ColumnName] = []interface{}{model, v}
						if len(v) != 0 {
							hasValue = true
						}
					}
				}
				isMatched = true
				break
			}
		}

		if isMatched {
			if *field.IsNecessary && !hasValue {
				ds.ResponseError(ctx, FAIL, "必填项没有值")
				return
			}
			continue
		}

		//处理级联选择控件数据
		for fieldName, value := range cascadeData {
			if field.Key == fieldName {
				if value == nil {
					continue
				}

				//查询级联选择控件信息
				isNotFound, fieldLinkage := formdataSrv.GetFieldLinkageByFieldId(uint32(field.ID))
				if isNotFound {
					ds.ResponseError(ctx, NOT_FOUND)
					return
				}

				var grades = make([][]string, 0)
				for _, grade := range strings.Split(fieldLinkage.Content, "@") { //表字段关联关系id#表字段id#显示文案@表字段关联关系id#表字段id#显示文案
					grades = append(grades, strings.Split(grade, "#"))
				}

				var itemValues = make([]string, 0)
				if subValues, ok := value.([]interface{}); !ok {
					ds.ResponseError(ctx, FAIL, "级联控件数据格式错误")
					return
				} else {
					for _, v := range subValues {
						if data, ok := v.(string); !ok {
							ds.ResponseError(ctx, FAIL, "级联控件数据格式错误")
							return
						} else {
							if len(data) == 0 {
								ds.ResponseError(ctx, FAIL, "级联控件数据格式错误")
								return
							}
							itemValues = append(itemValues, data)
						}
					}
				}
				hasValue = true

				for k, grade := range grades {
					relationId, _ := strconv.ParseUint(grade[0], 10, 32)
					isNotFound, datasourceColumnRelation := formdataSrv.GetDatasourceColumnRelationById(uint32(relationId))
					if isNotFound {
						ds.ResponseError(ctx, NOT_FOUND)
						return
					}

					isNotFound, sourceDatasourceColumn := formdataSrv.GetDatasourceColumnById(datasourceColumnRelation.SourceColumnID)
					if isNotFound {
						ds.ResponseError(ctx, NOT_FOUND)
						return
					}

					isNotFound, targetDatasourceTable := formdataSrv.GetDatasourceTableById(datasourceColumnRelation.TargetTableID)
					if isNotFound {
						ds.ResponseError(ctx, NOT_FOUND)
						return
					}

					isNotFound, sourceDatasourceTable := formdataSrv.GetDatasourceTableById(datasourceColumnRelation.SourceTableID)
					if isNotFound {
						ds.ResponseError(ctx, NOT_FOUND)
						return
					}

					if grade[0] == "0" {
						//本表单的字段

						colId, _ := strconv.ParseUint(grade[1], 10, 64)
						isNotFound, datasourceColumn := formdataSrv.GetDatasourceColumnById(uint32(colId))
						if isNotFound {
							ds.ResponseError(ctx, NOT_FOUND)
							return
						}

						err, data := formdataSrv.GetColValues(ds.businessDb, sourceDatasourceTable.TableName, datasourceColumn.ColumnName, fmt.Sprintf("id = %s", itemValues[k]))
						if err != nil {
							logrus.Error("获取业务数据出错：", err.Error())
							ds.InternalServiceError(ctx, err.Error())
							return
						}

						switch datasourceColumn.DataType {
						case "int", "tinyint":
							model := "int"
							value, _ := strconv.Atoi(data[0][1])
							colData[datasourceColumn.ColumnName] = []interface{}{model, value}
						case "char", "varchar":
							model := "string"
							colData[datasourceColumn.ColumnName] = []interface{}{model, data[0][1]}
						case "datetime":
							model := "string"
							colData[datasourceColumn.ColumnName] = []interface{}{model, data[0][1]}
						default:
							model := "string"
							colData[datasourceColumn.ColumnName] = []interface{}{model, data[0][1]}
						}
					} else {
						//关联表的字段

						if sourceDatasourceColumn.FieldType == "Records" {
							//本表的字段为关联字段，只需记录关联数据的id
							model := "string"
							colData[sourceDatasourceColumn.ColumnName] = []interface{}{model, itemValues[k]}
						} else {
							//本表的字段非关联类字段，则存储对应的值

							colId, _ := strconv.ParseUint(grade[1], 10, 64)
							isNotFound, targetDatasourceColumn := formdataSrv.GetDatasourceColumnById(uint32(colId))
							if isNotFound {
								ds.ResponseError(ctx, NOT_FOUND)
								return
							}

							err, data := formdataSrv.GetColValues(ds.businessDb, targetDatasourceTable.TableName, targetDatasourceColumn.ColumnName, fmt.Sprintf("id = %s", itemValues[k]))
							if err != nil {
								logrus.Error("获取业务数据出错：", err.Error())
								ds.InternalServiceError(ctx, err.Error())
								return
							}

							switch sourceDatasourceColumn.DataType {
							case "int", "tinyint":
								model := "int"
								value, _ := strconv.Atoi(data[0][1])
								colData[sourceDatasourceColumn.ColumnName] = []interface{}{model, value}
							case "char", "varchar":
								model := "string"
								colData[sourceDatasourceColumn.ColumnName] = []interface{}{model, data[0][1]}
							case "datetime":
								model := "string"
								colData[sourceDatasourceColumn.ColumnName] = []interface{}{model, data[0][1]}
							default:
								model := "string"
								colData[sourceDatasourceColumn.ColumnName] = []interface{}{model, data[0][1]}
							}
						}
					}
				}
				isMatched = true
				break
			}
		}

		if isMatched {
			if *field.IsNecessary && !hasValue {
				ds.ResponseError(ctx, FAIL, "必填项没有值")
				return
			}
			continue
		}

		//todo: 处理多表单控件数据
		//处理多表单控件数据
		//for fieldName, value := range multiFormData {
		//	a := value.([]interface{})
		//}
	}

	if len(fields) == 0 {
		ds.ResponseError(ctx, FAIL, "当前表单无控件")
		return
	}

	logrus.Info("colData = ", colData)

	//todo:操作多张表时要加入事务
	//写入业务数据表
	err = formdataSrv.InsertBusinessData(ds.businessDb, datasourceTable.TableName, colData)
	if err != nil {
		logrus.Error("新增业务数据出错：", err.Error())
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	ds.ResponseSuccess(ctx, nil)
	return
}
