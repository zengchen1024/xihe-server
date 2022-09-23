package repository

import (
	"github.com/opensourceways/xihe-server/domain"
)

type ModelPropertyUpdateInfo struct {
	Owner    domain.Account
	Id       string
	Version  int
	Property domain.ModelModifiableProperty
}

type Model interface {
	Save(*domain.Model) (domain.Model, error)
	Get(domain.Account, string) (domain.Model, error)
	GetByName(domain.Account, domain.ModelName) (domain.Model, error)
	List(domain.Account, ResourceListOption) ([]domain.Model, error)
	FindUserModels([]UserResourceListOption) ([]domain.Model, error)

	AddLike(domain.Account, string) error
	RemoveLike(domain.Account, string) error

	AddRelatedDataset(*RelatedResourceInfo) error
	RemoveRelatedDataset(*RelatedResourceInfo) error

	UpdateProperty(*ModelPropertyUpdateInfo) error
}
