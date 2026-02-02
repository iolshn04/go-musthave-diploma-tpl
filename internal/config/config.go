package config

import (
	"flag"
	"os"
)

type Config struct {
	RunAddress           string
	DatabaseURI          string
	AccrualSystemAddress string
	SecretKey            string
}

func New() *Config {
	const (
		defaultRunAddress  = "localhost:8080"
		defaultDatabaseURI = "postgres://gophermart:secret@localhost:5432/gophermart?sslmode=disable"
		defaultSecretKey   = "super-secret-key"
	)

	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	flagRunAddr := fs.String("a", "", "service run address")
	flagDB := fs.String("d", "", "database uri")
	flagAccrual := fs.String("r", "", "accrual system address")

	_ = fs.Parse(os.Args[1:])

	cfg := &Config{}

	if v, ok := os.LookupEnv("RUN_ADDRESS"); ok {
		cfg.RunAddress = v
	} else if *flagRunAddr != "" {
		cfg.RunAddress = *flagRunAddr
	} else {
		cfg.RunAddress = defaultRunAddress
	}

	if v, ok := os.LookupEnv("DATABASE_URI"); ok {
		cfg.DatabaseURI = v
	} else if *flagDB != "" {
		cfg.DatabaseURI = *flagDB
	} else {
		cfg.DatabaseURI = defaultDatabaseURI
	}

	if v, ok := os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS"); ok {
		cfg.AccrualSystemAddress = v
	} else if *flagAccrual != "" {
		cfg.AccrualSystemAddress = *flagAccrual
	} else {
		cfg.AccrualSystemAddress = "http://localhost:8081"
	}

	if val, ok := os.LookupEnv("SECRET_KEY"); ok {
		cfg.SecretKey = val
	} else {
		cfg.SecretKey = defaultSecretKey
	}

	return cfg
}
