package repository

import (
	"github.com/opensourceways/xihe-server/domain"
)

type ModelListOption struct {
	Name domain.ProjName
}

type Model interface {
	Save(*domain.Model) (domain.Model, error)
	Get(domain.Account, string) (domain.Model, error)
	List(domain.Account, ModelListOption) ([]domain.Model, error)
}
