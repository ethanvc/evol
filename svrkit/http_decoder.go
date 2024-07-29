package svrkit

import (
	"context"
	"encoding/json"
	"github.com/ethanvc/evol/base"
	"github.com/ethanvc/evol/xlog"
	"google.golang.org/grpc/codes"
	"io"
)

type HttpDecoder struct {
}

func NewHttpDecoder() *HttpDecoder {
	return &HttpDecoder{}
}

func (decoder *HttpDecoder) Intercept(c context.Context, req any, nexter Nexter) (any, error) {
	httpReq := GetHttpRequestContext(c)
	if httpReq == nil {
		return nil, base.New(codes.Internal, "HttpDecoderMustMusedWithHttpProtocol")
	}
	if nexter.Chain().NewReq == nil {
		return nexter.Next(c, req)
	}
	limiterR := &io.LimitedReader{
		R: httpReq.Request.Body,
		N: 1024 * 1024 * 2,
	}
	content, err := io.ReadAll(limiterR)
	if err != nil {
		return nil, err
	}
	GetAccessInfo(c).SetReq(content)
	req = nexter.Chain().NewReq()
	err = json.Unmarshal(content, req)
	if err != nil {
		return nil, xlog.New(c, err).Error()
	}
	return nexter.Next(c, req)
}
