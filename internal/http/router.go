package http

import (
	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/config"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/auth"
	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/http/handlers"
	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/http/middleware"
	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/repository/postgres"
)

func NewRouter(repo *postgres.Repository, cfg *config.Config) http.Handler {
	r := chi.NewRouter()

	userRepo := postgres.NewUserRepository(repo.DB)
	authService := auth.NewService(userRepo)
	authHandler := handlers.NewAuthHandler(authService, cfg.SecretKey)

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Post("/api/user/register", authHandler.Register)
	r.Post("/api/user/login", authHandler.Login)

	r.Route("/api/user", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(cfg.SecretKey))
		// orders / balance / withdrawals
	})

	return r
}
