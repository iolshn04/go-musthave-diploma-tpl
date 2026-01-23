package app

import (
	"context"
	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/repository/postgres"
	"net/http"
	"time"

	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/config"
	apphttp "github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/http"
)

type App struct {
	httpServer *http.Server
}

func New(cfg *config.Config, repo *postgres.Repository) *App {
	router := apphttp.NewRouter(repo, cfg)
	srv := &http.Server{
		Addr:         cfg.RunAddress,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	return &App{httpServer: srv}
}

func (a *App) Run() error {
	return a.httpServer.ListenAndServe()
}

func (a *App) Shutdown(ctx context.Context) error {
	return a.httpServer.Shutdown(ctx)
}
