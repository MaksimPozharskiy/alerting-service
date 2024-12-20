package usecases

import (
	repositories "alerting-service/internal/repository"
	"fmt"
	"math"
	"strconv"
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
		return fmt.Errorf("incorrect metric value")
	}

	metricValueIsInt := metricValue == math.Trunc(metricValue)

	if metricType == counterMetric {

		if metricValueIsInt {
			if v, err := strconv.Atoi(valueStr); err == nil {
				usecase.storageRepository.UpdateCounterMetic(metricType, v)
			} else {
				return fmt.Errorf("incorrect metric value")
			}
		} else {
			return fmt.Errorf("incorrect metric value")
		}
	}

	if metricType == gaugeMetric {
		if metricValueIsInt {
			return fmt.Errorf("incorrect metric value")
		} else {
			usecase.storageRepository.UpdateGaugeMetic(metricType, metricValue)
		}
	}

	return nil
}
