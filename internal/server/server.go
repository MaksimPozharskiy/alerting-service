package server

import (
	"net/http"
)

type Server interface {
	Start(*http.ServeMux) error
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

func (s *ServerImpl) Start(mux *http.ServeMux) error {
	if err := http.ListenAndServe(s.httpServer.Addr, mux); err != http.ErrServerClosed {
		return err
	}
	return nil
}
