package svrkit

import (
	"golang.org/x/net/context"
	"sync"
)

type ContextStore struct {
	mux   sync.Mutex
	store map[any]any
}

type contextKeyContextStore struct{}

func WithContextStore(c context.Context) context.Context {
	return context.WithValue(c, contextKeyContextStore{}, &ContextStore{})
}

func GetContextStore(c context.Context) *ContextStore {
	store, _ := c.Value(contextKeyContextStore{}).(*ContextStore)
	return store
}
