package web

import (
	"crypto/tls"
	"io/ioutil"
	"net"
	"net/http"
)

// Server 表示一个服务
type Server interface {
	Serve() error
}

// NewServer 返回一个在 addr 监听，使用 handler 的 Server
func NewServer(addr string, handler http.Handler) Server {
	s := new(server)
	s.Server.Addr = addr
	s.Server.Handler = handler
	return s
}

// NewTLSServer 返回一个在 addr 监听 tls，使用 handler 的 Server ，certFile 和 keyFile 表示证书路径。
func NewTLSServer(addr, certFile, keyFile string, handler http.Handler) (Server, error) {
	// 读取证书
	certPEM, err := ioutil.ReadFile(certFile)
	if err != nil {
		return nil, err
	}
	keyPEM, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}
	// 初始化返回
	s := new(server)
	s.Server.Addr = addr
	s.Server.Handler = handler
	s.certPEM = certPEM
	s.certPEM = keyPEM
	return s, nil
}

// NewTLSServerWithKeyPair 返回一个在 addr 监听 tls，使用 handler 的 Server ，certPEM 和 keyPEM 表示证书的数据。
func NewTLSServerWithKeyPair(addr string, certPEM, keyPEM []byte, handler http.Handler) Server {
	s := new(server)
	s.Server.Addr = addr
	s.Server.Handler = handler
	s.certPEM = certPEM
	s.keyPEM = keyPEM
	return s
}

// server 实现 Server 接口
type server struct {
	http.Server
	certPEM []byte
	keyPEM  []byte
}

// Serve 实现 Server 接口
func (s *server) Serve() error {
	// 如果有证书，监听 https
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
	// 监听 http
	return s.Server.ListenAndServe()
}
