package xlog

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

type LogContext struct {
	mux        sync.Mutex
	method     string
	startTime  time.Time
	attributes []slog.Attr
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

func (lc *LogContext) SetAttribute(key string, val slog.Value) {
	if lc == nil || key == "" {
		return
	}
	lc.mux.Lock()
	defer lc.mux.Unlock()

	const maxExtraLen = 100
	if len(lc.attributes) > maxExtraLen {
		return
	}
	lc.attributes = append(lc.attributes, slog.Attr{
		Key:   key,
		Value: val,
	})
}

func (lc *LogContext) TraverseAttributes(f func(attributes []slog.Attr)) {
	if lc == nil {
		return
	}
	lc.mux.Lock()
	defer lc.mux.Unlock()
	f(lc.attributes)
}
