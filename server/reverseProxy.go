package server

import (
	"github.com/Lenora-Z/low-code/utils"
	"github.com/gin-gonic/gin"
	"net/url"
)

// CallInterface
// @Summary 转发请求第三方接口
// @Description 请求方法这里写作get，实际是可以各种请求方法；第三方接口的参数正常设置
// @Tags interface
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param redirect query string true "第三方接口链接"
// @Success 200 {string} json ""
// @Router /api/interface/call [get]
func (ds *defaultServer) CallInterface(ctx *gin.Context) {
	destination := ctx.Query("redirect")
	if len(destination) == 0 {
		ds.InvalidParametersError(ctx)
		return
	}

	destUrl, err := url.Parse(destination)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	// 创建反向代理处理方法
	proxy := utils.NewReverseProxy(destUrl)
	proxy.ServeHTTP(ctx.Writer, ctx.Request)
}
