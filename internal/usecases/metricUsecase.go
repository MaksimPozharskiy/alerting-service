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
	GetMetrics() map[string]string
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
	switch metric.Type {
	case counterMetric:
		value, err := strconv.Atoi(metric.Value)
		if err != nil {
			return v.ErrInvalidMetricValue
		}
		usecase.storageRepository.UpdateCounterMetric(metric.Name, value)
	case gaugeMetric:
		value, err := strconv.ParseFloat(metric.Value, 64)
		if err != nil {
			return v.ErrInvalidMetricValue
		}
		usecase.storageRepository.UpdateGaugeMetric(metric.Name, value)
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

func (usecase *MetricUsecaseImpl) GetMetrics() map[string]string {
	return usecase.storageRepository.GetMetrics()
}
