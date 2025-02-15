package migrations

import "database/sql"

func InitDB(db *sql.DB) error {
	schema := `
    CREATE TABLE IF NOT EXISTS metrics (
        id SERIAL PRIMARY KEY,
        name TEXT NOT NULL UNIQUE,
        type TEXT CHECK (type IN ('gauge', 'counter')) NOT NULL,
        value DOUBLE PRECISION,
        delta BIGINT,
        updated_at TIMESTAMP DEFAULT NOW()
    );`
	_, err := db.Exec(schema)
	return err
}
