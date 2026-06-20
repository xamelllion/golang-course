package config

import "os"

type Config struct {
	GithubToken    string
	Port           string
	SubscriberAddr string
	KafkaBroker    string
}

func Load() Config {
	var config Config

	config.GithubToken = getenv("GITHUB_TOKEN", "")
	config.Port = getenv("COLLECTOR_PORT", "50052")
	config.SubscriberAddr = getenv("SUBSCRIBER_ADDR", "localhost:50053")
	config.KafkaBroker = getenv("KAFKA_BROKER", "localhost:9092")

	return config
}

func getenv(key, defaultValue string) string {
	var v string

	if v = os.Getenv(key); v == "" {
		v = defaultValue
	}

	return v
}
