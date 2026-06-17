package usecase

import (
	"log/slog"

	"github.com/xamelllion/golang-course/gateway/internal/domain"
	apperror "github.com/xamelllion/golang-course/internal/errors"
)

type Pinger interface {
	Ping() (string, error)
}

type PingUsecase struct {
	pingers map[string]Pinger
	log     *slog.Logger
}

func NewPingUsecase(pingers map[string]Pinger, log *slog.Logger) *PingUsecase {
	return &PingUsecase{pingers: pingers, log: log}
}

func (p *PingUsecase) PingAll() (domain.ServicesInfo, error) {
	result := domain.ServicesInfo{Status: domain.ServicesStatusOk}

	for key, pinger := range p.pingers {
		_, err := pinger.Ping()

		newServiceStatus := domain.ServiceStatus{Name: key, Status: domain.PingStatusUp}
		if err != nil {
			if apperror.CodeOf(err) != apperror.CodeUnavailable {
				return domain.ServicesInfo{}, err
			}

			newServiceStatus.Status = domain.PingStatusDown
			result.Status = domain.ServicesStatusDegraded
		}

		result.Services = append(result.Services, newServiceStatus)
	}

	return result, nil
}
