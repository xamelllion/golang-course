package usecase

import "github.com/xamelllion/golang-course/gateway/internal/domain"

type Collector interface {
	GetRepository(owner, repo string) (domain.Repository, error)
}

type RepositoryUseCase struct {
	Collector Collector
}

func NewRepositoryUseCase(collector Collector) *RepositoryUseCase {
	return &RepositoryUseCase{Collector: collector}
}

func (r *RepositoryUseCase) GetRepository(owner, repo string) (domain.Repository, error) {
	return r.Collector.GetRepository(owner, repo)
}
