package accrual

import (
	"context"
	"log"
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
	Number string `db:"number"`
	UserID string `db:"user_id"`
}

type Worker struct {
	client  *Client
	orders  OrdersRepo
	balance BalanceRepo
}

func NewWorker(c *Client, o OrdersRepo, b BalanceRepo) *Worker {
	return &Worker{
		client:  c,
		orders:  o,
		balance: b,
	}
}

func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("accrual worker stopped")
			return
		case <-ticker.C:
			w.process()
		}
	}
}

func (w *Worker) process() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	orders, err := w.orders.ListForProcessing(ctx)
	if err != nil {
		log.Printf("worker: failed to list orders: %v", err)
		return
	}

	for _, o := range orders {
		resp, err := w.client.Get(ctx, o.Number)
		if err != nil {
			log.Printf("worker: failed to get order %s from accrual: %v", o.Number, err)
			continue
		}

		if resp.Status == "PROCESSED" && resp.Accrual != nil {
			if err := w.balance.Add(ctx, o.UserID, *resp.Accrual); err != nil {
				log.Printf("worker: failed to add balance for user %s: %v", o.UserID, err)
			}
		}

		if err := w.orders.Update(ctx, o.Number, resp.Status, resp.Accrual); err != nil {
			log.Printf("worker: failed to update order %s: %v", o.Number, err)
		}
	}
}
