package user

import (
	"context"
	"database/sql"
	"errors"
	"literate-barnacle/database"

	_ "embed"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Get(ctx context.Context, login string) (DbUser, error)
	Add(ctx context.Context, user DbUser) error
}

type SqlRepository struct {
	db *sqlx.DB
}

func NewSqlRepository(db *sqlx.DB) SqlRepository {
	return SqlRepository{db: db}
}

//go:embed get.sql
var getSql string

//go:embed insert.sql
var insertSql string

func (s SqlRepository) Get(ctx context.Context, login string) (DbUser, error) {
	user := DbUser{}
	err := s.db.GetContext(ctx, &user, getSql, login)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return DbUser{}, database.ErrNotFound
	case err != nil:
		return DbUser{}, err
	default:
		return user, err
	}
}

func (s SqlRepository) Add(ctx context.Context, user DbUser) error {
	_, err := s.db.NamedExecContext(ctx, insertSql, user)
	if err != nil {
		return err
	}

	return nil
}
