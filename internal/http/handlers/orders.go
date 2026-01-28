package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/http/middleware"
	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/orders"
)

type OrdersHandler struct {
	service orders.ServiceInterface
}

func NewOrdersHandler(s orders.ServiceInterface) *OrdersHandler {
	return &OrdersHandler{service: s}
}

func (h *OrdersHandler) Upload(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserID(r.Context())

	body, _ := io.ReadAll(r.Body)
	number := string(body)

	err := h.service.Upload(r.Context(), userID, number)

	switch err {
	case nil:
		w.WriteHeader(http.StatusAccepted) // 202
	case orders.ErrEmptyBody:
		w.WriteHeader(http.StatusBadRequest) // 400
	case orders.ErrInvalidNumber:
		w.WriteHeader(http.StatusUnprocessableEntity) // 422
	case orders.ErrAlreadyMine:
		w.WriteHeader(http.StatusOK) // 200
	case orders.ErrAlreadyExists:
		w.WriteHeader(http.StatusConflict) // 409
	default:
		w.WriteHeader(http.StatusInternalServerError) // 500
	}
}

func (h *OrdersHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserID(r.Context())

	ordersList, err := h.service.List(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(ordersList) == 0 {
		w.WriteHeader(http.StatusNoContent) // 204
		return
	}

	type OrderResp struct {
		Number     string   `json:"number"`
		Status     string   `json:"status"`
		Accrual    *float64 `json:"accrual,omitempty"`
		UploadedAt string   `json:"uploaded_at"`
	}

	resp := make([]OrderResp, len(ordersList))
	for i, o := range ordersList {
		loc := time.Now().Location()
		resp[i] = OrderResp{
			Number:     o.Number,
			Status:     o.Status,
			Accrual:    o.Accrual,
			UploadedAt: o.UploadedAt.In(loc).Format(time.RFC3339),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
