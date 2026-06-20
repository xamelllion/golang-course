package config

import "os"

type Config struct {
	Port           string
	CollectorAddr  string
	SubscriberAddr string
	KafkaBroker    string
	DB_DSN         string
}

func Load() Config {
	var config Config

	config.Port = getenv("PROCESSOR_PORT", "50051")
	config.CollectorAddr = getenv("COLLECTOR_ADDR", "localhost:50052")
	config.SubscriberAddr = getenv("SUBSCRIBER_ADDR", "localhost:50053")
	config.KafkaBroker = getenv("KAFKA_BROKER", "localhost:9092")
	config.DB_DSN = getenv("DB_DSN", "postgres://postgres:postgres@localhost:5433/repositories?sslmode=disable")

	return config
}

func getenv(key, defaultValue string) string {
	var v string

	if v = os.Getenv(key); v == "" {
		v = defaultValue
	}

	return v
}
