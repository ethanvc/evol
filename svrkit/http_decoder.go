package svrkit

import (
	"context"
	"encoding/json"
	"github.com/ethanvc/evol/xlog"
	"io"
	"log/slog"

	"github.com/ethanvc/evol/base"
	"google.golang.org/grpc/codes"
)

type HttpDecoder struct {
}

func NewHttpDecoder() *HttpDecoder {
	return &HttpDecoder{}
}

func (decoder *HttpDecoder) Intercept(c context.Context, req any, nexter Nexter) (any, error) {
	httpReq := GetHttpRequestContext(c)
	if httpReq == nil {
		return nil, base.New(codes.Internal).
			SetEvent("HttpDecoderMustMusedWithHttpProtocol")
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
	req = nexter.Chain().NewReq()
	if len(content) == 0 {
		return nexter.Next(c, req)
	}
	GetAccessInfo(c).SetReq(req)
	err = json.Unmarshal(content, req)
	if err != nil {
		GetAccessInfo(c).SetReq(string(content))
		return nil, base.NewComposer(c, err).Error()
	}
	resp, err := nexter.Next(c, req)
	lc := xlog.GetLogContext(c)
	lc.SetAttribute(xlog.AttributeKeyHttpRequestPath, slog.StringValue(httpReq.Request.URL.Path))
	lc.SetAttribute(xlog.AttributeKeyHttpRequestHost,
		slog.StringValue(httpReq.Request.Host))
	lc.SetAttribute(
		xlog.AttributeKeyHttpRequestHeaders,
		slog.AnyValue(httpReq.Request.Header))
	return resp, err
}
