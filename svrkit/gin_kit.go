package svrkit

import (
	"context"
	"github.com/gin-gonic/gin"
	"slices"
)

type GinKit struct {
	interceptors []InterceptorFunc
}

func NewGinKit() *GinKit {
	return &GinKit{}
}

func (kit *GinKit) AddInterceptor(fn ...InterceptorFunc) *GinKit {
	kit.interceptors = append(kit.interceptors, fn...)
	return kit
}

func (kit *GinKit) Clone(fn ...InterceptorFunc) *GinKit {
	newKit := &GinKit{
		interceptors: slices.Clone(kit.interceptors),
	}
	newKit.interceptors = append(newKit.interceptors, fn...)
	return newKit
}

func (kit *GinKit) Handlers(handler any) gin.HandlerFunc {
	if handler == nil {
		panic("handler must not be nil")
	}
	interceptors := slices.Clone(kit.interceptors)
	interceptors = append(interceptors, func(c context.Context, req any, nexter Nexter) (any, error) {
		return nil, nil
	})
	nexter := NewNexter(interceptors...)
	return func(c *gin.Context) {
		ctx, req := WithHttpRequestContext(c.Request.Context())
		req.Pattern = c.FullPath()
		req.Request = c.Request
		req.PathParams = c.Params
		req.Writer = c.Writer
		nexter.Next(ctx, nil)
	}
}
