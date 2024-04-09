package carsfetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/danblok/em-assessment/domain"
)

// CarsFetcher retrieves a car information from some API.
type CarsFetcher struct {
	baseURL string
}

// New creates new CarsFetcher with base URL provided.
func New(baseURL string) *CarsFetcher {
	return &CarsFetcher{
		baseURL: baseURL,
	}
}

// FetchOneByRegNum fetches a Car information by its reg num.
func (s *CarsFetcher) FetchOneByRegNum(ctx context.Context, regNum string) (*domain.Car, error) {
	res, err := http.Get(fmt.Sprintf("%s/info?regNum=%s", s.baseURL, regNum))
	if err != nil {
		return nil, err
	}

	car := new(domain.Car)
	err = json.NewDecoder(res.Body).Decode(car)
	if err != nil {
		return nil, err
	}
	res.Body.Close()

	return car, nil
}
