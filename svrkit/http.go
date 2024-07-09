package svrkit

import (
	"context"
	"net/http"
)

type HttpRequestContext struct {
	Pattern string
	Request *http.Request
	Writer  http.ResponseWriter
}

type contextKeyHttpRequestContext struct{}

func WithHttpRequestContext(c context.Context) (context.Context, *HttpRequestContext) {
	hrc := &HttpRequestContext{}
	return context.WithValue(c, &contextKeyHttpRequestContext{}, hrc), hrc
}

func GetHttpRequestContext(c context.Context) *HttpRequestContext {
	hrc, _ := c.Value(&contextKeyHttpRequestContext{}).(*HttpRequestContext)
	return hrc
}
