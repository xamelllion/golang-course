package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/xamelllion/golang-course/gateway/internal/adapter/collector"
	"github.com/xamelllion/golang-course/gateway/internal/config"
	controller "github.com/xamelllion/golang-course/gateway/internal/controller/http"
	"github.com/xamelllion/golang-course/gateway/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.Load()

	conn, error := grpc.NewClient(cfg.CollectorAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if error != nil {
		panic(error)
	}
	defer conn.Close()

	collector := collector.NewCollector(conn)

	repositoryUseCase := usecase.NewRepositoryUseCase(collector)

	handler := controller.NewHandler(repositoryUseCase)

	http.HandleFunc("/", handler.GetRepository)
	log.Printf("run server on %s port\n", cfg.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), nil))
}
