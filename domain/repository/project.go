package repository

import (
	"github.com/opensourceways/xihe-server/domain"
)

type ProjectListOption struct {
	Name domain.ProjName
}

type Project interface {
	Save(*domain.Project) (domain.Project, error)
	Get(string, string) (domain.Project, error)
	List(string, ProjectListOption) ([]domain.Project, error)
}
