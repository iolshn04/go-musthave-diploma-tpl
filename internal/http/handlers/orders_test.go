package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/http/middleware"
	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/orders"
)

type mockService struct {
	uploadFunc func(ctx context.Context, userID, number string) error
	listFunc   func(ctx context.Context, userID string) ([]models.Order, error)
}

func (m *mockService) Upload(ctx context.Context, userID, number string) error {
	return m.uploadFunc(ctx, userID, number)
}

func (m *mockService) List(ctx context.Context, userID string) ([]models.Order, error) {
	return m.listFunc(ctx, userID)
}

func fakeUserID() string { return "user123" }

func TestOrdersHandlerUpload(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		mockErr    error
		wantStatus int
	}{
		{"empty body", "", orders.ErrEmptyBody, http.StatusBadRequest},
		{"invalid number", "abc123", orders.ErrInvalidNumber, http.StatusUnprocessableEntity},
		{"already mine", "1234567890", orders.ErrAlreadyMine, http.StatusOK},
		{"already exists", "1234567890", orders.ErrAlreadyExists, http.StatusConflict},
		{"success", "79927398713", nil, http.StatusAccepted},
		{"other error", "79927398713", context.DeadlineExceeded, http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockService{
				uploadFunc: func(ctx context.Context, userID, number string) error {
					return tt.mockErr
				},
			}
			h := NewOrdersHandler(mock)

			req := httptest.NewRequest(http.MethodPost, "/api/user/orders", bytes.NewBufferString(tt.body))
			req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, fakeUserID()))

			w := httptest.NewRecorder()
			h.Upload(w, req)
			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("got status %d, want %d", resp.StatusCode, tt.wantStatus)
			}
		})
	}
}

func TestOrdersHandlerList(t *testing.T) {
	accrual := 100.0
	now := time.Now()

	tests := []struct {
		name       string
		mockOrders []models.Order
		mockErr    error
		wantStatus int
		wantBody   bool
	}{
		{"empty list", []models.Order{}, nil, http.StatusNoContent, false},
		{"list with orders", []models.Order{
			{Number: "123", Status: "NEW", UploadedAt: now, Accrual: &accrual},
		}, nil, http.StatusOK, true},
		{"error", nil, context.DeadlineExceeded, http.StatusInternalServerError, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockService{
				listFunc: func(ctx context.Context, userID string) ([]models.Order, error) {
					return tt.mockOrders, tt.mockErr
				},
			}
			h := NewOrdersHandler(mock)

			req := httptest.NewRequest(http.MethodGet, "/api/user/orders", nil)
			req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, fakeUserID()))
			w := httptest.NewRecorder()

			h.List(w, req)
			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("got status %d, want %d", resp.StatusCode, tt.wantStatus)
			}

			if tt.wantBody {
				var resp []map[string]interface{}
				err := json.NewDecoder(w.Body).Decode(&resp)
				if err != nil {
					t.Errorf("failed to decode response: %v", err)
				}
				if len(resp) != len(tt.mockOrders) {
					t.Errorf("got %d orders, want %d", len(resp), len(tt.mockOrders))
				}
			}
		})
	}
}
