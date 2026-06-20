package adapter

import (
	"context"
	"encoding/json"
	"log"

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

func (ka KafkaAdapter) SendTaskRequest(owner, repo string) error {
	req := kafkaClient.RepoTaskRequest{
		Owner: owner,
		Repo:  repo,
	}

	err := ka.producer.Send(context.Background(), "repo_request", req)
	if err != nil {
		return apperror.WrapCode(apperror.CodeInternal, "send kafka task requests", err)
	}
	log.Printf("sended task request %v", req)

	return nil
}

func (ka KafkaAdapter) RunGetTaskResponse(ctx context.Context, handle func(ctx context.Context, task kafkaClient.RepoTaskResponse) error) error {
	return ka.consumer.Run(ctx, func(_ context.Context, data []byte) error {
		var task kafkaClient.RepoTaskResponse

		err := json.Unmarshal(data, &task)
		if err != nil {
			return apperror.Wrap("unmarshal kafka task", err)
		}

		log.Printf("get task for proceed: %v", task)

		return handle(ctx, task)
	})
}
