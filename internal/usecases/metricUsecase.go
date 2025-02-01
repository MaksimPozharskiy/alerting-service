package usecases

import (
	"alerting-service/internal/models"

	v "alerting-service/internal/validation"
)

type MemStorage interface {
	GetCounterMetric(string) (int, bool)
	GetGaugeMetric(string) (float64, bool)
	UpdateGaugeMetric(string, float64)
	UpdateCounterMetric(string, int)
	GetMetrics() []models.Metrics
	SetMetrics([]models.Metrics)
}

type MetricUsecase interface {
	MetricDataProcessing(models.Metrics) error
	GetMetricDataProcessing(models.Metrics) (float64, error)
	GetMetrics() []models.Metrics
}

type MetricUsecaseImpl struct {
	storageRepository MemStorage
}

func NewMetricUsecase(storageRepository MemStorage) MetricUsecase {
	return &MetricUsecaseImpl{
		storageRepository: storageRepository,
	}
}

func (usecase *MetricUsecaseImpl) MetricDataProcessing(metric models.Metrics) error {
	switch metric.MType {
	case models.CounterMetric:
		usecase.storageRepository.UpdateCounterMetric(metric.ID, int(*metric.Delta))
	case models.GaugeMetric:
		usecase.storageRepository.UpdateGaugeMetric(metric.ID, *metric.Value)
	}

	return nil
}

func (usecase *MetricUsecaseImpl) GetMetricDataProcessing(metric models.Metrics) (float64, error) {
	switch metric.MType {
	case models.CounterMetric:
		if value, ok := usecase.storageRepository.GetCounterMetric(metric.ID); ok {
			return float64(value), nil
		} else {
			return 0, v.ErrMetricNotFound
		}
	case models.GaugeMetric:
		if value, ok := usecase.storageRepository.GetGaugeMetric(metric.ID); ok {
			return value, nil
		} else {
			return 0, v.ErrMetricNotFound
		}
	}
	return 0, v.ErrInvalidMetricValue
}

func (usecase *MetricUsecaseImpl) GetMetrics() []models.Metrics {
	return usecase.storageRepository.GetMetrics()
}
