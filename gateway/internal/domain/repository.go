package domain

import "time"

type Subscription struct {
	Owner string `json:"owner" example:"octocat"`
	Repo  string `json:"repo" example:"Hello-World"`
}

type Repository struct {
	Name        string    `json:"name" example:"xamelllion"`
	Description string    `json:"description" example:"xamelllion's repo"`
	Stars       int32     `json:"stars" example:"122"`
	Forks       int32     `json:"forks" example:"2"`
	CreateDate  time.Time `json:"create_date" example:"2026-03-03T16:45:04Z"`
}
