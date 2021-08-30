package sms

import (
	"context"
	smsProto "github.com/Lenora-Z/low-code/service/sms/proto"
)

type SmsService interface {
	SendSms(mobile string, useType int64) (*smsProto.SendCodeResponse, error)
}

type smsService struct {
	rpcClient smsProto.SmsService
}

func NewSmsService(rpcClient smsProto.SmsService) SmsService {
	s := new(smsService)
	s.rpcClient = rpcClient
	return s
}

func (s *smsService) SendSms(mobile string, useType int64) (*smsProto.SendCodeResponse, error) {
	return s.rpcClient.SendCode(context.TODO(), &smsProto.SendCodeRequest{
		Mobile:  mobile,
		UseType: useType,
	})
}
