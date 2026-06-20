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
	UpdateRepo(ctx context.Context, repo domain.Repository, owner, repoName string) error
}

type Tasker interface {
	SendTaskRequest(owner, repo string) error
}

type Subscriber interface {
	GetSubscriptions() ([]domain.Subscription, error)
}

type RepositoryUseCase struct {
	repoKeeper RepoKeeper
	tasker     Tasker
	subscriber Subscriber
}

func NewRepositoryUsecase(repoKeeper RepoKeeper, tasker Tasker, subscriber Subscriber) *RepositoryUseCase {
	return &RepositoryUseCase{repoKeeper: repoKeeper, tasker: tasker, subscriber: subscriber}
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
	subs, err := r.subscriber.GetSubscriptions()
	if err != nil {
		return nil, err
	}

	result := []*domain.Repository{}
	for _, sub := range subs {
		repo, err := r.GetRepository(sub.Owner, sub.Repo)
		if err != nil {
			return nil, err
		}

		if repo == nil {
			continue
		}

		result = append(result, repo)
	}

	return result, nil
}

func (r *RepositoryUseCase) ProcessTaskReponse(task domain.TaskResponse) error {
	ok, err := r.repoKeeper.IsExistsRepo(context.Background(), task.Owner, task.Repo)
	if err != nil {
		return err
	}

	repo := domain.Repository{
		Name:        task.Repo,
		Description: task.Description,
		Stars:       task.Stars,
		Forks:       task.Forks,
		CreateDate:  task.CreateDate,
	}

	if ok {
		return r.repoKeeper.UpdateRepo(context.Background(), repo, task.Owner, task.Repo)
	}

	return r.repoKeeper.CreateRepo(context.Background(), repo, task.Owner, task.Repo)
}
