package handlers

import (
	"alerting-service/internal/logger"
	"alerting-service/internal/models"
	"alerting-service/internal/usecases"
	"alerting-service/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"

	v "alerting-service/internal/validation"

	"go.uber.org/zap"
)

type metricHandler struct {
	metricUsecase usecases.MetricUsecase
}

func NewMetricHandler(metricUsecase usecases.MetricUsecase) *metricHandler {
	return &metricHandler{metricUsecase: metricUsecase}
}
func (handler *metricHandler) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		handleError(w, v.ErrMethodNotAllowed)
		return
	}

	logger.Log.Debug("decoding request")
	var req models.Metrics

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	metric := models.Metrics{
		ID:    req.ID,
		MType: req.MType,
	}

	if req.MType == "gauge" {
		metric.Value = req.Value
	} else {
		metric.Delta = req.Delta
	}

	handler.metricUsecase.MetricDataProcessing(metric)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	if err := enc.Encode(metric); err != nil {
		logger.Log.Debug("error encoding response", zap.Error(err))
		return
	}
	logger.Log.Debug("sending HTTP 200 response")
}

func (handler *metricHandler) GetMetric(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		handleError(w, v.ErrMethodNotAllowed)
		return
	}

	logger.Log.Debug("decoding request")
	var req models.Metrics

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	metric := models.Metrics{
		ID:    req.ID,
		MType: req.MType,
	}

	value, err := handler.metricUsecase.GetMetricDataProcessing(metric)
	if err != nil {
		handleError(w, err)
		return
	}

	if metric.MType == "gauge" {
		metric.Value = &value
	} else {
		val := int64(value)
		metric.Delta = &val
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	if err := enc.Encode(metric); err != nil {
		logger.Log.Debug("error encoding response", zap.Error(err))
		return
	}
	logger.Log.Debug("sending HTTP 200 response")
}

func (handler *metricHandler) UpdateURLMetric(w http.ResponseWriter, req *http.Request) {
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

func (handler *metricHandler) GetURLMetric(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		handleError(w, v.ErrMethodNotAllowed)
		return
	}

	metric, err := utils.ParseGetMetricURL(req.URL.Path)
	if err != nil {
		handleError(w, err)
		return
	}

	value, err := handler.metricUsecase.GetMetricDataProcessing(metric)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte(fmt.Sprint(value)))
}

func (handler *metricHandler) GetAllMetrics(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		handleError(w, v.ErrMethodNotAllowed)
		return
	}

	allMetrics := handler.metricUsecase.GetMetrics()

	w.Header().Set("Content-type", "text/html; charset=utf-8")
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
