package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/xamelllion/golang-course/collector/internal/adapter"
	"github.com/xamelllion/golang-course/collector/internal/config"
	apperror "github.com/xamelllion/golang-course/internal/errors"
	pb "github.com/xamelllion/golang-course/proto/collector"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type server struct {
	pb.UnimplementedCollectorServiceServer
	gh adapter.GithubRepositoryAdapter
}

func (s *server) GetRepository(ctx context.Context, req *pb.GetRepositoryRequest) (*pb.GetRepositoryResponse, error) {
	if req.Owner == "" || req.Repo == "" {
		return &pb.GetRepositoryResponse{}, status.Errorf(codes.InvalidArgument, "missing required fields: owner or repo")
	}

	repo, error := s.gh.GetRepository(req.Owner, req.Repo)
	if error != nil {
		switch apperror.CodeOf(error) {
		case apperror.CodeNotFound:
			return nil, status.Error(codes.NotFound, error.Error())
		case apperror.CodeInvalidArgument:
			return nil, status.Error(codes.InvalidArgument, error.Error())
		case apperror.CodeUnavailable:
			return nil, status.Error(codes.Unavailable, error.Error())
		default:
			return nil, status.Error(codes.Internal, error.Error())
		}
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
	cfg := config.Load()
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Port))
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterCollectorServiceServer(grpcServer, &server{gh: adapter.NewGithubRepositoryAdapter()})

	log.Printf("grpc listen on %s port", cfg.Port)

	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}
