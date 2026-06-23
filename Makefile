SHELL := /bin/zsh

.PHONY: help build run docker-build docker-up docker-down clean test

help:
	@echo "Available targets:"
	@echo "  make build        Build the project with Maven."
	@echo "  make run          Run the application locally."
	@echo "  make docker-build Build the Docker image."
	@echo "  make docker-up    Start app and PostgreSQL via Docker Compose."
	@echo "  make docker-down  Stop Docker Compose services."
	@echo "  make clean        Clean Maven build artifacts."
	@echo "  make test         Run unit tests."

build:
	./mvnw clean package

run: build
	java -jar target/curs3ProjectBack-0.0.1-SNAPSHOT.jar

docker-build:
	docker build -t curs3projectback .

docker-up:
	docker compose up --build

docker-down:
	docker compose down

clean:
	./mvnw clean

test:
	./mvnw test
