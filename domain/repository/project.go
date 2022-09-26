package repository

import (
	"github.com/opensourceways/xihe-server/domain"
)

type UserResourceListOption struct {
	Owner domain.Account
	Ids   []string
}

type ResourceListOption struct {
	// can't define Name as domain.ResourceName
	// because the Name can be subpart of the real resource name
	Name     string
	RepoType domain.RepoType

	CountPerPage int
	PageNum      int
}

type RelatedResourceInfo struct {
	ResourceToUpdate

	RelatedResource domain.ResourceIndex
}

type ProjectPropertyUpdateInfo struct {
	ResourceToUpdate

	Property domain.ProjectModifiableProperty
}

type ResourceToUpdate struct {
	Owner     domain.Account
	Id        string
	Version   int
	UpdatedAt int64
}

type Project interface {
	Save(*domain.Project) (domain.Project, error)
	Get(domain.Account, string) (domain.Project, error)
	GetByName(domain.Account, domain.ProjName) (domain.Project, error)
	FindUserProjects([]UserResourceListOption) ([]domain.Project, error)

	List(domain.Account, *ResourceListOption) ([]domain.Project, error)
	ListAndSortByUpdateTime(domain.Account, *ResourceListOption) ([]domain.Project, error)
	ListAndSortByFirtLetter(domain.Account, *ResourceListOption) ([]domain.Project, error)
	ListAndSortByDownloadCount(domain.Account, *ResourceListOption) ([]domain.Project, error)

	AddLike(domain.Account, string) error
	RemoveLike(domain.Account, string) error

	AddRelatedModel(*RelatedResourceInfo) error
	RemoveRelatedModel(*RelatedResourceInfo) error

	AddRelatedDataset(*RelatedResourceInfo) error
	RemoveRelatedDataset(*RelatedResourceInfo) error

	UpdateProperty(*ProjectPropertyUpdateInfo) error

	IncreaseFork(*domain.ResourceIndex) error
}
