package signature

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
)

var HashSHA256 = "HashSHA256"

var hashKey []byte

func SetServerHashKey(key string) {
	hashKey = []byte(key)
}

func GetHash(data []byte, hashKey []byte) string {
	hash := hmac.New(sha256.New, hashKey)
	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}

func HashMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if string(hashKey) == "" || req.Header.Get(HashSHA256) == "" {
			next.ServeHTTP(w, req)
			return
		}

		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			http.Error(w, "Unable to read request body", http.StatusInternalServerError)
			return
		}

		req.Body = io.NopCloser(bytes.NewReader(bodyBytes))

		expectedHash := GetHash(bodyBytes, hashKey)

		receivedHash := req.Header.Get("HashSHA256")

		if hmac.Equal([]byte(expectedHash), []byte(receivedHash)) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, req)
	})
}
