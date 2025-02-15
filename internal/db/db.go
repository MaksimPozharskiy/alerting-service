package db

import (
	"alerting-service/internal/migrations"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connect(connectParams string) (*sql.DB, error) {
	db, err := sql.Open("pgx", connectParams)
	if err != nil {
		return nil, err
	}

	err = migrations.InitDB(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}
