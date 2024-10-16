package svrkit

import (
	"context"
	"encoding/json"
	"github.com/ethanvc/evol/base"
	"github.com/ethanvc/evol/xlog"
	"google.golang.org/grpc/codes"
	"log/slog"
)

type HttpEncoder struct{}

func NewHttpEncoder() *HttpEncoder {
	return &HttpEncoder{}
}

func (e *HttpEncoder) Intercept(c context.Context, req any, nexter Nexter) (any, error) {
	resp, businessErr := nexter.Next(c, req)
	httpReq := GetHttpRequestContext(c)
	if httpReq == nil {
		return nil, base.New(codes.Internal).SetEvent("HttpEncoderMustUsedWithHttpProtocol")
	}
	accInfo := GetAccessInfo(c)
	statusErr := e.convertToStatus(c, businessErr)
	var httpResp HttpResponse
	httpResp.Code = statusErr.GetCode()
	httpResp.Msg = statusErr.GetMsg()
	httpResp.Data = resp
	content, err := json.Marshal(httpResp)
	if err != nil {
		return nil, base.New(codes.Internal)
	}
	accInfo.SetResp(resp)
	httpReq.Writer.Header().Set("Content-Type", "application/json")
	n, err := httpReq.Writer.Write(content)
	if err != nil {
		return nil, base.New(codes.Unknown).SetErrEvent(err)
	}
	if n != len(content) {
		return nil, base.New(codes.Unknown).SetEvent("ContentLenError")
	}
	xlog.GetLogContext(c).SetAttribute(
		xlog.AttributeKeyHttpResponseHeaders,
		slog.AnyValue(httpReq.Writer.Header()))
	return resp, businessErr
}

func (e *HttpEncoder) convertToStatus(c context.Context, err error) *base.Status {
	if err == nil {
		return nil
	}
	if s, ok := err.(*base.Status); ok {
		return s
	}
	return base.New(c, err).Status()
}

type HttpResponse struct {
	Code codes.Code `json:"code"`
	Msg  string     `json:"msg"`
	Data any        `json:"data"`
}
