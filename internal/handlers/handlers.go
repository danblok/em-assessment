package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/danblok/em-assessment/domain"
)

// CarsService is an interface to manage Car records.
type CarsService interface {
	FetchMany(context.Context, domain.FetchCarsFilter) ([]domain.Car, error)
	AddMany(context.Context, []domain.Car) error
	UpdateOne(context.Context, *domain.Car) error
	DeleteOneByID(context.Context, string) error
}

// CarsFetcher is an interface to fetch information about Car.
type CarsFetcher interface {
	FetchOneByRegNum(context.Context, string) (*domain.Car, error)
}

// HandlerFunc is a custom handler func that can return error.
type HandlerFunc func(http.ResponseWriter, *http.Request) error

// Handler is a custom handler.
type Handler struct {
	log *slog.Logger
	hf  HandlerFunc
}

// Wraps Handler with logger.
func WithLogger(log *slog.Logger) func(HandlerFunc) *Handler {
	return func(hf HandlerFunc) *Handler {
		return &Handler{log: log, hf: hf}
	}
}

// ServeHTTP implements http.HandlerFunc interface.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.hf(w, r)
	if e := new(APIError); errors.As(err, e) {
		h.log.Error(
			"handler error",
			"time", time.Now(),
			"err", e.Err.Error(),
			"handler", e.Handler,
			"status", e.Status,
			"method", r.Method,
			"path", r.URL.Path,
		)
		_ = JSON(w, e.Status, ResponseError{Msg: e.Msg})
		return
	} else if errors.Is(err, domain.ErrNoItemFound) {
		h.log.Error(
			"internal error",
			"time", time.Now(),
			"err", err.Error(),
			"status", http.StatusNotFound,
			"method", r.Method,
			"path", r.URL.Path,
		)
		_ = JSON(w, http.StatusNotFound, ResponseError{Msg: e.Msg})
		return
	} else if err != nil {
		h.log.Error(
			"internal error",
			"time", time.Now(),
			"err", e.Error(),
		)
		_ = JSON(w, http.StatusInternalServerError, ResponseError{Msg: "internal error"})
		return
	}
}

// ResponseError represents a HTTP response when an error occured.
type ResponseError struct {
	Msg string `json:"message"`
}

// APIError is a custom error that contains context info about an error occured in a handler.
type APIError struct {
	Err     error
	Status  int
	Msg     string
	Handler string
}

// Error implements error interface on APIError.
func (e APIError) Error() string {
	return e.Msg
}

// JSON writes status code and encodes body to JSON format.
func JSON(w http.ResponseWriter, code int, body any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(body)
}

// Status writes status code to response.
func Status(w http.ResponseWriter, code int) error {
	w.WriteHeader(code)
	return nil
}
