.PHONY: test clean all

build:
	@go build -o bin/cars cmd/cars/main.go

up:build
	@bin/cars
