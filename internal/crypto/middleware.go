package crypto

import (
	"bytes"
	"crypto/rsa"
	"io"
	"net/http"

	"alerting-service/internal/logger"

	"go.uber.org/zap"
)

// DecryptionMiddleware decrypts incoming requests if they are encrypted.
func DecryptionMiddleware(privateKey *rsa.PrivateKey) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if privateKey == nil || r.Header.Get("Content-Encryption") != "RSA" {
				next.ServeHTTP(w, r)
				return
			}

			logger.Log.Debug("Decrypting incoming request")

			encryptedData, err := io.ReadAll(r.Body)
			if err != nil {
				logger.Log.Error("Failed to read encrypted body", zap.Error(err))
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}
			defer r.Body.Close()

			decryptedData, err := DecryptData(encryptedData, privateKey)
			if err != nil {
				logger.Log.Error("Failed to decrypt body", zap.Error(err))
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}

			r.Body = io.NopCloser(bytes.NewReader(decryptedData))
			r.ContentLength = int64(len(decryptedData))
			r.Header.Del("Content-Encryption")

			next.ServeHTTP(w, r)
		})
	}
}
