package server

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/thoohv5/person/internal/provide/logger"
	"github.com/thoohv5/person/pkg/log"
	"github.com/thoohv5/person/pkg/transport"
	"github.com/thoohv5/person/pkg/util/endpoint"
	"github.com/thoohv5/person/pkg/util/host"
)

// Option is an HTTP server option.
type Option func(*server)

// Network with server network.
func Network(network string) Option {
	return func(s *server) {
		s.network = network
	}
}

// Address with server address.
func Address(addr string) Option {
	return func(s *server) {
		s.address = addr
	}
}

// Timeout with server timeout.
func Timeout(timeout time.Duration) Option {
	return func(s *server) {
		s.timeout = timeout
	}
}

// Logger with server log.
func Logger(log log.Logger) Option {
	return func(s *server) {
		s.log = log
	}
}

// TLSConfig with TLS config.
func TLSConfig(c *tls.Config) Option {
	return func(o *server) {
		o.tlsConf = c
	}
}

// Listener with server lis
func Listener(lis net.Listener) Option {
	return func(s *server) {
		s.lis = lis
	}
}

// Handler with server handler
func Handler(handler http.Handler) Option {
	return func(s *server) {
		s.handler = handler
	}
}

type server struct {
	*http.Server

	lis      net.Listener
	tlsConf  *tls.Config
	endpoint *url.URL
	err      error
	network  string
	address  string
	timeout  time.Duration
	log      log.Logger

	handler http.Handler
}

func New(opts ...Option) transport.Server {

	srv := &server{
		network: "tcp",
		address: ":0",
		timeout: 1 * time.Second,
		log:     NewDevNullLog(os.Stdout),
		handler: func() http.Handler {
			return http.NewServeMux()
		}(),
	}
	for _, o := range opts {
		o(srv)
	}

	srv.Server = &http.Server{
		Handler:   srv.handler,
		TLSConfig: srv.tlsConf,
	}
	srv.err = srv.listenAndEndpoint()

	return srv
}

func (s *server) Start(ctx context.Context) error {
	if s.err != nil {
		return s.err
	}
	s.BaseContext = func(net.Listener) context.Context {
		return ctx
	}
	s.log.Infoc(ctx, "[HTTP] server listening on", logger.FieldString("addr", s.lis.Addr().String()))
	var err error
	if s.tlsConf != nil {
		err = s.ServeTLS(s.lis, "", "")
	} else {
		err = s.Serve(s.lis)
	}
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *server) Stop(ctx context.Context) error {
	s.log.Infoc(ctx, "[HTTP] server stopping")
	if err := s.Shutdown(ctx); err != nil {
		s.log.Errorc(ctx, "[HTTP] server stop", logger.FieldError(err))
		return err
	}
	s.log.Infoc(ctx, "[HTTP] server stop")
	return nil
}

func (s *server) listenAndEndpoint() error {
	if s.lis == nil {
		lis, err := net.Listen(s.network, s.address)
		if err != nil {
			return err
		}
		s.lis = lis
	}
	addr, err := host.Extract(s.address, s.lis)
	if err != nil {
		_ = s.lis.Close()
		return err
	}
	s.endpoint = endpoint.NewEndpoint(endpoint.Scheme("server", s.tlsConf != nil), addr)
	return nil
}
