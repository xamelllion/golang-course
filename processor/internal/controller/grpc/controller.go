package controller

import (
	"context"

	apperror "github.com/xamelllion/golang-course/internal/errors"
	"github.com/xamelllion/golang-course/processor/internal/domain"
	pbProcessor "github.com/xamelllion/golang-course/proto/processor"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type RepositoryUsecase interface {
	GetRepository(owner, repo string) (domain.Repository, error)
}

type PingUsecase interface {
	Ping() (string, error)
}

type Server struct {
	pbProcessor.UnimplementedProcessorServiceServer
	RepositoryUsecase RepositoryUsecase
	PingUsecase       PingUsecase
}

func NewServer(repositoryUsecase RepositoryUsecase, pingUsecase PingUsecase) *Server {
	return &Server{RepositoryUsecase: repositoryUsecase, PingUsecase: pingUsecase}
}

func (s *Server) GetRepository(_ context.Context, req *pbProcessor.GetRepositoryRequest) (*pbProcessor.GetRepositoryResponse, error) {
	if req.Owner == "" || req.Repo == "" {
		return &pbProcessor.GetRepositoryResponse{}, status.Errorf(codes.InvalidArgument, "missing required fields: owner or repo")
	}

	repo, error := s.RepositoryUsecase.GetRepository(req.Owner, req.Repo)
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

	return &pbProcessor.GetRepositoryResponse{
		Name:        repo.Name,
		Description: repo.Description,
		Stars:       repo.Stars,
		Forks:       repo.Forks,
		CreateDate:  timestamppb.New(repo.CreateDate),
	}, nil
}

func (s *Server) Ping(_ context.Context, req *pbProcessor.PingRequest) (*pbProcessor.PingResponse, error) {
	pong, err := s.PingUsecase.Ping()
	if err != nil {
		return nil, err
	}

	return &pbProcessor.PingResponse{Reply: pong}, nil
}
