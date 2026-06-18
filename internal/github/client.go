package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	apperror "github.com/xamelllion/golang-course/internal/errors"
)

type Client struct {
	Client http.Client
}

func NewClient(http http.Client) Client {
	return Client{Client: http}
}

func (c Client) makeGithubRequest(method, endpoint string) (*http.Response, error) {
	url := "https://api.github.com" + endpoint

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, apperror.WrapCode(apperror.CodeInternal, "making request", err)
	}

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, apperror.WrapCode(apperror.CodeUnavailable, "newtwork problem", err)
	}

	return res, nil
}

func (c Client) getRepo(owner, repo string) (githubRepo, error) {
	res, err := c.makeGithubRequest(http.MethodGet, fmt.Sprintf("/repos/%v/%v", owner, repo))
	if err != nil {
		return githubRepo{}, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return githubRepo{}, apperror.New(apperror.CodeNotFound, "not found repo")
	}

	if res.StatusCode != http.StatusOK {
		return githubRepo{}, apperror.New(apperror.CodeInternal, fmt.Sprintf("unexpected error while getting repo with status code %v\nres: %v", res.StatusCode, res))
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return githubRepo{}, apperror.Wrap("get repo", err)
	}

	result := githubRepo{}
	if err := json.Unmarshal(data, &result); err != nil {
		return githubRepo{}, apperror.Wrap("get repo:", err)
	}

	return result, nil
}

func (c Client) GetRepository(owner, repo string) (RepositoryDTO, error) {
	repoData, err := c.getRepo(owner, repo)
	if err != nil {
		return RepositoryDTO{}, err
	}

	createDate, err := time.Parse(time.RFC3339, repoData.CreationDate)
	if err != nil {
		return RepositoryDTO{}, apperror.Wrap("unexpected error while parsing date", err)
	}

	return RepositoryDTO{
		Name:         repoData.Name,
		Description:  repoData.Description,
		Stars:        repoData.StarsCount,
		Forks:        repoData.ForksCount,
		CreationDate: createDate,
	}, nil
}
