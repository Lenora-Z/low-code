package server

import (
	"github.com/gin-gonic/gin"
)

type MobileArg struct {
	Mobile  string `json:"mobile"`
	UseType int64  `json:"useType"`
	Content string `json:"content"`
}

// @Summary 发送短信
// @Description 调用短信业务服务
// @Tags client
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param MobileArg body MobileArg true "请求体"
// @Success 200 {object} ApiResponse{result=object}
// @Router /api/sms/send [post]
func (ds *defaultServer) SendSms(ctx *gin.Context) {
	//args := new(MobileArg)
	//if err := ctx.BindJSON(args); err != nil {
	//	ds.InvalidParametersError(ctx)
	//	return
	//}
	//
	//smsClient := sms.NewSmsService(ds.rpcSmsClient)
	//resp, err := smsClient.SendSms(args.Mobile, args.UseType)
	//if err != nil {
	//	logrus.Error("get rpc response err: ", err.Error())
	//	ds.InternalServiceError(ctx, err.Error())
	//	return
	//}
	//
	//ds.ResponseSuccess(ctx, resp)
}
