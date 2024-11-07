package mitm

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/ethanvc/evol/base"
	"github.com/ethanvc/evol/xlog"
	"google.golang.org/grpc/codes"
)

type HttpHandler struct {
	httpSvr  *HttpServer
	httpsSvr *HttpsServer
	store    *PacketStore
}

func NewHttpHandler(store *PacketStore) *HttpHandler {
	return &HttpHandler{
		store: store,
	}
}

func (h *HttpHandler) SetHttpSvr(httpSvr *HttpServer, httpsSvr *HttpsServer) {
	h.httpSvr = httpSvr
	h.httpsSvr = httpsSvr
}

func (h *HttpHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if h.httpSvr == nil || h.httpsSvr == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("initializing..."))
		return
	}
	err := h.serveHTTPInternal(w, req)
	if err != nil {
		slog.ErrorContext(req.Context(), "ServeError", xlog.Error(err))
	}
}

func (h *HttpHandler) serveHTTPInternal(w http.ResponseWriter, req *http.Request) error {
	if req.Method == http.MethodConnect {
		return h.serveConnect(w, req)
	}
	newReq, newResp, err := h.forwardToRemote(req)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		msg := fmt.Sprintf("mitm: forward error(%s)", err.Error())
		if _, err := w.Write([]byte(msg)); err != nil {
			slog.ErrorContext(req.Context(), "ForwardError", xlog.Error(err))
		}
		h.appendToStore(errors.New(msg), newReq, newResp)
		return nil
	}
	for k, vv := range newResp.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(newResp.StatusCode)
	if _, err := io.Copy(w, newResp.Body); err != nil {
		slog.ErrorContext(req.Context(), "CopyContentToRespError", xlog.Error(err))
	}
	h.appendToStore(nil, newReq, newResp)
	return nil
}

func (h *HttpHandler) serveConnect(w http.ResponseWriter, req *http.Request) error {
	hij, ok := w.(http.Hijacker)
	if !ok {
		return base.New(codes.Internal).SetEvent("ConvertToHijackerFailed")
	}
	conn, _, err := hij.Hijack()
	if err != nil {
		return base.ErrWithCaller(err)
	}
	conn.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
	h.httpsSvr.AddConn(conn)
	return nil
}

func (h *HttpHandler) forwardToRemote(req *http.Request) (*http.Request, *http.Response, error) {
	newReq, err := http.NewRequestWithContext(req.Context(), req.Method, req.URL.String(), req.Body)
	if err != nil {
		return nil, nil, err
	}
	if newReq.URL.Host == "" {
		if req.TLS != nil {
			newReq.URL.Scheme = "https"
		} else {
			newReq.URL.Scheme = "http"
		}
		newReq.URL.Host = req.Host
	}
	newReq.Header = req.Header
	newResp, err := http.DefaultClient.Do(newReq)
	if err != nil {
		return newReq, nil, err
	}
	return newReq, newResp, nil
}

func (h *HttpHandler) buildRequest(req *http.Request) (*http.Request, error) {
	newReq, err := http.NewRequest(req.Method, req.URL.String(), req.Body)
	if err != nil {
		return nil, err
	}
	newReq.Header = req.Header
	return newReq, nil
}

func (h *HttpHandler) requestRemoteServer(req *http.Request) (*http.Response, error) {
	newReq, err := http.NewRequestWithContext(req.Context(), req.Method, req.URL.String(), req.Body)
	if err != nil {
		return nil, err
	}
	newReq.Header = req.Header
	resp, err := http.DefaultClient.Do(newReq)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (h *HttpHandler) appendToStore(err error, req *http.Request, resp *http.Response) {
	packet := &HttpPacket{}
	packet.Req.Build(req)
	packet.Resp.Build(resp)
	packet.Err = err
	h.store.AppendHttpPacket(packet)
}
