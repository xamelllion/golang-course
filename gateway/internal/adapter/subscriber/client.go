package subscriber

import (
	"context"
	"log/slog"

	"github.com/xamelllion/golang-course/gateway/internal/config"
	"github.com/xamelllion/golang-course/gateway/internal/domain"
	apperror "github.com/xamelllion/golang-course/internal/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pbSubscriber "github.com/xamelllion/golang-course/proto/subscriber"
)

type Subscriber struct {
	conn   *grpc.ClientConn
	client pbSubscriber.SubscriberServiceClient
	log    *slog.Logger
}

func NewSubscriber(cfg config.Config, log *slog.Logger) Subscriber {
	conn, error := grpc.NewClient(cfg.SubscriberAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if error != nil {
		panic(error)
	}

	client := pbSubscriber.NewSubscriberServiceClient(conn)

	return Subscriber{
		conn:   conn,
		client: client,
		log:    log,
	}
}

func (s Subscriber) Subscribe(owner, repo string) error {
	_, err := s.client.Subscribe(context.Background(), &pbSubscriber.SubscribeRequest{Owner: owner, Repo: repo})
	if err != nil {
		return apperror.FromGRPC(err, "subscriber subscribe")
	}

	return nil
}

func (s Subscriber) Unsubscribe(owner, repo string) error {
	_, err := s.client.Unsubscribe(context.Background(), &pbSubscriber.UnsubscribeRequest{Owner: owner, Repo: repo})
	if err != nil {
		return apperror.FromGRPC(err, "subscriber unsubscribe")
	}

	return nil
}

func (s Subscriber) GetSubscriptions() ([]domain.Subscription, error) {
	subs, err := s.client.GetSubscriptions(context.Background(), &pbSubscriber.GetSubscriptionsRequest{})
	if err != nil {
		return nil, apperror.FromGRPC(err, "subscriber get subscriptions")
	}

	result := make([]domain.Subscription, len(subs.Subscriptions))
	for k, sub := range subs.Subscriptions {
		result[k] = domain.Subscription{
			Owner: sub.Owner,
			Repo:  sub.Repo,
		}
	}

	return result, nil
}

func (s Subscriber) Ping() (string, error) {
	pong, err := s.client.Ping(context.Background(), &pbSubscriber.PingRequest{})
	if err != nil {
		return "", apperror.New(apperror.CodeUnavailable, "processor unvailable")
	}

	return pong.Reply, nil
}

func (s Subscriber) Close() {
	s.conn.Close()
}
