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

	var metric models.Metrics

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&req)

	if err != nil {
		logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
		metric, err = utils.ParseUpdateMetricURL(r.URL.Path)
		if err != nil {
			handleError(w, err)
			return
		}
	} else {
		metric = models.Metrics{
			ID:    req.ID,
			MType: req.MType,
		}

		if req.MType == "gauge" {
			metric.Value = req.Value
		} else {
			metric.Delta = req.Delta
		}
	}

	handler.metricUsecase.MetricDataProcessing(metric)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (handler *metricHandler) GetMetric(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handleError(w, v.ErrMethodNotAllowed)
		return
	}

	logger.Log.Debug("decoding request")
	var metric models.Metrics

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&r)

	if err != nil {
		logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
		metric, err = utils.ParseGetMetricURL(r.URL.Path)
		if err != nil {
			handleError(w, err)
			return
		}
	} else {
		metric = models.Metrics{
			ID:    metric.ID,
			MType: metric.MType,
		}
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

	if metric.MType == "gauge" {
		v := *metric.Value
		w.Write([]byte(fmt.Sprint(v)))
	} else {
		v := *metric.Delta
		w.Write([]byte(fmt.Sprint(v)))
	}
	logger.Log.Debug("sending HTTP 200 response")
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
