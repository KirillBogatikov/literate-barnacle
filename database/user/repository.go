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
	GetById(ctx context.Context, id string) (DbUser, error)
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

//go:embed get_by_id.sql
var getByIdSql string

//go:embed insert.sql
var insertSql string

func (s SqlRepository) Get(ctx context.Context, login string) (DbUser, error) {
	return s.getOne(ctx, getSql, login)
}

func (s SqlRepository) GetById(ctx context.Context, id string) (DbUser, error) {
	return s.getOne(ctx, getByIdSql, id)
}

func (s SqlRepository) Add(ctx context.Context, user DbUser) error {
	_, err := s.db.NamedExecContext(ctx, insertSql, user)
	if err != nil {
		return err
	}

	return nil
}

func (s SqlRepository) getOne(ctx context.Context, script string, params ...interface{}) (DbUser, error) {
	user := DbUser{}
	err := s.db.GetContext(ctx, &user, script, params...)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return DbUser{}, database.ErrNotFound
	case err != nil:
		return DbUser{}, err
	default:
		return user, err
	}
}
