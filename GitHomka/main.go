package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type GithubRepo struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	ForksCount   int    `json:"forks_count"`
	StarsCount   int    `json:"stargazers_count"`
	CreationDate string `json:"created_at"`
}

func (gr GithubRepo) String() string {
	return fmt.Sprintf("%v:\n"+
		"\tdescription: %v\n"+
		"\tstart count: %v\n"+
		"\tforks count: %v\n"+
		"\tdate creation: %v\n", gr.Name, gr.Description, gr.StarsCount, gr.ForksCount, gr.CreationDate)
}

func makeGithubRequest(method, endpoint string) (*http.Response, error) {
	client := http.Client{
		Timeout: time.Second * 2,
	}
	url := "https://api.github.com" + endpoint

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func getRepo(repoName string) (GithubRepo, error) {
	repoSplited := strings.Split(repoName, "/")
	if len(repoSplited) < 2 || len(repoSplited[len(repoSplited)-1]) == 0 || len(repoSplited[len(repoSplited)-2]) == 0 {
		return GithubRepo{}, fmt.Errorf(`wrong repository name %v, must be "owner/repo"`, repoName)
	}
	owner, name := repoSplited[len(repoSplited)-2], repoSplited[len(repoSplited)-1]

	res, err := makeGithubRequest(http.MethodGet, fmt.Sprintf("/repos/%v/%v", owner, name))
	if err != nil {
		return GithubRepo{}, err
	}
	if res.StatusCode == http.StatusNotFound {
		return GithubRepo{}, fmt.Errorf(`there is no repo "%v"`, repoName)
	}

	if res.StatusCode != http.StatusOK {
		return GithubRepo{}, fmt.Errorf("unexpected error while getting repo with status code %v\nres: %v", res.StatusCode, res)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return GithubRepo{}, err
	}
	defer res.Body.Close()

	repo := GithubRepo{}
	if err := json.Unmarshal(data, &repo); err != nil {
		return GithubRepo{}, err
	}

	return repo, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stdout, "Usage: %v owner/repo\n\n"+
			"TestApp is a cli tool that provide you information about repo\n", os.Args[0])
		return
	}

	repo, err := getRepo(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "[Error] %v\n", err)
		return
	}

	fmt.Println(repo)
}
