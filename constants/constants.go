package constants

import "github.com/Hilina-t/go-kafka-analytics-pipeline/utils"

var (
	// Get configuration from environment variables
	KafkaBrokers  = utils.GetEnv("KAFKA_BROKERS", "localhost:9092")
	KafkaTopic    = utils.GetEnv("KAFKA_TOPIC", "analytics-events")
	ServerPort    = utils.GetEnv("SERVER_PORT", "8080")
	ConsumerGroup = utils.GetEnv("CONSUMER_GROUP", "analytics-consumer-group")
)
