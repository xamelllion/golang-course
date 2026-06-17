package collector

import (
	"context"

	"github.com/xamelllion/golang-course/processor/internal/domain"
	apperror "github.com/xamelllion/golang-course/internal/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/xamelllion/golang-course/proto/collector"
)

type Collector struct {
	conn *grpc.ClientConn
}

func NewCollector(conn *grpc.ClientConn) Collector {
	return Collector{
		conn: conn,
	}
}

func (c Collector) GetRepository(owner, repo string) (domain.Repository, error) {
	client := pb.NewCollectorServiceClient(c.conn)

	repository, error := client.GetRepository(context.Background(), &pb.GetRepositoryRequest{Owner: owner, Repo: repo})
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
