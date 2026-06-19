package kafkaClient

import "time"

type RepoTaskRequest struct {
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
}

type RepoTaskResponse struct {
	Owner       string    `json:"owner"`
	Repo        string    `json:"repo"`
	Description string    `json:"description"`
	Stars       int       `json:"stars"`
	Forks       int       `json:"forks"`
	CreatedAt   time.Time `json:"created_at"`
}
