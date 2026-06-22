package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/xamelllion/golang-course/gateway/internal/adapter/cache"
	"github.com/xamelllion/golang-course/gateway/internal/adapter/collector"
	"github.com/xamelllion/golang-course/gateway/internal/adapter/processor"
	"github.com/xamelllion/golang-course/gateway/internal/adapter/ratelimiter"
	"github.com/xamelllion/golang-course/gateway/internal/adapter/subscriber"
	"github.com/xamelllion/golang-course/gateway/internal/config"
	controller "github.com/xamelllion/golang-course/gateway/internal/controller/http"
	"github.com/xamelllion/golang-course/gateway/internal/middleware"
	"github.com/xamelllion/golang-course/gateway/internal/usecase"
	"github.com/xamelllion/golang-course/internal/logger"
	httpSwagger "github.com/swaggo/http-swagger"
)

//	@title		TestApp API
//	@version	0.1

//	@contact.name	Xamelllion
//	@contact.url	t.me/xamelllion
//	@contact.email	me@xamelllion.ru

//	@license.name	MIT
//	@license.url	https://mit-license.org/

//	@host		localhost:8080
//	@BasePath	/

func main() {
	cfg := config.Load()
	logger := logger.Load("DEBUG")

	processor := processor.NewProcessor(cfg, logger)
	collector := collector.NewCollector(cfg, logger)
	subscriber := subscriber.NewSubscriber(cfg, logger)

	repositoryUseCase := usecase.NewRepositoryUseCase(processor, logger)
	pingUseCase := usecase.NewPingUsecase(map[string]usecase.Pinger{
		"collector":  collector,
		"processor":  processor,
		"subscriber": subscriber,
	}, logger)
	subscribeUseCase := usecase.NewSubscribeUseCase(subscriber, logger)

	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	redisRateLimiter := ratelimiter.NewRedisRateLimiter(redisClient, float64(cfg.RateLimitReqPerSecond), cfg.RateLimitCapacity)
	memoryRateLimiter := ratelimiter.NewMemoryBucketRateLimiter(cfg.RateLimitCapacity, cfg.RateLimitReqPerSecond)
	ratelimiter := ratelimiter.NewFallbackRateLimiter(redisRateLimiter, memoryRateLimiter)
	rateLimitMiddleware := middleware.RateLimitMiddleware(ratelimiter)

	cacher := cache.NewRedisCacher(redisClient)
	cacherMiddleware := middleware.CacheMiddleware(cacher, cfg.CacheTTL)

	loggerMiddleware := middleware.LoggerMiddleware()

	handler := controller.NewHandler(repositoryUseCase, pingUseCase, subscribeUseCase, logger)

	mux := http.NewServeMux()
	swaggerHandleFunc := httpSwagger.Handler(httpSwagger.URL(fmt.Sprintf("http://localhost:%s/docs/swagger/doc.json", cfg.Port)))
	mux.Handle("/docs/swagger/", middleware.Chain(swaggerHandleFunc, rateLimitMiddleware))
	mux.Handle("GET /api/repositories/info", middleware.Chain(handler.GetRepository(), loggerMiddleware, rateLimitMiddleware, cacherMiddleware))
	mux.Handle("GET /api/ping", middleware.Chain(handler.PingServices(), loggerMiddleware, rateLimitMiddleware))
	mux.Handle("POST /api/subscriptions", middleware.Chain(handler.Subscribe(), loggerMiddleware, rateLimitMiddleware))
	mux.Handle("DELETE /api/subscriptions/{owner}/{repo}", middleware.Chain(handler.Unsubscribe(), loggerMiddleware, rateLimitMiddleware))
	mux.Handle("GET /api/subscriptions", middleware.Chain(handler.GetSubscriptions(), loggerMiddleware, rateLimitMiddleware))
	mux.Handle("GET /api/subscriptions/info", middleware.Chain(handler.GetSubscribedRepositories(), loggerMiddleware, rateLimitMiddleware, cacherMiddleware))

	log.Printf("run server on %s port\n", cfg.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), mux))
}
