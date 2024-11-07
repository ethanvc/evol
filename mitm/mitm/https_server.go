package mitm

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
)

type HttpsServer struct {
	certMgr *CertManager
	ln      *ConnListener
	tlsLn   net.Listener
	handler *HttpHandler
	svr     *http.Server
}

func NewHttpsServer(certMgr *CertManager, handler *HttpHandler) *HttpsServer {
	return &HttpsServer{
		certMgr: certMgr,
		handler: handler,
	}
}

func (s *HttpsServer) AddConn(conn net.Conn) {
	s.ln.Add(conn)
}

func (s *HttpsServer) Listen() error {
	s.ln = NewConnListener()
	tlsConf := &tls.Config{
		GetCertificate: s.certMgr.GetCertificate,
	}
	s.tlsLn = tls.NewListener(s.ln, tlsConf)
	return nil
}

func (s *HttpsServer) Shutdown() error {
	return s.svr.Shutdown(context.Background())
}

func (s *HttpsServer) Serve() error {
	s.svr = &http.Server{
		Handler: http.HandlerFunc(s.handler.ServeHTTP),
	}
	return s.svr.Serve(s.tlsLn)
}
