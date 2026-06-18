package adapter

import (
	apperror "github.com/xamelllion/golang-course/internal/errors"
	"github.com/xamelllion/golang-course/internal/github"
)

type GithubRepositoryAdapter struct {
	GithubClient github.Client
}

func NewGithubRepositoryAdapter(client github.Client) GithubRepositoryAdapter {
	return GithubRepositoryAdapter{GithubClient: client}
}

func (gra GithubRepositoryAdapter) IsExistsRepository(owner, repo string) (bool, error) {
	_, err := gra.GithubClient.GetRepository(owner, repo)
	if err != nil {
		if apperror.CodeOf(err) == apperror.CodeNotFound {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
