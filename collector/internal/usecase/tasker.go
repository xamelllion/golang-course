package usecase

import (
	"log"

	"github.com/xamelllion/golang-course/collector/internal/domain"
	apperror "github.com/xamelllion/golang-course/internal/errors"
)

type GithubAdapter interface {
	GetRepository(owner, repo string) (domain.GithubRepository, error)
}

type KafkaAdapter interface {
	SendTaskResponse(task domain.Task, repo domain.GithubRepository) error
	SendTaskResponseError(task domain.Task, status string, err error) error
}

type TaskerUsecase struct {
	gh GithubAdapter
	ka KafkaAdapter
}

func NewTaskerUsecase(gh GithubAdapter, ka KafkaAdapter) TaskerUsecase {
	return TaskerUsecase{gh: gh, ka: ka}
}

func (tu TaskerUsecase) ProcessTask(task domain.Task) error {
	repo, err := tu.gh.GetRepository(task.Owner, task.Repo)
	if err != nil {
		log.Printf("process task error %v", err)
		switch apperror.CodeOf(err) {
		case apperror.CodeNotFound:
			return tu.ka.SendTaskResponseError(task, "not_found", err)
		case apperror.CodeInvalidArgument:
			return tu.ka.SendTaskResponseError(task, "invalid_argument", err)
		case apperror.CodeUnavailable:
			return tu.ka.SendTaskResponseError(task, "unavailable", err)
		default:
			return tu.ka.SendTaskResponseError(task, "internal", err)
		}
	}

	log.Printf("get repo for kafka task: %v", repo)

	return tu.ka.SendTaskResponse(task, repo)
}
