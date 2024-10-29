package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/ethanvc/evol/base"
	"github.com/ethanvc/evol/mitm/mitm"
	"github.com/ethanvc/evol/xlog"
	"go.uber.org/fx"
	"google.golang.org/grpc/codes"
	"io"
	"log/slog"
	"net"
	"net/http"
)

func main() {
	c := context.Background()
	err := xlog.InitDefaultLogger()
	base.PanicIfErr(c, err)
	app := fx.New(
		fx.NopLogger,
		fx.Provide(NewTcpServer),
		fx.Provide(mitm.NewPacketStore),
		fx.Provide(NewHttpServer),
		fx.Invoke(func(server *TcpMuxServer) {}),
	)
	base.PanicIfErr(c, app.Err())
	app.Run()
}

func NewTcpServer(lc fx.Lifecycle, httpSvr *HttpServer) (*TcpMuxServer, error) {
	svr := &TcpMuxServer{
		httpSvr: httpSvr,
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			err := svr.Listen(ctx)
			base.PanicIfErr(ctx, err)
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

type TcpMuxServer struct {
	ln      net.Listener
	httpSvr *HttpServer
}

func (s *TcpMuxServer) Listen(c context.Context) error {
	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		return base.ErrWithCaller(err)
	}
	s.ln = ln
	return nil
}

func (s *TcpMuxServer) Serve() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			return base.ErrWithCaller(err)
		}
		c := context.Background()
		go func() {
			err := s.serveNewConnection(c, conn)
			if err != nil {
				slog.ErrorContext(c, "ServeConnectionError", xlog.Error(err))
			}
		}()
	}
}

func (s *TcpMuxServer) Shutdown() error {
	return nil
}

func (s *TcpMuxServer) serveNewConnection(c context.Context, conn net.Conn) error {
	bufReader := bufio.NewReaderSize(conn, 1024)
	buf, err := bufReader.Peek(10)
	if err != nil {
		conn.Close()
		return base.ErrWithCaller(err)
	}
	switch typ := s.guessProtocolType(buf); typ {
	case ProtocolTypeHttp:
		s.httpSvr.AddConn(mitm.NewBufferIoConn(conn, bufReader))
		return nil
	default:
		conn.Close()
		return base.New(codes.Internal).SetEvent("ProtocolTypeError")
	}
}

func (s *TcpMuxServer) guessProtocolType(buf []byte) ProtocolType {
	return ProtocolTypeHttp
}

type ProtocolType int

const (
	ProtocolTypeTcp ProtocolType = iota
	ProtocolTypeHttp
)

type HttpServer struct {
	svr   *http.Server
	store *mitm.PacketStore
	ln    *mitm.ConnListener
}

func NewHttpServer(lc fx.Lifecycle, store *mitm.PacketStore) (*HttpServer, error) {
	svr := &HttpServer{
		store: store,
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
	s.ln = mitm.NewConnListener()
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
	packet := &mitm.HttpPacket{}
	packet.Req.Build(req)
	packet.Resp.Build(resp)
	packet.Err = err
	s.store.AppendHttpPacket(packet)
}

func (s *HttpServer) serveHTTP(w http.ResponseWriter, req *http.Request) {
	newReq, newResp, err := s.forwardToRemote(req)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		msg := fmt.Sprintf("mitm: forward error(%s)", err.Error())
		if _, err := w.Write([]byte(msg)); err != nil {
			slog.ErrorContext(req.Context(), "ForwardError", xlog.Error(err))
		}
		s.appendToStore(errors.New(msg), newReq, newResp)
		return
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
