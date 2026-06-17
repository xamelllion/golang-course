package collector

import (
	"context"
	"log/slog"

	"github.com/xamelllion/golang-course/gateway/internal/domain"
	apperror "github.com/xamelllion/golang-course/internal/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pbProcessor "github.com/xamelllion/golang-course/proto/processor"
)

type Collector struct {
	conn *grpc.ClientConn
	log  *slog.Logger
}

func NewCollector(conn *grpc.ClientConn, log *slog.Logger) Collector {
	return Collector{
		conn: conn,
		log:  log,
	}
}

func (c Collector) GetRepository(owner, repo string) (domain.Repository, error) {
	client := pbProcessor.NewProcessorServiceClient(c.conn)

	c.log.Debug("adapter: get repositon of %s/%s", owner, repo)
	repository, error := client.GetRepository(context.Background(), &pbProcessor.GetRepositoryRequest{Owner: owner, Repo: repo})
	if error != nil {
		c.log.Debug("adapter: error %v", error)
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
