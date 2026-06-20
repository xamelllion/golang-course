package adapter

import (
	"context"
	"log"

	"github.com/xamelllion/golang-course/collector/internal/domain"
	apperror "github.com/xamelllion/golang-course/internal/errors"
	kafkaClient "github.com/xamelllion/golang-course/internal/kafka"
)

type KafkaAdapter struct {
	producer kafkaClient.KafkaProducer
	consumer kafkaClient.KafkaConsumer
}

func NewKafkaAdapter(producer kafkaClient.KafkaProducer, consumer kafkaClient.KafkaConsumer) KafkaAdapter {
	return KafkaAdapter{producer: producer, consumer: consumer}
}

func (ka KafkaAdapter) SendTaskRequest(sub domain.Subscription) error {
	req := kafkaClient.RepoTaskRequest{
		Owner: sub.Owner,
		Repo:  sub.Repo,
	}

	err := ka.producer.Send(context.Background(), "repo_request", req)
	if err != nil {
		return apperror.WrapCode(apperror.CodeInternal, "send kafka task requests", err)
	}
	log.Printf("sended task request %v", req)

	return nil
}

func (ka KafkaAdapter) SendTaskResponse(task domain.Task, repo domain.GithubRepository) error {
	res := kafkaClient.RepoTaskResponse{
		Owner:       task.Owner,
		Repo:        task.Repo,
		Description: repo.Description,
		Stars:       int(repo.Stars),
		Forks:       int(repo.Forks),
		CreatedAt:   repo.Create_date,
		Status:      kafkaClient.STATUS_OK,
		Error:       "",
	}

	log.Printf("start send kafka response %v", res)

	err := ka.producer.Send(context.Background(), "repo_response", res)
	if err != nil {
		return apperror.WrapCode(apperror.CodeInternal, "send kafka task response", err)
	}

	log.Printf("send kafka response %v", res)

	return nil
}

func (ka KafkaAdapter) SendTaskResponseError(task domain.Task, status string, err error) error {
	res := kafkaClient.RepoTaskResponse{
		Owner:  task.Owner,
		Repo:   task.Repo,
		Status: status,
		Error:  err.Error(),
	}

	if err := ka.producer.Send(context.Background(), "repo_response", res); err != nil {
		return apperror.WrapCode(apperror.CodeInternal, "send kafka task response", err)
	}

	return nil
}
