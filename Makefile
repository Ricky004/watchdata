# Variables
BINARY_NAME=watchdata
SERVER_DIR=cmd/server
CLIENT_DIR=cmd/client
COLLECTOR_CONFIG=configs/otel-collector-config.yaml
DOCKER_COMPOSE_FILE=docker-compose.yaml

# Default target
.DEFAULT_GOAL := help

# Build the server binary
build-server:
	@echo "Building server..."
	go build -o bin/$(BINARY_NAME)-server $(SERVER_DIR)

# Build the client binary
build-client:
	@echo "Building client..."
	go build -o bin/$(BINARY_NAME)-client $(CLIENT_DIR)

# Build all binaries
build: build-server build-client

# Run the server
run-server:
	@echo "Running server..."
	go run $(SERVER_DIR)/main.go

# Run the client
run-client:
	@echo "Running client..."
	go run $(CLIENT_DIR)/main.go

# Start all services with Docker Compose
up:
	@echo "Starting all services..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d

# Stop all services
down:
	@echo "Stopping all services..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) down

# View logs from all services
logs:
	docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f

# Build and start services
start: build up

# Format Go code
fmt:
	@echo "Formatting Go code..."
	go fmt ./...

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run benchmarks
bench:
	@echo "Running benchmarks..."
	go test -bench=. ./...

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Create bin directory
bin:
	mkdir -p bin

# Clean build files and containers
clean:
	@echo "Cleaning up..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	docker-compose -f $(DOCKER_COMPOSE_FILE) down --volumes --remove-orphans

# Development setup
dev-setup: deps
	@echo "Setting up development environment..."
	@if ! command -v golangci-lint &> /dev/null; then \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi

# Check if ClickHouse is ready
check-clickhouse:
	@echo "Checking ClickHouse connection..."
	@timeout 30s bash -c 'until docker-compose exec clickhouse clickhouse-client --user=default --password=pass --query="SELECT 1"; do sleep 2; done'

# Initialize database schema
init-db: check-clickhouse
	@echo "Initializing database schema..."
	docker-compose exec clickhouse clickhouse-client --user=default --password=pass --multiquery < scripts/init.sql

# Show help
help:
	@echo "Available targets:"
	@echo "  build-server    - Build the server binary"
	@echo "  build-client    - Build the client binary" 
	@echo "  build          - Build all binaries"
	@echo "  run-server     - Run the server"
	@echo "  run-client     - Run the client"
	@echo "  up             - Start all services with Docker Compose"
	@echo "  down           - Stop all services"
	@echo "  logs           - View logs from all services"
	@echo "  start          - Build and start services"
	@echo "  fmt            - Format Go code"
	@echo "  lint           - Run linter"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage"
	@echo "  bench          - Run benchmarks"
	@echo "  deps           - Install dependencies"
	@echo "  clean          - Clean build files and containers"
	@echo "  dev-setup      - Setup development environment"
	@echo "  check-clickhouse - Check ClickHouse connection"
	@echo "  init-db        - Initialize database schema"
	@echo "  help           - Show this help message"

.PHONY: build-server build-client build run-server run-client up down logs start fmt lint test test-coverage bench deps clean dev-setup check-clickhouse init-db help bin
