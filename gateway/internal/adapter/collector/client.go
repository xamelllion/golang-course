package collector

import (
	"context"

	"github.com/xamelllion/golang-course/gateway/internal/domain"
	"google.golang.org/grpc"

	pb "github.com/xamelllion/golang-course/proto"
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
		return domain.Repository{}, error
	}

	return domain.Repository{
		Name:        repository.Name,
		Description: repository.Description,
		Stars:       repository.Stars,
		Forks:       repository.Forks,
		CreateDate: repository.CreateDate.AsTime(),
	}, nil
}
