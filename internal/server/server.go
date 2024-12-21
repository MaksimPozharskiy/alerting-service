package server

import (
	"fmt"
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
	fmt.Printf("Starting server at port %s\n", s.httpServer.Addr)
	err := http.ListenAndServe(s.httpServer.Addr, mux)
	if err != nil {
		panic(err)
	}
}
