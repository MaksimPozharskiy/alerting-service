package db

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connect(connectParams string) (*sql.DB, error) {
	db, err := sql.Open("pgx", connectParams)
	if err != nil {
		return nil, err
	}

	err = initDB(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func initDB(db *sql.DB) error {
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
