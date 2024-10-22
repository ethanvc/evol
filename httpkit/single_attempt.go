package httpkit

import "net/http"

type SingleAttempt struct {
	Request  *http.Request
	Response *http.Response
	err      error
}
