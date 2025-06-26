package compressor

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGzipMiddleware_CompressesResponse(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})

	gzHandler := GzipMiddleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rec := httptest.NewRecorder()
	gzHandler.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if !strings.Contains(res.Header.Get("Content-Encoding"), "gzip") {
		t.Errorf("expected Content-Encoding: gzip, got %q", res.Header.Get("Content-Encoding"))
	}

	gr, err := gzip.NewReader(res.Body)
	if err != nil {
		t.Fatal("failed to create gzip reader:", err)
	}
	defer gr.Close()

	body, err := io.ReadAll(gr)
	if err != nil {
		t.Fatal("failed to read response body:", err)
	}

	if string(body) != "hello" {
		t.Errorf("expected body 'hello', got %q", string(body))
	}
}
