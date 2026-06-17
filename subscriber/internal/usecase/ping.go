package usecase

type PingUsecase struct {
}

func NewPingUsecase() PingUsecase {
	return PingUsecase{}
}

func (r PingUsecase) Ping() (string, error) {
	return "pong", nil
}
