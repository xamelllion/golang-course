package config

import "os"

type Config struct {
	Port   string
	DB_DSN string
}

func Load() Config {
	var config Config

	config.Port = getenv("SUBSCRIBER_PORT", "50053")
	config.DB_DSN = getenv("DB_DSN", "postgres://postgres:postgres@localhost:5432/subscribers?sslmode=disable")

	return config
}

func getenv(key, defaultValue string) string {
	var v string

	if v = os.Getenv(key); v == "" {
		v = defaultValue
	}

	return v
}
