package usecase

import (
	"context"

	"github.com/xamelllion/golang-course/processor/internal/domain"
)

type RepoKeeper interface {
	GetRepositories(ctx context.Context) ([]domain.Repository, error)
	GetRepository(ctx context.Context, owner, repoName string) (domain.Repository, error)
	IsExistsRepo(ctx context.Context, owner, repo string) (bool, error)
	CreateRepo(ctx context.Context, repo domain.Repository, owner, repoName string) error
}

type Tasker interface {
	SendTaskRequest(owner, repo string) error
}

type RepositoryUseCase struct {
	repoKeeper RepoKeeper
	tasker     Tasker
}

func NewRepositoryUsecase(repoKeeper RepoKeeper, tasker Tasker) *RepositoryUseCase {
	return &RepositoryUseCase{repoKeeper: repoKeeper, tasker: tasker}
}

func (r *RepositoryUseCase) GetRepository(owner, repoName string) (*domain.Repository, error) {
	ctx := context.Background()

	ok, err := r.repoKeeper.IsExistsRepo(ctx, owner, repoName)
	if err != nil {
		return nil, err
	}

	if !ok {
		if err := r.tasker.SendTaskRequest(owner, repoName); err != nil {
			return nil, err
		}

		return nil, nil
	}

	repo, err := r.repoKeeper.GetRepository(ctx, owner, repoName)
	if err != nil {
		return nil, err
	}

	return &repo, nil
}

func (r *RepositoryUseCase) GetSubscribedRepository() ([](*domain.Repository), error) {
	// TODO
	return nil, nil
	// return r.Collector.GetSubscribedRepository()
}
