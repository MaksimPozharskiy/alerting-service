package signature

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
)

var HashSHA256 = "HashSHA256"

var hashKey []byte

func SetServerHashKey(key string) {
	hashKey = []byte(key)
}

func GetHash(hashKey string) string {
	hash := hmac.New(sha256.New, []byte(hashKey))
	signature := hex.EncodeToString(hash.Sum(nil))

	return signature
}

func SignRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		signature := r.Header.Get(HashSHA256)
		hash := hmac.New(sha256.New, hashKey)

		expectedSignature := hex.EncodeToString(hash.Sum(nil))

		if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func SignResponse(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		hash := hmac.New(sha256.New, hashKey)
		signature := hex.EncodeToString(hash.Sum(nil))

		w.Header().Set(HashSHA256, signature)

		h.ServeHTTP(w, req)
	})
}
