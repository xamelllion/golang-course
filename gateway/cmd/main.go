package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/xamelllion/golang-course/gateway/internal/adapter/collector"
	"github.com/xamelllion/golang-course/gateway/internal/adapter/processor"
	"github.com/xamelllion/golang-course/gateway/internal/adapter/subscriber"
	"github.com/xamelllion/golang-course/gateway/internal/config"
	controller "github.com/xamelllion/golang-course/gateway/internal/controller/http"
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

	handler := controller.NewHandler(repositoryUseCase, pingUseCase, subscribeUseCase, logger)

	mux := http.NewServeMux()
	mux.HandleFunc("/docs/swagger/", httpSwagger.Handler(httpSwagger.URL(fmt.Sprintf("http://localhost:%s/docs/swagger/doc.json", cfg.Port))))
	mux.HandleFunc("GET /api/repositories/info", handler.GetRepository)
	mux.HandleFunc("GET /api/ping", handler.PingServices)
	mux.HandleFunc("POST /api/subscriptions", handler.Subscribe)
	mux.HandleFunc("DELETE /api/subscriptions/{owner}/{repo}", handler.Unsubscribe)
	mux.HandleFunc("GET /api/subscriptions", handler.GetSubscriptions)
	mux.HandleFunc("GET /api/subscriptions/info", handler.GetSubscribedRepositories)

	log.Printf("run server on %s port\n", cfg.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), mux))
}
