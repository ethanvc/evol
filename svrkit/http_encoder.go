package svrkit

import (
	"context"
	"encoding/json"
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
		return nil, base.New(codes.Internal, "HttpEncoderMustUsedWithHttpProtocol")
	}
	var httpResp HttpResponse
	httpResp.Data = resp
	content, err := json.Marshal(httpResp)
	if err != nil {
		return nil, base.New(codes.Internal, "MarshalHttpResponseError")
	}
	GetAccessInfo(c).SetResp(content)
	n, err := httpReq.Writer.Write(content)
	if err != nil {
		return nil, base.New(codes.Unknown, "WriteHttpResponseError")
	}
	if n != len(content) {
		return nil, base.New(codes.Unknown, "WriteContentPartial")
	}
	return resp, nil
}

type HttpResponse struct {
	Code codes.Code `json:"code"`
	Data any        `json:"data"`
}
