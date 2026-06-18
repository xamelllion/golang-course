package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/xamelllion/golang-course/collector/internal/adapter"
	"github.com/xamelllion/golang-course/collector/internal/config"
	"github.com/xamelllion/golang-course/collector/internal/domain"
	apperror "github.com/xamelllion/golang-course/internal/errors"
	"github.com/xamelllion/golang-course/internal/github"
	pb "github.com/xamelllion/golang-course/proto/collector"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type server struct {
	pb.UnimplementedCollectorServiceServer
	gh  adapter.GithubRepositoryAdapter
	rp  adapter.Subscriber
	cfg config.Config
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

func (s *server) GetSubscribedRepository(ctx context.Context, req *pb.GetSubscribedRepositoryRequest) (*pb.GetSubscribedRepositoryResponse, error) {
	subs, err := s.rp.GetSubscriptions()
	if err != nil {
		return nil, err
	}

	repos := make([]domain.GithubRepository, len(subs))
	for k, sub := range subs {
		repo, err := s.gh.GetRepository(sub.Owner, sub.Repo)
		if err != nil {
			return nil, status.Error(apperror.ToGRPCCode(err), err.Error())
		}
		repos[k] = repo
	}

	results := make([](*pb.GetRepositoryResponse), len(subs))
	for k, repo := range repos {
		results[k] = &pb.GetRepositoryResponse{
			Name:        repo.Name,
			Description: repo.Description,
			Stars:       repo.Stars,
			Forks:       repo.Forks,
			CreateDate:  timestamppb.New(repo.Create_date),
		}
	}

	return &pb.GetSubscribedRepositoryResponse{Repositories: results}, nil
}

func (s *server) Ping(_ context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	return &pb.PingResponse{Reply: "pong"}, nil
}

func main() {
	cfg := config.Load()
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Port))
	if err != nil {
		panic(err)
	}

	conn, err := grpc.NewClient(
		cfg.SubscriberAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	subscriptionAdapter := adapter.NewSubscriber(conn)

	client := http.Client{Timeout: 10 * time.Second}
	ghClient := github.NewClient(client)

	grpcServer := grpc.NewServer()
	pb.RegisterCollectorServiceServer(grpcServer, &server{
		gh:  adapter.NewGithubRepositoryAdapter(ghClient),
		rp:  subscriptionAdapter,
		cfg: cfg,
	})

	log.Printf("grpc listen on %s port", cfg.Port)

	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}
