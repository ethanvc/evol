package mitm

import (
	"context"
	"net"
	"net/http"

	"github.com/ethanvc/evol/base"
	"go.uber.org/fx"
)

type HttpServer struct {
	svr     *http.Server
	store   *PacketStore
	ln      *ConnListener
	handler *HttpHandler
}

func NewHttpServer(lc fx.Lifecycle, store *PacketStore, handler *HttpHandler) (*HttpServer, error) {
	svr := &HttpServer{
		store:   store,
		handler: handler,
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
		Handler: http.HandlerFunc(s.handler.ServeHTTP),
	}
	return s.svr.Serve(s.ln)
}

func (s *HttpServer) Shutdown() error {
	return s.svr.Shutdown(context.Background())
}
