# Real-Time Website Analytics Dashboard with Kafka and Go

A comprehensive real-time analytics platform built with Apache Kafka and Go that captures, processes, and visualizes website analytics events. Features a beautiful web dashboard with live charts, real-time event streaming, and intelligent alerting.

## Features

### 🎯 Core Analytics
- **Event Producer API**: HTTP API endpoint to receive analytics events
- **Event Consumer**: Background service to process and aggregate events  
- **Real-time Processing**: Events are processed in real-time using Apache Kafka
- **Event Types**: Support for page views, clicks, sessions, and custom events
- **Scalable Architecture**: Kafka-based architecture allows horizontal scaling

### 📊 Real-Time Dashboard
- **Interactive Web Dashboard**: Beautiful, responsive dashboard with live updates
- **Real-time Charts**: Line charts, doughnuts, bar charts with Chart.js
- **Live Event Stream**: See events as they happen with WebSocket updates
- **Key Metrics**: Total events, unique users, active sessions, performance metrics
- **Device & Browser Analytics**: Detailed breakdowns of user environments

### 🔔 Intelligence & Alerts
- **Smart Alerts**: Configurable threshold-based alerting system
- **Performance Monitoring**: Track page load times and performance metrics
- **Traffic Source Analysis**: Understand where your traffic comes from
- **Time-windowed Analytics**: Hourly breakdowns and historical data

### 🛠 DevOps & Deployment
- **WebSocket Support**: Real-time bidirectional communication
- **Docker Support**: Complete Docker Compose setup for easy deployment
- **Graceful Shutdown**: Proper cleanup and resource management
- **Health Monitoring**: Built-in health checks and monitoring endpoints

## Architecture

```
┌─────────────┐        ┌──────────────┐        ┌────────┐        ┌──────────────┐
│   Website   │───────>│   Producer   │───────>│ Kafka  │───────>│   Consumer   │
│  (Clients)  │  HTTP  │   Service    │        │ Topics │        │   Service    │
└─────────────┘        └──────────────┘        └────────┘        └──────────────┘
                            │:8080                                       │
                            │                                            ▼
                            ▼                                    ┌──────────────┐
                ┌──────────────────┐                           │  Analytics   │
                │  Web Dashboard   │◄──────────────────────────│  Processing  │
                │   (Real-time)    │      WebSocket             │   & Alerts   │
                └──────────────────┘                           └──────────────┘
                       │                                               │
                   Browser                                       Real-time
                   Client                                        Aggregation
```

**Components:**
- **Producer Service**: HTTP API + WebSocket server + Dashboard hosting
- **Consumer Service**: Event processing + Analytics engine + Alerting
- **Kafka**: Message broker for scalable event streaming
- **Web Dashboard**: Real-time visualization with live updates
- **Analytics Engine**: Time-windowed aggregations and intelligence

## Prerequisites

- Go 1.23 or higher
- Docker and Docker Compose (for running Kafka)
- Make (optional, for using Makefile commands)

## Quick Start

### 1. Clone the repository

```bash
git clone https://github.com/Hilina-t/go-kafka-analytics-pipeline.git
cd go-kafka-analytics-pipeline
```

### 2. Start services with Docker Compose

```bash
docker-compose up -d
```

This will start:
- Zookeeper (port 2181)
- Kafka (port 9092)
- Producer service (port 8080)
- Consumer service

### 3. Send test events

```bash
./examples/send_events.sh
```

### 4. Access the Real-Time Dashboard

Open your browser and navigate to:
```
http://localhost:8080
```

You'll see:
- 📊 **Live Analytics Dashboard** with real-time charts and metrics
- 📈 **Interactive Visualizations** using Chart.js
- 🔴 **Live Event Stream** showing events as they happen
- 📱 **Device & Browser Analytics**
- ⚡ **Performance Monitoring**

### 5. Test with Enhanced Events

```bash
./examples/test_dashboard.sh
```

This sends realistic test data including page views, clicks, and session events with performance metrics.

### 6. View consumer logs to see processed events

```bash
docker-compose logs -f consumer
```

## Development

### Install dependencies

```bash
go mod download
```

### Build locally

```bash
make build
```

Or manually:

```bash
go build -o producer ./cmd/producer
go build -o consumer ./cmd/consumer
```

### Run services locally

First, make sure Kafka is running:

```bash
docker-compose up -d zookeeper kafka
```

Then run the producer:

```bash
make run-producer
# or
go run ./cmd/producer
```

In another terminal, run the consumer:

```bash
make run-consumer
# or
go run ./cmd/consumer
```

## API Endpoints

### GET / (Dashboard)

Access the real-time analytics dashboard.

**Response:** HTML dashboard with live WebSocket updates

### GET /analytics

Get current analytics snapshot as JSON.

**Response:**

```json
{
  "timestamp": "2024-01-01T12:00:00Z",
  "total_events": 1500,
  "unique_users": 245,
  "active_sessions": 12,
  "events_by_type": {
    "page_view": 1200,
    "click": 250,
    "session": 50
  },
  "top_pages": [...],
  "traffic_sources": [...],
  "device_stats": {...},
  "browser_stats": {...},
  "hourly_page_views": [...],
  "performance_metrics": {...}
}
```

### WebSocket /ws

Real-time WebSocket endpoint for live dashboard updates.

**Message Types:**
- `analytics_snapshot`: Complete analytics data
- `analytics_update`: Incremental updates (every 5s)
- `real_time_event`: Individual events as they happen
- `alert`: System alerts and notifications

### POST /event

Send an analytics event to be processed.

**Request Body:**

```json
{
  "type": "page_view",
  "user_id": "user123",
  "session_id": "session456",
  "url": "https://example.com/home",
  "path": "/home",
  "referrer": "https://google.com",
  "user_agent": "Mozilla/5.0",
  "ip_address": "192.168.1.1",
  "metadata": {
    "page_title": "Home Page",
    "load_time": 1200
  }
}
```

**Response:**

```json
{
  "status": "accepted",
  "id": "550e8400-e29b-41d4-a716-446655440000"
}
```

### GET /health

Health check endpoint.

**Response:**

```json
{
  "status": "healthy",
  "service": "analytics-producer"
}
```

## Event Types

### Page View Event

Tracks when a user views a page.

```json
{
  "type": "page_view",
  "user_id": "user123",
  "session_id": "session456",
  "url": "https://example.com/home",
  "path": "/home",
  "metadata": {
    "page_title": "Home Page",
    "load_time": 1200
  }
}
```

### Click Event

Tracks user clicks on elements.

```json
{
  "type": "click",
  "user_id": "user123",
  "session_id": "session456",
  "url": "https://example.com/products",
  "path": "/products",
  "metadata": {
    "element_id": "buy-button",
    "element_type": "button",
    "element_text": "Buy Now"
  }
}
```

### Session Event

Tracks user session information.

```json
{
  "type": "session",
  "user_id": "user123",
  "session_id": "session456",
  "url": "https://example.com",
  "path": "/",
  "metadata": {
    "duration": 300,
    "page_count": 5,
    "device": "desktop",
    "browser": "Chrome"
  }
}
```

## Configuration

Both services can be configured using environment variables:

### Producer Service

| Variable | Default | Description |
|----------|---------|-------------|
| `KAFKA_BROKERS` | `localhost:9092` | Kafka broker addresses |
| `KAFKA_TOPIC` | `analytics-events` | Kafka topic name |
| `SERVER_PORT` | `8080` | HTTP server port |

### Consumer Service

| Variable | Default | Description |
|----------|---------|-------------|
| `KAFKA_BROKERS` | `localhost:9092` | Kafka broker addresses |
| `KAFKA_TOPIC` | `analytics-events` | Kafka topic name |
| `CONSUMER_GROUP` | `analytics-consumer-group` | Consumer group ID |

## Available Make Commands

```bash
make build           # Build producer and consumer binaries
make clean           # Remove build artifacts
make test            # Run tests
make run-producer    # Run producer locally
make run-consumer    # Run consumer locally
make docker-up       # Start all services with Docker Compose
make docker-down     # Stop all services
make docker-restart  # Rebuild and restart Docker services
make docker-logs     # View logs from Docker services
make deps            # Install and tidy dependencies
make fmt             # Format code
make lint            # Run linter
make help            # Show help message
```

## Project Structure

```
.
├── cmd/
│   ├── producer/          # Producer service (HTTP API)
│   └── consumer/          # Consumer service (event processor)
├── pkg/
│   ├── kafka/             # Kafka producer and consumer wrappers
│   └── models/            # Event data models
├── examples/
│   └── send_events.sh     # Script to send test events
├── docker-compose.yml     # Docker Compose configuration
├── Dockerfile.producer    # Producer Dockerfile
├── Dockerfile.consumer    # Consumer Dockerfile
├── Makefile              # Build and development commands
└── README.md             # This file
```

## Testing

Send test events using curl:

```bash
curl -X POST http://localhost:8080/event \
  -H "Content-Type: application/json" \
  -d '{
    "type": "page_view",
    "user_id": "user123",
    "session_id": "session456",
    "url": "https://example.com/home",
    "path": "/home"
  }'
```

## Monitoring

The consumer service prints analytics statistics every 30 seconds, showing:
- Unique user count
- Events by type
- Top pages by view count

## Troubleshooting

### Kafka connection issues

If services can't connect to Kafka, ensure Kafka is running:

```bash
docker-compose ps
```

Check Kafka logs:

```bash
docker-compose logs kafka
```

### Port conflicts

If port 8080 or 9092 is already in use, you can change the ports in `docker-compose.yml`.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
