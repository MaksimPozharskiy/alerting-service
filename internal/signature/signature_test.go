package signature

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHashMiddleware_ValidHash(t *testing.T) {
	SetServerHashKey("secret")

	body := []byte(`{"test":"data"}`)
	hash := GetHash(body, []byte("secret"))

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set(HashSHA256, hash)

	rec := httptest.NewRecorder()
	HashMiddleware(handler).ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", res.StatusCode)
	}
}

func TestHashMiddleware_InvalidHash(t *testing.T) {
	SetServerHashKey("secret")

	body := []byte(`{"test":"data"}`)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set(HashSHA256, "invalidhash")

	rec := httptest.NewRecorder()
	HashMiddleware(handler).ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	respBody, _ := io.ReadAll(res.Body)

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request, got %d", res.StatusCode)
	}
	if string(respBody) != "invalid hash\n" {
		t.Errorf("unexpected response body: %s", string(respBody))
	}
}

func TestGetHash_Deterministic(t *testing.T) {
	key := []byte("key")
	data := []byte("some data")

	hash1 := GetHash(data, key)
	hash2 := GetHash(data, key)

	if hash1 != hash2 {
		t.Errorf("expected hashes to match: %s != %s", hash1, hash2)
	}
}
