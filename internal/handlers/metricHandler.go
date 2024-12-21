package handlers

import (
	"alerting-service/internal/usecases"
	"errors"
	"fmt"
	"net/http"

	utils "alerting-service/internal/utils"
)

type metricHandler struct {
	metricUsecase usecases.MetricUsecase
}

var ErrInvalidMetricValue = errors.New("invalid metric type")
var ErrMethodNotAllowred = errors.New("method not allowed")

var errMap = map[error]int{
	utils.ErrMetricNotFound:    http.StatusNotFound,
	utils.ErrInvalidMetricType: http.StatusBadRequest,
	ErrInvalidMetricValue:      http.StatusBadRequest,
	ErrMethodNotAllowred:       http.StatusMethodNotAllowed,
}

func NewMetricHandler(metricUsecase usecases.MetricUsecase) *metricHandler {
	return &metricHandler{metricUsecase: metricUsecase}
}

func (handler *metricHandler) UpdateMetric(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		handleError(w, ErrMethodNotAllowred)
		return
	}

	metric, err := utils.ParseMetricURL(req.URL.Path)
	if err != nil {
		handleError(w, err)
		return
	}

	err = handler.metricUsecase.MetricDataProcessing(metric.Type, metric.Name, metric.Value)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func handleError(w http.ResponseWriter, err error) {
	statusCode, ok := errMap[err]
	if !ok {
		statusCode = http.StatusInternalServerError
	}

	http.Error(w, fmt.Sprint(err), statusCode)
}
