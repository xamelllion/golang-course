package main

import (
	"log"
	"net/http"

	"github.com/xamelllion/golang-course/gateway/internal/adapter/collector"
	controller "github.com/xamelllion/golang-course/gateway/internal/controller/http"
	"github.com/xamelllion/golang-course/gateway/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, error := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if error != nil {
		panic(error)
	}
	defer conn.Close()

	collector := collector.NewCollector(conn)

	repositoryUseCase := usecase.NewRepositoryUseCase(collector)

	handler := controller.NewHandler(repositoryUseCase)

	http.HandleFunc("/", handler.GetRepository)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
