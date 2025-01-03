package server

import (
	"net/http"
	"reflect"
	"testing"
)

func TestNewServer(t *testing.T) {
	adress := "localhost:8080"

	tests := []struct {
		name   string
		want   Server
		adress string
	}{
		{
			name:   "new server test",
			adress: adress,
			want: &ServerImpl{
				httpServer: http.Server{
					Addr: adress,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			want := NewServer(test.adress)

			if !reflect.DeepEqual(want, test.want) {
				t.Errorf("want: %v, got: %v", test.want, want)
			}
		})
	}
}

// @TODO Не понял как написать тест на Start
// func TestStart(t *testing.T) {
// 	mux := http.NewServeMux()

// 	tests := []struct {
// 		name string
// 	}{
// 		{
// 			name: "start server test",
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			server := NewServer("8080")
// 			go func() {
// 				time.Sleep(1 * time.Second)
// 				panic("")
// 			}()

// 			err := server.Start(mux)

// 			if err != nil {
// 				t.Error("unexpected error:", err)
// 			}
// 		})
// 	}
// }
