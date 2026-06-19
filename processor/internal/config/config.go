package config

import "os"

type Config struct {
	Port          string
	CollectorAddr string
	KafkaBrokers  string
}

func Load() Config {
	var config Config

	config.Port = getenv("PROCESSOR_PORT", "50051")
	config.CollectorAddr = getenv("COLLECTOR_ADDR", "localhost:50052")
	config.KafkaBrokers = getenv("KAFKA_BROKERS", "localhost:9092")

	return config
}

func getenv(key, defaultValue string) string {
	var v string

	if v = os.Getenv(key); v == "" {
		v = defaultValue
	}

	return v
}
