package postgres

import (
	"context"
	"github.com/google/uuid"
	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/balance"

	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/models"
	"github.com/jmoiron/sqlx"
)

type BalanceRepository struct {
	db *sqlx.DB
}

func NewBalanceRepository(db *sqlx.DB) *BalanceRepository {
	return &BalanceRepository{db: db}
}

func (r *BalanceRepository) Get(ctx context.Context, userID string) (models.Balance, error) {
	const q = `
	SELECT
	  COALESCE(SUM(o.accrual),0) -
	  COALESCE((SELECT SUM(sum) FROM withdrawals WHERE user_id=$1),0) as current,
	  COALESCE((SELECT SUM(sum) FROM withdrawals WHERE user_id=$1),0) as withdrawn
	FROM orders o
	WHERE o.user_id=$1;
	`

	var b models.Balance
	err := r.db.GetContext(ctx, &b, q, userID)
	return b, err
}

func (r *BalanceRepository) Withdraw(ctx context.Context, userID, order string, sum float64) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var current float64

	err = tx.GetContext(ctx, &current,
		`SELECT current FROM balances WHERE user_id=$1 FOR UPDATE`, userID)
	if err != nil {
		return err
	}

	if current < sum {
		return balance.ErrNotEnoughFunds
	}

	// insert withdrawal
	_, err = tx.ExecContext(ctx, `
		INSERT INTO withdrawals (id, user_id, order_number, sum)
		VALUES ($1,$2,$3,$4)
	`, uuid.New(), userID, order, sum)
	if err != nil {
		return err
	}

	// update balance
	_, err = tx.ExecContext(ctx, `
		UPDATE balances
		SET current = current - $2,
		    withdrawn = withdrawn + $2
		WHERE user_id=$1
	`, userID, sum)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *BalanceRepository) Withdrawals(ctx context.Context, userID string) ([]models.Withdrawal, error) {
	var res []models.Withdrawal

	err := r.db.SelectContext(ctx, &res, `
		SELECT order_number, sum, processed_at
		FROM withdrawals
		WHERE user_id=$1
		ORDER BY processed_at DESC
	`, userID)

	return res, err
}
