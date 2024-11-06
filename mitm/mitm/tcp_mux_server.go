package mitm

import (
	"bufio"
	"context"
	"log/slog"
	"net"

	"github.com/ethanvc/evol/base"
	"github.com/ethanvc/evol/xlog"
	"go.uber.org/fx"
	"google.golang.org/grpc/codes"
)

func NewTcpMuxServer(lc fx.Lifecycle, httpSvr *HttpServer, httpsSvr *HttpsServer) (*TcpMuxServer, error) {
	svr := &TcpMuxServer{
		httpSvr:  httpSvr,
		httpsSvr: httpsSvr,
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
	ln       net.Listener
	httpSvr  *HttpServer
	httpsSvr *HttpsServer
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
		s.httpSvr.AddConn(NewBufferIoConn(conn, bufReader))
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
