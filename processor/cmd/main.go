package main

import (
	"context"
	"fmt"
	"log"
	"net"

	kafkaClient "github.com/xamelllion/golang-course/internal/kafka"
	"github.com/xamelllion/golang-course/processor/internal/adapter"
	"github.com/xamelllion/golang-course/processor/internal/config"
	controller "github.com/xamelllion/golang-course/processor/internal/controller/grpc"
	"github.com/xamelllion/golang-course/processor/internal/usecase"
	pbProcessor "github.com/xamelllion/golang-course/proto/processor"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.Load()
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Port))
	if err != nil {
		panic(err)
	}

	conn, err := grpc.NewClient(cfg.CollectorAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	pool, err := pgxpool.New(context.Background(), cfg.DB_DSN)
	if err != nil {
		log.Fatalf("create pg pool: %v", err)
	}
	defer pool.Close()

	// collector := adapter.NewCollector(conn)
	postgres := adapter.NewPostgresAdapter(pool)

	kafkaRequestProducer := kafkaClient.NewKafkaProducer([]string{cfg.KafkaBroker}, "repo.tasks.request")
	kafkaResponseConsumer := kafkaClient.NewKafkaReader([]string{cfg.KafkaBroker}, "repo.tasks.response", "processor")
	kafka := adapter.NewKafkaAdapter(kafkaRequestProducer, kafkaResponseConsumer)

	repositoryUsecase := usecase.NewRepositoryUsecase(postgres, kafka)
	pingUsecase := usecase.NewPingUsecase()
	server := controller.NewServer(repositoryUsecase, pingUsecase)

	grpcServer := grpc.NewServer()
	pbProcessor.RegisterProcessorServiceServer(grpcServer, server)

	log.Printf("grpc listen on %s port", cfg.Port)

	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}
