package config

import "os"

type Config struct {
	Port string
	SubscriberAddr string
	KafkaBrokers string
}

func Load() Config {
	var config Config

	config.Port = getenv("COLLECTOR_PORT", "50052")
	config.SubscriberAddr = getenv("SUBSCRIBER_ADDR", "localhost:50053")
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
