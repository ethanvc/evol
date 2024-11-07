package main

import (
	"context"

	"github.com/ethanvc/evol/base"
	"github.com/ethanvc/evol/mitm/mitm"
	"github.com/ethanvc/evol/xlog"
	"go.uber.org/fx"
)

func main() {
	c := context.Background()
	err := xlog.InitDefaultLogger()
	base.PanicIfErr(c, err)
	initializeHandler := func(handler *mitm.HttpHandler, httpSvr *mitm.HttpServer, httpsSvr *mitm.HttpsServer) {
		handler.SetHttpSvr(httpSvr, httpsSvr)
	}
	app := fx.New(
		fx.NopLogger,
		fx.Provide(mitm.NewHttpHandler),
		fx.Provide(mitm.NewCertManager),
		fx.Provide(mitm.NewTcpMuxServer),
		fx.Provide(mitm.NewPacketStore),
		fx.Provide(mitm.NewHttpServer),
		fx.Provide(NewHttpsServer),
		fx.Invoke(func(server *mitm.TcpMuxServer) {}),
		fx.Invoke(initializeHandler),
	)
	base.PanicIfErr(c, app.Err())
	app.Run()
}

func NewHttpsServer(lc fx.Lifecycle, certMgr *mitm.CertManager, handler *mitm.HttpHandler) *mitm.HttpsServer {
	svr := mitm.NewHttpsServer(certMgr, handler)
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
	return svr
}
