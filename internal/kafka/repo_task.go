package kafkaClient

import "time"

type RepoTaskRequest struct {
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
}

const (
	STATUS_OK               = "ok"
	STATUS_NOT_FOUND        = "not_found"
	STATUS_INVALID_ARGUMENT = "invalid_argument"
	STATUS_UNAVAILABLE      = "unavailable"
	STATUS_INTERNAL         = "internal"
)

type RepoTaskResponse struct {
	Owner       string    `json:"owner"`
	Repo        string    `json:"repo"`
	Description string    `json:"description"`
	Stars       int       `json:"stars"`
	Forks       int       `json:"forks"`
	CreatedAt   time.Time `json:"created_at"`
	Status      string    `json:"status"`
	Error       string    `json:"error"`
}
