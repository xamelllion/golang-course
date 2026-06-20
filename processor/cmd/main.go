package main

import (
	"context"
	"fmt"
	"log"
	"net"

	kafkaClient "github.com/xamelllion/golang-course/internal/kafka"
	"github.com/xamelllion/golang-course/processor/internal/adapter"
	"github.com/xamelllion/golang-course/processor/internal/config"
	controllerGRPC "github.com/xamelllion/golang-course/processor/internal/controller/grpc"
	controllerKafka "github.com/xamelllion/golang-course/processor/internal/controller/kafka"
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

	conn, err := grpc.NewClient(cfg.SubscriberAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	pool, err := pgxpool.New(context.Background(), cfg.DB_DSN)
	if err != nil {
		log.Fatalf("create pg pool: %v", err)
		panic(err)
	}
	defer pool.Close()

	subscriber := adapter.NewSubscriber(conn)
	postgres := adapter.NewPostgresAdapter(pool)

	kafkaRequestProducer := kafkaClient.NewKafkaProducer([]string{cfg.KafkaBroker}, "repo.tasks.request")
	kafkaResponseConsumer := kafkaClient.NewKafkaReader([]string{cfg.KafkaBroker}, "repo.tasks.response", "processor")
	kafkaAdapter := adapter.NewKafkaAdapter(kafkaRequestProducer, kafkaResponseConsumer)

	repositoryUsecase := usecase.NewRepositoryUsecase(postgres, kafkaAdapter, subscriber)
	pingUsecase := usecase.NewPingUsecase()
	server := controllerGRPC.NewServer(repositoryUsecase, pingUsecase)

	kafkaController := controllerKafka.NewKafkaController(repositoryUsecase, kafkaResponseConsumer)

	go func() {
		if err := kafkaController.Run(context.Background()); err != nil {
			log.Print("kafka error $v", err)
		}
	}()

	grpcServer := grpc.NewServer()
	pbProcessor.RegisterProcessorServiceServer(grpcServer, server)

	log.Printf("grpc listen on %s port", cfg.Port)

	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}
