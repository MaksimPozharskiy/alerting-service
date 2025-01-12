package usecases

import (
	"alerting-service/internal/domain"
	repositories "alerting-service/internal/repository"
	"strconv"

	v "alerting-service/internal/validation"
)

const counterMetric = "counter"
const gaugeMetric = "gauge"

type MetricUsecase interface {
	MetricDataProcessing(domain.Metric) error
	GetMetricDataProcessing(string, string) (float64, error)
	GetAllMetricsProcessing() map[string]string
}

type MetricUsecaseImpl struct {
	storageRepository repositories.MemStorage
}

func NewMetricUsecase(storageRepository repositories.MemStorage) MetricUsecase {
	return &MetricUsecaseImpl{
		storageRepository: storageRepository,
	}
}

func (usecase *MetricUsecaseImpl) MetricDataProcessing(metric domain.Metric) error {
	metricValue, err := strconv.ParseFloat(metric.Value, 64)
	if err != nil {
		return v.ErrInvalidMetricValue
	}

	switch metric.Type {
	case counterMetric:
		if value, err := strconv.Atoi(metric.Value); err == nil {
			usecase.storageRepository.UpdateCounterMetric(metric.Name, value)
		} else {
			return v.ErrInvalidMetricValue
		}
	case gaugeMetric:
		usecase.storageRepository.UpdateGaugeMetric(metric.Name, metricValue)
	}

	return nil
}

func (usecase *MetricUsecaseImpl) GetMetricDataProcessing(metricType, metricName string) (float64, error) {
	switch metricType {
	case counterMetric:
		if value, ok := usecase.storageRepository.GetCounterMetric(metricName); ok {
			return float64(value), nil
		} else {
			return 0, v.ErrMetricNotFound
		}
	case gaugeMetric:
		if value, ok := usecase.storageRepository.GetGaugeMetric(metricName); ok {
			return value, nil
		} else {
			return 0, v.ErrMetricNotFound
		}
	}
	return 0, v.ErrInvalidMetricValue
}

func (usecase *MetricUsecaseImpl) GetAllMetricsProcessing() map[string]string {
	allMetrics := make(map[string]string)

	gauges := usecase.storageRepository.GetAllGaugeMetrics()
	counter := usecase.storageRepository.GetAllCounterMetrics()

	for key, value := range gauges {
		allMetrics[key] = strconv.FormatFloat(value, 'f', -1, 64)
	}

	for key, value := range counter {
		allMetrics[key] = strconv.Itoa(value)
	}

	return allMetrics
}
