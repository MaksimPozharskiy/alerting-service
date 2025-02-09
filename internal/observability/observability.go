package observability

import (
	"alerting-service/internal/logger"
	v "alerting-service/internal/validation"
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"
)

type obsHandler struct {
	db *sql.DB
}

func NewObsHandler(db *sql.DB) *obsHandler {
	return &obsHandler{db: db}
}

func (h *obsHandler) HealthCheckDB(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handleError(w, v.ErrDBNotAvailable)
		return
	}

	logger.Log.Debug("ping request")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := h.db.PingContext(ctx); err != nil {
		handleError(w, v.ErrDBNotAvailable)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))

	logger.Log.Debug("sending HTTP 200 response")
}

func handleError(w http.ResponseWriter, err error) {
	statusCode, ok := v.ErrMap[err]

	if !ok {
		statusCode = http.StatusInternalServerError
	}

	http.Error(w, fmt.Sprint(err), statusCode)
}
