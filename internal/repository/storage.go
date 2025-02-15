package repository

import (
	"alerting-service/internal/models"
)

type StorageRepository interface {
	GetCounterMetric(string) (int, bool, error)
	GetGaugeMetric(string) (float64, bool, error)
	UpdateGaugeMetric(string, float64) error
	UpdateCounterMetric(string, int) error
	GetMetrics() ([]models.Metrics, error)
	SetMetrics([]models.Metrics)
	UpdateMetrics([]models.Metrics) error
}
