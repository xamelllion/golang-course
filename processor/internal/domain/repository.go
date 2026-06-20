package domain

import "time"

type Subscription struct {
	Owner string
	Repo string
}

type Repository struct {
	Name        string
	Description string
	Stars       int32
	Forks       int32
	CreateDate time.Time
}
