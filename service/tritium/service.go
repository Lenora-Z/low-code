package tritium

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Lenora-Z/low-code/utils"
	"github.com/sirupsen/logrus"
	"time"
	"unicode/utf8"
)

type TritiumService interface {
	// BatchGetTriUserDetail 批量获取用户详情
	BatchGetTriUserDetail(group []uint64) (error, []UserDetail)
	//用户详情
	GetTriUserDetail(uid uint64) (error, *UserDetail)
	//检验模块权限
	CheckModule(uid uint64) (bool, error)
	//发送日志
	SendLog(arg *LogArg) error
	//提交路由
	ResourceSubmit(path string)
	//拥有权限的账户列表
	GetTriCheckModule(module string) (error, []int)
	//获取当前机构名称
	GetTritiumConf() (error, string)
	//发送消息
	SendNotice(arg *NoticeArg) error
}

type tritiumService struct {
	api string
}

type ResponseCommon struct {
	Code    int8   `json:"code"`
	Message string `json:"message"`
}

type UserDetail struct {
	Id           uint64 `json:"id"`
	TrueName     string `json:"true_name"`
	Mobile       string `json:"mobile"`
	DepartmentID uint64 `json:"department_id"`
}

type LogArg struct {
	UId     uint64
	Action  string
	Message string
	Ip      string
	Hash    string
	Status  uint8
}

type NoticeArg struct {
	Title      string   `json:"title"`
	Message    string   `json:"message"`
	UserIds    []string `json:"userIds"`
	Module     string   `json:"module"`
	NotifyType int8     `json:"notify_type"` // 通知方式,0:系统 1:短信
	BaseUrl    string   `json:"base_url"`    //为空默认写入本地平台氚
	//Extra      ExtraData `json:"extra"`
}

func NewTritiumService(api string) TritiumService {
	s := new(tritiumService)
	s.api = "http://" + api
	return s
}

func (srv *tritiumService) CheckModule(uid uint64) (bool, error) {
	url := srv.api + CHECK_MODULE
	param := fmt.Sprintf(`{"userIds":["%d"], "platform": "%s", "module": "%s"}`, uid, PLATFORM, "cms")

	r, err := utils.SendRequest(&utils.HeaderRequest{
		Method: "POST",
		Url:    url,
		Body:   param,
	})

	if err != nil {
		return false, err
	}

	type response struct {
		ResponseCommon
		Result []map[string]string
	}

	var result response
	err = json.Unmarshal([]byte(r), &result)
	if err != nil {
		logrus.Error("middleware => entrance => Unmarshal: ", err)
		return false, err
	}
	if result.Code == 1 {
		logrus.Error("middleware => entrance => GetTriCheckModule: ", result)
		return false, errors.New(result.Message)
	}
	for _, value := range result.Result[0] {
		if value != "1" {
			return false, errors.New("Insufficient user rights")
		}
	}
	return true, nil

}

func (srv *tritiumService) GetTriCheckModule(module string) (error, []int) {
	url := srv.api + MODULE_USER
	param := fmt.Sprintf(`{"platform":"%s","module":"%s"}`, PLATFORM, module)
	r, err := utils.SendRequest(&utils.HeaderRequest{
		Method: "POST",
		Url:    url,
		Body:   param,
	})
	if err != nil {
		return err, nil
	}

	type res struct {
		List []int
	}

	type ResBody struct {
		ResponseCommon
		Result res
	}
	var result ResBody
	err = json.Unmarshal([]byte(r), &result)
	if err != nil {
		logrus.Error("common => SendUserLogs => Unmarshal: ", err)
		return err, nil
	}

	if result.Code != 0 {
		logrus.Error("common =>SendUserLogs => WrongCode:", result.Code, "---", result.Message)
		return errors.New("WrongCode"), nil
	}
	return nil, result.Result.List

}

func (srv *tritiumService) GetTritiumConf() (error, string) {
	url := srv.api + CONF

	r, err := utils.SendRequest(&utils.HeaderRequest{
		Method: "POST",
		Url:    url,
		Body:   "",
	})
	if err != nil {
		logrus.Error("main => InitSetting => getTritiumName => wrong:", err)
		return err, ""
	}

	logrus.Info(r)
	type ResBody struct {
		ResponseCommon
		Result struct {
			NodeName string `json:"node_name"`
		}
	}
	var result ResBody
	err = json.Unmarshal([]byte(r), &result)
	if err != nil {
		logrus.Error("main => InitSetting => getTritiumName => Unmarshal: ", err)
	}

	return nil, result.Result.NodeName
}

func (srv *tritiumService) GetTriUserDetail(uid uint64) (error, *UserDetail) {
	url := srv.api + USER_DETAIL
	param := fmt.Sprintf(`{"userId":"%d"}`, uid)

	r, err := utils.SendRequest(&utils.HeaderRequest{
		Method: "POST",
		Url:    url,
		Body:   param,
	})

	if err != nil {
		return err, nil
	}
	type response struct {
		ResponseCommon
		Result UserDetail
	}

	var result response
	err = json.Unmarshal([]byte(r), &result)
	if err != nil {
		logrus.Error("middleware => entrance => Unmarshal: ", err)
		return err, nil
	}

	if result.Code == 1 {
		logrus.Error("middleware => entrance => GetTriUserDetail: ", result)
		return errors.New(result.Message), nil
	}

	return nil, &result.Result

}

func (srv *tritiumService) BatchGetTriUserDetail(group []uint64) (error, []UserDetail) {
	param := struct {
		UserIds []uint64 `json:"user_ids"`
	}{
		group,
	}
	byt, err := json.Marshal(param)
	if err != nil {
		return err, nil
	}
	url := srv.api + GROUP_USER_DETAIL
	r, err := utils.SendRequest(&utils.HeaderRequest{
		Method: "POST",
		Url:    url,
		Body:   string(byt),
	})
	if err != nil {
		return err, nil
	}
	type response struct {
		ResponseCommon
		Result []UserDetail `json:"result"`
	}
	var result response
	err = json.Unmarshal([]byte(r), &result)
	if err != nil {
		logrus.Error("middleware => entrance => Unmarshal: ", err)
		return err, nil
	}

	if result.Code == 1 {
		logrus.Error("middleware => entrance => GetTriUserDetail: ", result)
		return errors.New(result.Message), nil
	}

	return nil, result.Result
}

func (srv *tritiumService) SendLog(arg *LogArg) error {
	url := srv.api + SEND_LOG
	param := fmt.Sprintf(`{"message":"%s","time":"%d","user_id":"%d","ip":"%s","platform":"%s","action_name":"%s"}`, arg.Message, time.Now().Unix(), arg.UId, arg.Ip, PLATFORM, arg.Action)

	r, err := utils.SendRequest(&utils.HeaderRequest{
		Method: "POST",
		Url:    url,
		Body:   param,
	})
	if err != nil {
		return err
	}

	type ResBody struct {
		ResponseCommon
		Result interface{}
	}
	var result ResBody
	err = json.Unmarshal([]byte(r), &result)
	if err != nil {
		logrus.Error("common => SendUserLogs => Unmarshal: ", err)
		return err
	}

	if result.Code != 0 {
		logrus.Error("common =>SendUserLogs => WrongCode:", result.Code, "---", result.Message)
		return errors.New("WrongCode")
	}
	return nil
}

func (srv *tritiumService) SendNotice(arg *NoticeArg) error {
	var url string
	var userIds []string
	if arg.BaseUrl == "" {
		url = srv.api + SEND_NOTICE
		err, users := srv.GetTriCheckModule("ims")
		if err != nil {
			return err
		}
		for _, id := range users {
			userIds = append(userIds, fmt.Sprintf("%d", id))
		}
	} else {
		url = arg.BaseUrl + "/ims/notice/submit"
	}

	if utf8.RuneCountInString(arg.Title) > 30 {
		r := []rune(arg.Title)
		arg.Title = fmt.Sprintf("%s...", string(r[:27]))
	}

	param := fmt.Sprintf(`{"message":"%s","userIds":"%s","time":"%d","platform":"%s","module":"%s","title":"%s","notify_type":"%d"}`, arg.Message, userIds, time.Now().Unix(), PLATFORM, arg.Module, arg.Title, arg.NotifyType)

	r, err := utils.SendRequest(&utils.HeaderRequest{
		Method: "POST",
		Url:    url,
		Body:   param,
	})
	if err != nil {
		return err
	}

	type ResBody struct {
		ResponseCommon
		Result interface{}
	}
	var result ResBody
	err = json.Unmarshal([]byte(r), &result)
	if err != nil {
		logrus.Error("common => SendNotice => Unmarshal: ", err)
		return err
	}

	if result.Code != 0 {
		logrus.Error("common =>SendNotice => WrongCode:", result.Code, "---", result.Message)
		return errors.New("WrongCode")
	}
	return nil

}

func (srv *tritiumService) ResourceSubmit(path string) {
	requestBody, err := utils.ReadFile(path)
	if err != nil {
		logrus.Error("read file failed:", err)
		return
	}
	logrus.Info("route register:", requestBody)
	ret, err := utils.SendRequest(&utils.HeaderRequest{
		Method: "POST",
		Url:    srv.api + RESOURCE_SUBMIT,
		Body:   requestBody,
	})
	if err != nil {
		logrus.Error("initTritiumRoute--Wrong:", err)
		return
	}
	logrus.Info("initTritiumRoute--:", ret)
	return
}
