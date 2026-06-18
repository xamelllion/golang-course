package usecase

import (
	"log/slog"

	"github.com/xamelllion/golang-course/gateway/internal/domain"
)

type Subscriber interface {
	Subscribe(owner, repo string) error
	Unsubscribe(owner, repo string) error
	GetSubscriptions() ([]domain.Subscription, error)
}

type SubscribeUseCase struct {
	Subscriber Subscriber
}

func NewSubscribeUseCase(subscriber Subscriber, log *slog.Logger) *SubscribeUseCase {
	return &SubscribeUseCase{Subscriber: subscriber}
}

func (u *SubscribeUseCase) Subscribe(owner, repo string) error {
	return u.Subscriber.Subscribe(owner, repo)
}

func (u *SubscribeUseCase) Unsubscribe(owner, repo string) error {
	return u.Subscriber.Unsubscribe(owner, repo)
}

func (u *SubscribeUseCase) GetSubscriptions() ([]domain.Subscription, error) {
	return u.Subscriber.GetSubscriptions()
}
