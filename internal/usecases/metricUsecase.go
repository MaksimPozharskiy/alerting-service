package usecases

import (
	repositories "alerting-service/internal/repository"
	"math"
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

	metricValueIsInt := metricValue == math.Trunc(metricValue)

	if metricType == counterMetric {
		if metricValueIsInt {
			if value, err := strconv.Atoi(valueStr); err == nil {
				usecase.storageRepository.UpdateCounterMetic(metricType, value)
			} else {
				return v.ErrInvalidMetricValue
			}
		} else {
			return v.ErrInvalidMetricValue
		}
	}

	if metricType == gaugeMetric {
		if metricValueIsInt {
			return v.ErrInvalidMetricValue
		} else {
			usecase.storageRepository.UpdateGaugeMetic(metricType, metricValue)
		}
	}

	return nil
}
