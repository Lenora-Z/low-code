//Created by Goland
//@User: lenora
//@Date: 2021/3/8
//@Time: 10:52 上午
package server

const (
	AppTitle   = "低代码平台"
	AppVersion = "v1.4"
	AppName    = "low-code"
	NickName   = "lowCode-backend"
)

const CLAIMS = "claims"

var TRUE, FALSE = true, false

type IdStruct struct {
	Id uint64 `json:"id" validate:"required"`
}

type CustomClaims struct {
	Name      string `json:"name"`
	AppId     uint64 `json:"app_id"`     //应用d
	UserId    uint64 `json:"user_id"`    //用户
	UUid      string `json:"user_uuid"`  //uuid
	AppStatus int8   `json:"app_status"` //应用状态
}

type UserClaims struct {
	AppId   uint64 `json:"app_id"` //应用id
	UserId  uint64 `json:"id"`
	Account string `json:"account"` //账号
	GroupId uint64 `json:"group_id"`
	AppName string `json:"app_name"`
}
