package svrkit

import (
	"github.com/gin-gonic/gin"
	"slices"
)

type GinKit struct {
	interceptors []InterceptorFunc
}

func NewGinKit() *GinKit {
	return &GinKit{}
}

func (kit *GinKit) Clone(fn ...InterceptorFunc) *GinKit {
	newKit := &GinKit{
		interceptors: slices.Clone(kit.interceptors),
	}
	newKit.interceptors = append(newKit.interceptors, fn...)
	return newKit
}

func (kit *GinKit) Handlers(handler InterceptorFunc) gin.HandlerFunc {
	if handler == nil {
		panic("handler must not be nil")
	}
	interceptors := slices.Clone(kit.interceptors)
	interceptors = append(interceptors, handler)
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
