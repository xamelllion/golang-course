package domain

import "time"

type Repository struct {
	Name        string
	Description string
	Stars       int32
	Forks       int32
	CreateDate time.Time
}
