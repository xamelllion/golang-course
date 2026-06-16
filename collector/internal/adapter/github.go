package adapter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/xamelllion/golang-course/collector/internal/domain"
	apperror "github.com/xamelllion/golang-course/internal/errors"
)

type GithubRepositoryAdapter interface {
	GetRepository(owner, repo string) (domain.GithubRepository, error)
}

func NewGithubRepositoryAdapter() GithubRepositoryAdapter {
	return &githubRepositoryAdapter{}
}

type githubRepo struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	ForksCount   int    `json:"forks_count"`
	StarsCount   int    `json:"stargazers_count"`
	CreationDate string `json:"created_at"`
}

func (gr githubRepo) String() string {
	return fmt.Sprintf("%v:\n"+
		"\tdescription: %v\n"+
		"\tstart count: %v\n"+
		"\tforks count: %v\n"+
		"\tdate creation: %v\n", gr.Name, gr.Description, gr.StarsCount, gr.ForksCount, gr.CreationDate)
}

func makeGithubRequest(method, endpoint string) (*http.Response, error) {
	client := http.Client{
		Timeout: time.Second * 10,
	}
	url := "https://api.github.com" + endpoint

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, apperror.WrapCode(apperror.CodeInternal, "making request", err)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, apperror.WrapCode(apperror.CodeInternal, "newtwork problem", err)
	}

	return res, nil
}

func getRepo(repoName string) (githubRepo, error) {
	repoSplited := strings.Split(repoName, "/")
	if len(repoSplited) < 2 || len(repoSplited[len(repoSplited)-1]) == 0 || len(repoSplited[len(repoSplited)-2]) == 0 {
		return githubRepo{}, apperror.New(apperror.CodeInvalidArgument, "bad repo and owner arguments")
	}
	owner, name := repoSplited[len(repoSplited)-2], repoSplited[len(repoSplited)-1]

	res, err := makeGithubRequest(http.MethodGet, fmt.Sprintf("/repos/%v/%v", owner, name))
	if err != nil {
		return githubRepo{}, err
	}
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
	defer res.Body.Close()

	repo := githubRepo{}
	if err := json.Unmarshal(data, &repo); err != nil {
		return githubRepo{}, apperror.Wrap("get repo:", err)
	}

	return repo, nil
}

type githubRepositoryAdapter struct {
}

func (gra *githubRepositoryAdapter) GetRepository(owner, repo string) (domain.GithubRepository, error) {
	repoData, err := getRepo(fmt.Sprintf("%v/%v", owner, repo))
	if err != nil {
		return domain.GithubRepository{}, err
	}

	createDate, err := time.Parse(time.RFC3339, repoData.CreationDate)
	if err != nil {
		return domain.GithubRepository{}, apperror.Wrap("unexpected error while parsing date", err)
	}

	return domain.GithubRepository{
		Name:        repoData.Name,
		Description: repoData.Description,
		Stars:       int32(repoData.StarsCount),
		Forks:       int32(repoData.ForksCount),
		Create_date: createDate,
	}, nil
}
