package postgres

import (
	"context"
	"errors"

	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/models"
	"github.com/jmoiron/sqlx"
)

var ErrUserExists = errors.New("user already exists")
var ErrUserNotFound = errors.New("user not found")

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	const q = `
		INSERT INTO users (id, login, password_hash)
		VALUES ($1, $2, $3)
	`

	_, err := r.db.ExecContext(ctx, q, user.ID, user.Login, user.PasswordHash)

	if err != nil {
		return ErrUserExists
	}

	return nil
}

func (r *UserRepository) GetByLogin(ctx context.Context, login string) (*models.User, error) {
	const q = `
		SELECT id, login, password_hash, created_at
		FROM users
		WHERE login = $1
	`

	var u models.User
	if err := r.db.GetContext(ctx, &u, q, login); err != nil {
		return nil, ErrUserNotFound
	}

	return &u, nil
}
