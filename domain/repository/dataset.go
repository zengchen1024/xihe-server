package repository

import (
	"github.com/opensourceways/xihe-server/domain"
)

type DatasetPropertyUpdateInfo struct {
	ResourceToUpdate

	Property domain.DatasetModifiableProperty
}

type Dataset interface {
	Save(*domain.Dataset) (domain.Dataset, error)
	Get(domain.Account, string) (domain.Dataset, error)
	GetByName(domain.Account, domain.DatasetName) (domain.Dataset, error)
	List(domain.Account, ResourceListOption) ([]domain.Dataset, error)
	FindUserDatasets([]UserResourceListOption) ([]domain.Dataset, error)

	AddLike(domain.Account, string) error
	RemoveLike(domain.Account, string) error

	UpdateProperty(*DatasetPropertyUpdateInfo) error
}
