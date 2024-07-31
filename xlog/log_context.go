package xlog

import (
	"context"
	"sync"
	"time"
)

type LogContext struct {
	mux       sync.Mutex
	method    string
	startTime time.Time
}

type contextKeyLogContext struct{}

func WithLogContext(c context.Context, method string) context.Context {
	logCtx := &LogContext{
		method: method,
	}
	return context.WithValue(c, contextKeyLogContext{}, logCtx)
}

func GetLogContext(c context.Context) *LogContext {
	lc, _ := c.Value(contextKeyLogContext{}).(*LogContext)
	return lc
}

func (lc *LogContext) GetMethod() string {
	if lc == nil {
		return "Global"
	}
	lc.mux.Lock()
	defer lc.mux.Unlock()
	return lc.method
}

func (lc *LogContext) GetStartTime() time.Time {
	if lc == nil {
		return time.Now()
	}
	lc.mux.Lock()
	defer lc.mux.Unlock()
	return lc.startTime

}
