package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	"github.com/danblok/em-assessment/domain"
)

type carsHandler struct {
	svc         CarsService
	carsFetcher CarsFetcher
}

// RegisterCarsHandler registers Car handlers.
func RegisterCarsHandler(r *chi.Mux, svc CarsService, carsFetcher CarsFetcher, log *slog.Logger) {
	h := &carsHandler{
		svc:         svc,
		carsFetcher: carsFetcher,
	}

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	handler := WithLogger(log)
	r.Method(http.MethodGet, "/cars", handler(h.handleGetMany))
	r.Method(http.MethodPost, "/cars", handler(h.handleAddMany))
	r.Method(http.MethodPatch, "/cars/{id}", handler(h.handleUpdateOne))
	r.Method(http.MethodDelete, "/cars/{id}", handler(h.handleDeleteOne))
}

func (h *carsHandler) handleGetMany(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query()
	offset := 0
	if query.Has("offset") {
		var err error
		offset, err = strconv.Atoi(query.Get("offset"))
		if err != nil || offset < 0 {
			return APIError{
				Err:     err,
				Status:  http.StatusBadRequest,
				Msg:     "invalid offset",
				Handler: "handleGetMany",
			}
		}
	}

	limit := 20
	if query.Has("limit") {
		var err error
		limit, err = strconv.Atoi(query.Get("limit"))
		if err != nil || limit < 0 {
			return APIError{
				Err:     err,
				Status:  http.StatusBadRequest,
				Msg:     "invalid limit",
				Handler: "handleGetMany",
			}
		}
	}

	years := make([]int, 0, len(query["year"]))
	for _, year := range query["year"] {
		n, err := strconv.Atoi(year)
		if err != nil {
			return APIError{
				Err:     err,
				Status:  http.StatusBadRequest,
				Msg:     "invalid year",
				Handler: "handleGetMany",
			}
		}
		years = append(years, n)
	}

	filter := domain.FetchCarsFilter{
		Offset:  offset,
		Limit:   limit,
		IDs:     query["id"],
		RegNums: query["regNum"],
		Marks:   query["mark"],
		Models:  query["model"],
		Years:   years,
	}

	cars, err := h.svc.FetchMany(r.Context(), filter)
	if err != nil {
		return err
	}

	return JSON(w, http.StatusOK, cars)
}

func (h *carsHandler) handleAddMany(w http.ResponseWriter, r *http.Request) error {
	var body struct {
		RegNums []string `json:"regNums"`
	}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return APIError{
			Err:     err,
			Status:  http.StatusBadRequest,
			Msg:     "invalid requestBody",
			Handler: "handleAddMany",
		}
	}
	r.Body.Close()

	if len(body.RegNums) == 0 {
		return APIError{
			Err:     fmt.Errorf("empty reg nums: len: %d", len(body.RegNums)),
			Status:  http.StatusBadRequest,
			Msg:     "empty reg nums",
			Handler: "handleAddMany",
		}
	}

	cars := make([]domain.Car, len(body.RegNums))
	g, ctx := errgroup.WithContext(r.Context()) // parallelize fetching
	for i, regNum := range body.RegNums {
		g.Go(func() error {
			car, err := h.carsFetcher.FetchOneByRegNum(ctx, regNum)
			if err == nil {
				cars[i] = *car
			}
			return err
		})
	}
	if err := g.Wait(); err != nil { // if some fetch fails, its error gets returned
		return err
	}

	err = h.svc.AddMany(r.Context(), cars)
	if err != nil {
		return err
	}

	return Status(w, http.StatusCreated)
}

func (h *carsHandler) handleUpdateOne(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	if err := uuid.Validate(id); err != nil {
		return APIError{
			Err:     err,
			Status:  http.StatusBadRequest,
			Msg:     "invalid path param car id",
			Handler: "handleUpdateOne",
		}
	}

	car := new(domain.Car)
	err := json.NewDecoder(r.Body).Decode(car)
	car.ID = id
	if err != nil {
		return APIError{
			Err:     err,
			Status:  http.StatusBadRequest,
			Msg:     "invalid request body",
			Handler: "handleUpdateOne",
		}
	}
	r.Body.Close()

	err = h.svc.UpdateOne(r.Context(), car)
	if err != nil {
		return err
	}

	return Status(w, http.StatusOK)
}

func (h *carsHandler) handleDeleteOne(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	if err := uuid.Validate(id); err != nil {
		return APIError{
			Err:     err,
			Status:  http.StatusBadRequest,
			Msg:     "invalid path param car id",
			Handler: "handleDeleteOne",
		}
	}

	err := h.svc.DeleteOneByID(r.Context(), id)
	if err != nil {
		return err
	}

	return Status(w, http.StatusOK)
}
