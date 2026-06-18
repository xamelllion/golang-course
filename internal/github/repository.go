package github

import (
	"fmt"
	"time"
)

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

type RepositoryDTO struct {
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Forks        int       `json:"forks_count"`
	Stars        int       `json:"stargazers_count"`
	CreationDate time.Time `json:"created_at"`
}
