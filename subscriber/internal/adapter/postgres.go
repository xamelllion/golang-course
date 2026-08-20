package adapter

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	subscriber "github.com/xamelllion/golang-course/subscriber/internal/adapter/postgres/sqlc"
	"github.com/xamelllion/golang-course/subscriber/internal/domain"
)

type SubscriptionPostgresAdapter struct {
	Pool  *pgxpool.Pool
	Query *subscriber.Queries
}

func NewSubscriptionPostgresAdapter(pool *pgxpool.Pool) SubscriptionPostgresAdapter {
	return SubscriptionPostgresAdapter{Pool: pool, Query: subscriber.New(pool)}
}

func (a SubscriptionPostgresAdapter) GetSubscriptions(ctx context.Context) ([]domain.Subscription, error) {
	subs, err := a.Query.GetSubscriptions(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]domain.Subscription, len(subs))
	for k, v := range subs {
		result[k] = domain.Subscription{Owner: v.Owner, Repo: v.Repo}
	}

	return result, nil
}

func (a SubscriptionPostgresAdapter) IsExistsSubscription(ctx context.Context, subi domain.Subscription) (bool, error) {
	result, err := a.Query.IsExistsSubscription(ctx, subscriber.IsExistsSubscriptionParams{Owner: subi.Owner, Repo: subi.Repo})
	if err != nil {
		return false, err
	}

	return result, nil
}

func (a SubscriptionPostgresAdapter) CreateSubscription(ctx context.Context, sub domain.Subscription) error {
	_, err := a.Query.CreateSubscription(ctx, subscriber.CreateSubscriptionParams{Owner: sub.Owner, Repo: sub.Repo})

	return err
}

func (a SubscriptionPostgresAdapter) DeleteSubscription(ctx context.Context, sub domain.Subscription) error {
	return a.Query.DeleteSubscription(ctx, subscriber.DeleteSubscriptionParams{Owner: sub.Owner, Repo: sub.Repo})
}
