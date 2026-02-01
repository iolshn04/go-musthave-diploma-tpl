package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/http/middleware"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/balance"
	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/models"
)

type mockBalanceService struct {
	getRes         models.Balance
	getErr         error
	withdrawErr    error
	withdrawalsRes []models.Withdrawal
	withdrawalsErr error
}

func (m *mockBalanceService) Get(ctx context.Context, userID string) (models.Balance, error) {
	return m.getRes, m.getErr
}

func (m *mockBalanceService) Withdraw(ctx context.Context, userID, order string, sum float64) error {
	return m.withdrawErr
}

func (m *mockBalanceService) Withdrawals(ctx context.Context, userID string) ([]models.Withdrawal, error) {
	return m.withdrawalsRes, m.withdrawalsErr
}

func ctxWithUser(r *http.Request) *http.Request {
	ctx := context.WithValue(r.Context(), middleware.UserIDKey, "user1")
	return r.WithContext(ctx)
}

func TestBalanceHandler_Get_OK(t *testing.T) {
	svc := &mockBalanceService{
		getRes: models.Balance{
			Current:   100,
			Withdrawn: 10,
		},
	}

	h := NewBalanceHandler(svc)

	req := ctxWithUser(httptest.NewRequest(http.MethodGet, "/balance", nil))
	w := httptest.NewRecorder()

	h.Get(w, req)

	if w.Code != http.StatusOK {
		t.Fatal("status not 200")
	}
}

func TestBalanceHandler_Get_Error(t *testing.T) {
	svc := &mockBalanceService{
		getErr: context.DeadlineExceeded,
	}

	h := NewBalanceHandler(svc)

	req := ctxWithUser(httptest.NewRequest(http.MethodGet, "/balance", nil))
	w := httptest.NewRecorder()

	h.Get(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatal("expected 500")
	}
}

func TestBalanceHandler_Withdraw_OK(t *testing.T) {
	svc := &mockBalanceService{}

	h := NewBalanceHandler(svc)

	body := []byte(`{"order":"79927398713","sum":10}`)

	req := ctxWithUser(httptest.NewRequest(http.MethodPost, "/withdraw", bytes.NewBuffer(body)))
	w := httptest.NewRecorder()

	h.Withdraw(w, req)

	if w.Code != http.StatusOK {
		t.Fatal("expected 200")
	}
}

func TestBalanceHandler_Withdraw_NotEnough(t *testing.T) {
	svc := &mockBalanceService{
		withdrawErr: balance.ErrNotEnoughFunds,
	}

	h := NewBalanceHandler(svc)

	body := []byte(`{"order":"79927398713","sum":10}`)

	req := ctxWithUser(httptest.NewRequest(http.MethodPost, "/withdraw", bytes.NewBuffer(body)))
	w := httptest.NewRecorder()

	h.Withdraw(w, req)

	if w.Code != http.StatusPaymentRequired {
		t.Fatal("expected 402")
	}
}

func TestBalanceHandler_Withdraw_Invalid(t *testing.T) {
	svc := &mockBalanceService{
		withdrawErr: balance.ErrInvalidOrder,
	}

	h := NewBalanceHandler(svc)

	body := []byte(`{"order":"123","sum":10}`)

	req := ctxWithUser(httptest.NewRequest(http.MethodPost, "/withdraw", bytes.NewBuffer(body)))
	w := httptest.NewRecorder()

	h.Withdraw(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatal("expected 422")
	}
}

func TestBalanceHandler_Withdrawals_OK(t *testing.T) {
	svc := &mockBalanceService{
		withdrawalsRes: []models.Withdrawal{
			{
				Order:       "1",
				Sum:         10,
				ProcessedAt: time.Now(),
			},
		},
	}

	h := NewBalanceHandler(svc)

	req := ctxWithUser(httptest.NewRequest(http.MethodGet, "/withdrawals", nil))
	w := httptest.NewRecorder()

	h.Withdrawals(w, req)

	if w.Code != http.StatusOK {
		t.Fatal("expected 200")
	}

	var resp []models.Withdrawal
	_ = json.Unmarshal(w.Body.Bytes(), &resp)

	if len(resp) != 1 {
		t.Fatal("wrong withdrawals count")
	}
}

func TestBalanceHandler_Withdrawals_Empty(t *testing.T) {
	svc := &mockBalanceService{}

	h := NewBalanceHandler(svc)

	req := ctxWithUser(httptest.NewRequest(http.MethodGet, "/withdrawals", nil))
	w := httptest.NewRecorder()

	h.Withdrawals(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatal("expected 204")
	}
}
