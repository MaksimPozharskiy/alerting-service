package repository

import (
	"alerting-service/internal/models"
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type DBStorageImp struct {
	db *sql.DB
}

func NewDBStorageRepository(db *sql.DB) StorageRepository {
	return &DBStorageImp{db: db}
}

func (d *DBStorageImp) GetCounterMetric(key string) (int, bool, error) {
	var delta int

	row := d.db.QueryRow("SELECT delta FROM metrics WHERE type = 'counter' AND name = $1", key)
	err := row.Scan(&delta)
	if err != nil {
		return 0, false, err
	}

	return delta, true, nil
}

func (d *DBStorageImp) GetGaugeMetric(key string) (float64, bool, error) {
	var value float64

	row := d.db.QueryRow("SELECT delta FROM metrics WHERE type = 'gauge' AND name = $1", key)
	err := row.Scan(&value)
	if err != nil {
		return 0, false, err
	}

	return value, true, nil
}

func (d *DBStorageImp) UpdateGaugeMetric(metricName string, value float64) error {
	sqlStatement := `
INSERT INTO metrics (name, type, value, delta)
VALUES ($1, $2, $3, NULL)
ON CONFLICT (name) DO UPDATE
SET value = EXCLUDED.value, updated_at = NOW();`

	_, err := d.db.Exec(sqlStatement, metricName, "gauge", value)
	if err != nil {
		return err
	}

	return nil
}

func (d *DBStorageImp) UpdateCounterMetric(metricName string, value int) error {
	sqlStatement := `
INSERT INTO metrics (name, type, value, delta)
VALUES ($1, $2, $3, $4)
ON CONFLICT (name) DO UPDATE
SET delta = metrics.delta + EXCLUDED.delta, updated_at = NOW();`

	_, err := d.db.Exec(sqlStatement, metricName, "counter", 0, value)
	if err != nil {
		return err
	}

	return nil
}

func (d *DBStorageImp) GetMetrics() ([]models.Metrics, error) {
	rows, err := d.db.Query("SELECT name, type, value, delta FROM metrics")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allMetrics []models.Metrics

	for rows.Next() {
		var metric models.Metrics
		var value sql.NullFloat64
		var delta sql.NullInt64

		err := rows.Scan(&metric.ID, &metric.MType, &value, &delta)
		if err != nil {
			return nil, err
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
		return nil, err
	}

	return allMetrics, nil
}

func (d *DBStorageImp) SetMetrics(allMetrics []models.Metrics) {

}

func (d *DBStorageImp) UpdateMetrics(metrics []models.Metrics) error {
	if len(metrics) == 0 {
		return nil
	}

	sqlStr := "INSERT INTO metrics (name, type, value, delta) VALUES "
	var values []interface{}
	var placeholders []string

	for i, metric := range metrics {
		startIdx := i*4 + 1
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d)", startIdx, startIdx+1, startIdx+2, startIdx+3))

		if metric.MType == models.GaugeMetric {
			values = append(values, metric.ID, metric.MType, metric.Value, nil)
		} else {
			values = append(values, metric.ID, metric.MType, nil, metric.Delta)
		}
	}

	sqlStr += strings.Join(placeholders, ", ")

	sqlStr += ` ON CONFLICT (name) DO UPDATE
	SET delta = COALESCE(metrics.delta, 0) + COALESCE(EXCLUDED.delta, 0),
	    value = COALESCE(EXCLUDED.value, metrics.value);`

	_, err := d.db.Exec(sqlStr, values...)
	if err != nil {
		return err
	}

	return nil
}

func (d *DBStorageImp) PingContext(ctx context.Context) error {
	return d.db.PingContext(ctx)
}
