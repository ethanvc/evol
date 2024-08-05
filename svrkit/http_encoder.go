package svrkit

import (
	"context"
	"encoding/json"
	"github.com/ethanvc/evol/obs"

	"github.com/ethanvc/evol/base"
	"google.golang.org/grpc/codes"
)

type HttpEncoder struct{}

func NewHttpEncoder() *HttpEncoder {
	return &HttpEncoder{}
}

func (e *HttpEncoder) Intercept(c context.Context, req any, nexter Nexter) (any, error) {
	resp, err := nexter.Next(c, req)
	httpReq := GetHttpRequestContext(c)
	if httpReq == nil {
		return nil, base.New(codes.Internal).SetEvent("HttpEncoderMustUsedWithHttpProtocol")
	}
	statusErr := e.convertToStatus(c, err)
	var httpResp HttpResponse
	httpResp.Code = statusErr.GetCode()
	httpResp.Msg = statusErr.GetMsg()
	httpResp.Data = resp
	content, err := json.Marshal(httpResp)
	if err != nil {
		return nil, base.New(codes.Internal)
	}
	GetAccessInfo(c).SetResp(content)
	httpReq.Writer.Header().Set("Content-Type", "application/json")
	n, err := httpReq.Writer.Write(content)
	if err != nil {
		return nil, base.New(codes.Unknown).SetErrAsEvent(err)
	}
	if n != len(content) {
		return nil, base.New(codes.Unknown).SetEvent("ContentLenError")
	}
	return resp, nil
}

func (e *HttpEncoder) convertToStatus(c context.Context, err error) *base.Status {
	if err == nil {
		return nil
	}
	if s, ok := err.(*base.Status); ok {
		return s
	}
	return obs.New(c, err).Status()
}

type HttpResponse struct {
	Code codes.Code `json:"code"`
	Msg  string     `json:"msg"`
	Data any        `json:"data"`
}
