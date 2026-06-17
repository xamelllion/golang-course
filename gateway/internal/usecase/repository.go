package usecase

import (
	"log/slog"

	"github.com/xamelllion/golang-course/gateway/internal/domain"
)

type Collector interface {
	GetRepository(owner, repo string) (domain.Repository, error)
}

type RepositoryUseCase struct {
	Collector Collector
	log       *slog.Logger
}

func NewRepositoryUseCase(collector Collector, log *slog.Logger) *RepositoryUseCase {
	return &RepositoryUseCase{Collector: collector, log: log}
}

func (r *RepositoryUseCase) GetRepository(owner, repo string) (domain.Repository, error) {
	r.log.Debug("usecase: get repository", "owner", owner, "repo", repo)
	return r.Collector.GetRepository(owner, repo)
}
