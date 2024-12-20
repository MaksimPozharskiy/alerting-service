package handlers

import (
	"alerting-service/internal/usecases"
	"fmt"
	"net/http"
	"slices"
	"strings"
)

const counterMetric = "counter"
const gaugeMetric = "gauge"

var validMetricTypes = []string{counterMetric, gaugeMetric}
var validCountURLParts = 5

type metricHandler struct {
	metricUsecase usecases.MetricUsecase
}

func NewMetricHandler(metricUsecase usecases.MetricUsecase) *metricHandler {
	return &metricHandler{metricUsecase: metricUsecase}
}

func (handler *metricHandler) UpdateMetric(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "MethodNotAllowed", http.StatusMethodNotAllowed)
		return
	}

	urlData := strings.Split(req.URL.Path, "/")
	if len(urlData) != validCountURLParts {
		http.Error(w, "Metric not found", http.StatusNotFound)
		return
	}

	metricType := urlData[2]
	if !slices.Contains(validMetricTypes, metricType) {
		fmt.Println(metricType)
		http.Error(w, "Incorrect metrics type", http.StatusBadRequest)
		return
	}

	metricName := urlData[3]

	metricStr := urlData[4]

	err := handler.metricUsecase.MetricDataProcessing(metricType, metricName, metricStr)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
	}

	// @TODO
	// Вынести парсинг метрик в функци
	// сделать валидацию отдельно
	// подумать че сделать с моделями, куда девать типа counter и gauge
	// сделать создание сервера по чистой архитектуре

	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("2222"))
}
