package routing

import (
	"encoding/json"
	"fmt"
	"github.com/Lenora-Z/low-code/utils"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// 导航栏中的菜单项
type RoutingService interface {
	// 通过前端菜单json字符串,创建导航栏中的菜单项
	//CreateRouting(navId uint64, status bool, content RouterItemList) ([]byte, error)
	// 通过前端菜单json字符串,新增、更新、删除旧导航栏中的菜单项
	UpdateRouting(appId uint64, content RouterItemList) error
	//获取路由树
	GetRouterListService(appId uint64) RouterList
	//批量获取路由
	GetRouteGroup(ids []uint64) (error, RouterList)
}

type routingService struct {
	db *gorm.DB
}

func NewRoutingService(db *gorm.DB) RoutingService {
	u := new(routingService)
	u.db = db
	return u
}

type RouterItem struct {
	Id            uint64         `json:"id"`
	Title         string         `json:"title" validate:"required"` //路由名称
	Editor        bool           `json:"editor"`
	Key           string         `json:"key" validate:"required"` //key值
	FormID        uint64         `json:"formId"`                  //绑定的表单id
	ParentId      string         `json:"parentId"`
	Icon          string         `json:"icon"`
	OrderNum      uint32         `json:"order_num"`          //排序号
	Action        uint8          `json:"action_type"`        //事件类型
	ActionContent string         `json:"action_content"`     //时间配置
	Children      RouterItemList `json:"children,omitempty"` //子路由
}

type RouterItemList []RouterItem

func (srv *routingService) UpdateRouting(appId uint64, itemList RouterItemList) error {

	// 获取当前的菜单项
	currentRouterList, err := ListRouterWithAppId(srv.db, appId)
	if err != nil {
		return fmt.Errorf("ListRouterWithAppId: %s", err.Error())
	}

	// 更新提交上来的菜单项
	err = srv.updateRouting(appId, 0, itemList)
	if err != nil {
		return err
	}
	str, err := json.Marshal(itemList)
	if err != nil {
		return err
	}
	logrus.Info(string(str))

	// 删除“菜单项”已被删除部分
	// 没有覆盖到部分的菜单项
	if len(currentRouterList) > 0 {
		oldIds := currentRouterList.Ids()
		newIds := itemList.ids()
		deleteIds := utils.DifferenceUInt64(oldIds, newIds)
		if len(deleteIds) > 0 {
			err := DeleteWithIds(srv.db, deleteIds)
			if err != nil {
				return fmt.Errorf("DeleteWithIds:%s", err.Error())
			}
		}
	}

	return nil
}

func (srv *routingService) updateRouting(appId, parentId uint64, itemList RouterItemList) error {
	for k, item := range itemList {

		var router *Router
		var err error

		if item.Action < 1 || item.Action > 3 {
			logrus.Error("wrong action type:", item.Action)
			return errors.New("wrong action type")
		}
		if item.Id == 0 {
			// id为0或空, 新增条目,插入当前菜单项
			err, router = createRouter(srv.db, Router{
				AppID:         appId,
				Title:         item.Title,
				ParentID:      parentId,
				Key:           item.Key,
				FormID:        item.FormID,
				Icon:          item.Icon,
				Status:        &TRUE,
				OrderNum:      item.OrderNum,
				Action:        item.Action,
				ActionContent: item.ActionContent,
			})
		} else {
			// id不为0或空,更新条目,更新当前菜单项
			err, router = updateRouter(srv.db, Router{
				ID:            item.Id,
				AppID:         appId,
				Title:         item.Title,
				ParentID:      parentId,
				Key:           item.Key,
				FormID:        item.FormID,
				Icon:          item.Icon,
				OrderNum:      item.OrderNum,
				Action:        item.Action,
				ActionContent: item.ActionContent,
			})
		}

		if err != nil {
			return err
		}

		// 更新插入当前菜单项
		itemList[k].Id = router.ID

		// 插入下级菜单项
		if item.Children != nil {
			if err := srv.updateRouting(appId, router.ID, item.Children); err != nil {
				return err
			}
		}

	}
	return nil
}

func (srv *routingService) GetRouterListService(appId uint64) RouterList {
	column := make(map[string]interface{})
	column["status"] = 1
	if appId > 0 {
		column[RouterColumns.AppID] = appId
	}
	return GetRouterList(srv.db, column)
}

func (srv *routingService) GetRouteGroup(ids []uint64) (error, RouterList) {
	return GetRouterGroupWithIds(srv.db, ids)
}

func (ril RouterItemList) ids() []uint64 {
	ids := make([]uint64, 0, len(ril))
	for _, v := range ril {
		ids = append(ids, v.Id)
		if v.Children != nil {
			ids = append(ids, v.Children.ids()...)
		}
	}
	return ids
}
