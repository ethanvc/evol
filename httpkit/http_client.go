package httpkit

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type HttpClient struct {
	interceptors []InterceptorFunc
}

func (cli *HttpClient) SendRequest(sa *SingleAttempt, req, resp any) error {
	if sa.Err != nil {
		return sa.Err
	}
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
	if realResp, ok := resp.(*io.ReadCloser); ok {
		*realResp = sa.Response.Body
		return nil
	}

	defer sa.Response.Body.Close()
	sa.RespBody, err = io.ReadAll(sa.Response.Body)
	if err != nil {
		return err
	}
	switch realResp := resp.(type) {
	case *string:
		*realResp = string(sa.RespBody)
	case *[]byte:
		*realResp = sa.RespBody
	default:
		if sa.Response.Header.Get("Content-Type") == "" {
			err := json.Unmarshal(sa.RespBody, resp)
			if err != nil {
				return err
			}
		} else {
			return errors.New("RespTypeNotSupport")
		}
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
