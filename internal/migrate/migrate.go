package migrate

import (
	"errors"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Do applies migrations from 'pathToMigrations' to database in 'dbURL'.
func Do(pathToMigrations, dbURL string) error {
	m, err := migrate.New("file://"+pathToMigrations, dbURL)
	if err != nil {
		return err
	}
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) { // pass if no change happened
		return err
	}

	return nil
}
