# Makefile for VentureX Backend

# Load environment variables from .env file.
# Create a .env file from .env.example and fill in your details.
-include .env

# Database configuration
DB_URL := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

# Migration commands
.PHONY: migrate-up migrate-down migrate-status migrate-create migrate-reset migrate-version

migrate-up:
	@echo "Running migrations up..."
	goose -dir migrations postgres "$(DB_URL)" up

migrate-down:
	@echo "Rolling back one migration..."
	goose -dir migrations postgres "$(DB_URL)" down

migrate-status:
	@echo "Migration status:"
	goose -dir migrations postgres "$(DB_URL)" status

migrate-create:
	@echo "Creating new migration: $(name)"
	@if [ -z "$(name)" ]; then echo "Usage: make migrate-create name=migration_name"; exit 1; fi
	goose -dir migrations create $(name) sql

migrate-reset:
	@echo "Resetting database (WARNING: This will drop all data!)"
	goose -dir migrations postgres "$(DB_URL)" reset

migrate-version:
	@echo "Current migration version:"
	goose -dir migrations postgres "$(DB_URL)" version

# Development commands
.PHONY: run build test clean

run:
	CONFIG_FILE=config/app.example.yaml go run ./cmd/api

build:
	go build -o bin/api ./cmd/api

test:
	go test ./...

clean:
	rm -rf bin/

# Docker commands
.PHONY: docker-up docker-down docker-logs

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

# Database setup
.PHONY: db-setup

db-setup: migrate-up
	@echo "Database setup complete!"

# Help
.PHONY: help

help:
	@echo "Available commands:"
	@echo "  migrate-up       - Run all pending migrations"
	@echo "  migrate-down     - Rollback one migration"
	@echo "  migrate-status   - Show migration status"
	@echo "  migrate-create   - Create new migration (use: make migrate-create name=migration_name)"
	@echo "  migrate-reset    - Reset database (WARNING: drops all data)"
	@echo "  migrate-version  - Show current migration version"
	@echo "  run             - Run the application"
	@echo "  build           - Build the application"
	@echo "  test            - Run tests"
	@echo "  clean           - Clean build artifacts"
	@echo "  docker-up       - Start docker services"
	@echo "  docker-down     - Stop docker services"
	@echo "  docker-logs     - Show docker logs"
	@echo "  db-setup        - Setup database with migrations" 

debug:
	dlv debug ./cmd/api -- --config=config/app.example.yaml 