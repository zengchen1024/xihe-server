package repository

import (
	"github.com/opensourceways/xihe-server/domain"
)

type UserResourceListOption struct {
	Owner domain.Account
	Ids   []string
}

type ProjectListOption struct {
	Name     domain.ProjName
	RepoType domain.RepoType
}

type Project interface {
	Save(*domain.Project) (domain.Project, error)
	Get(domain.Account, string) (domain.Project, error)
	GetByName(domain.Account, domain.ProjName) (domain.Project, error)
	List(domain.Account, ProjectListOption) ([]domain.Project, error)
	FindUserProjects([]UserResourceListOption) ([]domain.Project, error)
}
