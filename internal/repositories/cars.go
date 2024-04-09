package repositories

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/danblok/em-assessment/domain"
)

// CarsRepository represents repository for managing Car information in postgresql.
type CarsRepository struct {
	db *sqlx.DB
}

// NewCarsRepository returns a new instance of CarsRepository.
func NewCarsRepository(db *sqlx.DB) *CarsRepository {
	return &CarsRepository{db: db}
}

// FetchMany queries cars from DB by provided filter and returns them.
func (r *CarsRepository) FetchMany(ctx context.Context, filter domain.FetchCarsFilter) ([]domain.Car, error) {
	query := sq.Select("*").From("cars")
	if len(filter.IDs) > 0 {
		query = query.Where(sq.Eq{"car_id": filter.IDs})
	}
	if len(filter.RegNums) > 0 {
		query = query.Where(sq.Eq{"reg_num": filter.RegNums})
	}
	if len(filter.Marks) > 0 {
		query = query.Where(sq.Eq{"mark": filter.Marks})
	}
	if len(filter.Models) > 0 {
		query = query.Where(sq.Eq{"model": filter.Models})
	}
	if len(filter.Years) > 0 {
		query = query.Where(sq.Eq{"year": filter.Years})
	}

	queryStr, args, err := query.Limit(uint64(filter.Limit)).
		Offset(uint64(filter.Offset)).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	cars := make([]domain.Car, 0)
	err = r.db.SelectContext(ctx, &cars, queryStr, args...)
	if err != nil {
		return nil, err
	}
	return cars, nil
}

// AddMany adds many cars to DB.
func (r *CarsRepository) AddMany(ctx context.Context, cars []domain.Car) error {
	query := sq.Insert("cars").Columns("reg_num", "mark", "model", "year")
	for _, car := range cars {
		query = query.Values(car.RegNum, car.Mark, car.Model, car.Year)
	}

	queryStr, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, queryStr, args...)
	return err
}

// UpdateOne updates a car in DB with provided payload.
func (r *CarsRepository) UpdateOne(ctx context.Context, payload *domain.Car) error {
	query := sq.Update("cars")
	if payload.RegNum != "" {
		query = query.Set("reg_num", payload.RegNum)
	}
	if payload.Model != "" {
		query = query.Set("model", payload.Model)
	}
	if payload.Mark != "" {
		query = query.Set("mark", payload.Mark)
	}
	if payload.Year > 0 {
		query = query.Set("year", payload.Year)
	}
	query = query.Where(sq.Eq{"car_id": payload.ID})

	queryStr, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}

	res, err := r.db.ExecContext(ctx, queryStr, args...)
	if err != nil {
		return nil
	}

	if rowsAffected, err := res.RowsAffected(); err != nil {
		return err
	} else if rowsAffected == 0 {
		return domain.ErrNoItemFound
	}

	return nil
}

// DeleteOneByID deletes a car by provided ID.
func (r *CarsRepository) DeleteOneByID(ctx context.Context, id string) error {
	query := sq.Delete("cars").Where(sq.Eq{"car_id": id})
	queryStr, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}

	res, err := r.db.ExecContext(ctx, queryStr, args...)
	if err != nil {
		return err
	}

	if rowsAffected, err := res.RowsAffected(); err != nil {
		return err
	} else if rowsAffected == 0 {
		return domain.ErrNoItemFound
	}

	return nil
}
