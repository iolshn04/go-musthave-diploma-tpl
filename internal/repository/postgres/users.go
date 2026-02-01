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
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	const userQ = `
		INSERT INTO users (id, login, password_hash)
		VALUES ($1, $2, $3)
	`

	if _, err := tx.ExecContext(ctx, userQ, user.ID, user.Login, user.PasswordHash); err != nil {
		return ErrUserExists
	}

	const balanceQ = `
		INSERT INTO balances (user_id, current, withdrawn)
		VALUES ($1, 0, 0)
	`

	if _, err := tx.ExecContext(ctx, balanceQ, user.ID); err != nil {
		return err
	}

	return tx.Commit()
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
