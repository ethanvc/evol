package main

import (
	"context"
	"fmt"
	"github.com/ethanvc/evol/xlog"
	"github.com/miekg/dns"
	"go.uber.org/fx"
	"log"
)

func main() {
	if err := xlog.InitDefaultLogger(); err != nil {
		panic(err)
	}
	app := fx.New(
		fx.NopLogger,
		fx.Provide(NewDnsServer),
		fx.Invoke(func(server *dns.Server) {}),
	)
	app.Run()
}

func NewDnsServer(lc fx.Lifecycle) (*dns.Server, error) {
	svr := &dns.Server{
		Addr:    ":5999",
		Net:     "udp",
		Handler: dns.HandlerFunc(ServeDNS),
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				err := svr.ListenAndServe()
				if err != nil {
					panic(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return svr.ShutdownContext(ctx)
		},
	})
	return svr, nil
}

func ServeDNS(w dns.ResponseWriter, req *dns.Msg) {
	printRequest(req)

	resp, err := dns.ExchangeContext(context.Background(), req, "8.8.8.8:53")
	if err != nil {
		fmt.Printf("err:%s\n", err)
		return
	}
	err = w.WriteMsg(resp)
	if err != nil {
		fmt.Printf("err:%s\n", err)
		return
	}
}

func printRequest(m *dns.Msg) {
	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeA:
			log.Printf("Query for %s\n", q.Name)
		}
	}
}
