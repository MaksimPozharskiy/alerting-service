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

	if req.MType == models.GaugeMetric {
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

	if metric.MType == models.GaugeMetric {
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

	allMetrics, err := handler.metricUsecase.GetMetrics()
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	for _, metric := range allMetrics {
		mType := metric.MType
		if mType == models.GaugeMetric {
			w.Write([]byte(fmt.Sprintf("%s: %f\n", metric.ID, *metric.Value)))
		} else {
			w.Write([]byte(fmt.Sprintf("%s: %d\n", metric.ID, *metric.Delta)))
		}
	}
}
func (handler *metricHandler) UpdateMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logger.Log.Warn("Invalid request method", zap.String("method", r.Method))
		handleError(w, v.ErrMethodNotAllowed)
		return
	}

	logger.Log.Debug("Decoding request body for batch metric update")
	var metrics []models.Metrics

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&metrics); err != nil {
		logger.Log.Error("Cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.Log.Debug("Received metrics for update", zap.Int("metrics_count", len(metrics)), zap.Any("metrics", metrics))

	err := handler.metricUsecase.UpdateMetrics(metrics)
	if err != nil {
		logger.Log.Error("Error updating metrics", zap.Error(err))
		handleError(w, err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	logger.Log.Debug("Successfully processed batch update, sending HTTP 200 response")
}

func handleError(w http.ResponseWriter, err error) {
	statusCode, ok := v.ErrMap[err]

	if !ok {
		statusCode = http.StatusInternalServerError
	}

	http.Error(w, fmt.Sprint(err), statusCode)
}
