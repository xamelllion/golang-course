package controller

import (
	"context"

	pbSubscriber "github.com/xamelllion/golang-course/proto/subscriber"
)

type PingUsecase interface {
	Ping() (string, error)
}

type Server struct {
	pbSubscriber.UnimplementedSubscriberServiceServer
	PingUsecase PingUsecase
}

func NewServer(usecase PingUsecase) *Server {
	return &Server{PingUsecase: usecase}
}

func (s *Server) Ping(_ context.Context, req *pbSubscriber.PingRequest) (*pbSubscriber.PingResponse, error) {
	pong, err := s.PingUsecase.Ping()

	return &pbSubscriber.PingResponse{Reply: pong}, err
}
