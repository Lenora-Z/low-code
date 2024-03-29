// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: service/sms/proto/sms.proto

package proto

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "github.com/asim/go-micro/v3/api"
	client "github.com/asim/go-micro/v3/client"
	server "github.com/asim/go-micro/v3/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for Sms service

func NewSmsEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Sms service

type SmsService interface {
	SmsInValid(ctx context.Context, in *SmsInValidRequest, opts ...client.CallOption) (*SmsInvalidResponse, error)
	SendCode(ctx context.Context, in *SendCodeRequest, opts ...client.CallOption) (*SendCodeResponse, error)
	VerifyCodeUnChange(ctx context.Context, in *VerifyCodeUnChangeRequest, opts ...client.CallOption) (*VerifyCodeUnChangeResponse, error)
	VerifyCode(ctx context.Context, in *VerifyCodeRequest, opts ...client.CallOption) (*VerifyCodeResponse, error)
}

type smsService struct {
	c    client.Client
	name string
}

func NewSmsService(name string, c client.Client) SmsService {
	return &smsService{
		c:    c,
		name: name,
	}
}

func (c *smsService) SmsInValid(ctx context.Context, in *SmsInValidRequest, opts ...client.CallOption) (*SmsInvalidResponse, error) {
	req := c.c.NewRequest(c.name, "Sms.SmsInValid", in)
	out := new(SmsInvalidResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *smsService) SendCode(ctx context.Context, in *SendCodeRequest, opts ...client.CallOption) (*SendCodeResponse, error) {
	req := c.c.NewRequest(c.name, "Sms.SendCode", in)
	out := new(SendCodeResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *smsService) VerifyCodeUnChange(ctx context.Context, in *VerifyCodeUnChangeRequest, opts ...client.CallOption) (*VerifyCodeUnChangeResponse, error) {
	req := c.c.NewRequest(c.name, "Sms.VerifyCodeUnChange", in)
	out := new(VerifyCodeUnChangeResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *smsService) VerifyCode(ctx context.Context, in *VerifyCodeRequest, opts ...client.CallOption) (*VerifyCodeResponse, error) {
	req := c.c.NewRequest(c.name, "Sms.VerifyCode", in)
	out := new(VerifyCodeResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Sms service

type SmsHandler interface {
	SmsInValid(context.Context, *SmsInValidRequest, *SmsInvalidResponse) error
	SendCode(context.Context, *SendCodeRequest, *SendCodeResponse) error
	VerifyCodeUnChange(context.Context, *VerifyCodeUnChangeRequest, *VerifyCodeUnChangeResponse) error
	VerifyCode(context.Context, *VerifyCodeRequest, *VerifyCodeResponse) error
}

func RegisterSmsHandler(s server.Server, hdlr SmsHandler, opts ...server.HandlerOption) error {
	type sms interface {
		SmsInValid(ctx context.Context, in *SmsInValidRequest, out *SmsInvalidResponse) error
		SendCode(ctx context.Context, in *SendCodeRequest, out *SendCodeResponse) error
		VerifyCodeUnChange(ctx context.Context, in *VerifyCodeUnChangeRequest, out *VerifyCodeUnChangeResponse) error
		VerifyCode(ctx context.Context, in *VerifyCodeRequest, out *VerifyCodeResponse) error
	}
	type Sms struct {
		sms
	}
	h := &smsHandler{hdlr}
	return s.Handle(s.NewHandler(&Sms{h}, opts...))
}

type smsHandler struct {
	SmsHandler
}

func (h *smsHandler) SmsInValid(ctx context.Context, in *SmsInValidRequest, out *SmsInvalidResponse) error {
	return h.SmsHandler.SmsInValid(ctx, in, out)
}

func (h *smsHandler) SendCode(ctx context.Context, in *SendCodeRequest, out *SendCodeResponse) error {
	return h.SmsHandler.SendCode(ctx, in, out)
}

func (h *smsHandler) VerifyCodeUnChange(ctx context.Context, in *VerifyCodeUnChangeRequest, out *VerifyCodeUnChangeResponse) error {
	return h.SmsHandler.VerifyCodeUnChange(ctx, in, out)
}

func (h *smsHandler) VerifyCode(ctx context.Context, in *VerifyCodeRequest, out *VerifyCodeResponse) error {
	return h.SmsHandler.VerifyCode(ctx, in, out)
}
