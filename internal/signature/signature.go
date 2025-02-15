package signature

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"strings"
)

var HashSHA256 = "HashSHA256"

var hashKey []byte

func SetServerHashKey(key string) {
	hashKey = []byte(key)
}

func GetHash(data []byte, hashKey string) string {
	hash := hmac.New(sha256.New, []byte(hashKey))
	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}

func SignRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Unable to read request body", http.StatusInternalServerError)
			return
		}

		r.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))
		receivedHash := r.Header.Get("HashSHA256")
		expectedSignature := GetHash(bodyBytes, string(hashKey))

		if !hmac.Equal([]byte(receivedHash), []byte(expectedSignature)) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func SignResponse(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if string(hashKey) == "" || req.Header.Get(HashSHA256) == "" {
			h.ServeHTTP(w, req)
			return
		}

		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			http.Error(w, "Unable to read request body", http.StatusInternalServerError)
			return
		}

		req.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))

		signature := GetHash(bodyBytes, string(hashKey))

		w.Header().Set(HashSHA256, signature)

		h.ServeHTTP(w, req)
	})
}
