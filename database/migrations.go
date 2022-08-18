package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose"
)

func Migrate(db *sqlx.DB) error {
	goose.SetTableName("db_version")
	err := goose.Up(db.DB, "database/migrations/")
	if err != nil {
		return fmt.Errorf("can't migrate up: %v", err)
	}

	return nil
}

func NewPgx(ctx context.Context, connectionString string) (*sqlx.DB, error) {
	connConfig, err := pgx.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("bad connection string: %v", err)
	}
	if connConfig == nil {
		return nil, errors.New("connection config is nil")
	}

	db := sqlx.NewDb(stdlib.OpenDB(*connConfig), "pgx")
	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("bad ping: %s", err)
	}

	return db, nil
}
