# ---------- ENV ----------
DB_HOST ?= localhost
DB_PORT ?= 5435
DB_USER ?= postgres
DB_PASSWORD ?= postgres
DB_NAME ?= auth_db
MIGRATIONS_DIR ?= migrations
DOCKER_COMPOSE ?= docker-compose.yml

# ---------- HELP ----------
.PHONY: help
help: ## Show this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# ---------- DOCKER ----------
.PHONY: up down logs
up: ## Start postgres container
	docker compose -f $(DOCKER_COMPOSE) up -d postgres

down: ## Stop & remove containers
	docker compose -f $(DOCKER_COMPOSE) down

logs: ## Tail postgres logs
	docker logs -f postgres-auth

# ---------- DATABASE ----------
DB_URL := postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

.PHONY: migrateup migratedown new-migration
migrateup: ## Apply pending migrations
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" -verbose up

migratedown: ## Roll back last migration
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" -verbose down 1

new-migration: ## Usage: make new-migration name=add_foo
	@echo "Creating migration: $(name)"
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

# ---------- UTILS ----------
.PHONY: psql
psql: ## Open psql shell to local container
	docker exec -it postgres-auth psql -U $(DB_USER) -d $(DB_NAME)