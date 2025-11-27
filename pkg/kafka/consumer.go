package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Hilina-t/go-kafka-analytics-pipeline/pkg/models"
	"github.com/segmentio/kafka-go"
)

// Consumer represents a Kafka consumer
type Consumer struct {
	reader  *kafka.Reader
	topic   string
	groupID string
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(brokers []string, topic, groupID string) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	return &Consumer{
		reader:  reader,
		topic:   topic,
		groupID: groupID,
	}
}

// ConsumeEvents consumes and processes events from Kafka
func (c *Consumer) ConsumeEvents(ctx context.Context, handler func(*models.AnalyticsEvent) error) error {
	log.Printf("Starting consumer for topic: %s, group: %s", c.topic, c.groupID)

	for {
		select {
		case <-ctx.Done():
			log.Println("Consumer context cancelled, shutting down")
			return ctx.Err()
		default:
			msg, err := c.reader.FetchMessage(ctx)
			if err != nil {
				return fmt.Errorf("failed to fetch message: %w", err)
			}

			var event models.AnalyticsEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				log.Printf("Failed to unmarshal event: %v", err)
				// Commit message even if unmarshal fails to avoid reprocessing
				if err := c.reader.CommitMessages(ctx, msg); err != nil {
					log.Printf("Failed to commit message: %v", err)
				}
				continue
			}

			log.Printf("Processing event - Type: %s, ID: %s, User: %s", event.Type, event.ID, event.UserID)

			if err := handler(&event); err != nil {
				log.Printf("Failed to process event: %v", err)
				continue
			}

			if err := c.reader.CommitMessages(ctx, msg); err != nil {
				log.Printf("Failed to commit message: %v", err)
			}
		}
	}
}

// Close closes the consumer
func (c *Consumer) Close() error {
	return c.reader.Close()
}
