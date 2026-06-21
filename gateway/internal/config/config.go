package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port                  string
	ProcessorAddr         string
	CollectorAddr         string
	SubscriberAddr        string
	RedisAddr             string
	RateLimitCapaciry     int
	RateLimitReqPerSecond float64
}

func Load() Config {
	var config Config

	config.Port = getenv("GATEWAY_PORT", "8080")
	config.ProcessorAddr = getenv("PROCESSOR_ADDR", "localhost:50051")
	config.CollectorAddr = getenv("COLLECTOR_ADDR", "localhost:50052")
	config.SubscriberAddr = getenv("SUBSCRIBER_ADDR", "localhost:50053")
	config.RedisAddr = getenv("REDIS_ADDR", "localhost:6379")

	capacity, err := strconv.Atoi(getenv("RATELIMIT_CAPACIRY", "10"))
	if err != nil {
		panic(fmt.Sprintf("wrong value type %v", err))
	}
	reqPerSec, err := strconv.ParseFloat(getenv("RATELIMIT_REQ_PER_SECOND", "5.0"), 64)
	if err != nil {
		panic(fmt.Sprintf("wrong value type %v", err))
	}

	config.RateLimitCapaciry = capacity
	config.RateLimitReqPerSecond = reqPerSec

	return config
}

func getenv(key, defaultValue string) string {
	var v string

	if v = os.Getenv(key); v == "" {
		v = defaultValue
	}

	return v
}
