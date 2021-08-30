package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/Lenora-Z/low-code/conf"
	"github.com/Lenora-Z/low-code/service/data_db"
	"github.com/Lenora-Z/low-code/service/email"
	"github.com/Lenora-Z/low-code/service/form"
	"github.com/Lenora-Z/low-code/service/form_data"
	"github.com/Lenora-Z/low-code/service/source"
	"github.com/Lenora-Z/low-code/utils"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

type PublicFormDataArg struct {
	FormNo string `json:"form_no" validate:"required"` //表单编号
}

type DataCreateArg struct {
	TableName string                 `json:"table_name" validate:"required"` //表名
	Data      map[string]interface{} `json:"data"`                           //数据对象  传递object,key为字段名称(string) value为数据值(泛型)
}

type DataUpdateArg struct {
	DataCreateArg
	Condition []data_db.ConditionGroup `json:"condition"` //筛选条件
}

type SignArg struct {
	AccessKey string `json:"access_key"` // access_key
	SecretKey string `json:"secret_key"` //secret_key
	Method    string `json:"method"`     //请求方式
	Path      string `json:"path"`       //请求地址
}

type SignVO struct {
	Date      string `json:"date"`      //日期
	Rand      int    `json:"rand"`      //随机数
	Signature string `json:"signature"` //签名
}

type EmailVo struct {
	Sender   string   `json:"sender" validate:"required"`
	Receiver []string `json:"receiver"`
	Subject  string   `json:"subject" validate:"required"`
	Content  string   `json:"content" validate:"required"`
}

func (ds *defaultServer) publicFormDataList(ctx *gin.Context) {
	claims := ctx.MustGet(CLAIMS).(*CustomClaims)
	var arg PublicFormDataArg
	if err := ctx.BindJSON(&arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}
	if err := validatorInstance().Struct(arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	srv := form.NewFormService(ds.db)
	notFound, item := srv.GetFormDetailByNo(arg.FormNo, claims.AppId)
	if notFound {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}

	dataSrv := form_data.NewFormDataService(ds.mongoDb)
	err, _, list := dataSrv.GetListByFormId(1, utils.MAX_LIMIT, item.ID)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	ds.ResponseSuccess(ctx, list)
	return
}

// CreateData
// @Summary 新增数据
// @Tags public
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param DataCreateArg body DataCreateArg true "请求体"
// @Success 200 {object} ApiResponse{result=object}
// @Router /public/data/new [post]
func (ds *defaultServer) CreateData(ctx *gin.Context) {
	var arg DataCreateArg
	if err := ctx.BindJSON(&arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}
	if err := validatorInstance().Struct(arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	if len(arg.Data) <= 0 {
		ds.InvalidParametersError(ctx)
		return
	}

	srv := source.NewService(ds.db)
	notFound, tb := srv.GetTableByName(arg.TableName, 0)
	if notFound {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}

	dataSrv := data_db.NewService(ds.dataDb, tb.TableName)
	if err := dataSrv.NewItem(arg.Data); err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	ds.ResponseSuccess(ctx, nil)
	return
}

// UpdateData
// @Summary 修改数据
// @Tags public
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param DataUpdateArg body DataUpdateArg true "请求体"
// @Success 200 {object} ApiResponse{result=object}
// @Router /public/data/edit [post]
func (ds *defaultServer) UpdateData(ctx *gin.Context) {
	var arg DataUpdateArg
	if err := ctx.BindJSON(&arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}
	if err := validatorInstance().Struct(arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	if len(arg.Data) <= 0 {
		ds.InvalidParametersError(ctx)
		return
	}

	srv := source.NewService(ds.db)
	notFound, tb := srv.GetTableByName(arg.TableName, 0)
	if notFound {
		ds.ResponseError(ctx, NOT_FOUND)
		return
	}

	dataSrv := data_db.NewService(ds.dataDb, tb.TableName)
	if err := dataSrv.UpdateData(arg.Data, arg.Condition); err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}
	ds.ResponseSuccess(ctx, nil)
	return
}

// Signature
// @Summary 生成签名
// @Tags public
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param SignArg body SignArg true "请求体"
// @Success 200 {object} ApiResponse{result=SignVO}
// @Router /public/package/sign [post]
func (ds *defaultServer) Signature(ctx *gin.Context) {
	var arg SignArg
	if err := ctx.BindJSON(&arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}
	if err := validatorInstance().Struct(arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	loc, _ := time.LoadLocation("GMT")
	date := time.Now().In(loc).Format(time.RFC1123)
	hash := hmac.New(sha256.New, []byte(arg.SecretKey)) // 创建对应的sha256哈希加密算法
	rand := utils.RandInt(10000, 99999)
	message := []string{
		strings.ToUpper(arg.Method),
		arg.Path,
		"",
		arg.AccessKey,
		date,
		fmt.Sprintf("hmacrand:%d", rand),
	}
	signedMessage := strings.Join(message, "\n")
	signedMessage += "\n"
	hash.Write([]byte(signedMessage))
	sign := base64.StdEncoding.EncodeToString(hash.Sum(nil))

	ds.ResponseSuccess(ctx, SignVO{
		Date:      date,
		Rand:      rand,
		Signature: sign,
	})
	return
}

// SendEmail
// @Summary 发送邮件
// @Tags public
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param EmailVo body EmailVo true "请求体"
// @Success 200 {object} ApiResponse{result=object}
// @Router /public/package/email [post]
func (ds *defaultServer) SendEmail(ctx *gin.Context) {
	var arg EmailVo
	if err := ctx.BindJSON(&arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}
	if err := validatorInstance().Struct(arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	var config *conf.SmtpConfig
	for _, x := range ds.conf.SmtpGroupConf {
		if x.Sender == arg.Sender {
			config = &x
			break
		}
	}
	if config == nil {
		ds.InternalServiceError(ctx, "wrong sender")
		return
	}
	srv := email.NewEmailService(config.Sender, config.SmtpAddr, config.Password, config.SmtpPort)
	if err := srv.SendMail(arg.Subject, arg.Content, arg.Receiver); err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	ds.ResponseSuccess(ctx, nil)
	return
}

// ContractFiling
// @Summary 文件归档
// @Tags public
// @Accept  json
// @Produce  json
// @param lang header string false "语言(zh/en),默认en"
// @param IdStruct body IdStruct true "请求体"
// @Success 200 {object} ApiResponse{result=object}
// @Router /public/package/filing [post]
func (ds *defaultServer) ContractFiling(ctx *gin.Context) {
	var arg IdStruct
	if err := ctx.BindJSON(&arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}
	if err := validatorInstance().Struct(arg); err != nil {
		ds.InvalidParametersError(ctx)
		return
	}

	srv := data_db.NewService(ds.dataDb, "hrm_labor_contract")
	err, res := srv.GetItem(arg.Id)
	if err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	sum, _ := res["sum"].(int64)
	if sum != 0 {
		sum = sum - 1
	}
	if err := srv.ContractFiling(arg.Id, uint64(sum)); err != nil {
		ds.InternalServiceError(ctx, err.Error())
		return
	}

	ds.ResponseSuccess(ctx, nil)
	return

}
