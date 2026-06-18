package postgres

import (
	"context"

	"github.com/xamelllion/golang-course/subscriber/internal/adapter/postgres/sqlc"
	"github.com/xamelllion/golang-course/subscriber/internal/domain"
	"github.com/jackc/pgx/v5"
)

type SubscriptionPostgresAdapter struct {
	Conn  *pgx.Conn
	Query *subscriber.Queries
}

func NewSubscriptionPostgresAdapter(conn *pgx.Conn) SubscriptionPostgresAdapter {
	return SubscriptionPostgresAdapter{Conn: conn, Query: subscriber.New(conn)}
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

func (a SubscriptionPostgresAdapter) CreateSubscription(ctx context.Context, sub domain.Subscription) error {
	_, err := a.Query.CreateSubscription(ctx, subscriber.CreateSubscriptionParams{Owner: sub.Owner, Repo: sub.Repo})

	return err
}

func (a SubscriptionPostgresAdapter) DeleteSubscription(ctx context.Context, sub domain.Subscription) error {
	return a.Query.DeleteSubscription(ctx, subscriber.DeleteSubscriptionParams{Owner: sub.Owner, Repo: sub.Repo})
}
