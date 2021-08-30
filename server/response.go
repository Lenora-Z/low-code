//Created by Goland
//@User: lenora
//@Date: 2021/1/15
//@Time: 10:34 上午

package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ApiResponse struct {
	Code    int16       `json:"code"`    // 状态码
	Message string      `json:"message"` // 状态短语
	Result  interface{} `json:"result"`  // 数据结果集
}

const (
	SUCCESS int16 = 1000
	FAIL    int16 = 4000

	ACCESS_DENY int16 = 2003

	NOT_FOUND int16 = 3001

	WRONG_PARAM int16 = 4001
)

var WrongMessageEn = map[int16]string{
	1000: "success",

	2001: "user licence expired",
	2002: "user wrong",
	2003: "permission denied",
	2004: "not support",

	3001: "record not found",
	3002: "record has been exist",
	3003: "record can't be changed",
	3101: "there is already a valid navigation",

	4000: "fail",
	4001: "param error",
}

var WrongMessageZh = map[int16]string{
	1000: "请求成功",

	2001: "用户凭证过期",
	2002: "用户错误",
	2003: "权限验证失败",
	2004: "不支持该操作",

	3001: "记录未找到",
	3002: "记录已存在",
	3003: "记录禁止修改",
	3101: "一个应用下只能存在一个生效的导航栏",

	4000: "请求失败",
	4001: "参数错误",
}

func (ds *defaultServer) InvalidParametersError(c *gin.Context) {
	responseError(c, WRONG_PARAM, nil, ds.getLanguage(c))
}

func (ds *defaultServer) InternalServiceError(c *gin.Context, message ...string) {
	responseError(c, FAIL, message, ds.getLanguage(c))
}

func (ds *defaultServer) ResponseError(c *gin.Context, code int16, message ...string) {
	responseError(c, code, message, ds.getLanguage(c))
}

func (ds *defaultServer) ResponseSuccess(c *gin.Context, result interface{}, msg ...string) {
	if len(msg) == 0 {
		lang := ds.getLanguage(c)
		msg = append(msg, getResponseMsgWithLang(SUCCESS, lang))
	}
	responseOutput(c, SUCCESS, msg[0], result)
}

func responseError(c *gin.Context, code int16, message []string, lang string) {
	var msg string
	if len(message) == 0 {
		msg = getResponseMsgWithLang(code, lang)
	} else {
		msg = message[0]
	}
	responseOutput(c, code, msg, nil)
}

// @Summary 返回码
// @Tags response
// @Accept  json
// @Produce  json
// @Success 1000 {object} ApiResponse{result=object} "请求成功"
// @Failure 2001 {string} ApiResponse{} "用户凭证过期"
// @Failure 2002 {string} ApiResponse{} "用户错误"
// @Failure 2003 {string} ApiResponse{} "权限验证失败"
// @Failure 2004 {string} ApiResponse{} "不支持该操作"
// @Failure 3001 {string} ApiResponse{} "记录未找到"
// @Failure 3002 {string} ApiResponse{} "记录已存在"
// @Failure 3003 {string} ApiResponse{} "记录禁止修改"
// @Failure 3101 {string} ApiResponse{} "一个应用下只能存在一个生效的导航栏"
// @Failure 4000 {string} ApiResponse{} "请求失败"
// @Failure 4001 {string} ApiResponse{} "参数错误"
// @Router / [get]
func responseOutput(c *gin.Context, code int16, message string, result interface{}) {
	if result == nil {
		result = ""
	}
	c.JSON(http.StatusOK, ApiResponse{
		Code:    code,
		Message: message,
		Result:  result,
	})
	return
}

func getResponseMsgWithLang(code int16, lang string) string {
	var msg string
	switch lang {
	case "zh":
		msg = WrongMessageZh[code]
	default:
		msg = WrongMessageEn[code]
	}
	return msg
}
