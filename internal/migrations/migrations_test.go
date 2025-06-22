package migrations

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestInitDB(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_DSN")
	if dsn == "" {
		t.Skip("TEST_DATABASE_DSN not set")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("failed to connect to PostgreSQL: %v", err)
	}
	defer db.Close()

	_, _ = db.Exec("DROP TABLE IF EXISTS metrics")

	if err := InitDB(db); err != nil {
		t.Fatalf("InitDB returned error: %v", err)
	}

	rows, err := db.Query("SELECT id, name, type, value, delta, updated_at FROM metrics")
	if err != nil {
		t.Fatalf("table 'metrics' does not exist or query failed: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, name, mtype string
		var value sql.NullFloat64
		var delta sql.NullInt64
		var updatedAt sql.NullTime

		if err := rows.Scan(&id, &name, &mtype, &value, &delta, &updatedAt); err != nil {
			t.Errorf("failed to scan row: %v", err)
		}
	}

	if err := rows.Err(); err != nil {
		t.Errorf("rows iteration error: %v", err)
	}
}
