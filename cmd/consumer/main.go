package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Hilina-t/go-kafka-analytics-pipeline/constants"
	"github.com/Hilina-t/go-kafka-analytics-pipeline/pkg/analytics"
	"github.com/Hilina-t/go-kafka-analytics-pipeline/pkg/kafka"
	"github.com/Hilina-t/go-kafka-analytics-pipeline/pkg/models"
)

// ConsumerService handles event processing and analytics
type ConsumerService struct {
	consumer         *kafka.Consumer
	analyticsService *analytics.Service
}

// NewConsumerService creates a new consumer service
func NewConsumerService(consumer *kafka.Consumer, analyticsService *analytics.Service) *ConsumerService {
	return &ConsumerService{
		consumer:         consumer,
		analyticsService: analyticsService,
	}
}

// processEvent handles incoming events from Kafka
func (cs *ConsumerService) processEvent(event *models.AnalyticsEvent) error {
	log.Printf("Processing %s event for user %s on %s", event.Type, event.UserID, event.URL)

	// Process the event through analytics service
	if err := cs.analyticsService.ProcessEvent(event); err != nil {
		log.Printf("Error processing analytics event: %v", err)
		return err
	}

	// Check for alerts
	alerts := cs.analyticsService.CheckAlerts()
	for _, alert := range alerts {
		log.Printf("ALERT [%s]: %s", alert.Severity, alert.Message)
	}

	return nil
}

// printStats prints current analytics statistics
func (cs *ConsumerService) printStats() {
	snapshot := cs.analyticsService.GetSnapshot()

	fmt.Println("\n=== Real-Time Analytics Summary ===")
	fmt.Printf("Last Updated: %s\n", snapshot.Timestamp.Format(time.RFC3339))
	fmt.Printf("Total Events: %d\n", snapshot.TotalEvents)
	fmt.Printf("Unique Users: %d\n", snapshot.UniqueUsers)
	fmt.Printf("Active Sessions: %d\n", snapshot.ActiveSessions)

	fmt.Println("\nEvents by Type:")
	for eventType, count := range snapshot.EventsByType {
		fmt.Printf("  %s: %d\n", eventType, count)
	}

	if len(snapshot.TopPages) > 0 {
		fmt.Println("\nTop Pages:")
		for i, page := range snapshot.TopPages {
			if i >= 10 {
				break
			}
			fmt.Printf("  %s: %d views (%d unique visitors)\n",
				page.Path, page.Views, page.UniqueVisitors)
		}
	}

	if len(snapshot.TrafficSources) > 0 {
		fmt.Println("\nTop Traffic Sources:")
		for i, source := range snapshot.TrafficSources {
			if i >= 5 {
				break
			}
			fmt.Printf("  %s: %d visits (%.1f%%)\n",
				source.Source, source.Count, source.Percent)
		}
	}

	fmt.Printf("\nPerformance Metrics:")
	fmt.Printf("  Average Load Time: %.1fms\n", snapshot.PerformanceMetrics.AverageLoadTime)
	fmt.Printf("  Fast Pages: %d, Slow Pages: %d\n",
		snapshot.PerformanceMetrics.FastPagesCount,
		snapshot.PerformanceMetrics.SlowPagesCount)

	fmt.Println("===================================")
}

func main() {

	log.Printf("Starting enhanced consumer with brokers: %s, topic: %s, group: %s",
		constants.KafkaBrokers, constants.KafkaTopic, constants.ConsumerGroup)

	// Create analytics service
	analyticsService := analytics.NewService()

	// Add some default alert configurations
	analyticsService.AddAlert(models.AlertConfig{
		Name:          "High Load Time Alert",
		Type:          "performance",
		Metric:        "average_load_time",
		Threshold:     5000, // 5 seconds
		Operator:      "gt",
		Enabled:       true,
		WindowMinutes: 5,
	})

	analyticsService.AddAlert(models.AlertConfig{
		Name:          "Traffic Surge Alert",
		Type:          "traffic",
		Metric:        "total_events",
		Threshold:     1000, // 1000 events
		Operator:      "gt",
		Enabled:       true,
		WindowMinutes: 5,
	})

	// Create Kafka consumer
	consumer := kafka.NewConsumer([]string{constants.KafkaBrokers}, constants.KafkaTopic, constants.ConsumerGroup)
	defer consumer.Close()

	// Create consumer service
	consumerService := NewConsumerService(consumer, analyticsService)

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("\nReceived shutdown signal, printing final stats...")
		consumerService.printStats()
		cancel()
	}()

	// Start periodic stats reporting
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				consumerService.printStats()
			case <-ctx.Done():
				return
			}
		}
	}()

	// Start consuming events
	log.Println("Enhanced consumer started, waiting for events...")
	log.Println("Real-time analytics processing enabled with alerts")
	if err := consumer.ConsumeEvents(ctx, consumerService.processEvent); err != nil {
		if err == context.Canceled {
			log.Println("Consumer stopped gracefully")
		} else {
			log.Fatalf("Consumer error: %v", err)
		}
	}
}
