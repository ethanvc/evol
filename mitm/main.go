package main

import (
	"bufio"
	"bytes"
	"context"
	"github.com/ethanvc/evol/base"
	"github.com/ethanvc/evol/xlog"
	"go.uber.org/fx"
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
		fx.Invoke(func(server *TcpMuxServer) {}),
	)
	app.Run()
}

func NewTcpServer(lc fx.Lifecycle) (*TcpMuxServer, error) {
	svr := &TcpMuxServer{}
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
	ln net.Listener
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
	defer conn.Close()
	bufReader := bufio.NewReaderSize(conn, 1024)
	buf, err := bufReader.Peek(10)
	if err != nil {
		return base.ErrWithCaller(err)
	}
	switch typ := s.guessProtocolType(buf); typ {
	case ProtocolTypeHttp:
	case ProtocolTypeTcp:
	case ProtocolTypeHttps:
	}
	return nil
}

func (s *TcpMuxServer) guessProtocolType(buf []byte) ProtocolType {
	if bytes.HasPrefix(buf, []byte("GET")) {
		return ProtocolTypeHttp
	}
	return ProtocolTypeTcp
}

type ProtocolType int

const (
	ProtocolTypeTcp ProtocolType = iota
	ProtocolTypeHttp
	ProtocolTypeHttps
)

type HttpServer struct {
	svr *http.Server
}

func NewHttpServer(lc fx.Lifecycle) (*HttpServer, error) {
	svr := &HttpServer{}
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

func (s *HttpServer) Listen() error {

}

func (s *HttpServer) Serve() error {

}

func (s *HttpServer) Shutdown() error {

}
