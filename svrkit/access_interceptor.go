package svrkit

import (
	"context"
	"github.com/ethanvc/evol/obs"
	"google.golang.org/grpc/codes"
)

type AccessInterceptor struct {
}

func (interceptor *AccessInterceptor) Intercept(c context.Context, req any) (any, error) {
	return nil, obs.New(codes.Unimplemented, "")
}

type AccessInfo struct {
	Req any
}
