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

	return db, nil
}
