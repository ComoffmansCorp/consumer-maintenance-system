SHELL := /bin/bash

APP_NAME := consumer-maintenance-system
DATABASE_URL ?= postgres://postgres:postgres@localhost:5432/curs3?sslmode=disable

.PHONY: help run build seed test lint generate migrate-up migrate-down \
	docker-build docker-up docker-down docker-fresh clean

help:
	@echo "Available targets:"
	@echo "  make run           Run the API server locally (needs local Postgres + env vars)"
	@echo "  make seed          Fill the local database with demo data (idempotent)"
	@echo "  make build         Build the Go binary into bin/api"
	@echo "  make test          Run Go tests"
	@echo "  make lint          Run golangci-lint"
	@echo "  make generate      Regenerate sqlc code from sql/ and migrations/"
	@echo "  make migrate-up    Apply all pending migrations"
	@echo "  make migrate-down  Roll back the last migration"
	@echo "  make docker-up     Start the whole stack (db+migrate+seed+api+frontend) via Docker Compose"
	@echo "  make docker-down   Stop Docker Compose services"
	@echo "  make docker-fresh  Recreate Docker Compose services with a clean volume"
	@echo "  make clean         Remove build artifacts"

run:
	go run ./cmd/api

seed:
	go run ./cmd/seed

build:
	go build -o bin/api ./cmd/api

test:
	go test ./...

lint:
	golangci-lint run ./...

generate:
	sqlc generate

migrate-up:
	migrate -database "$(DATABASE_URL)" -path migrations up

migrate-down:
	migrate -database "$(DATABASE_URL)" -path migrations down 1

docker-build:
	docker build -t $(APP_NAME) .

docker-up:
	docker compose up --build

docker-down:
	docker compose down

docker-fresh:
	docker compose down -v
	docker compose up --build

clean:
	rm -rf bin
