package svrkit

import (
	"context"
)

type InterceptorFunc func(c context.Context, req any, next Nexter) (any, error)

type Nexter struct {
	interceptors []InterceptorFunc
	index        int
}

func NewNexter(interceptors ...InterceptorFunc) Nexter {
	return Nexter{
		interceptors: interceptors,
	}
}

func (n Nexter) Next(c context.Context, req any) (any, error) {
	if n.index >= len(n.interceptors) {
		return nil, ErrNexterEnded
	}
	newNext := Nexter{
		interceptors: n.interceptors,
		index:        n.index + 1,
	}
	return n.interceptors[n.index](c, req, newNext)
}
