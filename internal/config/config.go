package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config contains the app configuration settings.
type Config struct {
	DbHost             string
	DbPort             uint16
	DbUser             string
	DbPassword         string
	DbName             string
	AppPort            uint16
	MigrationsDir      string
	ExternalCarsAPIURL string
	Env                string
}

// New loads env variables from .env file and the environment and assigns them to the Config.
func MustLoad() *Config {
	_ = godotenv.Load() // if .env doesn't exist, we still can try to retrieve env vars from the environment

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		panic("failed to parse DB_PORT")
	}
	appPort, err := strconv.Atoi(os.Getenv("APP_PORT"))
	if err != nil {
		panic("failed to parse APP_PORT")
	}
	dbHost, ok := os.LookupEnv("DB_HOST")
	if !ok {
		panic("failed to retrieve DB_HOST env variable")
	}
	dbUser, ok := os.LookupEnv("DB_USER")
	if !ok {
		panic("failed to retrieve DB_HOST env variable")
	}
	dbPassword, ok := os.LookupEnv("DB_PASSWORD")
	if !ok {
		panic("failed to retrieve DB_PASSWORD env variable")
	}
	dbName, ok := os.LookupEnv("DB_NAME")
	if !ok {
		panic("failed to retrieve DB_PASSWORD env variable")
	}
	migrationsDir, ok := os.LookupEnv("MIGRATIONS_DIR")
	if !ok {
		panic("failed to retrieve MIGRATIONS_DIR env variable")
	}
	externalCarsAPIURL, ok := os.LookupEnv("EXTERNAL_CARS_API_URL")
	if !ok {
		panic("failed to retrieve EXTERNAL_CARS_API_URL")
	}
	env, ok := os.LookupEnv("ENV")
	if !ok {
		panic("failed to retrive ENV")
	}

	return &Config{
		DbHost:             dbHost,
		DbPort:             uint16(dbPort),
		DbUser:             dbUser,
		DbPassword:         dbPassword,
		DbName:             dbName,
		AppPort:            uint16(appPort),
		MigrationsDir:      migrationsDir,
		ExternalCarsAPIURL: externalCarsAPIURL,
		Env:                env,
	}
}
