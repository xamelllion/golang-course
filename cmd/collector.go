package main

import (
	"context"
	"errors"
	"log"
	"net"
	"time"

	pb "github.com/xamelllion/golang-course/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type server struct {
	pb.UnimplementedCollectorServiceServer
}

func (s *server) GetRepository(ctx context.Context, req *pb.GetRepositoryRequest) (*pb.GetRepositoryResponse, error) {
	if req.Owner == "" || req.Repo == "" {
		return &pb.GetRepositoryResponse{}, errors.New("there is no owner or repo")
	}

	return &pb.GetRepositoryResponse{
		Name:        "Homka",
		Description: "Homka",
		Stars:       2,
		Forks:       2,
		CreateDate:  timestamppb.New(time.Now()),
	}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterCollectorServiceServer(grpcServer, &server{})

	log.Printf("grpc listen on 50051 port")

	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}
