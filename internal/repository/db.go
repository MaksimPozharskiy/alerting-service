package repository

import (
	"alerting-service/internal/logger"
	"alerting-service/internal/models"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
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

	row := d.db.QueryRow("SELECT value FROM metrics WHERE type = 'gauge' AND name = $1", key)
	err := row.Scan(&value)

	if err != nil {
		return 0, false, err
	}

	return value, true, nil
}

func (d *DBStorageImp) UpdateGaugeMetric(metricName string, value float64) error {
	query := `
INSERT INTO metrics (name, type, value, delta)
VALUES ($1, $2, $3, NULL)
ON CONFLICT (name) DO UPDATE
SET value = EXCLUDED.value, updated_at = NOW();`

	stmt, err := d.db.PrepareContext(context.Background(), query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = d.retryExecute(context.Background(), stmt, metricName, "gauge", value)
	return err
}

func (d *DBStorageImp) UpdateCounterMetric(metricName string, value int) error {
	query := `
INSERT INTO metrics (name, type, value, delta)
VALUES ($1, $2, $3, $4)
ON CONFLICT (name) DO UPDATE
SET delta = metrics.delta + EXCLUDED.delta, updated_at = NOW();`

	stmt, err := d.db.PrepareContext(context.Background(), query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = d.retryExecute(context.Background(), stmt, metricName, "counter", nil, value)
	return err

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
		logger.Log.Debug("No metrics provided for update")
		return nil
	}

	tx, err := d.db.Begin()
	if err != nil {
		logger.Log.Error("Failed to begin transaction", zap.Error(err))
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `INSERT INTO metrics (name, type, value, delta) 
              VALUES ($1, $2, $3, $4) 
              ON CONFLICT (name) DO UPDATE 
              SET delta = COALESCE(metrics.delta, 0) + COALESCE(EXCLUDED.delta, 0), 
                  value = COALESCE(EXCLUDED.value, metrics.value);`

	stmt, err := d.db.PrepareContext(context.Background(), query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, metric := range metrics {
		logger.Log.Debug("Processing metric", zap.String("id", metric.ID), zap.String("type", metric.MType))

		var value, delta interface{}
		if metric.MType == models.GaugeMetric {
			value = metric.Value
			delta = nil
		} else {
			value = nil
			delta = metric.Delta
		}

		_, err := d.retryExecute(context.Background(), stmt, metric.ID, metric.MType, value, delta)
		if err != nil {
			logger.Log.Error("Error executing SQL query", zap.String("metric_id", metric.ID), zap.Error(err))
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		logger.Log.Error("Failed to commit transaction", zap.Error(err))
		return err
	}

	logger.Log.Debug("Successfully updated all metrics in database")
	return nil
}

func (d *DBStorageImp) PingContext(ctx context.Context) error {
	return d.db.PingContext(ctx)
}

func (d *DBStorageImp) retryExecute(ctx context.Context, stmt *sql.Stmt, args ...any) (sql.Result, error) {
	var err error
	var result sql.Result

	retryDelays := []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}

	for attempt := 0; attempt <= len(retryDelays); attempt++ {
		result, err = stmt.ExecContext(ctx, args...)
		if err == nil {
			return result, nil
		}

		if isRetriableError(err) {
			logger.Log.Warn("Retriable error occurred, retrying...",
				zap.Error(err), zap.Int("attempt", attempt+1))

			if attempt < len(retryDelays) {
				time.Sleep(retryDelays[attempt])
				continue
			}
		}

		logger.Log.Error("Non-retriable error occurred", zap.Error(err))
		return nil, err
	}

	return nil, err
}

func isRetriableError(err error) bool {
	var pqErr *pgconn.PgError
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case pgerrcode.ConnectionException, pgerrcode.SerializationFailure,
			pgerrcode.DeadlockDetected, pgerrcode.StatementCompletionUnknown:
			return true
		}
	}
	return false
}
