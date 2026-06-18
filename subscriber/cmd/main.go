package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	pbSubscriber "github.com/xamelllion/golang-course/proto/subscriber"
	"github.com/xamelllion/golang-course/subscriber/internal/config"
	controller "github.com/xamelllion/golang-course/subscriber/internal/controller/grpc"
	"github.com/xamelllion/golang-course/subscriber/internal/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_DSN"))
	if err != nil {
		log.Fatalf("create pg pool: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("ping postgres: %v", err)
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Port))
	if err != nil {
		panic(err)
	}

	// postgresAdapter := adapter.NewSubscriptionPostgresAdapter(pool)

	usecase := usecase.NewPingUsecase()
	server := controller.NewServer(usecase)

	grpcServer := grpc.NewServer()
	pbSubscriber.RegisterSubscriberServiceServer(grpcServer, server)

	log.Printf("grpc listen on %s port", cfg.Port)

	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}
