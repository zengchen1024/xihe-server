package repository

import (
	"github.com/opensourceways/xihe-server/domain"
)

type ModelListOption struct {
	Name domain.ProjName
}

type Model interface {
	Save(*domain.Model) (domain.Model, error)
	Get(string, string) (domain.Model, error)
	List(string, ModelListOption) ([]domain.Model, error)
}
