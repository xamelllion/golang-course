package processor

import (
	"context"
	"log/slog"

	"github.com/xamelllion/golang-course/gateway/internal/config"
	"github.com/xamelllion/golang-course/gateway/internal/domain"
	apperror "github.com/xamelllion/golang-course/internal/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	pbProcessor "github.com/xamelllion/golang-course/proto/processor"
)

type Processor struct {
	conn   *grpc.ClientConn
	client pbProcessor.ProcessorServiceClient
	log    *slog.Logger
}

func NewProcessor(cfg config.Config, log *slog.Logger) Processor {
	conn, error := grpc.NewClient(cfg.ProcessorAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if error != nil {
		panic(error)
	}

	client := pbProcessor.NewProcessorServiceClient(conn)

	return Processor{
		conn:   conn,
		client: client,
		log:    log,
	}
}

func (p Processor) GetRepository(owner, repo string) (domain.Repository, error) {
	p.log.Debug("adapter: get repositon of %s/%s", owner, repo)
	repository, error := p.client.GetRepository(context.Background(), &pbProcessor.GetRepositoryRequest{Owner: owner, Repo: repo})
	if error != nil {
		switch status.Code(error) {
		case codes.NotFound:
			return domain.Repository{}, apperror.New(apperror.CodeNotFound, error.Error())
		case codes.InvalidArgument:
			return domain.Repository{}, apperror.New(apperror.CodeInvalidArgument, error.Error())
		case codes.Unavailable:
			return domain.Repository{}, apperror.New(apperror.CodeUnavailable, error.Error())
		default:
			return domain.Repository{}, apperror.New(apperror.CodeInternal, error.Error())
		}
	}

	return domain.Repository{
		Name:        repository.Name,
		Description: repository.Description,
		Stars:       repository.Stars,
		Forks:       repository.Forks,
		CreateDate:  repository.CreateDate.AsTime(),
	}, nil
}

func (p Processor) Ping() (string, error) {
	pong, err := p.client.Ping(context.Background(), &pbProcessor.PingRequest{})
	if err != nil {
		return "", apperror.New(apperror.CodeUnavailable, "processor unvailable")
	}

	return pong.Reply, nil
}

func (p Processor) Close() {
	p.conn.Close()
}
