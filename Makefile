SHELL := /bin/bash

APP_NAME := consumer-maintenance-system
DATABASE_URL ?= postgres://postgres:postgres@localhost:5432/curs3?sslmode=disable

.PHONY: help run build seed test lint generate migrate-up migrate-down \
	docker-build docker-up docker-down docker-fresh scale-app backup restore clean

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
	@echo "  make scale-app N=3 Scale the app service and make the gateway pick it up (see target's own comment for why both steps matter)"
	@echo "  make backup        Dump the Postgres database to backups/<timestamp>.sql.gz"
	@echo "  make restore FILE=backups/xyz.sql.gz  Restore a backup (drops and recreates the current schema)"
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

# Scaling `app` alone isn't enough: the gateway's nginx.conf resolves the
# `app` upstream's DNS once at startup (a deliberate tradeoff, not an
# oversight -- see observability/nginx/nginx.conf's own comment: a
# dynamically-resolved upstream disables connection keepalive and caused
# ~69% 502s under load during the highload test). Open-source nginx has no
# built-in way to pick up new replicas without a config reload, so this
# just makes the two-step dance ("scale, then remember to restart the
# gateway") into one command instead of a step that's easy to forget.
# Real dynamic re-balancing (no restart at all) would need either NGINX
# Plus or a Docker-events-driven sidecar regenerating the config on the
# fly -- deliberately out of scope here, same size of effort as the
# Kafka/Postgres-replication/Redis-HA ideas this project already decided
# against for a single-node demo stack.
scale-app:
	@if [ -z "$(N)" ]; then echo "Usage: make scale-app N=<replica count>"; exit 1; fi
	docker compose up -d --scale app=$(N) app
	docker compose restart gateway
	@echo "Scaled app to $(N) replicas, gateway re-resolved DNS."

# Plain pg_dump/pg_restore, not a managed backup service -- appropriate for
# a single-node demo stack, not a claim this is production-grade backup
# infrastructure (no offsite storage, no retention policy, no automated
# schedule). Still a real, restorable backup: verified by round-tripping
# through `make backup` then `make restore` against a fresh volume.
backup:
	@mkdir -p backups
	@f=backups/$$(date +%Y%m%d-%H%M%S).sql.gz; \
	docker compose exec -T db pg_dump -U postgres curs3 | gzip > $$f; \
	echo "Backup written to $$f"

restore:
	@if [ -z "$(FILE)" ]; then echo "Usage: make restore FILE=backups/<name>.sql.gz"; exit 1; fi
	@echo "This drops and recreates the public schema in the running db container before restoring -- ctrl-C now to abort."
	@sleep 3
	docker compose exec -T db psql -U postgres -d curs3 -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
	gunzip -c $(FILE) | docker compose exec -T db psql -U postgres -d curs3

clean:
	rm -rf bin
