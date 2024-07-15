package svrkit

import (
	"context"
	"github.com/gin-gonic/gin"
	"slices"
)

type InterceptorFunc func(c context.Context, req any, next Nexter) (any, error)

type Nexter struct {
	chain *Chain
	index int
}

func (n Nexter) Next(c context.Context, req any) (any, error) {
	if n.index >= len(n.chain.interceptors) {
		if n.chain.handler == nil {
			return nil, ErrNexterEnded
		}
		return n.chain.handler(c, req)
	}
	newNext := Nexter{
		chain: n.chain,
		index: n.index + 1,
	}
	return n.chain.interceptors[n.index](c, req, newNext)
}

func (n Nexter) Chain() *Chain {
	return n.chain
}

type Chain struct {
	interceptors []InterceptorFunc
	NewReq       func() any
	handler      func(ctx context.Context, req any) (resp any, err error)
}

func NewChain[Req, Resp any](interceptors []InterceptorFunc,
	handler func(ctx context.Context, req *Req) (resp *Resp, err error)) *Chain {
	return &Chain{
		interceptors: slices.Clone(interceptors),
		NewReq: func() any {
			return new(Req)
		},
		handler: func(ctx context.Context, req any) (resp any, err error) {
			realReq := req.(*Req)
			resp, err = handler(ctx, realReq)
			return resp, err
		},
	}
}

func NewGinChain[Req, Resp any](interceptors []InterceptorFunc,
	handler func(ctx context.Context, req *Req) (resp *Resp, err error)) gin.HandlerFunc {
	chain := NewChain(interceptors, handler)
	return func(c *gin.Context) {
		ctx, req := WithHttpRequestContext(c.Request.Context())
		req.Pattern = c.FullPath()
		req.Request = c.Request
		req.PathParams = c.Params
		req.Writer = c.Writer
		chain.GetNexter().Next(ctx, nil)
	}
}

func (ch *Chain) GetNexter() Nexter {
	return Nexter{
		chain: ch,
	}
}
