package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port                  string
	ProcessorAddr         string
	CollectorAddr         string
	SubscriberAddr        string
	RedisAddr             string
	RateLimitCapacity     int
	RateLimitReqPerSecond float64
	CacheTTL              time.Duration
}

func Load() Config {
	var config Config

	config.Port = getenv("GATEWAY_PORT", "8080")
	config.ProcessorAddr = getenv("PROCESSOR_ADDR", "localhost:50051")
	config.CollectorAddr = getenv("COLLECTOR_ADDR", "localhost:50052")
	config.SubscriberAddr = getenv("SUBSCRIBER_ADDR", "localhost:50053")
	config.RedisAddr = getenv("REDIS_ADDR", "localhost:6379")

	config.RateLimitCapacity = getenvInt("RATELIMIT_CAPACITY", 10)
	config.RateLimitReqPerSecond = getenvFloat64("RATELIMIT_REQ_PER_SECOND", 5.0)
	config.CacheTTL = time.Duration(getenvInt("CACHE_TTL_SECONDS", 60)) * time.Second

	return config
}

func getenv(key, defaultValue string) string {
	var v string

	if v = os.Getenv(key); v == "" {
		v = defaultValue
	}

	return v
}

func getenvConvert[T any](key string, defaultValue T, converter func(string) (T, error)) T {
	if str := os.Getenv(key); str == "" {
		return defaultValue
	} else {
		result, err := converter(str)
		if err != nil {
			panic(fmt.Sprintf("wrong value type %v", err))
		}
		return result
	}
}

func getenvInt(key string, defaultValue int) int {
	return getenvConvert(key, defaultValue, strconv.Atoi)
}

func getenvFloat64(key string, defaultValue float64) float64 {
	converter := func(v string) (float64, error) {
		return strconv.ParseFloat(v, 64)
	}

	return getenvConvert(key, defaultValue, converter)
}
