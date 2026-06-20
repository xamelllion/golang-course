package controllerKafka

import (
	"context"
	"encoding/json"
	"log"

	apperror "github.com/xamelllion/golang-course/internal/errors"
	kafkaClient "github.com/xamelllion/golang-course/internal/kafka"
	"github.com/xamelllion/golang-course/processor/internal/domain"
)

type Tasker interface {
	ProcessTaskReponse(domain.TaskResponse) error
}

type KafkaController struct {
	tasker   Tasker
	consumer kafkaClient.KafkaConsumer
}

func NewKafkaController(tasker Tasker, consumer kafkaClient.KafkaConsumer) KafkaController {
	return KafkaController{
		tasker:   tasker,
		consumer: consumer,
	}
}

func (kc KafkaController) Run(ctx context.Context) error {
	return kc.consumer.Run(ctx, func(_ context.Context, data []byte) error {
		var task kafkaClient.RepoTaskResponse

		err := json.Unmarshal(data, &task)
		if err != nil {
			return apperror.Wrap("unmarshal kafka task", err)
		}

		log.Printf("get task for proceed: %v/%v", task.Owner, task.Repo)

		switch task.Status {
		case kafkaClient.STATUS_OK:
			return kc.tasker.ProcessTaskReponse(domain.TaskResponse{
				Owner:       task.Owner,
				Repo:        task.Repo,
				Description: task.Description,
				Stars:       int32(task.Stars),
				Forks:       int32(task.Forks),
				CreateDate:  task.CreatedAt,
			})
		case kafkaClient.STATUS_NOT_FOUND:
			return apperror.New(kafkaClient.STATUS_NOT_FOUND, "repo not found")
		case kafkaClient.STATUS_INVALID_ARGUMENT:
			return apperror.New(apperror.CodeInvalidArgument, "invalid repo owner or name")
		case kafkaClient.STATUS_UNAVAILABLE:
			return apperror.New(apperror.CodeUnavailable, "service unavaliable")
		default:
			log.Printf("kafka adapter internal error %v", task)
			return apperror.New(apperror.CodeInternal, "internal error")
		}
	})
}
