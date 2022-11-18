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
	Delete(*domain.ResourceIndex) error
	Get(domain.Account, string) (domain.Dataset, error)
	GetByName(domain.Account, domain.ResourceName) (domain.Dataset, error)
	GetSummaryByName(domain.Account, domain.ResourceName) (domain.ResourceSummary, error)

	FindUserDatasets([]UserResourceListOption) ([]domain.DatasetSummary, error)
	ListSummary([]ResourceSummaryListOption) ([]domain.ResourceSummary, error)

	ListAndSortByUpdateTime(domain.Account, *ResourceListOption) (UserDatasetsInfo, error)
	ListAndSortByFirstLetter(domain.Account, *ResourceListOption) (UserDatasetsInfo, error)
	ListAndSortByDownloadCount(domain.Account, *ResourceListOption) (UserDatasetsInfo, error)

	ListGlobalAndSortByUpdateTime(*GlobalResourceListOption) (UserDatasetsInfo, error)
	ListGlobalAndSortByFirstLetter(*GlobalResourceListOption) (UserDatasetsInfo, error)
	ListGlobalAndSortByDownloadCount(*GlobalResourceListOption) (UserDatasetsInfo, error)

	Search(*ResourceSearchOption) (ResourceSearchResult, error)

	AddLike(*domain.ResourceIndex) error
	RemoveLike(*domain.ResourceIndex) error

	AddRelatedProject(*domain.ReverselyRelatedResourceInfo) error
	RemoveRelatedProject(*domain.ReverselyRelatedResourceInfo) error

	AddRelatedModel(*domain.ReverselyRelatedResourceInfo) error
	RemoveRelatedModel(*domain.ReverselyRelatedResourceInfo) error

	UpdateProperty(*DatasetPropertyUpdateInfo) error
}
