package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/balance"
	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/http/middleware"
)

type BalanceHandler struct {
	service balance.ServiceInterface
}

func NewBalanceHandler(s balance.ServiceInterface) *BalanceHandler {
	return &BalanceHandler{service: s}
}

func (h *BalanceHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserID(r.Context())

	res, err := h.service.Get(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *BalanceHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserID(r.Context())

	type reqBody struct {
		Order string  `json:"order"`
		Sum   float64 `json:"sum"`
	}

	var body reqBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := h.service.Withdraw(r.Context(), userID, body.Order, body.Sum)

	switch err {
	case nil:
		w.WriteHeader(http.StatusOK)
	case balance.ErrInvalidOrder:
		w.WriteHeader(http.StatusUnprocessableEntity)
	case balance.ErrNotEnoughFunds:
		w.WriteHeader(http.StatusPaymentRequired)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *BalanceHandler) Withdrawals(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserID(r.Context())

	list, err := h.service.Withdrawals(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(list) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}
