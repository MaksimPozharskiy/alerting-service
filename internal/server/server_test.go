package server

import (
	"testing"
)

func TestNewServer(t *testing.T) {
	addr := "localhost:8080"
	s := NewServer(addr)

	impl, ok := s.(*ServerImpl)
	if !ok {
		t.Fatalf("expected *ServerImpl, got %T", s)
	}

	if impl.httpServer.Addr != addr {
		t.Errorf("expected addr %s, got %s", addr, impl.httpServer.Addr)
	}
}
