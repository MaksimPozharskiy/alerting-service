package handlers

import (
	"alerting-service/internal/logger"
	"alerting-service/internal/models"
	"alerting-service/internal/usecases"
	"alerting-service/internal/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	if err := enc.Encode(metric); err != nil {
		logger.Log.Debug("error encoding response", zap.Error(err))
		return
	}
	logger.Log.Debug("sending HTTP 200 response")
}

func (handler *metricHandler) GetMetric(w http.ResponseWriter, r *http.Request) {
	logger.Log.Debug("Received request for GetMetric",
		zap.String("method", r.Method),
		zap.String("url", r.URL.Path))

	// Проверка метода запроса
	if r.Method != http.MethodPost {
		logger.Log.Warn("Method not allowed", zap.String("received_method", r.Method))
		handleError(w, v.ErrMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Логируем тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Log.Error("Failed to read request body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	logger.Log.Debug("Request body received", zap.String("body", string(body)))

	// Восстанавливаем тело запроса для JSON-декодера
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	// Декодируем JSON
	var req models.Metrics
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		logger.Log.Warn("Cannot decode request JSON body", zap.Error(err), zap.String("body", string(body)))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	logger.Log.Debug("Decoded JSON", zap.Any("metric", req))

	// Проверяем, передан ли ID
	if req.ID == "" {
		logger.Log.Warn("Missing metric ID in request", zap.Any("metric", req))
		handleError(w, v.ErrInvalidMetricValue)
		return
	}

	// Проверяем тип метрики
	if req.MType != models.GaugeMetric && req.MType != models.CounterMetric {
		logger.Log.Warn("Invalid metric type", zap.String("received_type", req.MType))
		handleError(w, v.ErrInvalidMetricType)
		return
	}

	// Создаём объект метрики
	metric := models.Metrics{
		ID:    req.ID,
		MType: req.MType,
	}

	// Вызываем обработку данных
	logger.Log.Debug("Calling GetMetricDataProcessing", zap.Any("metric", metric))
	value, err := handler.metricUsecase.GetMetricDataProcessing(metric)
	if err != nil {
		logger.Log.Error("GetMetricDataProcessing returned an error", zap.Error(err))
		handleError(w, err)
		return
	}

	// Записываем значение в метрику
	if metric.MType == models.GaugeMetric {
		metric.Value = &value
		logger.Log.Debug("Metric is gauge", zap.Float64("value", value))
	} else {
		val := int64(value)
		metric.Delta = &val
		logger.Log.Debug("Metric is counter", zap.Int64("delta", val))
	}

	// Отправляем ответ
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	if err := enc.Encode(metric); err != nil {
		logger.Log.Error("Error encoding response", zap.Error(err))
		return
	}

	logger.Log.Debug("Successfully sent HTTP 200 response", zap.Any("metric", metric))
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

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
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

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
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

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
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

	w.Header().Set("Content-Type", "application/json")
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
