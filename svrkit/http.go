package svrkit

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HttpRequestContext struct {
	Pattern    string
	Request    *http.Request
	PathParams gin.Params
	Writer     http.ResponseWriter
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
