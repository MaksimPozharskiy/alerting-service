package usecases

import (
	"alerting-service/internal/models"
	repositories "alerting-service/internal/repository"

	v "alerting-service/internal/validation"
)

const counterMetric = "counter"
const gaugeMetric = "gauge"

type MetricUsecase interface {
	MetricDataProcessing(models.Metrics) error
	GetMetricDataProcessing(models.Metrics) (float64, error)
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

func (usecase *MetricUsecaseImpl) MetricDataProcessing(metric models.Metrics) error {
	switch metric.MType {
	case counterMetric:
		usecase.storageRepository.UpdateCounterMetric(metric.ID, int(*metric.Delta))
	case gaugeMetric:
		usecase.storageRepository.UpdateGaugeMetric(metric.ID, *metric.Value)
	}

	return nil
}

func (usecase *MetricUsecaseImpl) GetMetricDataProcessing(metric models.Metrics) (float64, error) {
	switch metric.MType {
	case counterMetric:
		if value, ok := usecase.storageRepository.GetCounterMetric(metric.ID); ok {
			return float64(value), nil
		} else {
			return 0, v.ErrMetricNotFound
		}
	case gaugeMetric:
		if value, ok := usecase.storageRepository.GetGaugeMetric(metric.ID); ok {
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
