package config

import "os"

type Config struct {
	Port string
}

func Load() Config {
	var config Config

	config.Port = getenv("PROCESSOR_PORT", "50051")

	return config
}

func getenv(key, defaultValue string) string {
	var v string

	if v = os.Getenv(key); v == "" {
		v = defaultValue
	}

	return v
}
