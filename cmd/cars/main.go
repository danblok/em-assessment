package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/danblok/em-assessment/internal/config"
	"github.com/danblok/em-assessment/internal/handlers"
	"github.com/danblok/em-assessment/internal/logger"
	"github.com/danblok/em-assessment/internal/migrate"
	"github.com/danblok/em-assessment/internal/repositories"
	"github.com/danblok/em-assessment/internal/services/carsfetcher"
)

func main() {
	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.Env)
	log.Debug("Configuration", "cfg", cfg)

	dbConnURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DbUser,
		cfg.DbPassword,
		cfg.DbHost,
		cfg.DbPort,
		cfg.DbName,
	)
	log.Debug("Establishing connection to DB", "url", dbConnURL)
	db, err := sqlx.Connect("postgres", dbConnURL)
	if err != nil {
		log.Error("Failed to connect to DB", "err", err)
		return
	}
	defer func() {
		log.Debug("Closing database connection...")
		err := db.Close()
		if err != nil {
			log.Error("Failed to close database", "err", err)
			return
		}
		log.Debug("Database connection closed")
	}()
	log.Info("DB connection established")

	log.Debug("Applying migrations")
	err = migrate.Do(cfg.MigrationsDir, dbConnURL)
	if err != nil {
		log.Error("Failed to apply migrations", "err", err)
		return
	}
	log.Info("Migrations applied")

	r := chi.NewRouter()
	svc := logger.NewLoggerService(log, repositories.NewCarsRepository(db))
	externalSvc := carsfetcher.New(cfg.ExternalCarsAPIURL)
	handlers.RegisterCarsHandler(r, svc, externalSvc, log)

	log.Info("Starting server")
	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.AppPort), r); err != nil {
		log.Error("Failed to start server", "err", err)
	}
}
