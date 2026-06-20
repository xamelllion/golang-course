package adapter

import (
	"context"

	apperror "github.com/xamelllion/golang-course/internal/errors"
	"github.com/xamelllion/golang-course/processor/internal/domain"
	"google.golang.org/grpc"

	pb "github.com/xamelllion/golang-course/proto/subscriber"
)

type Subscriber struct {
	conn *grpc.ClientConn
}

func NewSubscriber(conn *grpc.ClientConn) Subscriber {
	return Subscriber{
		conn: conn,
	}
}

func (s Subscriber) GetSubscriptions() ([]domain.Subscription, error) {
	client := pb.NewSubscriberServiceClient(s.conn)

	subs, error := client.GetSubscriptions(context.Background(), &pb.GetSubscriptionsRequest{})
	if error != nil {
		return nil, apperror.FromGRPC(error, "processor get subscriptions")
	}

	result := make([]domain.Subscription, len(subs.Subscriptions))
	for k, sub := range subs.Subscriptions {
		result[k] = domain.Subscription{Owner: sub.Owner, Repo: sub.Repo}
	}

	return result, nil
}
