package domain

import "time"

type TaskResponse struct {
	Owner       string
	Repo        string
	Description string
	Stars       int32
	Forks       int32
	CreateDate  time.Time
}

type Subscription struct {
	Owner string
	Repo  string
}

type Repository struct {
	Name        string
	Description string
	Stars       int32
	Forks       int32
	CreateDate  time.Time
}
