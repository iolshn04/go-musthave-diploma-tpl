package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/orders"
)

type OrdersRepository struct {
	db *sqlx.DB
}

func NewOrdersRepository(db *sqlx.DB) *OrdersRepository {
	return &OrdersRepository{db: db}
}

func (r *OrdersRepository) Save(ctx context.Context, userID, number string) error {
	const q = `
	INSERT INTO orders (id, number, user_id, status)
	VALUES ($1,$2,$3,'NEW')
	`

	_, err := r.db.ExecContext(ctx, q, uuid.New(), number, userID)
	return err
}

func (r *OrdersRepository) Owner(ctx context.Context, number string) (string, error) {
	const q = `SELECT user_id FROM orders WHERE number=$1`

	var uid string
	err := r.db.GetContext(ctx, &uid, q, number)
	if err != nil {
		return "", errors.New("not found")
	}
	return uid, nil
}

func (r *OrdersRepository) List(ctx context.Context, userID string) ([]orders.Order, error) {
	var res []orders.Order

	err := r.db.SelectContext(ctx, &res, `
		SELECT number, status, accrual, uploaded_at
		FROM orders
		WHERE user_id=$1
		ORDER BY uploaded_at DESC
	`, userID)

	return res, err
}
