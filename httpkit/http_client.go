package httpkit

import (
	"errors"
	"io"
	"net/http"
)

type HttpClient struct {
	interceptors []InterceptorFunc
}

func (cli *HttpClient) SendRequest(sa *SingleAttempt, req, resp any) error {
	invoker := Invoker{
		cli: cli,
	}
	return invoker.Invoke(sa, req, resp)
}

func (cli *HttpClient) sendHttpRequest(sa *SingleAttempt, req, resp any) error {
	if sa.Err != nil {
		return sa.Err
	}
	var err error
	sa.Response, err = http.DefaultClient.Do(sa.Request)
	if err != nil {
		return err
	}

	defer sa.Response.Body.Close()
	switch realResp := resp.(type) {
	case *string:
		buf, err := io.ReadAll(sa.Response.Body)
		if err != nil {
			return err
		}
		*realResp = string(buf)
	default:
		return errors.New("UnsupportResp")
	}
	return nil
}

type Invoker struct {
	cli   *HttpClient
	index int
}

func (invoker Invoker) Invoke(sa *SingleAttempt, req, resp any) error {
	if invoker.index >= len(invoker.cli.interceptors) {
		return invoker.cli.sendHttpRequest(sa, req, resp)
	}
	newNext := Invoker{
		cli:   invoker.cli,
		index: invoker.index + 1,
	}
	return invoker.cli.interceptors[invoker.index](sa, req, resp, newNext)
}

type InterceptorFunc func(sa *SingleAttempt, req, resp any, invoker Invoker) error
