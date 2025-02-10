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

func (d *DBStorageImp) GetCounterMetric(key string) (int, bool) {
	var delta int

	row := d.db.QueryRow("SELECT delta FROM metrics WHERE type = 'counter' AND name = $1", key)
	err := row.Scan(&delta)
	if err != nil {
		panic(err)
	}

	return delta, true
}

func (d *DBStorageImp) GetGaugeMetric(key string) (float64, bool) {
	var value float64

	row := d.db.QueryRow("SELECT delta FROM metrics WHERE type = 'gauge' AND name = $1", key)
	err := row.Scan(&value)
	if err != nil {
		panic(err)
	}

	return value, true
}

func (d *DBStorageImp) UpdateGaugeMetric(metricName string, value float64) {
	sqlStatement := `
INSERT INTO metrics (name, type, value, delta)
VALUES ($1, $2, $3, NULL)
ON CONFLICT (name) DO UPDATE
SET value = EXCLUDED.value, updated_at = NOW();`

	_, err := d.db.Exec(sqlStatement, metricName, "gauge", value)
	if err != nil {
		panic(err)
	}
}

func (d *DBStorageImp) UpdateCounterMetric(metricName string, value int) {
	sqlStatement := `
INSERT INTO metrics (name, type, value, delta)
VALUES ($1, $2, $3, $4)
ON CONFLICT (name) DO UPDATE
SET delta = metrics.delta + EXCLUDED.delta, updated_at = NOW();`

	_, err := d.db.Exec(sqlStatement, metricName, "counter", 0, value)
	if err != nil {
		panic(err)
	}
}

func (d *DBStorageImp) GetMetrics() []models.Metrics {
	rows, err := d.db.Query("SELECT name, type, value, delta FROM metrics")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var allMetrics []models.Metrics

	for rows.Next() {
		var metric models.Metrics
		var value sql.NullFloat64
		var delta sql.NullInt64

		err := rows.Scan(&metric.ID, &metric.MType, &value, &delta)
		if err != nil {
			panic(err)
		}

		if value.Valid {
			metric.Value = &value.Float64
		}
		if delta.Valid {
			deltaValue := delta.Int64
			metric.Delta = &deltaValue
		}

		allMetrics = append(allMetrics, metric)
	}

	if err = rows.Err(); err != nil {
		panic(err)
	}

	return allMetrics
}

func (d *DBStorageImp) SetMetrics(allMetrics []models.Metrics) {

}

func (d *DBStorageImp) PingContext(ctx context.Context) error {
	return d.db.PingContext(ctx)
}
