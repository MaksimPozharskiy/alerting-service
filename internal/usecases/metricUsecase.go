package usecases

import (
	repositories "alerting-service/internal/repository"
	"strconv"

	v "alerting-service/internal/validation"
)

const counterMetric = "counter"
const gaugeMetric = "gauge"

type MetricUsecase interface {
	MetricDataProcessing(string, string, string) error
}

type MetricUsecaseImpl struct {
	storageRepository repositories.MemStorage
}

func NewMetricUsecase(storageRepository repositories.MemStorage) MetricUsecase {
	return &MetricUsecaseImpl{
		storageRepository: storageRepository,
	}
}

func (usecase *MetricUsecaseImpl) MetricDataProcessing(metricType, metricName, valueStr string) error {
	metricValue, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return v.ErrInvalidMetricValue
	}

	switch metricType {
	case counterMetric:
		if value, err := strconv.Atoi(valueStr); err == nil {
			usecase.storageRepository.UpdateCounterMetic(metricType, value)
		} else {
			return v.ErrInvalidMetricValue
		}
	case gaugeMetric:
		usecase.storageRepository.UpdateGaugeMetic(metricType, metricValue)
	}

	return nil
}
