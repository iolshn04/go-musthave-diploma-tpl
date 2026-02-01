package main

import (
	"log"

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

	application := app.New(cfg, repo)

	log.Printf("starting server on %s", cfg.RunAddress)
	if err := application.Run(); err != nil {
		log.Fatal(err)
	}
}
