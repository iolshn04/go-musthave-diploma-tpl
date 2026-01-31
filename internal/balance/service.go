package balance

import (
	"context"
	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/models"
	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/orders"
)

type Repository interface {
	Get(ctx context.Context, userID string) (models.Balance, error)
	Withdraw(ctx context.Context, userID, order string, sum float64) error
	Withdrawals(ctx context.Context, userID string) ([]models.Withdrawal, error)
}

type ServiceInterface interface {
	Get(ctx context.Context, userID string) (models.Balance, error)
	Withdraw(ctx context.Context, userID, order string, sum float64) error
	Withdrawals(ctx context.Context, userID string) ([]models.Withdrawal, error)
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Get(ctx context.Context, userID string) (models.Balance, error) {
	return s.repo.Get(ctx, userID)
}

func (s *Service) Withdraw(ctx context.Context, userID, order string, sum float64) error {
	if !orders.ValidLuhn(order) {
		return ErrInvalidOrder
	}

	return s.repo.Withdraw(ctx, userID, order, sum)
}

func (s *Service) Withdrawals(ctx context.Context, userID string) ([]models.Withdrawal, error) {
	return s.repo.Withdrawals(ctx, userID)
}
