package controller

import (
	"context"

	apperror "github.com/xamelllion/golang-course/internal/errors"
	pbSubscriber "github.com/xamelllion/golang-course/proto/subscriber"
	"github.com/xamelllion/golang-course/subscriber/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PingUsecase interface {
	Ping() (string, error)
}

type SubscribeUsecase interface {
	Subscribe(subi domain.Subscription) error
	Unsubscribe(subi domain.Subscription) error
	GetSubscriptions() ([]domain.Subscription, error)
}

type Server struct {
	pbSubscriber.UnimplementedSubscriberServiceServer
	PingUsecase      PingUsecase
	SubscribeUsecase SubscribeUsecase
}

func NewServer(pingUsecase PingUsecase, subscribeUsecase SubscribeUsecase) *Server {
	return &Server{PingUsecase: pingUsecase, SubscribeUsecase: subscribeUsecase}
}

func (s *Server) Ping(_ context.Context, req *pbSubscriber.PingRequest) (*pbSubscriber.PingResponse, error) {
	pong, err := s.PingUsecase.Ping()

	return &pbSubscriber.PingResponse{Reply: pong}, err
}

func (s *Server) Subscribe(_ context.Context, req *pbSubscriber.SubscribeRequest) (*pbSubscriber.SubscribeResponse, error) {
	if req.Owner == "" || req.Repo == "" {
		return &pbSubscriber.SubscribeResponse{}, status.Error(codes.InvalidArgument, "empty owner or repo")
	}

	subi := domain.Subscription{Owner: req.Owner, Repo: req.Repo}

	err := s.SubscribeUsecase.Subscribe(subi)
	if err != nil {
		return nil, status.Error(apperror.ToGRPCCode(err), err.Error())
	}

	return &pbSubscriber.SubscribeResponse{}, err
}

func (s *Server) Unsubscribe(_ context.Context, req *pbSubscriber.UnsubscribeRequest) (*pbSubscriber.UnsubscribeResponse, error) {
	if req.Owner == "" || req.Repo == "" {
		return &pbSubscriber.UnsubscribeResponse{}, apperror.New(apperror.CodeInvalidArgument, "empty owner or repo")
	}

	subi := domain.Subscription{Owner: req.Owner, Repo: req.Repo}

	err := s.SubscribeUsecase.Unsubscribe(subi)
	if err != nil {
		return nil, status.Error(apperror.ToGRPCCode(err), err.Error())
	}

	return &pbSubscriber.UnsubscribeResponse{}, err
}

func (s *Server) GetSubscriptions(_ context.Context, req *pbSubscriber.GetSubscriptionsRequest) (*pbSubscriber.GetSubscriptionsResponse, error) {
	subs, err := s.SubscribeUsecase.GetSubscriptions()
	if err != nil {
		return nil, status.Error(apperror.ToGRPCCode(err), err.Error())
	}

	result := make([](*pbSubscriber.Subscription), len(subs))
	for key, sub := range subs {
		result[key] = &pbSubscriber.Subscription{Owner: sub.Owner, Repo: sub.Repo}
	}

	return &pbSubscriber.GetSubscriptionsResponse{Subscriptions: result}, err
}
