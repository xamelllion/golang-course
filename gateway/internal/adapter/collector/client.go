package collector

import (
	"context"
	"log/slog"

	"github.com/xamelllion/golang-course/gateway/internal/config"
	apperror "github.com/xamelllion/golang-course/internal/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pbCollector "github.com/xamelllion/golang-course/proto/collector"
)

type Collector struct {
	conn   *grpc.ClientConn
	client pbCollector.CollectorServiceClient
	log    *slog.Logger
}

func NewCollector(cfg config.Config, log *slog.Logger) Collector {
	conn, error := grpc.NewClient(cfg.CollectorAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if error != nil {
		panic(error)
	}

	client := pbCollector.NewCollectorServiceClient(conn)

	return Collector{
		conn:   conn,
		client: client,
		log:    log,
	}
}

func (c Collector) Ping() (string, error) {
	pong, err := c.client.Ping(context.Background(), &pbCollector.PingRequest{})
	if err != nil {
		return "", apperror.New(apperror.CodeUnavailable, "processor unvailable")
	}

	return pong.Reply, nil
}

func (c Collector) Close() {
	c.conn.Close()
}
