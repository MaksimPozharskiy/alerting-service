package repository

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func setupTestDB(t *testing.T) *sql.DB {
	dsn := os.Getenv("TEST_DATABASE_DSN")
	if dsn == "" {
		t.Skip("TEST_DATABASE_DSN not set")
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	_, _ = db.Exec("DROP TABLE IF EXISTS metrics")
	_, _ = db.Exec(`
	CREATE TABLE metrics (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		type TEXT CHECK (type IN ('gauge', 'counter')) NOT NULL,
		value DOUBLE PRECISION,
		delta BIGINT,
		updated_at TIMESTAMP DEFAULT NOW()
	)`)
	return db
}

func TestDBStorage_UpdateAndGetGauge(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewDBStorageRepository(db)

	err := repo.UpdateGaugeMetric("gauge_test", 3.14)
	if err != nil {
		t.Fatalf("update gauge failed: %v", err)
	}

	val, ok, err := repo.GetGaugeMetric("gauge_test")
	if err != nil || !ok {
		t.Fatalf("get gauge failed: %v", err)
	}
	if val != 3.14 {
		t.Errorf("expected 3.14, got %f", val)
	}
}

func TestDBStorage_UpdateAndGetCounter(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewDBStorageRepository(db)

	err := repo.UpdateCounterMetric("counter_test", 10)
	if err != nil {
		t.Fatalf("update counter failed: %v", err)
	}
	err = repo.UpdateCounterMetric("counter_test", 5)
	if err != nil {
		t.Fatalf("second update failed: %v", err)
	}

	val, ok, err := repo.GetCounterMetric("counter_test")
	if err != nil || !ok {
		t.Fatalf("get counter failed: %v", err)
	}
	if val != 15 {
		t.Errorf("expected 15, got %d", val)
	}
}

func TestDBStorage_GetMetrics(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewDBStorageRepository(db)

	_ = repo.UpdateGaugeMetric("g1", 1.23)
	_ = repo.UpdateCounterMetric("c1", 7)

	metrics, err := repo.GetMetrics()
	if err != nil {
		t.Fatalf("get metrics failed: %v", err)
	}
	if len(metrics) != 2 {
		t.Errorf("expected 2 metrics, got %d", len(metrics))
	}
}
