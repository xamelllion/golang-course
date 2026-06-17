package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/xamelllion/golang-course/gateway/internal/adapter/collector"
	"github.com/xamelllion/golang-course/gateway/internal/config"
	controller "github.com/xamelllion/golang-course/gateway/internal/controller/http"
	"github.com/xamelllion/golang-course/gateway/internal/usecase"
	"github.com/xamelllion/golang-course/internal/logger"
	httpSwagger "github.com/swaggo/http-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//	@title		TestApp API
//	@version	0.1

//	@contact.name	Xamelllion
//	@contact.url	t.me/xamelllion
//	@contact.email	me@xamelllion.ru

//	@license.name	MIT
//	@license.url	https://mit-license.org/

//	@host		localhost:8080
//	@BasePath	/api/v1

func main() {
	cfg := config.Load()
	logger := logger.Load("DEBUG")

	conn, error := grpc.NewClient(cfg.ProcessorAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if error != nil {
		panic(error)
	}
	defer conn.Close()

	collector := collector.NewCollector(conn, logger)

	repositoryUseCase := usecase.NewRepositoryUseCase(collector, logger)

	handler := controller.NewHandler(repositoryUseCase, logger)

	mux := http.NewServeMux()
	mux.HandleFunc("/docs/swagger/", httpSwagger.Handler(httpSwagger.URL(fmt.Sprintf("http://localhost:%s/docs/swagger/doc.json", cfg.Port))))
	mux.HandleFunc("/repo/{owner}/{repo}", handler.GetRepository)

	log.Printf("run server on %s port\n", cfg.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), mux))
}
