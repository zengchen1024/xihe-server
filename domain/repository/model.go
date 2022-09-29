package repository

import (
	"github.com/opensourceways/xihe-server/domain"
)

type ModelPropertyUpdateInfo struct {
	ResourceToUpdate

	Property domain.ModelModifiableProperty
}

type UserModelsInfo struct {
	Models []domain.ModelSummary
	Total  int
}

type Model interface {
	Save(*domain.Model) (domain.Model, error)
	Get(domain.Account, string) (domain.Model, error)
	GetByName(domain.Account, domain.ModelName) (domain.Model, error)

	FindUserModels([]UserResourceListOption) ([]domain.ModelSummary, error)

	List(domain.Account, *ResourceListOption) (UserModelsInfo, error)
	ListAndSortByUpdateTime(domain.Account, *ResourceListOption) (UserModelsInfo, error)
	ListAndSortByFirstLetter(domain.Account, *ResourceListOption) (UserModelsInfo, error)
	ListAndSortByDownloadCount(domain.Account, *ResourceListOption) (UserModelsInfo, error)

	AddLike(domain.Account, string) error
	RemoveLike(domain.Account, string) error

	AddRelatedDataset(*RelatedResourceInfo) error
	RemoveRelatedDataset(*RelatedResourceInfo) error

	UpdateProperty(*ModelPropertyUpdateInfo) error
}
