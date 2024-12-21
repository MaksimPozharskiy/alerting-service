package server

import (
	"net/http"
)

type Server interface {
	Start(*http.ServeMux)
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

func (s *ServerImpl) Start(mux *http.ServeMux) {
	err := http.ListenAndServe(s.httpServer.Addr, mux)
	if err != nil {
		panic(err)
	}
}
