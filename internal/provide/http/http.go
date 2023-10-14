package http

import (
	"net/http"
	"time"

	"github.com/thoohv5/person/pkg/log"
	"github.com/thoohv5/person/pkg/server"
	"github.com/thoohv5/person/pkg/transport"
)

// New 注册HTTP
func New(
	conf *Config,
	log log.Logger,
	flag bool,
	handler http.Handler,
) transport.Server {
	if !flag {
		return nil
	}

	var opts = []server.Option{
		server.Logger(log),
		// server.TLSConfig(generateTLSConfig()),
		server.Handler(handler),
	}

	if network := conf.Network; network != "" {
		opts = append(opts, server.Network(network))
	}
	if addr := conf.Addr; addr != "" {
		opts = append(opts, server.Address(addr))
	}
	if timeout := conf.Timeout; timeout > 0 {
		opts = append(opts, server.Timeout(time.Duration(timeout)*time.Second))
	}
	srv := server.New(opts...)

	return srv
}

//
// func generateTLSConfig() *tls.Config {
// 	pool := x509.NewCertPool()
// 	pool.AppendCertsFromPEM(ih.CACrt)
//
// 	tlsCert, err := tls.X509KeyPair(ih.ServerCrt, ih.ServerKey)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	// #nosec
// 	return &tls.Config{
// 		Certificates: []tls.Certificate{tlsCert},
// 		ClientCAs:    pool,
// 		ClientAuth:   tls.VerifyClientCertIfGiven,
// 		Rand:         rand.Reader,
// 	}
// }
