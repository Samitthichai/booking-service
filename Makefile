include .env
export

# Project name = network prefix for docker-run migrate
PROJECT = booking-service

# Run golang-migrate via docker image (no local install needed).
# Connects to the compose network; DB host = "db" (service name).
MIGRATE = docker run --rm \
	--network $(PROJECT)_default \
	-v $(PWD)/migrations:/migrations \
	migrate/migrate -path /migrations \
	-database "postgres://$(DB_USER):$(DB_PASSWORD)@db-postgres:5432/$(DB_NAME)?sslmode=$(DB_SSLMODE)"

.PHONY: run up up-build down logs api-logs db-up db-down db-logs psql migrate-up migrate-down migrate-force migrate-create

run:
	go run ./cmd/api

# --- Docker (full stack: api + Postgres + pgadmin) ---
up:
	docker compose up -d

up-build:
	docker compose up -d --build

down:
	docker compose down

logs:
	docker compose logs -f

api-logs:
	docker compose logs -f api

# --- Docker (Postgres only) ---
db-up:
	docker compose up -d db-postgres

db-down:
	docker compose down

db-logs:
	docker compose logs -f db-postgres

psql:
	docker compose exec db-postgres psql -U $(DB_USER) -d $(DB_NAME)

# --- Migration (golang-migrate) ---
migrate-up:
	$(MIGRATE) up

migrate-down:
	$(MIGRATE) down 1

migrate-force:
	$(MIGRATE) force $(V)

# make migrate-create name=add_foo
migrate-create:
	docker run --rm -v $(PWD)/migrations:/migrations \
		migrate/migrate create -ext sql -dir /migrations -seq $(name)
