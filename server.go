package web

import (
	"crypto/tls"
	"net"
	"net/http"
)

type Server interface {
	Serve() error
	ServeTLS(certFile, keyFile string) error
	ServeTLSWithKeyPair(certPEM, keyPEM []byte) error
}

func NewServer(addr string, handler http.Handler) Server {
	s := new(server)
	s.Server.Addr = addr
	s.Server.Handler = handler
	return s
}

type server struct {
	http.Server
}

func (s *server) Serve() error {
	return s.Server.ListenAndServe()
}

func (s *server) ServeTLS(certFile, keyFile string) error {
	return s.Server.ListenAndServeTLS(certFile, keyFile)
}

func (s *server) ServeTLSWithKeyPair(certFile, keyFile []byte) error {
	l, err := net.Listen("tcp", s.Server.Addr)
	if err != nil {
		return err
	}
	//
	s.Server.TLSConfig = &tls.Config{
		// http2.NextProtoTLS
		NextProtos: []string{"h2"},
	}
	s.Server.TLSConfig.Certificates = make([]tls.Certificate, 1)
	s.Server.TLSConfig.Certificates[0], err = tls.X509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}
	//
	return s.Server.Serve(tls.NewListener(l, s.Server.TLSConfig))
}
