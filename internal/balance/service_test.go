package balance

import (
	"context"
	"testing"

	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/models"
)

type mockRepo struct {
	getRes models.Balance
	getErr error
	err    error
}

func (m *mockRepo) Get(ctx context.Context, userID string) (models.Balance, error) {
	return m.getRes, m.getErr
}

func (m *mockRepo) Withdraw(ctx context.Context, userID, order string, sum float64) error {
	return m.err
}

func (m *mockRepo) Withdrawals(ctx context.Context, userID string) ([]models.Withdrawal, error) {
	return nil, m.err
}

//
// GET
//

func TestService_Get_OK(t *testing.T) {
	repo := &mockRepo{
		getRes: models.Balance{Current: 50},
	}

	s := New(repo)

	res, err := s.Get(context.Background(), "u")

	if err != nil || res.Current != 50 {
		t.Fatal("failed get")
	}
}

//
// WITHDRAW
//

func TestService_Withdraw_InvalidNumber(t *testing.T) {
	s := New(&mockRepo{})

	err := s.Withdraw(context.Background(), "u", "123", 10)

	if err != ErrInvalidOrder {
		t.Fatal("expected invalid order")
	}
}

func TestService_Withdraw_RepoError(t *testing.T) {
	repo := &mockRepo{err: ErrNotEnoughFunds}
	s := New(repo)

	err := s.Withdraw(context.Background(), "u", "79927398713", 10)

	if err != ErrNotEnoughFunds {
		t.Fatal("expected repo error")
	}
}

func TestService_Withdraw_OK(t *testing.T) {
	s := New(&mockRepo{})

	err := s.Withdraw(context.Background(), "u", "79927398713", 10)

	if err != nil {
		t.Fatal("unexpected error")
	}
}
