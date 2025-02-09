package repository

import (
	"alerting-service/internal/models"
	"context"
	"database/sql"
)

type DBStorageImp struct {
	db *sql.DB
}

func NewDBStorageRepository(db *sql.DB) StorageRepository {
	return &DBStorageImp{db: db}
}

func (s *DBStorageImp) GetCounterMetric(key string) (int, bool) {
	return 0, false
}

func (s *DBStorageImp) GetGaugeMetric(key string) (float64, bool) {
	return 0, false
}

func (s *DBStorageImp) UpdateGaugeMetric(metricName string, value float64) {

}

func (s *DBStorageImp) UpdateCounterMetric(metricName string, value int) {

}

func (s *DBStorageImp) GetMetrics() []models.Metrics {
	return []models.Metrics{}
}

func (s *DBStorageImp) SetMetrics(allMetrics []models.Metrics) {

}

func (s *DBStorageImp) PingContext(ctx context.Context) error {
	return s.db.PingContext(ctx)
}
