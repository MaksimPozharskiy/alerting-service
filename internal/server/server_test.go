package server

import (
	"net/http"
	"reflect"
	"testing"
)

func TestNewServer(t *testing.T) {
	address := "localhost:8080"

	tests := []struct {
		name    string
		want    Server
		address string
	}{
		{
			name:    "new server test",
			address: address,
			want: &ServerImpl{
				httpServer: http.Server{
					Addr: address,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			want := NewServer(test.address)

			if !reflect.DeepEqual(want, test.want) {
				t.Errorf("want: %v, got: %v", test.want, want)
			}
		})
	}
}
