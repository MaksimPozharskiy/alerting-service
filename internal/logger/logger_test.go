package logger

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoggerResponseWriter_WriteHeaderAndWrite(t *testing.T) {
	rec := httptest.NewRecorder()
	lrw := NewLoggerResponseWriter(rec)

	lrw.WriteHeader(http.StatusCreated)
	if lrw.statusCode != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, lrw.statusCode)
	}

	body := []byte("test body")
	n, err := lrw.Write(body)
	if err != nil {
		t.Fatalf("unexpected write error: %v", err)
	}
	if n != len(body) {
		t.Errorf("expected write size %d, got %d", len(body), n)
	}
	if lrw.size != len(body) {
		t.Errorf("expected recorded size %d, got %d", len(body), lrw.size)
	}
}

func TestRequestLogger(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	RequestLogger(handler).ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	if string(body) != "ok" {
		t.Errorf("expected body 'ok', got %s", string(body))
	}
}

func TestResponseLogger(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("accepted"))
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	ResponseLogger(handler).ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusAccepted {
		t.Errorf("expected status 202, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	if string(body) != "accepted" {
		t.Errorf("expected body 'accepted', got %s", string(body))
	}
}
