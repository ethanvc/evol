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

func (kit *GinKit) Handlers(handler InterceptorFunc) gin.HandlerFunc {
	if handler == nil {
		panic("handler must not be nil")
	}
	interceptors := slices.Clone(kit.interceptors)
	interceptors = append(interceptors, handler)
	nexter := NewNexter(interceptors...)
	return func(c *gin.Context) {
		nexter.Next(c.Request.Context(), nil)
	}
}
