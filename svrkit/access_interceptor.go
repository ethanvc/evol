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
	xlog.GetAccessLogger().LogAccess(c, 0, err, info.req, info.resp, info.extra...)
	return resp, err
}

type contextKeyAccessInfo struct{}

func WithAccessInfo(c context.Context) (*AccessInfo, context.Context) {
	info := &AccessInfo{}
	c = context.WithValue(c, contextKeyAccessInfo{}, info)
	return info, c
}

type AccessInfo struct {
	mux   sync.Mutex
	req   any
	resp  any
	extra []any
}
