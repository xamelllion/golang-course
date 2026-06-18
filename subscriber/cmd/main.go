package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pbSubscriber "github.com/xamelllion/golang-course/proto/subscriber"
	"github.com/xamelllion/golang-course/subscriber/internal/config"
	controller "github.com/xamelllion/golang-course/subscriber/internal/controller/grpc"
	"github.com/xamelllion/golang-course/subscriber/internal/usecase"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()

	conn, err := pgx.Connect(ctx, cfg.DB_DSN)
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Port))
	if err != nil {
		panic(err)
	}

	usecase := usecase.NewPingUsecase()
	server := controller.NewServer(usecase)

	grpcServer := grpc.NewServer()
	pbSubscriber.RegisterSubscriberServiceServer(grpcServer, server)

	log.Printf("grpc listen on %s port", cfg.Port)

	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}
