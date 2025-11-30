.PHONY: all build clean test run-producer run-consumer docker-up docker-down docker-restart docker-logs deps fmt lint test-dashboard help

# Variables
PRODUCER_BINARY=producer
CONSUMER_BINARY=consumer

all: build

# Build both services with enhanced features
build:
	@echo "🔨 Building enhanced producer with dashboard..."
	go build -o $(PRODUCER_BINARY) ./cmd/producer
	@echo "🔨 Building enhanced consumer with analytics..."
	go build -o $(CONSUMER_BINARY) ./cmd/consumer
	@echo "✅ Build complete! Dashboard available at http://localhost:8080"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -f $(PRODUCER_BINARY) $(CONSUMER_BINARY)
	go clean

# Install and tidy dependencies
deps:
	@echo "📦 Installing dependencies..."
	go mod download
	go mod tidy

# Run tests
test:
	@echo "🧪 Running tests..."
	go test -v ./...

# Format code
fmt:
	@echo "🎨 Formatting code..."
	go fmt ./...

# Run linter
lint:
	@echo "🔍 Running linter..."
	golangci-lint run

# Run producer locally
run-producer:
	@echo "🚀 Running producer with dashboard..."
	@echo "📊 Dashboard: http://localhost:8080"
	@echo "🔗 WebSocket: ws://localhost:8080/ws"
	@echo "📈 Analytics API: http://localhost:8080/analytics"
	go run ./cmd/producer

# Run consumer locally
run-consumer:
	@echo "🚀 Running enhanced consumer with analytics..."
	@echo "📊 Real-time analytics processing enabled"
	@echo "🔔 Smart alerts configured"
	go run ./cmd/consumer

# Start all services with Docker Compose
docker-up:
	@echo "🐳 Starting services with Docker Compose..."
	docker-compose up -d
	@echo "✅ Services started!"
	@echo "📊 Dashboard: http://localhost:8080"
	@echo "📈 Analytics API: http://localhost:8080/analytics"

# Stop all services
docker-down:
	@echo "🐳 Stopping services..."
	docker-compose down

# Rebuild and restart Docker services
docker-restart: docker-down
	@echo "🔄 Rebuilding and restarting services..."
	docker-compose up -d --build

# View logs from all services
docker-logs:
	@echo "📋 Viewing logs..."
	docker-compose logs -f

# Test the dashboard with sample data
test-dashboard:
	@echo "🧪 Testing Real-Time Analytics Dashboard..."
	@echo "📊 Sending sample events..."
	@chmod +x ./examples/test_dashboard.sh
	./examples/test_dashboard.sh

# Show help message
help:
	@echo "🚀 Real-Time Analytics Dashboard - Available Commands:"
	@echo ""
	@echo "  📦 Building & Dependencies:"
	@echo "    build            - Build enhanced producer and consumer with dashboard"
	@echo "    clean            - Remove build artifacts"
	@echo "    deps             - Install and tidy Go dependencies"
	@echo ""
	@echo "  🧪 Development & Testing:"
	@echo "    test             - Run all tests"
	@echo "    test-dashboard   - Test dashboard with realistic sample data"
	@echo "    fmt              - Format Go code"
	@echo "    lint             - Run code linter"
	@echo ""
	@echo "  🚀 Local Development:"
	@echo "    run-producer     - Run producer with dashboard locally (port 8080)"
	@echo "    run-consumer     - Run enhanced consumer with analytics locally"
	@echo ""
	@echo "  🐳 Docker Operations:"
	@echo "    docker-up        - Start all services with Docker Compose"
	@echo "    docker-down      - Stop all services"
	@echo "    docker-restart   - Rebuild and restart Docker services"
	@echo "    docker-logs      - View logs from Docker services"
	@echo ""
	@echo "  📊 Dashboard Features:"
	@echo "    • Real-time analytics with live charts"
	@echo "    • WebSocket-powered event streaming"
	@echo "    • Performance monitoring & alerts"
	@echo "    • Device/browser analytics"
	@echo "    • Traffic source analysis"
	@echo ""
	@echo "  🎯 Quick Start:"
	@echo "    make docker-up && make test-dashboard"
	@echo "    Then visit: http://localhost:8080"
