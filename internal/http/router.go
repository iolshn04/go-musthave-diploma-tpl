package http

import (
	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/balance"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/auth"
	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/config"
	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/http/handlers"
	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/http/middleware"
	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/orders"
	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/repository/postgres"
)

func NewRouter(repo *postgres.Repository, cfg *config.Config) http.Handler {
	r := chi.NewRouter()

	// auth
	userRepo := postgres.NewUserRepository(repo.DB)
	authService := auth.NewService(userRepo)
	authHandler := handlers.NewAuthHandler(authService, cfg.SecretKey)

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Post("/api/user/register", authHandler.Register)
	r.Post("/api/user/login", authHandler.Login)

	// orders
	ordersRepo := postgres.NewOrdersRepository(repo.DB)
	ordersService := orders.New(ordersRepo)
	ordersHandler := handlers.NewOrdersHandler(ordersService)

	r.Route("/api/user", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(cfg.SecretKey))

		r.Post("/orders", ordersHandler.Upload)
		r.Get("/orders", ordersHandler.List)
	})

	// balance
	balanceRepo := postgres.NewBalanceRepository(repo.DB)
	balanceService := balance.New(balanceRepo)
	balanceHandler := handlers.NewBalanceHandler(balanceService)

	r.Get("/api/user/balance", balanceHandler.Get)
	r.Post("/api/user/balance/withdraw", balanceHandler.Withdraw)
	r.Get("/api/user/withdrawals", balanceHandler.Withdrawals)

	return r
}
