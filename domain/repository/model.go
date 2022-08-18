package repository

import (
	"github.com/opensourceways/xihe-server/domain"
)

type ModelListOption struct {
	Name     domain.ModelName
	RepoType domain.RepoType
}

type Model interface {
	Save(*domain.Model) (domain.Model, error)
	Get(domain.Account, string) (domain.Model, error)
	GetByName(domain.Account, domain.ModelName) (domain.Model, error)
	List(domain.Account, ModelListOption) ([]domain.Model, error)
	FindUserModels([]UserResourceListOption) ([]domain.Model, error)

	AddLike(domain.Account, string) error
	RemoveLike(domain.Account, string) error
}
