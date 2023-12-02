package storage

import (
	"context"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func NewStorage(ctx context.Context, filepath string) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}

	return db, nil
}
