package user

import (
	"context"
	"literate-barnacle/models"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	CreateUser(ctx context.Context, user models.User) error
	GetUser(ctx context.Context, username, password string) (models.User, error)
}

type SqlRepository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return SqlRepository{db}
}

// go:embed insert.sql
var insertSql string

func (r SqlRepository) CreateUser(ctx context.Context, user models.User) error {
	u := toPostgresUser(user)
	_, err := r.db.NamedExecContext(ctx, insertSql, u)
	return err
}

// go:embed get.sql
var getUserSql string

func (r SqlRepository) GetUser(ctx context.Context, username, password string) (models.User, error) {
	user := User{}
	err := r.db.GetContext(ctx, &user, getUserSql, username, password)
	if err != nil {
		return models.User{}, err
	}

	return toModel(user), nil
}

func toPostgresUser(u models.User) User {
	return User{
		Username: u.Username,
		Password: u.Password,
	}
}

func toModel(u User) models.User {
	return models.User{
		ID:       u.ID,
		Username: u.Username,
		Password: u.Password,
	}
}
