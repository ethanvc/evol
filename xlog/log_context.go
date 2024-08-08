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
	extra     map[string]interface{}
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

func (lc *LogContext) SetExtra(key string, val any) {
	if lc == nil || key == "" || val == nil {
		return
	}
	lc.mux.Lock()
	defer lc.mux.Unlock()

	const maxExtraLen = 10
	if len(lc.extra) > maxExtraLen {
		return
	}
	if lc.extra == nil {
		lc.extra = make(map[string]interface{})
	}
	lc.extra[key] = val
}
