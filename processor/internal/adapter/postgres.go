package adapter

import (
	"context"
	"log"

	apperror "github.com/xamelllion/golang-course/internal/errors"
	processor "github.com/xamelllion/golang-course/processor/internal/adapter/postgres/sqlc"
	"github.com/xamelllion/golang-course/processor/internal/domain"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresAdapter struct {
	Pool  *pgxpool.Pool
	Query *processor.Queries
}

func repoToDomain(repo processor.Repository) domain.Repository {
	return domain.Repository{
		Name:        repo.Repo,
		Description: repo.Description,
		Stars:       repo.Stars,
		Forks:       repo.Forks,
		CreateDate:  repo.CreatedAt.Time}
}

func repoFromDomain(repo domain.Repository, owner string) processor.Repository {
	return processor.Repository{
		Owner:       owner,
		Repo:        repo.Name,
		Description: repo.Description,
		Stars:       repo.Stars,
		Forks:       repo.Forks,
		CreatedAt:   pgtype.Timestamptz{Time: repo.CreateDate},
	}
}

func NewPostgresAdapter(pool *pgxpool.Pool) PostgresAdapter {
	return PostgresAdapter{Pool: pool, Query: processor.New(pool)}
}

func (a PostgresAdapter) GetRepositories(ctx context.Context) ([]domain.Repository, error) {
	repos, err := a.Query.GetRepositories(ctx)
	if err != nil {
		return nil, apperror.Wrap("get repo postgres", err)
	}

	result := make([]domain.Repository, len(repos))
	for k, v := range repos {
		result[k] = repoToDomain(v)
	}

	return result, nil
}

func (a PostgresAdapter) GetRepository(ctx context.Context, owner, repoName string) (domain.Repository, error) {
	repo, err := a.Query.GetRepository(ctx, processor.GetRepositoryParams{Owner: owner, Repo: repoName})
	if err != nil {
		return domain.Repository{}, apperror.Wrap("get repo postgres", err)
	}

	return repoToDomain(repo), nil
}

func (a PostgresAdapter) IsExistsRepo(ctx context.Context, owner, repo string) (bool, error) {
	result, err := a.Query.IsExistsRepository(ctx, processor.IsExistsRepositoryParams{Owner: owner, Repo: repo})
	if err != nil {
		return false, apperror.Wrap("is exists repo postgres", err)
	}

	return result, nil
}

func (a PostgresAdapter) CreateRepo(ctx context.Context, repo domain.Repository, owner, repoName string) error {
	_, err := a.Query.CreateRepository(ctx, processor.CreateRepositoryParams{
		Owner:       owner,
		Repo:        repoName,
		Description: repo.Description,
		Stars:       repo.Stars,
		Forks:       repo.Forks,
		CreatedAt:   pgtype.Timestamptz{Time: repo.CreateDate, Valid: true},
	})

	if err != nil {
		return apperror.Wrap("create repo postgres", err)
	}

	log.Printf("postgres: create repo %v/%v", owner, repoName)

	return nil
}

func (a PostgresAdapter) UpdateRepo(ctx context.Context, repo domain.Repository, owner, repoName string) error {
	err := a.Query.UpdateRepository(ctx, processor.UpdateRepositoryParams{
		Owner:       owner,
		Repo:        repoName,
		Description: repo.Description,
		Stars:       repo.Stars,
		Forks:       repo.Forks,
		CreatedAt:   pgtype.Timestamptz{Time: repo.CreateDate, Valid: true},
	})

	if err != nil {
		return apperror.Wrap("update repo postgres", err)
	}

	log.Printf("postgres: update repo %v/%v", owner, repoName)

	return nil
}
