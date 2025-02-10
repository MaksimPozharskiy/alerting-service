package validation

import (
	"alerting-service/internal/models"
	"errors"
	"net/http"
)

var (
	ErrMetricNotFound     = errors.New("metric not found")
	ErrInvalidMetricType  = errors.New("invalid metric type")
	ErrInvalidMetricValue = errors.New("invalid metric value")
	ErrMethodNotAllowed   = errors.New("method not allowed")
	ErrDBNotAvailable     = errors.New("database is not available")
)

var ErrMap = map[error]int{
	ErrMetricNotFound:     http.StatusNotFound,
	ErrInvalidMetricType:  http.StatusBadRequest,
	ErrInvalidMetricValue: http.StatusBadRequest,
	ErrMethodNotAllowed:   http.StatusMethodNotAllowed,
	ErrDBNotAvailable:     http.StatusInternalServerError,
}

var ValidMetricTypes = []string{models.CounterMetric, models.GaugeMetric}
var ValidCountUpdateURLParts = 5
var ValidCountGetURLParts = 4
