package web

import (
	"crypto/tls"
	"io/ioutil"
	"net"
	"net/http"
)

type Server interface {
	Serve() error
}

func NewServer(addr string, handler http.Handler) Server {
	s := new(server)
	s.Server.Addr = addr
	s.Server.Handler = handler
	return s
}

func NewTSLServer(addr, certFile, keyFile string, handler http.Handler) (Server, error) {
	certPEM, err := ioutil.ReadFile(certFile)
	if err != nil {
		return nil, err
	}
	keyPEM, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}
	s := new(server)
	s.Server.Addr = addr
	s.Server.Handler = handler
	s.certPEM = certPEM
	s.certPEM = keyPEM
	return s, nil
}

func NewTSLServerWithPair(addr string, certPEM, keyPEM []byte, handler http.Handler) Server {
	s := new(server)
	s.Server.Addr = addr
	s.Server.Handler = handler
	s.certPEM = certPEM
	s.keyPEM = keyPEM
	return s
}

type server struct {
	http.Server
	certPEM []byte
	keyPEM  []byte
}

func (s *server) Serve() error {
	if len(s.keyPEM) > 0 && len(s.certPEM) > 0 {
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
		s.Server.TLSConfig.Certificates[0], err = tls.X509KeyPair(s.certPEM, s.keyPEM)
		if err != nil {
			return err
		}
		//
		return s.Server.Serve(tls.NewListener(l, s.Server.TLSConfig))
	}
	return s.Server.ListenAndServe()
}
