package repository

import (
	"github.com/opensourceways/xihe-server/domain"
)

type ModelSummaryListOption struct {
	Owner domain.Account
	Name  domain.ModelName
}

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
	GetSummaryByName(domain.Account, domain.ResourceName) (domain.ResourceSummary, error)

	FindUserModels([]UserResourceListOption) ([]domain.ModelSummary, error)
	ListSummary([]ModelSummaryListOption) ([]domain.ResourceSummary, error)

	List(domain.Account, *ResourceListOption) (UserModelsInfo, error)
	ListAndSortByUpdateTime(domain.Account, *ResourceListOption) (UserModelsInfo, error)
	ListAndSortByFirstLetter(domain.Account, *ResourceListOption) (UserModelsInfo, error)
	ListAndSortByDownloadCount(domain.Account, *ResourceListOption) (UserModelsInfo, error)

	AddLike(*domain.ResourceIndex) error
	RemoveLike(*domain.ResourceIndex) error

	AddRelatedDataset(*RelatedResourceInfo) error
	RemoveRelatedDataset(*RelatedResourceInfo) error

	AddRelatedProject(*domain.ReverselyRelatedResourceInfo) error
	RemoveRelatedProject(*domain.ReverselyRelatedResourceInfo) error

	UpdateProperty(*ModelPropertyUpdateInfo) error
}
