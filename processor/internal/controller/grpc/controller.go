package controllerGRPC

import (
	"context"

	apperror "github.com/xamelllion/golang-course/internal/errors"
	"github.com/xamelllion/golang-course/processor/internal/domain"
	pbProcessor "github.com/xamelllion/golang-course/proto/processor"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func fromDomainRepo(repo *domain.Repository) *pbProcessor.Repository {
	return &pbProcessor.Repository{
		Name:        repo.Name,
		Description: repo.Description,
		Stars:       repo.Stars,
		Forks:       repo.Forks,
		CreateDate:  timestamppb.New(repo.CreateDate),
		Status:      pbProcessor.RepositoryStatus_REPOSITORY_STATUS_READY,
	}
}

type RepositoryUsecase interface {
	GetRepository(owner, repo string) (*domain.Repository, error)
	GetSubscribedRepository() ([](*domain.Repository), error)
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
		return nil, status.Error(apperror.ToGRPCCode(error), error.Error())
	}

	// task created
	if repo == nil {
		return &pbProcessor.GetRepositoryResponse{
			Repository: &pbProcessor.Repository{
				Status: pbProcessor.RepositoryStatus_REPOSITORY_STATUS_PREPARING,
			},
		}, nil
	}

	return &pbProcessor.GetRepositoryResponse{
		Repository: fromDomainRepo(repo),
	}, nil
}

func (s *Server) GetSubscribedRepository(_ context.Context, _ *pbProcessor.GetSubscribedRepositoryRequest) (*pbProcessor.GetSubscribedRepositoryResponse, error) {
	repos, error := s.RepositoryUsecase.GetSubscribedRepository()
	if error != nil {
		return nil, status.Error(apperror.ToGRPCCode(error), error.Error())
	}

	result := make([]*pbProcessor.Repository, len(repos))
	for k, repo := range repos {
		if repo == nil {
			result[k] = &pbProcessor.Repository{
				Status: pbProcessor.RepositoryStatus_REPOSITORY_STATUS_PREPARING,
			}

			continue
		}

		result[k] = fromDomainRepo(repo)
	}

	return &pbProcessor.GetSubscribedRepositoryResponse{
		Repositories: result,
	}, nil
}

func (s *Server) Ping(_ context.Context, req *pbProcessor.PingRequest) (*pbProcessor.PingResponse, error) {
	pong, err := s.PingUsecase.Ping()
	if err != nil {
		return nil, err
	}

	return &pbProcessor.PingResponse{Reply: pong}, nil
}
