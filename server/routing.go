package server

import (
	"github.com/Lenora-Z/low-code/service/routing"
	"github.com/gin-gonic/gin"
)

type RouteItem struct {
	routing.Router
	Child []RouteItem `json:"child"`
}

// @Summary 路由列表(树形)
// @Tags route
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param app-id header string true "应用id"
// @Success 200 {object} ApiResponse{result=[]RouteItem}
// @Router /route/list [get]
func (ds *defaultServer) RouteList(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*CustomClaims)
	result := make([]RouterItem, 0)

	rs := routing.NewRoutingService(ds.db)
	list := rs.GetRouterListService(claims.AppId)

	if len(list) > 0 {
		for _, v := range list {
			//寻找孩子节点
			if v.ParentID == 0 {
				item := RouterItem{v, []RouterItem{}}
				item.getChildren(list)
				result = append(result, item)
			}
		}
	}
	ds.ResponseSuccess(ctx, result)
	return
}

type RouterItem struct {
	routing.Router              //路由信息
	Children       []RouterItem `json:"children"` //子路由
}

func (this *RouterItem) getChildren(group routing.RouterList) {
	for _, v := range group {
		if v.ParentID == this.ID {
			this.Children = append(this.Children, RouterItem{v, []RouterItem{}})
		}
	}
	if len(this.Children) > 0 {
		for i, _ := range this.Children {
			this.Children[i].getChildren(group)
		}
	}
	return
}
