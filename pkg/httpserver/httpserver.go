package httpserver

import (
	"context"
	"net"
	"net/http"
	"time"

	"log"
)

const (
	_defaultAddr            = ":80"
	_defaultReadTimeout     = 5 * time.Second
	_defaultWriteTimeout    = 5 * time.Second
	_defaultMaxHeaderBytes  = 1 << 20
	_defaultShutdownTimeout = 3 * time.Second
)

// Server - http server wrapper with custom logic.
type Server struct {
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
}

// Option - options for configuring http server.
type Option func(*Server)

func Port(port string) Option {
	return func(s *Server) {
		s.server.Addr = net.JoinHostPort("", port)
	}
}

func ReadTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.server.ReadTimeout = timeout
	}
}

func WriteTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.server.WriteTimeout = timeout
	}
}

func ShutdownTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = timeout
	}
}

func New(handler http.Handler, opts ...Option) *Server {
	httpServer := &http.Server{
		Addr:           _defaultAddr,
		Handler:        handler,
		ReadTimeout:    _defaultReadTimeout,
		WriteTimeout:   _defaultWriteTimeout,
		MaxHeaderBytes: _defaultMaxHeaderBytes,
	}

	s := &Server{
		server:          httpServer,
		notify:          make(chan error, 1),
		shutdownTimeout: _defaultShutdownTimeout,
	}

	for _, opt := range opts {
		opt(s)
	}

	s.start()

	return s
}

// start - bootstraps http server.
func (s *Server) start() {
	log.Printf("Starting HTTP server on port %s", s.server.Addr)
	go func() {
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()
}

// Notify - returns error notification channel.
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown - shuts down http server gracefully.
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
