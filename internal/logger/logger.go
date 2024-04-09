package logger

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/danblok/em-assessment/domain"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// CarsService is an interface to manage Car records.
type CarsService interface {
	FetchMany(context.Context, domain.FetchCarsFilter) ([]domain.Car, error)
	AddMany(context.Context, []domain.Car) error
	UpdateOne(context.Context, *domain.Car) error
	DeleteOneByID(context.Context, string) error
}

// LoggerService implements types.CarsService interface to manage logging.
type LoggerService struct {
	log *slog.Logger
	svc CarsService
}

// NewLoggerService creates new Logger
func NewLoggerService(log *slog.Logger, svc CarsService) *LoggerService {
	return &LoggerService{log: log, svc: svc}
}

// SetupLogger sets up different loggers depending on current environment.
func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

// FetchMany implements and wraps CarsService.FetchMany call to be able to log.
func (s *LoggerService) FetchMany(ctx context.Context, filter domain.FetchCarsFilter) (cars []domain.Car, err error) {
	s.log.Debug(
		"FetchMany",
		"time", time.Now(),
		"filter", filter,
	)
	defer func(t time.Time) {
		if err != nil {
			s.log.Error(
				"FetchMany",
				"time", t,
				"err", err,
			)
		}
	}(time.Now())

	return s.svc.FetchMany(ctx, filter)
}

// AddMany implements and wraps CarsService.AddMany call to be able to log.
func (s *LoggerService) AddMany(ctx context.Context, cars []domain.Car) (err error) {
	s.log.Debug(
		"AddMany",
		"time", time.Now(),
		"cars", cars,
	)
	defer func(t time.Time) {
		if err != nil {
			s.log.Error(
				"AddMany",
				"time", t,
				"err", err,
			)
		}
	}(time.Now())

	return s.svc.AddMany(ctx, cars)
}

// UpdateOne implements and wraps CarsService.UpdateOne call to be able to log.
func (s *LoggerService) UpdateOne(ctx context.Context, car *domain.Car) (err error) {
	s.log.Debug(
		"UpdateOne",
		"time", time.Now(),
		"car", car,
	)
	defer func(t time.Time) {
		if err != nil {
			s.log.Error(
				"UpdateOne",
				"time", t,
				"err", err,
				"car", car,
			)
		}
	}(time.Now())

	return s.svc.UpdateOne(ctx, car)
}

// DeleteOneByID implements and wraps CarsService.DeleteOneByID call to be able to log.
func (s *LoggerService) DeleteOneByID(ctx context.Context, id string) (err error) {
	s.log.Debug(
		"DeleteOnyByID",
		"time", time.Now(),
		"car_id", id,
	)
	defer func(t time.Time) {
		if err != nil {
			s.log.Error(
				"DeleteOneByID",
				"time", t,
				"err", err,
				"car_id", id,
			)
		}
	}(time.Now())

	return s.svc.DeleteOneByID(ctx, id)
}
