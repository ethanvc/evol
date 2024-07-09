package obs

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
	return context.WithValue(c, contextKeyLogContext{}, &LogContext{})
}

func GetLogContext(c context.Context) *LogContext {
	lc, _ := c.Value(contextKeyLogContext{}).(*LogContext)
	return lc
}
