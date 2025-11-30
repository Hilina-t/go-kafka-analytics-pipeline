package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/Hilina-t/go-kafka-analytics-pipeline/constants"
	"github.com/Hilina-t/go-kafka-analytics-pipeline/pkg/analytics"
	"github.com/Hilina-t/go-kafka-analytics-pipeline/pkg/kafka"
	"github.com/Hilina-t/go-kafka-analytics-pipeline/pkg/models"
	"github.com/Hilina-t/go-kafka-analytics-pipeline/pkg/websocket"
	"github.com/google/uuid"
)

type Server struct {
	producer         *kafka.Producer
	analyticsService *analytics.Service
	wsHub            *websocket.Hub
	port             string
}

func NewServer(producer *kafka.Producer, port string) *Server {
	analyticsService := analytics.NewService()
	wsHub := websocket.NewHub(analyticsService)

	return &Server{
		producer:         producer,
		analyticsService: analyticsService,
		wsHub:            wsHub,
		port:             port,
	}
}

func (s *Server) handleEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var event models.AnalyticsEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// Set ID and timestamp if not provided
	if event.ID == "" {
		event.ID = uuid.New().String()
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	ctx := context.Background()
	if err := s.producer.SendEvent(ctx, event.ID, event); err != nil {
		log.Printf("Failed to send event: %v", err)
		http.Error(w, "Failed to send event", http.StatusInternalServerError)
		return
	}

	// Process event for real-time analytics
	if err := s.analyticsService.ProcessEvent(&event); err != nil {
		log.Printf("Failed to process analytics event: %v", err)
	}

	// Broadcast event to WebSocket clients
	s.wsHub.BroadcastEvent(&event)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "accepted",
		"id":     event.ID,
	})
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "healthy",
		"service": "analytics-producer",
	})
}

func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	// Serve the dashboard HTML file
	dashboardPath := filepath.Join("web", "dashboard.html")
	http.ServeFile(w, r, dashboardPath)
}

func (s *Server) handleAnalytics(w http.ResponseWriter, r *http.Request) {
	snapshot := s.analyticsService.GetSnapshot()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(snapshot)
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	s.wsHub.ServeWS(w, r)
}

func (s *Server) Start(ctx context.Context) error {
	// Start WebSocket hub in a goroutine
	go s.wsHub.Run()

	mux := http.NewServeMux()
	mux.HandleFunc("/event", s.handleEvent)
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/", s.handleDashboard)
	mux.HandleFunc("/analytics", s.handleAnalytics)
	mux.HandleFunc("/ws", s.handleWebSocket)

	server := &http.Server{
		Addr:         ":" + s.port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Producer server starting on port %s", s.port)
		log.Printf("Dashboard available at http://localhost:%s", s.port)
		log.Printf("WebSocket endpoint: ws://localhost:%s/ws", s.port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	<-ctx.Done()

	// Graceful shutdown with 30 second timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Println("Shutting down server gracefully...")
	return server.Shutdown(shutdownCtx)
}

func main() {
	// Create Kafka producer
	producer := kafka.NewProducer([]string{constants.KafkaBrokers}, constants.KafkaTopic)
	defer producer.Close()

	// Create and start server
	server := NewServer(producer, constants.ServerPort)

	// Handle graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Received shutdown signal...")
		cancel()
	}()

	if err := server.Start(ctx); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}

	log.Println("Server stopped gracefully")
}
