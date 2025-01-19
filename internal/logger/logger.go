package logger

import (
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

var Log *zap.Logger = zap.NewNop()

type loggerResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func NewLoggerResponseWriter(w http.ResponseWriter) *loggerResponseWriter {
	return &loggerResponseWriter{
		ResponseWriter: w,
	}
}

func (lrw *loggerResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (r *loggerResponseWriter) Write(body []byte) (int, error) {
	var err error
	r.size, err = r.ResponseWriter.Write(body)

	return r.size, err
}

func Initialize(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl

	zl, err := cfg.Build()
	if err != nil {
		return err
	}

	Log = zl
	return nil
}

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(startTime)
		Log.Info("got incoming HTTP request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("duration", duration.String()),
		)
	})
}

func ResponseLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		lrw := NewLoggerResponseWriter(w)
		h.ServeHTTP(lrw, req)

		statusCode := lrw.statusCode
		size := lrw.size
		Log.Info("send outcoming HTTP response",
			zap.String("method", fmt.Sprint(statusCode)),
			zap.String("bodySize", fmt.Sprint(size)),
		)
	})
}
