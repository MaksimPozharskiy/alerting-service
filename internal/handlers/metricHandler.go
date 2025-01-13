package handlers

import (
	"alerting-service/internal/usecases"
	"fmt"
	"net/http"
	"strconv"

	utils "alerting-service/internal/utils"
	v "alerting-service/internal/validation"
)

type metricHandler struct {
	metricUsecase usecases.MetricUsecase
}

func NewMetricHandler(metricUsecase usecases.MetricUsecase) *metricHandler {
	return &metricHandler{metricUsecase: metricUsecase}
}

func (handler *metricHandler) UpdateMetric(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		handleError(w, v.ErrMethodNotAllowed)
		return
	}

	metric, err := utils.ParseUpdateMetricURL(req.URL.Path)
	if err != nil {
		handleError(w, err)
		return
	}

	err = handler.metricUsecase.MetricDataProcessing(metric)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func (handler *metricHandler) GetMetric(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		handleError(w, v.ErrMethodNotAllowed)
		return
	}

	metric, err := utils.ParseGetMetricURL(req.URL.Path)
	if err != nil {
		handleError(w, err)
		return
	}

	value, err := handler.metricUsecase.GetMetricDataProcessing(metric.Type, metric.Name)
	if err != nil {
		handleError(w, err)
		return
	}

	metric.Value = strconv.FormatFloat(value, 'f', -1, 64)

	w.Header().Set("Content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(metric.Value))
}

func (handler *metricHandler) GetAllMetrics(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		handleError(w, v.ErrMethodNotAllowed)
		return
	}

	allMetrics := handler.metricUsecase.GetMetrics()

	w.Header().Set("Content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	for key, value := range allMetrics {
		w.Write([]byte(fmt.Sprintf("%s: %s\n", key, value)))
	}
}

func handleError(w http.ResponseWriter, err error) {
	statusCode, ok := v.ErrMap[err]

	if !ok {
		statusCode = http.StatusInternalServerError
	}

	http.Error(w, fmt.Sprint(err), statusCode)
}
