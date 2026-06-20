package domain

import "time"

type Subscription struct {
	Owner string
	Repo  string
}

type Task struct {
	Owner string
	Repo  string
}

type GithubRepository struct {
	Name        string
	Description string
	Stars       int32
	Forks       int32
	Create_date time.Time
}
