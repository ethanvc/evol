package svrkit

import (
	"context"
	"sync"

	"github.com/ethanvc/evol/xlog"
)

// AccessInterceptor please make me the first interceptor to cover all report and monitor scene.
type AccessInterceptor struct {
	GenerateContext func(c context.Context) context.Context
}

func NewAccessInterceptor() *AccessInterceptor {
	return &AccessInterceptor{}
}

func (interceptor *AccessInterceptor) Intercept(c context.Context, req any, nexter Nexter) (any, error) {
	c = interceptor.generateContext(c)
	info, c := WithAccessInfo(c)
	resp, err := nexter.Next(c, nil)
	xlog.GetAccessLogger().LogAccess(c, 0, err, info.req, info.resp)
	return resp, err
}

func (interceptor *AccessInterceptor) generateContext(c context.Context) context.Context {
	if interceptor.GenerateContext == nil {
		return c
	}
	return interceptor.GenerateContext(c)
}

type contextKeyAccessInfo struct{}

func WithAccessInfo(c context.Context) (*AccessInfo, context.Context) {
	info := &AccessInfo{}
	c = context.WithValue(c, contextKeyAccessInfo{}, info)
	return info, c
}

func GetAccessInfo(c context.Context) *AccessInfo {
	info, _ := c.Value(contextKeyAccessInfo{}).(*AccessInfo)
	return info
}

type AccessInfo struct {
	mux  sync.Mutex
	req  any
	resp any
}

func (info *AccessInfo) SetReq(req any) {
	if info == nil {
		return
	}
	info.mux.Lock()
	defer info.mux.Unlock()
	info.req = req
}

func (info *AccessInfo) SetResp(resp any) {
	if info == nil {
		return
	}
	info.mux.Lock()
	defer info.mux.Unlock()
	info.resp = resp
}

func GenerateHttpContext(c context.Context) context.Context {
	httpCtx := GetHttpRequestContext(c)
	if httpCtx == nil {
		return nil
	}
	return xlog.WithLogContext(c, httpCtx.Request.Method+" "+httpCtx.Pattern)
}
