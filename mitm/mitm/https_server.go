package mitm

import (
	"context"
	"crypto/tls"
	"github.com/ethanvc/evol/xlog"
	"log/slog"
	"net"
	"net/http"
)

type HttpsServer struct {
	certMgr *CertManager
	ln      *ConnListener
	tlsLn   net.Listener
	svr     *http.Server
}

func NewHttpsServer(certMgr *CertManager) *HttpsServer {
	return &HttpsServer{
		certMgr: certMgr,
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
		Handler: http.HandlerFunc(s.serveHTTP),
	}
	return s.svr.Serve(s.tlsLn)
}

func (s *HttpsServer) serveHTTP(w http.ResponseWriter, r *http.Request) {
	if err := s.serveHTTPInternal(w, r); err != nil {
		slog.ErrorContext(r.Context(), "serveHTTPInternalError", xlog.Error(err))
	}
}

func (s *HttpsServer) serveHTTPInternal(w http.ResponseWriter, r *http.Request) error {
	return nil
}
