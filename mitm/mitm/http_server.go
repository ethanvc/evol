package mitm

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"

	"github.com/ethanvc/evol/base"
	"github.com/ethanvc/evol/xlog"
	"go.uber.org/fx"
	"google.golang.org/grpc/codes"
)

type HttpServer struct {
	svr      *http.Server
	store    *PacketStore
	ln       *ConnListener
	httpsSvr *HttpsServer
}

func NewHttpServer(lc fx.Lifecycle, store *PacketStore, httpsSvr *HttpsServer) (*HttpServer, error) {
	svr := &HttpServer{
		store:    store,
		httpsSvr: httpsSvr,
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			err := svr.Listen()
			if err != nil {
				return base.ErrWithCaller(err)
			}
			go func() {
				_ = svr.Serve()
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return svr.Shutdown()
		},
	})
	return svr, nil
}

func (s *HttpServer) AddConn(conn net.Conn) {
	s.ln.Add(conn)
}

func (s *HttpServer) GetListenerAddr() *net.TCPAddr {
	return s.ln.Addr().(*net.TCPAddr)
}

func (s *HttpServer) Listen() error {
	s.ln = NewConnListener()
	return nil
}

func (s *HttpServer) Serve() error {
	s.svr = &http.Server{
		Handler: http.HandlerFunc(s.serveHTTP),
	}
	return s.svr.Serve(s.ln)
}

func (s *HttpServer) Shutdown() error {
	return s.svr.Shutdown(context.Background())
}

func (s *HttpServer) appendToStore(err error, req *http.Request, resp *http.Response) {
	packet := &HttpPacket{}
	packet.Req.Build(req)
	packet.Resp.Build(resp)
	packet.Err = err
	s.store.AppendHttpPacket(packet)
}

func (s *HttpServer) serveHTTP(w http.ResponseWriter, req *http.Request) {
	err := s.serveHTTPInternal(w, req)
	if err != nil {
		slog.ErrorContext(req.Context(), "ServeError", xlog.Error(err))
	}
}

func (s *HttpServer) serveHTTPInternal(w http.ResponseWriter, req *http.Request) error {
	if req.Method == http.MethodConnect {
		return s.serveConnect(w, req)
	}
	newReq, newResp, err := s.forwardToRemote(req)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		msg := fmt.Sprintf("mitm: forward error(%s)", err.Error())
		if _, err := w.Write([]byte(msg)); err != nil {
			slog.ErrorContext(req.Context(), "ForwardError", xlog.Error(err))
		}
		s.appendToStore(errors.New(msg), newReq, newResp)
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
	s.appendToStore(nil, newReq, newResp)
	return nil
}

func (s *HttpServer) serveConnect(w http.ResponseWriter, req *http.Request) error {
	hij, ok := w.(http.Hijacker)
	if !ok {
		return base.New(codes.Internal).SetEvent("ConvertToHijackerFailed")
	}
	conn, _, err := hij.Hijack()
	if err != nil {
		return base.ErrWithCaller(err)
	}
	conn.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
	s.httpsSvr.AddConn(conn)
	return nil
}

func (s *HttpServer) forwardToRemote(req *http.Request) (*http.Request, *http.Response, error) {
	newReq, err := http.NewRequestWithContext(req.Context(), req.Method, req.URL.String(), req.Body)
	if err != nil {
		return nil, nil, err
	}
	newReq.Header = req.Header
	newResp, err := http.DefaultClient.Do(newReq)
	if err != nil {
		return newReq, nil, err
	}
	return newReq, newResp, nil
}

func (s *HttpServer) buildRequest(req *http.Request) (*http.Request, error) {
	newReq, err := http.NewRequest(req.Method, req.URL.String(), req.Body)
	if err != nil {
		return nil, err
	}
	newReq.Header = req.Header
	return newReq, nil
}

func (s *HttpServer) requestRemoteServer(req *http.Request) (*http.Response, error) {
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
