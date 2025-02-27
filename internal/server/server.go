package server

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Server interface {
	Start(*chi.Mux) error
	Shutdown(context.Context) error
}

type ServerImpl struct {
	httpServer http.Server
}

func NewServer(addr string) Server {
	return &ServerImpl{
		httpServer: http.Server{
			Addr: addr,
		},
	}
}

func (s *ServerImpl) Start(mux *chi.Mux) error {
	if err := http.ListenAndServe(s.httpServer.Addr, mux); err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *ServerImpl) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
