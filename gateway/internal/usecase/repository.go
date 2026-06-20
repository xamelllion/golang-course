package usecase

import (
	"log/slog"

	"github.com/xamelllion/golang-course/gateway/internal/domain"
)

type Processor interface {
	GetRepository(owner, repo string) (*domain.Repository, error)
	GetSubscribedRepository() ([](*domain.Repository), error)
}

type RepositoryUseCase struct {
	Processor Processor
	log       *slog.Logger
}

func NewRepositoryUseCase(processor Processor, log *slog.Logger) *RepositoryUseCase {
	return &RepositoryUseCase{Processor: processor, log: log}
}

func (r *RepositoryUseCase) GetRepository(owner, repo string) (*domain.Repository, error) {
	r.log.Debug("usecase: get repository", "owner", owner, "repo", repo)
	return r.Processor.GetRepository(owner, repo)
}

func (r *RepositoryUseCase) GetSubscribedRepository() ([](*domain.Repository), error) {
	r.log.Debug("usecase: get subscribed repositories")
	return r.Processor.GetSubscribedRepository()
}
