package controller

import (
	"log"
	"time"
)

type Refresher interface {
	Refresh() error
}

type Scheduler struct {
	refresher Refresher
	interval  time.Duration
}

func NewScheduler(refresher Refresher, interval time.Duration) Scheduler {
	return Scheduler{refresher: refresher, interval: interval}
}

func (sh Scheduler) Run() error {
	ticker := time.NewTicker(sh.interval)
	defer ticker.Stop()

	for range ticker.C {
		err := sh.refresher.Refresh()
		if err != nil {
			log.Printf("scheduler internal error %v", err)
		}
	}

	return nil
}
