package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Server interface {
	Start(*chi.Mux) error
}

type ServerImpl struct {
	httpServer http.Server
}

func NewServer(port string) Server {
	return &ServerImpl{
		httpServer: http.Server{
			Addr: ":" + port,
		},
	}
}

func (s *ServerImpl) Start(mux *chi.Mux) error {
	if err := http.ListenAndServe(s.httpServer.Addr, mux); err != http.ErrServerClosed {
		return err
	}
	return nil
}
