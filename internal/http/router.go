package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/http/handlers"
	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/http/middleware"
)

func NewRouter(
	secret string,
	auth *handlers.AuthHandler,
	orders *handlers.OrdersHandler,
	balance *handlers.BalanceHandler,
) http.Handler {

	r := chi.NewRouter()

	r.Post("/api/user/register", auth.Register)
	r.Post("/api/user/login", auth.Login)

	r.Route("/api/user", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(secret))

		r.Post("/orders", orders.Upload)
		r.Get("/orders", orders.List)

		r.Get("/balance", balance.Get)
		r.Post("/balance/withdraw", balance.Withdraw)
		r.Get("/withdrawals", balance.Withdrawals)
	})

	return r
}
