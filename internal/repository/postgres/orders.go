package postgres

import (
	"context"
	"errors"
	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/accrual"
	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
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

func (r *OrdersRepository) List(ctx context.Context, userID string) ([]models.Order, error) {
	var res []models.Order

	err := r.db.SelectContext(ctx, &res, `
		SELECT number, status, accrual, uploaded_at
		FROM orders
		WHERE user_id=$1
		ORDER BY uploaded_at DESC
	`, userID)

	return res, err
}

func (r *OrdersRepository) ListForProcessing(ctx context.Context) ([]accrual.OrderForUpdate, error) {
	var res []accrual.OrderForUpdate

	err := r.db.SelectContext(ctx, &res, `
		SELECT number, user_id
		FROM orders
		WHERE status IN ('NEW','PROCESSING')
	`)

	return res, err
}

func (r *OrdersRepository) Update(ctx context.Context, number, status string, accrual *float64) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE orders
		SET status=$2, accrual=$3
		WHERE number=$1
	`, number, status, accrual)

	return err
}
