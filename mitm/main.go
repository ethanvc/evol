package main

import (
	"context"
	"github.com/ethanvc/evol/base"
	"github.com/ethanvc/evol/xlog"
	"go.uber.org/fx"
	"net"
)

func main() {
	c := context.Background()
	err := xlog.InitDefaultLogger()
	base.PanicIfErr(c, err)
	app := fx.New(
		fx.NopLogger,
		fx.Provide(NewTcpServer),
		fx.Invoke(func(server *TcpServer) {}),
	)
	app.Run()
}

func NewTcpServer(lc fx.Lifecycle) (*TcpServer, error) {
	svr := &TcpServer{}
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

type TcpServer struct {
	ln net.Listener
}

func (s *TcpServer) Listen(c context.Context) error {
	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		return base.ErrWithCaller(err)
	}
	s.ln = ln
	return nil
}

func (s *TcpServer) Serve() error {
	for {
		conn, err := s.ln.Accept()
		c := context.Background()
		if err != nil {
			continue
		}
		go s.serveNewConnection(c, conn)
	}
	return nil
}

func (s *TcpServer) Shutdown() error {
	return nil
}

func (s *TcpServer) serveNewConnection(c context.Context, conn net.Conn) {
}
