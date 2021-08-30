//Package server
//Created by GoLand
//@User: lenora
//@Date: 2021/7/6
//@Time: 10:26

package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Lenora-Z/low-code/service/bpmn"
	"github.com/Lenora-Z/low-code/service/data_log"
	"github.com/Lenora-Z/low-code/service/flow"
	"github.com/Lenora-Z/low-code/service/form"
	"github.com/Lenora-Z/low-code/service/form_data"
	"github.com/Lenora-Z/low-code/service/mongodb"
	"github.com/Lenora-Z/low-code/service/service"
	"github.com/Lenora-Z/low-code/service/user"
	"github.com/Lenora-Z/low-code/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sort"
	"strings"
	"time"
)

type FormTreeList []FormItemVO

func (list FormTreeList) allFormDetail(srv form.FormService, tree form.FormTreeList, lang string) {
	for _, x := range tree {
		err, _, item := formDetail(srv, x.Id, lang)
		if err != nil {
			continue
		}
		list = append(list, *item)
		if len(x.Children) > 0 {
			list.allFormDetail(srv, x.Children, lang)
		}
	}

}

func formDetail(srv form.FormService, id uint64, lang string) (error, int16, *FormItemVO) {
	status, item := srv.GetFormDetail(id)
	if status || *item.IsDelete == true {
		return errors.New(getResponseMsgWithLang(NOT_FOUND, lang)), NOT_FOUND, nil
	}
	ret := make([]form.FieldFormat, 0)
	if err := json.Unmarshal(item.Content, &ret); err != nil {
		return err, FAIL, nil
	}
	var foot form.FieldFormat
	if err := json.Unmarshal(item.Footer, &foot); err != nil {
		return err, FAIL, nil
	}
	item.Content = []byte{}
	return nil, 0, &FormItemVO{
		Form:    item,
		Plugins: ret,
		Footer:  foot,
	}
}

type Srv struct {
	form     form.FormService
	flow     flow.FlowService
	formData form_data.FormDataService
	log      data_log.Service
}

func (ds *defaultServer) submitFormData(ctx *gin.Context, srv Srv, id, userId uint64, data []map[string]interface{},
	parent interface{}) (error, []interface{}) {
	lang := ds.getLanguage(ctx)
	status, tableItem := srv.form.GetFormDetail(id)
	if status || *tableItem.IsDelete == true {
		logrus.Error("table not found:", id)
		return errors.New(getResponseMsgWithLang(NOT_FOUND, lang)), nil
	}

	if *tableItem.Status != true {
		return errors.New("wrong table"), nil
	}

	err, fields := srv.form.FieldList(id)
	if err != nil {
		return err, nil
	}

	list := make([]interface{}, 0, cap(data))
	for i, x := range data {
		err, x = ds.validateItem(ctx, x, fields, tableItem.AppID, uint64(i), data_log.ADD, "")
		if err != nil {
			return err, nil
		}

		var instance string
		var flowId uint64

		//检索对应的流程并将提交上报到引擎中
		err, count, maps := srv.flow.GetMappingByFormId(tableItem.AppID, []uint64{id})
		flows := maps.Flows()
		logrus.Info("flows:", flows)
		if count > 0 {
			status, flowItem := srv.flow.GetFlowDetail(flows[0])
			if status {
				logrus.Error("flow not found:", flowItem)
				return errors.New(getResponseMsgWithLang(NOT_FOUND, lang)), nil
			}
			flowId = flowItem.ID

			srvSrv := service.NewOutsideService(ds.db)
			serviceGroup := strings.Split(flowItem.ServiceGroup, ",")

			var serviceMaps []map[string]string
			for _, s := range serviceGroup {
				srvId := utils.NewStr(s).Uint64()
				err, _, rely := srvSrv.GetParamRely([]uint64{srvId}, []uint64{flowItem.ID}, []uint64{id})
				if err != nil {
					return err, nil
				}
				mapItem := make(map[string]string, 0)
				for _, r := range rely {
					status, field := srv.form.GetFieldDetail(r.FieldID)
					if status {
						logrus.Error("field not found:", r.FieldID)
						return errors.New(getResponseMsgWithLang(NOT_FOUND, lang)), nil
					}

					status, param := srvSrv.GetParamDetail(r.ParamID)
					if status {
						logrus.Error("field not found:", r.FieldID)
						return errors.New(getResponseMsgWithLang(NOT_FOUND, lang)), nil
					}
					mapItem[param.Name] = field.Key
				}
				serviceMaps = append(serviceMaps, mapItem)
			}

			bpmnSrv := bpmn.NewBpmnService(ds.conf.Engine.Api)
			err, instance = bpmnSrv.ExecuteProcess(flowItem.Key, serviceMaps, x)
			if err != nil {
				return err, nil
			}
		}
		if parent != nil {
			x["parent_id"] = parent
		}

		err, objId := srv.formData.NewItem(userId, id, flowId, instance, x)
		if err != nil {
			return err, nil
		}
		if err, _ := srv.log.NewLog(userId, id, data_log.ADD, x); err != nil {
			return err, nil
		}

		list = append(list, objId)
	}
	return nil, list
}

/**
 * 新增数据控件格式校验
 * @params *gin.Context ctx http请求内容
 * @params map[string]interface{} data 原始数据
 * @params []form.Field list 控件列表
 * @params uint64 appId 应用id
 * @params uint64 index 数据序号
 * @params string method 校验方式 新增/编辑
 * @params string objId 记录id
 * @return error 错误信息
 * @return map[string]interface{} 校验完成的数据信息
 */
func (ds *defaultServer) validateItem(ctx *gin.Context, data map[string]interface{}, list []form.Field, appId, index uint64, method, objId string) (error, map[string]interface{}) {
	maps := make(map[string]interface{})
	mongoSrv := form_data.NewFormDataService(ds.mongoDb)
	for _, x := range list {
		//非收集信息类控件的值不做存储
		if form.FIELD_USE_TYPE[x.Type] != form.C {
			continue
		}

		//参数必填校验
		if (*x.IsNecessary == true) && (data[x.Key] == nil) {
			return errors.New(getResponseMsgWithLang(WRONG_PARAM, ds.getLanguage(ctx))), nil
		}

		//参数查重
		if *(x.IsOnly) == true {
			//mongo数据库查询
			found, d := mongoSrv.GetByFormAndKey(x.FormID, x.Key, data[x.Key])
			if found {
				if method == data_log.EDIT {
					itemId, ok := d["_id"].(primitive.ObjectID)
					if !ok {
						return errors.New("get log fail"), nil
					}
					if itemId.Hex() != objId {
						logrus.Error(itemId.Hex())
						return errors.New("data duplication:" + x.Title), nil
					}
				} else {
					return errors.New("data duplication:" + x.Title), nil
				}
			}
		}

		maps[x.Key] = data[x.Key]

		//文件收集类控件
		if x.Type == form.AUTOGRAPH || x.Type == form.FILE {
			if objId != "" {
				delete(maps, x.Key)
				continue
			}
			forms, _ := ctx.MultipartForm()
			if *x.IsNecessary == true && forms == nil {
				return errors.New(getResponseMsgWithLang(WRONG_PARAM, ds.getLanguage(ctx))), nil
			}
			if forms != nil {
				fileKey := fmt.Sprintf("%s_%d[]", x.Key, index)
				err, path := ds.fileUpload(forms.File[fileKey], appId)
				if err != nil {
					return err, nil
				}
				group := make([]fileFormat, 0, cap(path))
				for _, x := range path {
					if x.Hash != "" {
						group = append(group, fileFormat{
							Name: x.Name,
							Path: ds.getFullUrl(x.Hash),
						})
					}
				}
				maps[x.Key] = group
			}
		} else if x.Type == form.DATETIME_RANGE {
			//时间区间控件时间重叠校验
			var cont map[string]interface{}
			if err := json.Unmarshal(x.Content, &cont); err != nil {
				logrus.Error("parse dateRange field failed", err)
				return err, nil
			}
			lap, ok := cont["timeOverlap"].(bool)
			if !ok {
				return errors.New("parse timeOverlap failed"), nil
			}
			//判断时间重叠
			if lap {
				if err := checkTimeOverlap(x, data[x.Key], mongoSrv, objId, ds.getLanguage(ctx)); err != nil {
					return errors.New(err.Error()), nil
				}
			}
		} else if x.Type == form.LINKAGE {
			if data[x.Key] == nil {
				continue
			}
			v, ok := data[x.Key].([]interface{})
			if !ok {
				logrus.Error("linkage decode failed")
				return errors.New("params error"), nil
			}
			if len(v) > 2 {
				logrus.Error("linkage length too long")
				return errors.New("params error"), nil
			}
		} else if x.Type == form.SINGLE_CHOICE {
			var cont form.SingleChoice
			if err := json.Unmarshal(x.Content, &cont); err != nil {
				logrus.Error("parse singleChoice field failed", err)
				return err, nil
			}
			if cont.Hidden {
				maps[x.Key] = cont.DefaultValue
			}
		} else if x.Type == form.MULTIPLE_CHOICE {
			if data[x.Key] == nil {
				continue
			}
			mcStr, ok := data[x.Key].(string)
			if !ok {
				logrus.Error("format of multipleChoice is wrong,is not a string")
				return errors.New(getResponseMsgWithLang(WRONG_PARAM, ds.getLanguage(ctx))), nil
			}
			mcSlice := strings.Split(mcStr, "、")
			sort.Strings(mcSlice[:])
			maps[x.Key] = strings.Join(mcSlice, "、")
		}
	}
	return nil, maps
}

// 校验时间段是否重叠
// @params field form.Field 要检验的字段信息
// @params value interface{} 要检验的时间区间
// @params srv form_data.FormDataService 数据model对象
// @params objId string 数据编辑时必传,校验的时间区间对应的主键id
// @params lang string 接口语言
// @return error 返回nil时则未发生时间重叠
func checkTimeOverlap(filed form.Field, value interface{}, srv form_data.FormDataService, objId, lang string) error {
	ranges, ok := value.([]float64)
	if !ok {
		logrus.Error("format of dateRange is wrong,is not a string")
		return errors.New(getResponseMsgWithLang(WRONG_PARAM, lang))
	}
	now1 := utils.Time(ranges[0])
	now2 := utils.Time(ranges[1])
	err, _, list := srv.GetFullFormList(filed.FormID)
	if err != nil {
		return err
	}

	for _, v := range list {
		if objId != "" {
			_id, ok := v["_id"].(string)
			if !ok {
				logrus.Error("record id layout is wrong")
				return errors.New("record id layout is wrong")
			}
			if objId == _id {
				continue
			}
		}
		itemRange, ok := v[filed.Key].([]float64)
		if !ok {
			logrus.Error("record layout is wrong")
			return errors.New("record layout is wrong")
		}
		r1 := utils.Time(itemRange[0])
		r2 := utils.Time(itemRange[1])
		if now2.Before(r1) || now1.After(r2) {
			break
		} else {
			return errors.New("Time cannot be overlapped:" + filed.Title)
		}
	}
	return nil
}

// 格式化筛选条件
// @params appId uint64 应用id
// @params formId uint64 表单id
// @params search map[string]interface{} 筛选条件
// @params lang string 请求语言
// @return error 错误信息
// @return map 格式化后的筛选条件
func (ds *defaultServer) getFilter(appId, formId uint64, search map[string]interface{}, lang string) (error, map[string]interface{}) {
	//获取要检索的所有key
	keys := utils.GetStringKeys(search)
	filter := make(map[string]interface{})

	//获取要检索key切片的具体信息
	srv := form.NewFormService(ds.db)
	err, list := srv.FieldList(formId, keys...)
	if err != nil {
		return err, nil
	}

	//置入默认字段
	if _, ok := search["user_id"]; ok {
		list = append(list, form.Field{
			Key:  "user_id",
			Type: "USER_ID",
		})
	}
	if _, ok := search["user_name"]; ok {
		list = append(list, form.Field{
			Key:  "user_name",
			Type: "USER_NAME",
		})
	}
	if _, ok := search["created_at"]; ok {
		list = append(list, form.Field{
			Key:  "created_at",
			Type: "CREATED_AT",
		})
	}

	//格式化
	for _, v := range list {
		switch v.Type {
		case "USER_ID":
			filter[v.Key] = search[v.Key]
		case "USER_NAME":
			name, ok := search[v.Key].(string)
			if !ok {
				return errors.New(getResponseMsgWithLang(WRONG_PARAM, lang)), nil
			}
			//根据true name模糊查询用户id
			userSrv := user.NewUserService(ds.db)
			err, _, users := userSrv.GetAllUser(appId, name)
			if err != nil {
				return err, nil
			}
			//in方式获取数据
			filter["user_id"] = mongodb.In(users.UserIds())
		case "CREATED_AT":
			format, ok := search[v.Key].([]interface{})
			if !ok {
				return errors.New(getResponseMsgWithLang(WRONG_PARAM, lang)), nil
			}
			period := make([]time.Time, 0, cap(format))
			for _, x := range format {
				period = append(period, *ds.getTime(x))
			}
			filter[v.Key] = mongodb.Where([]bson.E{
				{">=", period[0]},
				{"<", period[1].AddDate(0, 0, 1)},
			})
		case form.INPUT:
			filter[v.Key] = mongodb.Like(search[v.Key])
		case form.AMOUNT:
			if _, ok := search[v.Key].([]interface{}); !ok {
				return errors.New(getResponseMsgWithLang(WRONG_PARAM, lang)), nil
			}
			filter[v.Key] = mongodb.Between(search[v.Key])
		case form.NUMBER:
			if _, ok := search[v.Key].([]interface{}); !ok {
				return errors.New(getResponseMsgWithLang(WRONG_PARAM, lang)), nil
			}
			filter[v.Key] = mongodb.Between(search[v.Key])
		case form.MAIL:
			filter[v.Key] = mongodb.Like(search[v.Key])
		case form.MULTIPLE_CHOICE:
			//TODO sort
			filter[v.Key] = search[v.Key]
		case form.SINGLE_CHOICE:
			filter[v.Key] = search[v.Key]
		case form.CERTIFICATE:
			filter[v.Key] = mongodb.Like(search[v.Key])
		case form.LEVEL:
			filter[v.Key] = search[v.Key]
		case form.PHONE:
			filter[v.Key] = mongodb.Like(search[v.Key])
		default:
			logrus.Info(v.Key, ":do not need filter")
		}
	}
	return nil, filter
}

/**
 * 筛选条件区分
 * @params ContentStruct config 查询控件配置信息
 * @params map[string]interface{} search 筛选条件原始数据
 * @return map[uint64]map[string]interface{} 处理后的筛选条件
 */
func searchData(config ContentStruct, search map[string]interface{}) map[uint64]map[string]interface{} {
	searchData := make(map[uint64]map[string]interface{})
	for k, s := range search {
		for _, c := range config.Columns {
			if k == c.Key {
				if c.FormId == config.FormId {
					searchData[0][k] = s
				} else {
					searchData[c.FormId][k] = s
				}
				break
			}
		}
	}
	return searchData
}

/**
 * 获取记录详情
 * @params int8 method 查询方式 1-_id 2-instance_id
 * @params string id 查询关键字
 * @params string types 结果返回方式 show-解析后的详情
 * @return int16 结果码
 * @return string 提示信息
 * @return form_data.FormDataItem 详情数据
 */
func (ds *defaultServer) formDataItemDetail(method int8, id, types string) (int16, string, form_data.FormDataItem) {
	dataSrv := form_data.NewFormDataService(ds.mongoDb)
	userSrv := user.NewUserService(ds.db)
	var status bool
	var item form_data.FormDataItem
	switch method {
	case 1:
		status, item = dataSrv.GetByObjId(id)
	case 2:
		status, item = dataSrv.GetByInstanceId(id)
	}
	if !status {
		return NOT_FOUND, "", nil
	}

	//编辑与查看时的详情返回须做区分,查看时的详情返回最终数据
	if types == "show" {
		formIdStr := fmt.Sprintf("%v", item["form_id"])
		formId := utils.NewStr(formIdStr).Uint64()
		if formId == 0 {
			logrus.Error("wrong form_id")
			return FAIL, "", nil
		}

		formSrv := form.NewFormService(ds.db)
		err, fields := formSrv.FieldList(formId)
		if err != nil {
			return FAIL, err.Error(), nil
		}
		for _, x := range fields {
			switch x.Type {
			case form.MEMBER:
				e, names := ds.getMembers(item[x.Key], userSrv)
				if e != nil && e.Error() != "" {
					return FAIL, e.Error(), nil
				}
				item[x.Key] = strings.Join(names, ",")
			case form.ORGANIZATION:
				e, names := ds.getOrganizations(item[x.Key], userSrv)
				if e != nil && e.Error() != "" {
					return FAIL, e.Error(), nil
				}
				item[x.Key] = strings.Join(names, ",")
			default:
				continue
			}
		}
	}
	t, ok := item["created_at"].(primitive.DateTime)
	if ok {
		item["created_at"] = t.Time().Format(utils.TimeFormatStr)
	} else {
		ti, ok := item["created_at"].(float64)
		if ok {
			item["created_at"] = utils.Time(ti).Format(utils.TimeFormatStr)
		}
	}

	idStr := fmt.Sprintf("%v", item["user_id"])
	userId := utils.NewStr(idStr).Uint64()
	if userId != 0 {
		_, users := userSrv.GetUser(userId)
		item["user_name"] = users.TrueName
	}
	return SUCCESS, "", item
}
