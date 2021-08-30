package form

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Lenora-Z/low-code/utils"
	"github.com/sirupsen/logrus"
	"strings"
)

func (srv *formService) formatField(formId, fieldId uint64, types string, item []byte) (error, []uint64) {
	if !inFieldFilter(types) {
		return nil, nil
	} else {
		var err error
		ret := make([]uint64, 0, 100)
		switch types {
		case NUMBER:
			err = srv.formNumber(item)
		case TEXT:
			err = srv.formText(item)
		case INPUT:
			err = srv.formInput(item)
		case MULTIFORM:
			err = srv.formatMultiForm(formId, fieldId, item)
		case LINKAGE:
			err = srv.formatLinkage(formId, fieldId, item)
		case RECORDS:
			err = srv.formRecords(formId, fieldId, item)
		case BUTTON:
			err, ret = srv.formButton(formId, fieldId, item)
		case TABLES:
			err, ret = srv.formTable(formId, fieldId, item)
		default:
			break
		}
		if err != nil {
			return err, nil
		}
		return nil, ret

	}
}

type Relations struct {
	Id     uint64 `json:"id"`
	Name   string `json:"name"`
	Status bool   `json:"status"`
}

type MultiForm struct {
	CommonField
	FormId     uint64        `json:"formId"` //关联主表id
	Relations  Relations     `json:"relationFormData"`
	ShowType   string        `json:"showType"`   // 呈现方式
	ImportType string        `json:"importType"` // 引入方式
	Columns    []FieldFormat `json:"columns"`    // 字段 {config} -> 控件 {config}
}

func (srv *formService) formatMultiForm(formId, fieldId uint64, content []byte) error {
	var item MultiForm
	if err := json.Unmarshal(content, &item); err != nil {
		return err
	}
	if item.FormId == 0 {
		logrus.Error("multi-form form_id can not be null")
		return errors.New("param error")
	}

	//呈现方式转换
	mode := multiFormShowType.searchKey(item.ShowType)
	if mode == 0 {
		logrus.Error("multi-form show_type wrong!!!")
		return errors.New("param error")
	}

	var method = multiFormImportType.searchKey(item.ImportType)
	if method == 0 {
		logrus.Error("multi-form importType wrong!!!")
		return errors.New("param error")
	}

	if item.Relations.Id == 0 {
		if method == COPY {
			logrus.Error("multi-form relation_id can not be null")
			return errors.New("param error")
		}

		if method == RELATION {
			item.Relations = Relations{Id: item.FormId}
		}
	}

	if err, _ := createMultiForm(srv.db, FieldMultiForm{
		FormID:  formId,
		FieldID: fieldId,
		ChildID: item.Relations.Id,
		Method:  method,
		Mode:    mode,
	}); err != nil {
		return err
	}
	return nil

}

type SingleLinkage struct {
	RelationId uint64 `json:"relation_id"` //关联id
	Id         uint64 `json:"id"`          //字段id
	Name       string `json:"name"`        // 字段名称
	Key        string `json:"key"`
	FormId     uint64 `json:"formId"`
	FormName   string `json:"formName"` // 原表单
	Title      string `json:"title"`    //显示文案
}

type LinkAge struct {
	CommonField
	FormId uint64          `json:"formId"` //关联主表id
	Ranks  []SingleLinkage `json:"ranks"`
}

func (srv *formService) formatLinkage(formId, fieldId uint64, content []byte) error {
	var item LinkAge
	if err := json.Unmarshal(content, &item); err != nil {
		return err
	}

	if len(item.Ranks) < 1 {
		logrus.Error("param error:rank length")
		return errors.New("param error")
	}

	con := linkageEncode(item.Ranks)

	if err, _ := createLinkage(srv.db, FieldLinkage{
		FormID:  formId,
		FieldID: fieldId,
		Content: con,
	}); err != nil {
		return err
	}
	return nil
}

type SingleChoice struct {
	DefaultValue string `json:"defaultValue"`
	Hidden       bool   `json:"hidden"`
}

type RecordsField struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
	Key  string `json:"key"`
}

type Records struct {
	CommonField
	Single         bool           `json:"single"`         //显示数量
	SheetId        uint64         `json:"sheetId"`        //关联id
	FormFields     []RecordsField `json:"formFields"`     //显示字段
	RelationFields []RecordsField `json:"relationFields"` //关联字段
	ShowType       string         `json:"showType"`       //'card'-卡片 | 'backfill'-自动回填
	IsShowDetail   bool           `json:"isShowDetail"`   //是否需要查看详情
}

func (srv *formService) formRecords(formId, fieldId uint64, content []byte) error {
	var item Records
	if err := json.Unmarshal(content, &item); err != nil {
		return err
	}

	mode := recordsModeType.searchKey(item.ShowType)
	if mode == 0 {
		logrus.Error("records show_type wrong!!!")
		return errors.New("param error")
	}

	var single uint8 = 1
	if !item.Single {
		single = 2
	}

	columns := make([]uint64, 0, cap(item.FormFields))
	for _, c := range item.FormFields {
		if c.Id != 0 {
			columns = append(columns, c.Id)
		}
	}
	relationColumns := make([]uint64, 0, cap(item.RelationFields))
	for _, rc := range item.RelationFields {
		if rc.Id != 0 {
			relationColumns = append(relationColumns, rc.Id)
		}
	}

	err, _ := createRecords(srv.db, FieldRecords{
		FormID:                     formId,
		FieldID:                    fieldId,
		DatasourceColumnRelationID: item.SheetId,
		DatasourceColumnIDs:        utils.UintJoin(relationColumns, ","),
		Mode:                       mode,
		CountType:                  single,
		DetailStatus:               item.IsShowDetail,
		Columns:                    utils.UintJoin(columns, ","),
	})
	return err
}

type TableColumn struct {
	Id         uint64 `json:"id"`
	Title      string `json:"title"`
	ShowName   string `json:"showName"`
	Filter     bool   `json:"filter"`
	RelationId uint64 `json:"relation_id"` //关联id
	ColumnId   uint64 `json:"column_id"`   //字段id
}

type FilterCondition struct {
	Id         uint64        `json:"id"`
	Condition  uint8         `json:"condition"`
	Types      string        `json:"type"`
	Value      []interface{} `json:"value"`
	RelationId uint64        `json:"relation_id"`
}

type Table struct {
	CommonField
	FormId             uint64            `json:"formId"`     //数据表id
	NeedExport         bool              `json:"needExport"` //是否导出 true => 是
	Columns            []TableColumn     `json:"columns"`
	InlineOperationBtn []string          `json:"inline_operation_btn"` //行内常规操作按钮
	Layers             []Button          `json:"layers"`               //操作按钮
	InlineBtn          []Button          `json:"inlineBtn"`            //行内操作按钮
	DataFilterConfig   bool              `json:"dataFilterConfig"`     //数据过滤配置
	Condition          []FilterCondition `json:"condition"`            //数据过滤
}

func (srv *formService) formTable(formId, fieldId uint64, content []byte) (error, []uint64) {
	var item Table
	if err := json.Unmarshal(content, &item); err != nil {
		return err, nil
	}

	flows := make([]uint64, 0, 100)

	//行内常规操作按钮解析
	for _, btn := range item.InlineOperationBtn {
		var event = ButtonType.searchKey(btn)
		if event == 0 {
			logrus.Error("button type wrong!!!")
			return errors.New("param error"), nil
		}
		if err, _ := createTableButton(srv.db, FieldTableButton{
			FormID:  formId,
			FieldID: fieldId,
			FlowID:  0,
			Name:    "",
			Event:   event,
		}); err != nil {
			return err, nil
		}
	}

	//行内按钮及其余操作按钮解析
	buttons := make([]Button, 0, cap(item.InlineBtn)+cap(item.Layers))
	buttons = append(item.InlineBtn, item.Layers...)
	for _, btn := range buttons {
		var event = ButtonType.searchKey(btn.Trigger)
		if event == 0 {
			logrus.Error("button type wrong!!!")
			return errors.New("param error"), nil
		}
		var flowId uint64
		if event == PROCESS {
			flowId = btn.Process.Id
		}
		if err, _ := createTableButton(srv.db, FieldTableButton{
			FormID:  formId,
			FieldID: fieldId,
			FlowID:  flowId,
			Name:    btn.Title,
			Event:   event,
		}); err != nil {
			return err, nil
		}

		if flowId != 0 {
			flows = append(flows, flowId)
		}
	}

	//解析存储表字段
	for _, col := range item.Columns {
		if err, _ := createTableColumn(srv.db, FieldTableColumn{
			FormID:                     formId,
			FieldID:                    fieldId,
			DatasourceColumnRelationID: col.RelationId,
			DatasourceColumnID:         col.ColumnId,
			ShowName:                   col.ShowName,
			IsCondition:                col.Filter,
		}); err != nil {
			return err, nil
		}
	}

	if item.DataFilterConfig {
		for _, cond := range item.Condition {
			switch cond.Types {
			case "num":
				switch cond.Condition {
				case EQ:
					cond.Condition = IS
				case NEQ:
					cond.Condition = IS_NOT
				case GT:
					cond.Condition = CONTAIN
				case LT:
					cond.Condition = EXCLUSIVE
				case GTE:
					cond.Condition = NULL
				case LTE:
					cond.Condition = NOT_NULL
				case NULL:
					cond.Condition = EQ
				case NOT_NULL:
					cond.Condition = NEQ
				}
			case "date":
				switch cond.Condition {
				case EARLY:
					cond.Condition = IS
				case LATER:
					cond.Condition = IS_NOT
				case EQ:
					cond.Condition = CONTAIN
				case NEQ:
					cond.Condition = EXCLUSIVE
				}
			}
			value := make([]string, 0, cap(cond.Value))
			for _, t := range cond.Value {
				if cond.Types == "date" {
					value = append(value, fmt.Sprintf("%.0f", t))
				} else {
					value = append(value, fmt.Sprintf("%v", t))
				}
			}
			if err, _ := createTableFilter(srv.db, FieldTableFilter{
				FormID:                     formId,
				FieldID:                    fieldId,
				DatasourceColumnRelationID: cond.RelationId,
				DatasourceColumnID:         cond.Id,
				FieldType:                  cond.Types,
				FieldTypeCondition:         cond.Condition,
				FieldTypeConditionValue:    strings.Join(value, ","),
			}); err != nil {
				return err, nil
			}
		}
	}

	err, _ := createTable(srv.db, FieldTable{
		FormID:            formId,
		FieldID:           fieldId,
		DatasourceTableID: item.FormId,
		IsExport:          item.NeedExport,
		IsFilter:          item.DataFilterConfig,
	})
	return err, flows

}

type ButtonProcess struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
}

type Button struct {
	CommonField
	Trigger string        `json:"trigger"`
	Process ButtonProcess `json:"process"`
}

func (srv *formService) formButton(formId, fieldId uint64, content []byte) (error, []uint64) {
	var item Button
	if err := json.Unmarshal(content, &item); err != nil {
		return err, nil
	}

	var event = ButtonType.searchKey(item.Trigger)
	if event == 0 {
		logrus.Error("button type wrong!!!")
		return errors.New("param error"), nil
	}
	var flowId uint64
	if event == PROCESS {
		flowId = item.Process.Id
	}

	err, _ := createButton(srv.db, FieldButton{
		FormID:  formId,
		FieldID: fieldId,
		FlowID:  flowId,
		Name:    item.Title,
		Event:   event,
	})
	if flowId == 0 {
		return err, nil
	} else {
		return err, []uint64{flowId}
	}
}

type Number struct {
	CommonField
	Digit       uint64      `json:"digit"`
	LimitAmount bool        `json:"limitAmount"`
	Max         interface{} `json:"max"` //float64
	Min         interface{} `json:"min"` //float64
	Unit        string      `json:"unit"`
}

func (srv *formService) formNumber(content []byte) error {
	var item Number
	if err := json.Unmarshal(content, &item); err != nil {
		return errors.New("param error:" + err.Error())
	}

	if item.LimitAmount {
		if item.Max == nil || item.Min == nil {
			logrus.Info("field number config without max or min")
			return errors.New("param error")
		}
	}
	return nil
}

type Text struct {
	CommonField
	Text          string `json:"text"`
	TitleSize     string `json:"titleSize"`
	PaddingBottom uint64 `json:"paddingBottom"`
	PaddingLeft   uint64 `json:"paddingLeft"`
	PaddingRight  uint64 `json:"paddingRight"`
	PaddingTop    uint64 `json:"paddingTop"`
}

func (srv *formService) formText(content []byte) error {
	var item Text
	if err := json.Unmarshal(content, &item); err != nil {
		return errors.New("param error:" + err.Error())
	}

	if item.Text == "" || item.TitleSize == "" {
		logrus.Info("field ext config without text or titleSize")
		return errors.New("param error")
	}

	if TextSize.searchKey(item.TitleSize) == 0 {
		logrus.Error("wrong title size")
		return errors.New("param error")
	}
	return nil
}

type Input struct {
	CommonField
	LimitWordNumber bool        `json:"limitWordNumber"`
	MaxWordNumber   interface{} `json:"maxWordNumber"`
	MinWordNumber   interface{} `json:"minWordNumber"`
}

func (srv *formService) formInput(content []byte) error {
	var item Input
	if err := json.Unmarshal(content, &item); err != nil {
		return errors.New("param error:" + err.Error())
	}

	if item.LimitWordNumber {
		if item.MaxWordNumber == nil || item.MinWordNumber == nil {
			logrus.Info("field input config without max or min")
			return errors.New("param error")
		}
	}
	return nil
}

func linkageEncode(group []SingleLinkage) string {
	items := make([]string, 0, cap(group))
	for _, x := range group {
		items = append(items, fmt.Sprintf("%d#%d#%s", x.RelationId, x.Id, x.Title))
	}
	return strings.Join(items, "@")
}

func linkageDecode(content string) []SingleLinkage {
	group := make([]SingleLinkage, 0, 50)
	items := strings.Split(content, "@")
	for _, x := range items {
		_temp := strings.Split(x, "#")
		group = append(group, SingleLinkage{
			RelationId: utils.NewStr(_temp[0]).Uint64(),
			Id:         utils.NewStr(_temp[1]).Uint64(),
			Title:      _temp[2],
		})
	}
	return group

}

func inFieldFilter(types string) bool {
	typeSlice := []string{MULTIFORM, LINKAGE, RECORDS, BUTTON, TABLES, NUMBER, TEXT, INPUT}
	for _, x := range typeSlice {
		if x == types {
			return true
		}
	}
	return false
}
