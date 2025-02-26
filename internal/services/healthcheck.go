package services

import (
	"ewallet-wallet/internal/interfaces/i_repository"
)

type Healthcheck struct {
	HealthcheckRepository i_repository.IHealthcheckRepo
}

func (s *Healthcheck) HealthcheckServices() (string, error) {
	return "service healty", nil
}
