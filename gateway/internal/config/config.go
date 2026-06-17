package config

import "os"

type Config struct {
	Port          string
	ProcessorAddr string
}

func Load() Config {
	var config Config

	config.Port = getenv("GATEWAY_PORT", "8080")
	config.ProcessorAddr = getenv("PROCESSOR_ADDR", "localhost:50051")

	return config
}

func getenv(key, defaultValue string) string {
	var v string

	if v = os.Getenv(key); v == "" {
		v = defaultValue
	}

	return v
}
