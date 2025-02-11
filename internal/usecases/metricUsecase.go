package usecases

import (
	"alerting-service/internal/logger"
	"alerting-service/internal/models"
	"alerting-service/internal/repository"

	v "alerting-service/internal/validation"

	"go.uber.org/zap"
)

type MetricUsecase interface {
	MetricDataProcessing(models.Metrics) error
	GetMetricDataProcessing(models.Metrics) (float64, error)
	GetMetrics() ([]models.Metrics, error)
	UpdateMetrics([]models.Metrics) error
}

type MetricUsecaseImpl struct {
	storageRepository repository.StorageRepository
}

func NewMetricUsecase(storageRepository repository.StorageRepository) MetricUsecase {
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
		if value, ok, _ := usecase.storageRepository.GetCounterMetric(metric.ID); ok {
			return float64(value), nil
		} else {
			return 0, v.ErrMetricNotFound
		}
	case models.GaugeMetric:
		if value, ok, _ := usecase.storageRepository.GetGaugeMetric(metric.ID); ok {
			return value, nil
		} else {
			return 0, v.ErrMetricNotFound
		}
	}
	return 0, v.ErrInvalidMetricValue
}

func (usecase *MetricUsecaseImpl) GetMetrics() ([]models.Metrics, error) {
	metrics, err := usecase.storageRepository.GetMetrics()
	if err != nil {
		return nil, err
	}
	return metrics, nil
}

func (usecase *MetricUsecaseImpl) UpdateMetrics(metrics []models.Metrics) error {
	logger.Log.Debug("Entering UpdateMetrics in usecase", zap.Int("metrics_count", len(metrics)))

	err := usecase.storageRepository.UpdateMetrics(metrics)
	if err != nil {
		logger.Log.Error("Error updating metrics in storage repository", zap.Error(err))
		return err
	}

	logger.Log.Debug("Successfully updated metrics in storage repository")
	return nil
}
