.PHONY: all build clean test run-producer run-consumer docker-up docker-down help

# Variables
PRODUCER_BINARY=producer
CONSUMER_BINARY=consumer

all: build

# Build both services
build:
	@echo "Building producer..."
	go build -o $(PRODUCER_BINARY) ./cmd/producer
	@echo "Building consumer..."
	go build -o $(CONSUMER_BINARY) ./cmd/consumer

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f $(PRODUCER_BINARY) $(CONSUMER_BINARY)
	go clean

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run producer locally
run-producer:
	@echo "Running producer..."
	go run ./cmd/producer

# Run consumer locally
run-consumer:
	@echo "Running consumer..."
	go run ./cmd/consumer

# Start all services with Docker Compose
docker-up:
	@echo "Starting services with Docker Compose..."
	docker-compose up -d

# Stop all services
docker-down:
	@echo "Stopping services..."
	docker-compose down

# Rebuild and restart Docker services
docker-restart: docker-down docker-up

# View logs
docker-logs:
	docker-compose logs -f

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run linter
lint:
	@echo "Running linter..."
	go vet ./...

# Help
help:
	@echo "Available targets:"
	@echo "  build            - Build producer and consumer binaries"
	@echo "  clean            - Remove build artifacts"
	@echo "  test             - Run tests"
	@echo "  run-producer     - Run producer locally"
	@echo "  run-consumer     - Run consumer locally"
	@echo "  docker-up        - Start all services with Docker Compose"
	@echo "  docker-down      - Stop all services"
	@echo "  docker-restart   - Rebuild and restart Docker services"
	@echo "  docker-logs      - View logs from Docker services"
	@echo "  deps             - Install and tidy dependencies"
	@echo "  fmt              - Format code"
	@echo "  lint             - Run linter"
	@echo "  help             - Show this help message"
