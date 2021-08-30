// Package flow
// Created by GoLand
// @User: lenora
// @Date: 2021/7/28
// @Time: 17:22

package flow

import (
	"encoding/json"
	"fmt"
	"github.com/Lenora-Z/low-code/utils"
	"regexp"
	"strings"
)

type StartNode struct {
	Name               string  `json:"name"`        //开始节点
	Description        string  `json:"description"` //节点描述
	Type               uint8   `json:"type"`        // 1: 固定时间 2: 表单时间
	DatasourceTableId  uint64  `json:"datasource_table_id"`
	DatasourceColumnId uint64  `json:"datasource_column_id"`
	TriggerTime        float64 `json:"trigger_time"`
	TriggerInterval    int64   `json:"trigger_interval"`
}

type ChainUpNode struct {
	Type     uint8  `json:"type"`      // 1: 节点 2: 表单
	TargetId uint64 `json:"target_id"` // 节点/表单id
	Columns  string `json:"columns"`   // 字段id逗号分割字符串
}

type EmailParams struct {
	ServiceId uint64 `json:"service_id"`
	SenderId  uint64 `json:"sender_id"`
	ReviverId uint64 `json:"reviver_id"`
	Subject   uint64 `json:"subject"`
	Content   uint64 `json:"content"`
}

type ChainUpParams struct {
	ServiceId uint64 `json:"service_id"`
	ParamId   uint64 `json:"param_id"`
}

type EmailTargetName struct {
	Type  uint8    `json:"type"`
	Value []uint64 `json:"value"`
}

type EmailTarget struct {
	Name EmailTargetName `json:"name"`
}

type EmailNode struct {
	Name        string        `json:"name"`        //开始节点
	Description string        `json:"description"` //节点描述
	Target      []EmailTarget `json:"target"`      // 发送对象
	Title       string        `json:"title"`       // 主题
	Content     string        `json:"content"`     //正文
	Address     string        `json:"address"`     // 地址
}

type AddDataColumn struct {
	Type  uint8       `json:"type"`
	Value interface{} `json:"value"`
}

type AddDataNode struct {
	Name        string                     `json:"name"`        //开始节点
	Description string                     `json:"description"` //节点描述
	TargetId    uint64                     `json:"target_id"`
	Columns     []map[string]AddDataColumn `json:"columns"`
}

type EditDataNode struct {
	AddDataNode
	ConditionList []ConditionGroup `json:"condition_list"`
}

type ConditionItem struct {
	OriginFormId   uint64 `json:"origin_form_id"`
	OriginId       uint64 `json:"origin_id"`
	TargetFormId   uint64 `json:"target_form_id"`
	TargetId       uint64 `json:"target_id"`
	Condition      uint8  `json:"condition"`
	TargetValue    string `json:"target_value"`
	OriginFormType uint8  `json:"origin_form_type"`
	TargetFormType uint8  `json:"target_form_type"`
}

type ConditionGroup []ConditionItem

type ConditionList struct {
	ConditionList []ConditionGroup `json:"condition_list"`
}

type CodeNode struct {
	TargetId uint64                     `json:"id"` //代码包id
	Columns  []map[uint64]AddDataColumn `json:"columns"`
}

type NodeCommon struct {
	Key   uint64      `json:"key"`         // 节点唯一标识
	Name  string      `json:"name"`        // 节点名称
	Desc  string      `json:"description"` // 节点描述
	Type  uint        `json:"type"`        // 节点类型
	Props interface{} `json:"props"`       //节点属性
	List  []Node      `json:"list"`
}

type Node []NodeCommon

func (srv *flowService) flowDecode(flowId, appId uint64, content []byte, chainUp ChainUpParams, email EmailParams) error {
	var node Node
	if err := json.Unmarshal(content, &node); err != nil {
		return err
	}

	for _, x := range node {
		err, act := createActivity(srv.db, FlowActivity{
			FlowID: flowId,
			NodeID: x.Key,
			Name:   x.Name,
			Desc:   x.Desc,
		})
		if err != nil {
			return err
		}
		acByte, err := json.Marshal(x.Props)
		if err != nil {
			return err
		}
		switch x.Type {
		case 0: // 开始节点
			var n StartNode
			var t uint8
			if err := json.Unmarshal(acByte, &n); err != nil {
				return err
			}
			if n.Type == 1 {
				t = START_TIME
				if err, _ := createActivityTime(srv.db, FlowActivityStartTime{
					FlowID:          flowId,
					FlowActivityID:  act.ID,
					TriggerTime:     utils.Time(n.TriggerTime),
					TriggerInterval: n.TriggerInterval,
				}); err != nil {
					return err
				}
			} else if n.Type == 2 {
				t = START_TABLE
				if err, _ := createActivityTable(srv.db, FlowActivityStartTable{
					FlowID:             flowId,
					FlowActivityID:     act.ID,
					DatasourceTableID:  n.DatasourceTableId,
					DatasourceColumnID: n.DatasourceColumnId,
					TriggerInterval:    n.TriggerInterval,
				}); err != nil {
					return err
				}
			}

			if err, _ := updateActivity(srv.db, FlowActivity{
				ID:   act.ID,
				Type: t,
			}); err != nil {
				return err
			}
		case 1: // 上链节点
			var n ChainUpNode
			if err := json.Unmarshal(acByte, &n); err != nil {
				return err
			}
			var com = FlowActivityService{
				AppID:                   appId,
				FlowID:                  flowId,
				FlowActivityID:          act.ID,
				ParamID:                 chainUp.ParamId,
				InputFieldTableColumnID: 0,
			}
			switch n.Type {
			case 1:
				for _, i := range strings.Split(n.Columns, ",") {
					item := com
					com.InputNodeID = n.TargetId
					com.InputParamID = utils.NewStr(i).Uint64()
					if err, _ := createActivityService(srv.db, item); err != nil {
						return err
					}
				}
			case 2:
				for _, i := range strings.Split(n.Columns, ",") {
					item := com
					com.InputFieldID = utils.NewStr(i).Uint64()
					if err, _ := createActivityService(srv.db, item); err != nil {
						return err
					}
				}
			case 4:
				for _, i := range strings.Split(n.Columns, ",") {
					item := com
					com.InputFieldTableColumnID = utils.NewStr(i).Uint64()
					if err, _ := createActivityService(srv.db, item); err != nil {
						return err
					}
				}
			}
			if err, _ := updateActivity(srv.db, FlowActivity{
				ID:        act.ID,
				Type:      SERVICE,
				ServiceID: chainUp.ServiceId,
			}); err != nil {
				return err
			}
		case 2: // 发送邮件节点
			var n EmailNode
			if err := json.Unmarshal(acByte, &n); err != nil {
				return err
			}

			//发送人
			if err, _ := createActivityService(srv.db, FlowActivityService{
				AppID:          appId,
				FlowID:         flowId,
				FlowActivityID: act.ID,
				ParamID:        email.SenderId,
				InputText:      n.Address,
			}); err != nil {
				return err
			}
			//主题
			if err, _ := createActivityService(srv.db, FlowActivityService{
				AppID:          appId,
				FlowID:         flowId,
				FlowActivityID: act.ID,
				ParamID:        email.Subject,
				InputText:      n.Title,
			}); err != nil {
				return err
			}

			//收件人
			for _, t := range n.Target {
				switch t.Name.Type {
				case 1:
					//节点
					if err, _ := createActivityService(srv.db, FlowActivityService{
						AppID:          appId,
						FlowID:         flowId,
						FlowActivityID: act.ID,
						ParamID:        email.ReviverId,
						InputNodeID:    t.Name.Value[0],
						InputParamID:   t.Name.Value[1],
					}); err != nil {
						return err
					}
				case 2:
					//表单
					if err, _ := createActivityService(srv.db, FlowActivityService{
						AppID:          appId,
						FlowID:         flowId,
						FlowActivityID: act.ID,
						ParamID:        email.ReviverId,
						InputFieldID:   t.Name.Value[1],
					}); err != nil {
						return err
					}
				case 3:
					//用户
					if err, _ := createActivityService(srv.db, FlowActivityService{
						AppID:              appId,
						FlowID:             flowId,
						FlowActivityID:     act.ID,
						ParamID:            email.ReviverId,
						InputLowcodeUserID: t.Name.Value[0],
					}); err != nil {
						return err
					}
				case 4: //列表字段
					if err, _ := createActivityService(srv.db, FlowActivityService{
						AppID:                   appId,
						FlowID:                  flowId,
						FlowActivityID:          act.ID,
						ParamID:                 email.ReviverId,
						InputFieldTableColumnID: t.Name.Value[1],
					}); err != nil {
						return err
					}
				}
			}

			//发送内容
			if err, _ := createActivityService(srv.db, FlowActivityService{
				AppID:          appId,
				FlowID:         flowId,
				FlowActivityID: act.ID,
				ParamID:        email.Content,
				InputText:      n.Content,
			}); err != nil {
				return err
			}
			reg := regexp.MustCompile(`{(\d|,)+}`)
			regSlice := reg.FindAllString(n.Content, -1)
			for _, i := range regSlice {
				str := i[1 : len(i)-1]
				strSlice := strings.Split(str, ",")
				switch strSlice[0] {
				case "1":
					if err, _ := createActivityService(srv.db, FlowActivityService{
						AppID:          appId,
						FlowID:         flowId,
						FlowActivityID: act.ID,
						ParamID:        email.Content,
						InputNodeID:    utils.NewStr(strSlice[1]).Uint64(),
						InputParamID:   utils.NewStr(strSlice[2]).Uint64(),
					}); err != nil {
						return err
					}
				case "2":
					if err, _ := createActivityService(srv.db, FlowActivityService{
						AppID:          appId,
						FlowID:         flowId,
						FlowActivityID: act.ID,
						ParamID:        email.Content,
						InputFieldID:   utils.NewStr(strSlice[2]).Uint64(),
					}); err != nil {
						return err
					}
				case "4":
					if err, _ := createActivityService(srv.db, FlowActivityService{
						AppID:                   appId,
						FlowID:                  flowId,
						FlowActivityID:          act.ID,
						ParamID:                 email.Content,
						InputFieldTableColumnID: utils.NewStr(strSlice[2]).Uint64(),
					}); err != nil {
						return err
					}
				}

			}

			if err, _ := updateActivity(srv.db, FlowActivity{
				ID:        act.ID,
				Type:      SERVICE,
				ServiceID: email.ServiceId,
			}); err != nil {
				return err
			}
		case 6: // 新增数据节点
			var n AddDataNode
			if err := json.Unmarshal(acByte, &n); err != nil {
				return err
			}

			for _, c := range n.Columns {
				for k, v := range c {
					value := make([]uint64, 0, 10)
					sourceId := utils.NewStr(k).Uint64()
					if v.Type != 3 {
						b, err := json.Marshal(v.Value)
						if err != nil {
							return err
						}
						if err := json.Unmarshal(b, &value); err != nil {
							return err
						}
					}
					switch v.Type {
					case 1:
						if err, _ := createActivityData(srv.db, FlowActivityData{
							AppID:              appId,
							FlowID:             flowId,
							FlowActivityID:     act.ID,
							DatasourceTableID:  n.TargetId,
							DatasourceColumnID: sourceId,
							Type:               ACT_TYPE,
							Expression:         IS,
							InputNodeID:        value[0],
							InputParamID:       value[1],
						}); err != nil {
							return err
						}
					case 2:
						if err, _ := createActivityData(srv.db, FlowActivityData{
							AppID:              appId,
							FlowID:             flowId,
							FlowActivityID:     act.ID,
							DatasourceTableID:  n.TargetId,
							DatasourceColumnID: sourceId,
							Type:               ACT_TYPE,
							Expression:         IS,
							InputFieldID:       value[1],
						}); err != nil {
							return err
						}
					case 3:
						if err, _ := createActivityData(srv.db, FlowActivityData{
							AppID:              appId,
							FlowID:             flowId,
							FlowActivityID:     act.ID,
							DatasourceTableID:  n.TargetId,
							DatasourceColumnID: sourceId,
							Type:               ACT_TYPE,
							Expression:         IS,
							InputText:          fmt.Sprintf(`%v`, v.Value),
						}); err != nil {
							return err
						}
					case 4:
						if err, _ := createActivityData(srv.db, FlowActivityData{
							AppID:                   appId,
							FlowID:                  flowId,
							FlowActivityID:          act.ID,
							DatasourceTableID:       n.TargetId,
							DatasourceColumnID:      sourceId,
							Type:                    ACT_TYPE,
							Expression:              IS,
							InputFieldTableColumnID: value[1],
						}); err != nil {
							return err
						}
					}
				}
			}

			if err, _ := updateActivity(srv.db, FlowActivity{
				ID:   act.ID,
				Type: DATA_ACT,
			}); err != nil {
				return err
			}
		case 5: // 编辑数据节点
			var n EditDataNode
			if err := json.Unmarshal(acByte, &n); err != nil {
				return err
			}

			for _, c := range n.Columns {
				for k, v := range c {
					value := make([]uint64, 0, 10)
					sourceId := utils.NewStr(k).Uint64()
					if v.Type != 3 {
						b, err := json.Marshal(v.Value)
						if err != nil {
							return err
						}
						if err := json.Unmarshal(b, &value); err != nil {
							return err
						}
					}
					switch v.Type {
					case 1:
						if err, _ := createActivityData(srv.db, FlowActivityData{
							AppID:              appId,
							FlowID:             flowId,
							FlowActivityID:     act.ID,
							DatasourceTableID:  n.TargetId,
							DatasourceColumnID: sourceId,
							Type:               ACT_TYPE,
							Expression:         IS,
							InputNodeID:        value[0],
							InputParamID:       value[1],
						}); err != nil {
							return err
						}
					case 2:
						if err, _ := createActivityData(srv.db, FlowActivityData{
							AppID:              appId,
							FlowID:             flowId,
							FlowActivityID:     act.ID,
							DatasourceTableID:  n.TargetId,
							DatasourceColumnID: sourceId,
							Type:               ACT_TYPE,
							Expression:         IS,
							InputFieldID:       value[1],
						}); err != nil {
							return err
						}
					case 3:
						if err, _ := createActivityData(srv.db, FlowActivityData{
							AppID:              appId,
							FlowID:             flowId,
							FlowActivityID:     act.ID,
							DatasourceTableID:  n.TargetId,
							DatasourceColumnID: sourceId,
							Type:               ACT_TYPE,
							Expression:         IS,
							InputText:          fmt.Sprintf(`%v`, v.Value),
						}); err != nil {
							return err
						}
					case 4:
						if err, _ := createActivityData(srv.db, FlowActivityData{
							AppID:                   appId,
							FlowID:                  flowId,
							FlowActivityID:          act.ID,
							DatasourceTableID:       n.TargetId,
							DatasourceColumnID:      sourceId,
							Type:                    ACT_TYPE,
							Expression:              IS,
							InputFieldTableColumnID: value[1],
						}); err != nil {
							return err
						}
					}
				}
			}

			//条件写入
			var i uint64
			for _, cond := range n.ConditionList {
				i = i + 1
				for _, c := range cond {
					var item = FlowActivityData{
						AppID:              appId,
						FlowID:             flowId,
						FlowActivityID:     act.ID,
						DatasourceTableID:  n.TargetId,
						DatasourceColumnID: c.OriginId,
						Type:               COND_TYPE,
						Expression:         c.Condition,
						GroupID:            i,
					}
					switch c.TargetFormType {
					case 1:
						item.InputNodeID = c.TargetFormId
						item.InputParamID = c.TargetId
					case 2:
						item.InputFieldID = c.TargetId
					case 3:
						item.InputText = c.TargetValue
					case 4:
						item.InputFieldTableColumnID = c.TargetId
					}

					if err, _ := createActivityData(srv.db, item); err != nil {
						return err
					}
				}

			}

			if err, _ := updateActivity(srv.db, FlowActivity{
				ID:   act.ID,
				Type: DATA_ACT,
			}); err != nil {
				return err
			}
		case 3: // 代码块节点
			var n CodeNode
			if err := json.Unmarshal(acByte, &n); err != nil {
				return err
			}
			for _, c := range n.Columns {
				value := make([]uint64, 0, 10)
				for k, v := range c {
					if v.Type != 3 {
						b, err := json.Marshal(v.Value)
						if err != nil {
							return err
						}
						if err := json.Unmarshal(b, &value); err != nil {
							return err
						}
					}
					switch v.Type {
					case 1:
						if err, _ := createActivityService(srv.db, FlowActivityService{
							AppID:          appId,
							FlowID:         flowId,
							FlowActivityID: act.ID,
							ParamID:        k,
							InputNodeID:    value[0],
							InputParamID:   value[1],
						}); err != nil {
							return err
						}
					case 2:
						if err, _ := createActivityService(srv.db, FlowActivityService{
							AppID:          appId,
							FlowID:         flowId,
							FlowActivityID: act.ID,
							ParamID:        k,
							InputFieldID:   value[1],
						}); err != nil {
							return err
						}
					case 3:
						if err, _ := createActivityService(srv.db, FlowActivityService{
							AppID:          appId,
							FlowID:         flowId,
							FlowActivityID: act.ID,
							ParamID:        k,
							InputText:      fmt.Sprintf(`%v`, v.Value),
						}); err != nil {
							return err
						}
					case 4:
						if err, _ := createActivityService(srv.db, FlowActivityService{
							AppID:                   appId,
							FlowID:                  flowId,
							FlowActivityID:          act.ID,
							ParamID:                 k,
							InputFieldTableColumnID: value[1],
						}); err != nil {
							return err
						}
					}
				}
			}
			if err, _ := updateActivity(srv.db, FlowActivity{
				ID:        act.ID,
				Type:      SERVICE,
				ServiceID: n.TargetId, //代码包id
			}); err != nil {
				return err
			}
		case 7: // 条件节点
			for _, l := range x.List {
				con, err := json.Marshal(l)
				if err != nil {
					return err
				}
				if err := srv.flowDecode(flowId, appId, con, chainUp, email); err != nil {
					return err
				}
			}
			if err, _ := updateActivity(srv.db, FlowActivity{
				ID:   act.ID,
				Type: BRANCH,
			}); err != nil {
				return err
			}
		case 8:
			var n ConditionList
			if err := json.Unmarshal(acByte, &n); err != nil {
				return err
			}
			var i uint64
			for _, cond := range n.ConditionList {
				i = i + 1
				for _, c := range cond {
					var item FlowActivityGateway
					item.FlowID = flowId
					item.FlowActivityID = act.ID
					item.NodeID = x.Key
					item.Expression = c.Condition
					item.GroupID = i
					switch c.OriginFormType {
					case 1:
						item.LeftInputParamID = c.OriginId
					case 2:
						item.LeftInputFieldID = c.OriginId
					case 4:
						item.LeftInputFieldTableColumnID = c.OriginId
					}
					switch c.TargetFormType {
					case 1:
						item.RightInputParamID = c.TargetId
					case 2:
						item.RightInputFieldID = c.TargetId
					case 3:
						item.RightInputText = c.TargetValue
					case 4:
						item.RightInputFieldTableColumnID = c.TargetId
					}
					if err, _ := createActivityGateway(srv.db, item); err != nil {
						return err
					}
				}
			}
		case 9:
			if err, _ := updateActivity(srv.db, FlowActivity{
				ID:   act.ID,
				Type: END,
			}); err != nil {
				return err
			}
		}
	}
	return nil
}
