package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Hilina-t/go-kafka-analytics-pipeline/pkg/kafka"
	"github.com/Hilina-t/go-kafka-analytics-pipeline/pkg/models"
)

// Analytics holds aggregated analytics data
type Analytics struct {
	mu              sync.RWMutex
	PageViews       map[string]int            // URL -> count
	UniqueUsers     map[string]bool           // UserID -> exists
	EventsByType    map[models.EventType]int  // EventType -> count
	TopPages        []PageStats
	LastUpdated     time.Time
}

type PageStats struct {
	URL   string
	Count int
}

func NewAnalytics() *Analytics {
	return &Analytics{
		PageViews:    make(map[string]int),
		UniqueUsers:  make(map[string]bool),
		EventsByType: make(map[models.EventType]int),
		LastUpdated:  time.Now(),
	}
}

func (a *Analytics) ProcessEvent(event *models.AnalyticsEvent) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Track event by type
	a.EventsByType[event.Type]++

	// Track unique users
	if event.UserID != "" {
		a.UniqueUsers[event.UserID] = true
	}

	// Track page views
	if event.Type == models.PageView && event.URL != "" {
		a.PageViews[event.URL]++
	}

	a.LastUpdated = time.Now()

	log.Printf("Processed %s event for user %s on %s", event.Type, event.UserID, event.URL)
	return nil
}

func (a *Analytics) PrintStats() {
	a.mu.RLock()
	defer a.mu.RUnlock()

	fmt.Println("\n=== Analytics Summary ===")
	fmt.Printf("Last Updated: %s\n", a.LastUpdated.Format(time.RFC3339))
	fmt.Printf("Unique Users: %d\n", len(a.UniqueUsers))
	
	fmt.Println("\nEvents by Type:")
	for eventType, count := range a.EventsByType {
		fmt.Printf("  %s: %d\n", eventType, count)
	}

	if len(a.PageViews) > 0 {
		fmt.Println("\nTop Pages:")
		count := 0
		for url, views := range a.PageViews {
			fmt.Printf("  %s: %d views\n", url, views)
			count++
			if count >= 10 {
				break
			}
		}
	}
	fmt.Println("========================")
}

func main() {
	// Get configuration from environment variables
	kafkaBrokers := getEnv("KAFKA_BROKERS", "localhost:9092")
	kafkaTopic := getEnv("KAFKA_TOPIC", "analytics-events")
	consumerGroup := getEnv("CONSUMER_GROUP", "analytics-consumer-group")

	log.Printf("Starting consumer with brokers: %s, topic: %s, group: %s", 
		kafkaBrokers, kafkaTopic, consumerGroup)

	// Create analytics processor
	analytics := NewAnalytics()

	// Create Kafka consumer
	consumer := kafka.NewConsumer([]string{kafkaBrokers}, kafkaTopic, consumerGroup)
	defer consumer.Close()

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("\nReceived shutdown signal, printing final stats...")
		analytics.PrintStats()
		cancel()
	}()

	// Start periodic stats reporting
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				analytics.PrintStats()
			case <-ctx.Done():
				return
			}
		}
	}()

	// Start consuming events
	log.Println("Consumer started, waiting for events...")
	if err := consumer.ConsumeEvents(ctx, analytics.ProcessEvent); err != nil {
		if err == context.Canceled {
			log.Println("Consumer stopped gracefully")
		} else {
			log.Fatalf("Consumer error: %v", err)
		}
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
