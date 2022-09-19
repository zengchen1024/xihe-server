package repository

import (
	"github.com/opensourceways/xihe-server/domain"
)

type UserResourceListOption struct {
	Owner domain.Account
	Ids   []string
}

type ResourceListOption struct {
	Name     string
	RepoType domain.RepoType
}

type RelatedResourceInfo struct {
	Owner      domain.Account
	ResourceId string
	Version    int

	domain.ResourceObj
}

type Project interface {
	Save(*domain.Project) (domain.Project, error)
	Get(domain.Account, string) (domain.Project, error)
	GetByName(domain.Account, domain.ProjName) (domain.Project, error)
	List(domain.Account, ResourceListOption) ([]domain.Project, error)
	FindUserProjects([]UserResourceListOption) ([]domain.Project, error)

	AddLike(domain.Account, string) error
	RemoveLike(domain.Account, string) error

	AddRelatedResource(*RelatedResourceInfo) error
	RemoveRelatedResource(*RelatedResourceInfo) error
}
