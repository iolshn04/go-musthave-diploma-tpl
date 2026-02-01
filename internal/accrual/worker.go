package accrual

import (
	"context"
	"time"
)

type OrdersRepo interface {
	ListForProcessing(ctx context.Context) ([]OrderForUpdate, error)
	Update(ctx context.Context, number, status string, accrual *float64) error
}

type BalanceRepo interface {
	Add(ctx context.Context, userID string, sum float64) error
}

type OrderForUpdate struct {
	Number string
	UserID string
}

type Worker struct {
	client  *Client
	orders  OrdersRepo
	balance BalanceRepo
}

func NewWorker(c *Client, o OrdersRepo, b BalanceRepo) *Worker {
	return &Worker{client: c, orders: o, balance: b}
}

func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second)

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			w.process(ctx)
		}
	}
}

func (w *Worker) process(ctx context.Context) {
	orders, err := w.orders.ListForProcessing(ctx)
	if err != nil {
		return
	}

	for _, o := range orders {
		resp, err := w.client.Get(ctx, o.Number)
		if err != nil {
			continue
		}

		_ = w.orders.Update(ctx, o.Number, resp.Status, resp.Accrual)

		if resp.Status == "PROCESSED" && resp.Accrual != nil {
			_ = w.balance.Add(ctx, o.UserID, *resp.Accrual)
		}
	}
}
