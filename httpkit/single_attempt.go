package httpkit

import (
	"context"
	"net/http"
)

type SingleAttempt struct {
	Request  *http.Request
	Response *http.Response
	Err      error
	RespBody []byte
}

func NewSingleAttempt(c context.Context, method, url string) *SingleAttempt {
	sa := &SingleAttempt{}
	sa.Request, sa.Err = http.NewRequestWithContext(c, method, url, nil)
	if sa.Err != nil {
		sa.Request, _ = http.NewRequestWithContext(c, http.MethodGet, "http://127.0.10.10/url_or_method_error", nil)
	}
	return sa
}
