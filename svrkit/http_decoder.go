package svrkit

import (
	"context"
)

type HttpDecoder struct {
}

func NewHttpDecoder() *HttpDecoder {
	return &HttpDecoder{}
}

func (decoder *HttpDecoder) Intercept(c context.Context, req any, nexter Nexter) (any, error) {
	return nil, nil
}
