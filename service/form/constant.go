// Package form Created by Goland
//@User: lenora
//@Date: 2021/3/12
//@Time: 7:39 下午
package form

const (
	DISABLE = iota
	ENABLE
)

const (
	FORM_PAGE = 1
	SHOW_PAGE = 2
)

const (
	INPUT           = "Input"          //文本
	AMOUNT          = "Amount"         //金额
	NUMBER          = "Number"         //数值
	MAIL            = "Mail"           //邮箱
	MULTIPLE_CHOICE = "MultipleChoice" //多选
	SINGLE_CHOICE   = "SingleChoice"   //单选
	AREA            = "Area"           //地区
	AUTOGRAPH       = "Autograph"      //签名
	CERTIFICATE     = "Certificates"   //证件
	DATETIME        = "DateTime"       //日期
	DATETIME_RANGE  = "DatetimeRange"  //起止时间
	FILE            = "File"           //附件
	LEVEL           = "Level"          //等级
	MEMBER          = "Member"         //成员
	PHONE           = "Phone"          //电话
	RECORDS         = "Records"        //关联记录
	REMARK          = "Remark"         //备注/说明
	BUTTON          = "Button"         //按钮
	CAROUSEL        = "Carousel"       //轮播
	IMAGE           = "Image"          //静态图片
	ORGANIZATION    = "Organization"   //组织
	TABLES          = "Table"          //列表
	DIVIDER         = "Divider"        //分段
	TEXT            = "Text"           //展示文本
	TABS            = "Tab"            //标签 //容器 //布局
	MULTIFORM       = "MultiForm"      //多表单
	LINKAGE         = "CascadeControl" //级联
	// SUBTABLE = "Subtable"
)

const (
	STANDARD int8 = iota + 1
	POPUP
	NAV
)

const (
	IS        uint8 = iota + 1 //是
	IS_NOT                     //不是
	CONTAIN                    //包含
	EXCLUSIVE                  //不包含
	NULL                       //为空
	NOT_NULL                   //不为空
	EQ                         //等于
	NEQ                        //不等于
	GT                         //大于
	LT                         //小于
	GTE                        //大于等于
	LTE                        //小于等于
	EARLY                      //早于
	LATER                      //晚于
)

var FIELD_USE_TYPE = map[string]string{
	INPUT:           C,
	AMOUNT:          C,
	NUMBER:          C,
	MAIL:            C,
	MULTIPLE_CHOICE: C,
	SINGLE_CHOICE:   C,
	AREA:            C,
	AUTOGRAPH:       C,
	CERTIFICATE:     C,
	DATETIME:        C,
	DATETIME_RANGE:  C,
	FILE:            C,
	LEVEL:           C,
	MEMBER:          C,
	PHONE:           C,
	RECORDS:         A,
	REMARK:          N,
	BUTTON:          N,
	CAROUSEL:        N,
	IMAGE:           N,
	ORGANIZATION:    C,
	TABLES:          D,
	DIVIDER:         N,
	TEXT:            N,
	TABS:            L,
	MULTIFORM:       A,
	LINKAGE:         A,
}

var TRUE, FALSE = true, false
var D, C, L, N, A = "display", "collect_data", "layout", "normal", "associated"

var TEXT_CONDITION = []uint8{IS, IS_NOT, CONTAIN, EXCLUSIVE, NULL, NOT_NULL}
var NUMBER_CONDITION = []uint8{EQ, NEQ, GT, LT, GTE, LTE, NULL, NOT_NULL}
var MAP_CONDITION = TEXT_CONDITION
var DATA_CONDITION = []uint8{EARLY, LATER, EQ, NEQ, NULL, NOT_NULL}

type FieldFormat struct {
	Id     string                   `json:"id"`
	Layers []map[string]interface{} `json:"layers"`
}

type ChildTabs struct {
	Id       uint64        `json:"id"`
	TabName  string        `json:"tabName"`
	Children []FieldFormat `json:"children"`
}

type CommonField struct {
	Type          string `json:"type"`
	TypeName      string `json:"typeName"`
	Name          string `json:"name"`
	Id            string `json:"id"`
	Description   string `json:"description"`
	Require       bool   `json:"require"`
	Repeat        bool   `json:"repeat"`
	Colspan       uint64 `json:"colspan"`
	Explain       string `json:"explain"`
	Title         string `json:"title"`
	Placeholder   string `json:"placeholder"`
	SourceFieldId uint64 `json:"sourceFieldId"`
}

type FormTree struct {
	Id       uint64       `json:"id"`
	Children FormTreeList `json:"children"`
}

type FormTreeList []FormTree
