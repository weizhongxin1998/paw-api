package database

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func Init() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./paw.db")
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)

	if err := Run(db); err != nil {
		return nil, err
	}

	return db, nil
}
