package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/app"
	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/config"
	"github.com/iolshn04/go-musthave-diploma-tpl/tree/master/internal/repository/postgres"
)

func main() {
	cfg := config.New()

	repo, err := postgres.New(cfg.DatabaseURI)
	if err != nil {
		log.Fatalf("db init failed: %v", err)
	}
	defer repo.DB.Close()

	log.Printf("accrual addr: %s", cfg.AccrualSystemAddress)
	application := app.New(cfg, repo)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("starting server on %s", cfg.RunAddress)
		if err := application.Run(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()

	log.Println("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = application.Shutdown(shutdownCtx)
}
