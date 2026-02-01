package app

import (
	"context"
	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/accrual"
	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/auth"
	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/balance"
	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/http/handlers"
	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/orders"
	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/repository/postgres"
	"net/http"
	"time"

	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/config"
	apphttp "github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/http"
)

type App struct {
	httpServer *http.Server
	worker     *accrual.Worker
}

func New(cfg *config.Config, repo *postgres.Repository) *App {

	// repos
	userRepo := postgres.NewUserRepository(repo.DB)
	ordersRepo := postgres.NewOrdersRepository(repo.DB)
	balanceRepo := postgres.NewBalanceRepository(repo.DB)

	// services
	authService := auth.NewService(userRepo)
	ordersService := orders.New(ordersRepo)
	balanceService := balance.New(balanceRepo)

	// handlers
	authHandler := handlers.NewAuthHandler(authService, cfg.SecretKey)
	ordersHandler := handlers.NewOrdersHandler(ordersService)
	balanceHandler := handlers.NewBalanceHandler(balanceService)

	// worker
	accrualClient := accrual.NewClient(cfg.AccrualSystemAddress)
	worker := accrual.NewWorker(accrualClient, ordersRepo, balanceRepo)

	router := apphttp.NewRouter(
		cfg.SecretKey,
		authHandler,
		ordersHandler,
		balanceHandler,
	)

	srv := &http.Server{
		Addr:         cfg.RunAddress,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &App{
		httpServer: srv,
		worker:     worker,
	}
}

func (a *App) Run(ctx context.Context) error {
	go a.worker.Run(ctx)

	return a.httpServer.ListenAndServe()
}

func (a *App) Shutdown(ctx context.Context) error {
	return a.httpServer.Shutdown(ctx)
}
