package domain

import "time"

type GithubRepository struct {
	Name        string
	Description string
	Stars       int32
	Forks       int32
	Create_date time.Time
}
