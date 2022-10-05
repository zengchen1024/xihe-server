package repository

import (
	"github.com/opensourceways/xihe-server/domain"
)

type DatasetPropertyUpdateInfo struct {
	ResourceToUpdate

	Property domain.DatasetModifiableProperty
}

type UserDatasetsInfo struct {
	Datasets []domain.DatasetSummary
	Total    int
}

type Dataset interface {
	Save(*domain.Dataset) (domain.Dataset, error)
	Get(domain.Account, string) (domain.Dataset, error)
	GetByName(domain.Account, domain.DatasetName) (domain.Dataset, error)

	FindUserDatasets([]UserResourceListOption) ([]domain.DatasetSummary, error)

	List(domain.Account, *ResourceListOption) (UserDatasetsInfo, error)
	ListAndSortByUpdateTime(domain.Account, *ResourceListOption) (UserDatasetsInfo, error)
	ListAndSortByFirstLetter(domain.Account, *ResourceListOption) (UserDatasetsInfo, error)
	ListAndSortByDownloadCount(domain.Account, *ResourceListOption) (UserDatasetsInfo, error)

	AddLike(*domain.ResourceIndex) error
	RemoveLike(*domain.ResourceIndex) error

	AddRelatedProject(*domain.ReverselyRelatedResourceInfo) error
	RemoveRelatedProject(*domain.ReverselyRelatedResourceInfo) error

	AddRelatedModel(*domain.ReverselyRelatedResourceInfo) error
	RemoveRelatedModel(*domain.ReverselyRelatedResourceInfo) error

	UpdateProperty(*DatasetPropertyUpdateInfo) error
}
