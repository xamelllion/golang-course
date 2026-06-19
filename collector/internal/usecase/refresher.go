package usecase

import (
	"log"

	"github.com/xamelllion/golang-course/collector/internal/domain"
)

type Subscriber interface {
	GetSubscriptions() ([]domain.Subscription, error)
}

type Tasker interface {
	SendTaskRequest(sub domain.Subscription) error
}

type RefresherUsecase struct {
	subscriber Subscriber
	tasker     Tasker
}

func NewRefresherUsecase(subscriber Subscriber, tasker Tasker) RefresherUsecase {
	return RefresherUsecase{subscriber: subscriber, tasker: tasker}
}

func (ru RefresherUsecase) Refresh() error {
	subs, err := ru.subscriber.GetSubscriptions()
	if err != nil {
		return err
	}

	for _, sub := range subs {
		err := ru.tasker.SendTaskRequest(sub)
		if err != nil {
			log.Printf("internal error %v", err)
			continue
		}
	}

	return nil
}
