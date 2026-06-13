package main

import (
	"context"
	"errors"
	"log"
	"net"

	"github.com/xamelllion/golang-course/internal/adapter"
	pb "github.com/xamelllion/golang-course/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type server struct {
	pb.UnimplementedCollectorServiceServer
	gh adapter.GithubRepositoryAdapter
}

func (s *server) GetRepository(ctx context.Context, req *pb.GetRepositoryRequest) (*pb.GetRepositoryResponse, error) {
	if req.Owner == "" || req.Repo == "" {
		return &pb.GetRepositoryResponse{}, errors.New("there is no owner or repo")
	}

	repo, error := s.gh.GetRepository(req.Owner, req.Repo)
	if error != nil {
		return &pb.GetRepositoryResponse{}, error
	}

	return &pb.GetRepositoryResponse{
		Name:        repo.Name,
		Description: repo.Description,
		Stars:       repo.Stars,
		Forks:       repo.Forks,
		CreateDate:  timestamppb.New(repo.Create_date),
	}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterCollectorServiceServer(grpcServer, &server{gh: adapter.NewGithubRepositoryAdapter()})

	log.Printf("grpc listen on 50051 port")

	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}
