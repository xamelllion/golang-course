package subscriber

import (
	"context"
	"log/slog"

	"github.com/xamelllion/golang-course/gateway/internal/config"
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
