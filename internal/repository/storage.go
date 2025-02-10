package repository

import (
	"alerting-service/internal/models"
)

type StorageRepository interface {
	GetCounterMetric(string) (int, bool)
	GetGaugeMetric(string) (float64, bool)
	UpdateGaugeMetric(string, float64)
	UpdateCounterMetric(string, int)
	GetMetrics() []models.Metrics
	SetMetrics([]models.Metrics)
}
