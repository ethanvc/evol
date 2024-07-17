package svrkit

import (
	"context"
	"github.com/ethanvc/evol/xlog"
	"sync"
)

// AccessInterceptor please make me the first interceptor to cover all report and monitor scene.
type AccessInterceptor struct {
}

func NewAccessInterceptor() *AccessInterceptor {
	return &AccessInterceptor{}
}

func (interceptor *AccessInterceptor) Intercept(c context.Context, req any, nexter Nexter) (any, error) {
	info, c := WithAccessInfo(c)
	resp, err := nexter.Next(c, nil)
	xlog.GetAccessLogger().LogAccess(c, 0, err, info.req, info.resp)
	return resp, err
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
	mux   sync.Mutex
	req   any
	resp  any
	extra map[string]any
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

func (info *AccessInfo) SetExtra(key string, value any) {
	if info == nil || key == "" || value == nil {
		return
	}
	info.mux.Lock()
	defer info.mux.Unlock()
	const maxAllowedExtra = 10
	if len(info.extra) > maxAllowedExtra {
		return
	}
	info.extra[key] = value
}
