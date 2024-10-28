package httpkit

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"slices"
	"strings"
)

type HttpClient struct {
	PreInterceptors  []InterceptorFunc
	PostInterceptors []InterceptorFunc
	EncodeRequest    func(c context.Context, sa *SingleAttempt, req any) error
	DecodeResponse   func(c context.Context, sa *SingleAttempt, resp any) error
}

func (cli *HttpClient) Clone() *HttpClient {
	newCli := &HttpClient{
		PreInterceptors:  slices.Clone(cli.PreInterceptors),
		PostInterceptors: slices.Clone(cli.PostInterceptors),
		EncodeRequest:    cli.EncodeRequest,
		DecodeResponse:   cli.DecodeResponse,
	}
	return newCli
}

func (cli *HttpClient) SendRequest(c context.Context, sa *SingleAttempt, req, resp any) error {
	if sa.Err != nil {
		return sa.Err
	}
	invoker := Invoker{
		cli: cli,
	}
	return invoker.Invoke(c, sa, req, resp)
}

func (cli *HttpClient) sendHttpRequestAfterAllInterceptors(c context.Context, sa *SingleAttempt, req, resp any) error {
	if sa.Err != nil {
		return sa.Err
	}
	if sa.UpdateContext && sa.Request.Context() != c {
		sa.Request = sa.Request.WithContext(c)
	}
	err := cli.encodeRequest(c, sa, req)
	if err != nil {
		return err
	}
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
	return cli.decodeResponse(c, sa, resp)
}

func (cli *HttpClient) encodeRequest(c context.Context, sa *SingleAttempt, req any) error {
	if sa.Request.Body != nil {
		return nil
	}
	if cli.EncodeRequest != nil {
		return cli.EncodeRequest(c, sa, req)
	}
	return nil
}

func (cli *HttpClient) decodeResponse(c context.Context, sa *SingleAttempt, resp any) error {
	if cli.DecodeResponse != nil {
		return cli.DecodeResponse(c, sa, resp)
	}
	switch realResp := resp.(type) {
	case *string:
		*realResp = string(sa.RespBody)
		return nil
	case *[]byte:
		*realResp = sa.RespBody
		return nil
	}
	contentType := strings.ToLower(sa.Request.Header.Get("Content-Type"))
	if strings.HasPrefix(contentType, "application/json") {
		err := json.Unmarshal(sa.RespBody, resp)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("RespTypeNotSupport")
}

type Invoker struct {
	cli       *HttpClient
	index     int
	indexType indexType
}

type indexType int

const (
	indexTypePreInterceptors indexType = iota
	indexTypeGlobalInterceptors
	indexTypePostInterceptors
)

func (invoker Invoker) Invoke(c context.Context, sa *SingleAttempt, req, resp any) error {
	// must be ordered
	indexTypes := []indexType{indexTypePreInterceptors, indexTypeGlobalInterceptors, indexTypePostInterceptors}
	interceptors := [][]InterceptorFunc{invoker.cli.PreInterceptors, globalInterceptors, invoker.cli.PostInterceptors}
	for i, typ := range indexTypes {
		if typ != invoker.indexType {
			continue
		}
		if invoker.index >= len(interceptors[i]) {
			invoker.indexType++
			continue
		}
		handler := interceptors[i][invoker.index]
		invoker.index++
		return handler(c, sa, req, resp, invoker)
	}
	return invoker.cli.sendHttpRequestAfterAllInterceptors(c, sa, req, resp)
}

type InterceptorFunc func(c context.Context, sa *SingleAttempt, req, resp any, invoker Invoker) error
