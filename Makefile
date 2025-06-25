# Makefile for E-commerce Microservices Project

# --- Configuration ---
# Define all service names. Used for looping in build commands.
SERVICES := order_service payment_service inventory_service notification_service
# Database connection URL for the migrate tool
DB_URL := "postgres://user:password@localhost:5432/ecommerce_db?sslmode=disable"
# Output directory for local builds
BIN_DIR := bin

# --- Phony Targets ---
.PHONY: all up down logs build clean migrate-create migrate-up migrate-down help

# --- Main Commands ---

# Default command: build and bring up the stack
all: build up

# Bring up all services in detached mode
up:
	@echo "Bringing up Docker containers..."
	docker compose up -d --build

# Stop and remove all containers and networks
down:
	@echo "Stopping Docker containers..."
	docker compose down

# Follow logs from all services. Use `make logs service=order-service` for specific logs.
logs:
	@echo "Following logs..."
	docker compose logs -f $(service)

# Build all service binaries locally into the ./bin directory
build:
	@echo "Building all service binaries locally..."
	@mkdir -p $(BIN_DIR)
	@for service in $(SERVICES); do \
		echo "--> Building $$service..."; \
		go build -o $(BIN_DIR)/$$service ./cmd/$$service/main.go; \
	done

# Remove the local build directory
clean:
	@echo "Cleaning local binaries..."
	@rm -rf $(BIN_DIR)

# --- Database Migrations ---
# These commands require golang-migrate to be installed locally.
# Example: `make migrate-create name=create_products_table`
migrate-create:
	@echo "Creating migration: $(name)"
	@migrate create -ext sql -dir migrations -seq $(name)

# Apply all up migrations
migrate-up:
	@echo "Applying UP migrations..."
	@migrate -path migrations -database $(DB_URL) -verbose up

# Roll back all down migrations
migrate-down:
	@echo "Applying DOWN migrations..."
	@migrate -path migrations -database $(DB_URL) -verbose down

# --- Help ---
help:
	@echo "Usage:"
	@echo "  make up                - Start all services with Docker Compose"
	@echo "  make down              - Stop all services"
	@echo "  make logs [service=...] - Tail logs from services (e.g., service=order-service)"
	@echo "  make build             - Build all Go binaries locally"
	@echo "  make clean             - Remove local binaries"
	@echo "  make migrate-create name=... - Create a new migration file"
	@echo "  make migrate-up        - Apply all database migrations"
	@echo "  make migrate-down      - Revert all database migrations"

