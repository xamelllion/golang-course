package usecase

import (
	"context"

	apperror "github.com/xamelllion/golang-course/internal/errors"
	"github.com/xamelllion/golang-course/subscriber/internal/domain"
)

type RepositoryAdapter interface {
	IsExistsRepository(owner, repo string) (bool, error)
}

type SubscriptionAdapter interface {
	IsExistsSubscription(ctx context.Context, sub domain.Subscription) (bool, error)
	GetSubscriptions(ctx context.Context) ([]domain.Subscription, error)
	CreateSubscription(ctx context.Context, sub domain.Subscription) error
	DeleteSubscription(ctx context.Context, sub domain.Subscription) error
}

type SubscriptionUsecase struct {
	RepositoryAdapter   RepositoryAdapter
	SubscriptionAdapter SubscriptionAdapter
}

func NewSubscriptionUsecase(repoAdapter RepositoryAdapter, subAdapter SubscriptionAdapter) SubscriptionUsecase {
	return SubscriptionUsecase{RepositoryAdapter: repoAdapter, SubscriptionAdapter: subAdapter}
}

func (u SubscriptionUsecase) Subscribe(subi domain.Subscription) error {
	ctx := context.Background()

	isExists, err := u.SubscriptionAdapter.IsExistsSubscription(ctx, subi)
	if err != nil {
		return err
	}

	if isExists {
		return apperror.New(apperror.CodeDublicate, "subscription already exists")
	}

	isExistsGithub, err := u.RepositoryAdapter.IsExistsRepository(subi.Owner, subi.Repo)
	if err != nil {
		return err
	}

	if !isExistsGithub {
		return apperror.New(apperror.CodeNotFound, "repo not found on github")
	}

	return u.SubscriptionAdapter.CreateSubscription(ctx, subi)
}

func (u SubscriptionUsecase) Unsubscribe(subi domain.Subscription) error {
	ctx := context.Background()

	isExists, err := u.SubscriptionAdapter.IsExistsSubscription(ctx, subi)
	if err != nil {
		return err
	}

	if !isExists {
		return apperror.New(apperror.CodeNotFound, "subscription not found")
	}

	return u.SubscriptionAdapter.DeleteSubscription(ctx, subi)
}

func (u SubscriptionUsecase) GetSubscriptions() ([]domain.Subscription, error) {
	ctx := context.Background()

	subs, err := u.SubscriptionAdapter.GetSubscriptions(ctx)
	if err != nil {
		return nil, err
	}

	return subs, nil
}
