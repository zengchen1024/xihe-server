package repositories

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

func NewProjectRepository(mapper ProjectMapper) repository.Project {
	return project{mapper}
}

type project struct {
	mapper ProjectMapper
}

func (impl project) Save(p *domain.Project) (r domain.Project, err error) {
	if p.Id != "" {
		return
	}

	do := ProjectDO{
		Owner:    p.Owner,
		Name:     p.Name.ProjName(),
		Desc:     p.Desc.ProjDesc(),
		Type:     p.Type.ProjType(),
		CoverId:  p.CoverId.CoverId(),
		RepoType: p.RepoType.RepoType(),
		Protocol: p.Protocol.ProtocolName(),
		Training: p.Training.TrainingPlatform(),
	}

	v, err := impl.mapper.Insert(do)
	if err == nil {
		r = *p
		r.Id = v
	}

	err = convertError(err)

	return
}

func (impl project) Get(pid string) (r domain.Project, err error) {
	return
}

type ProjectDO struct {
	Id       string
	Owner    string
	Name     string
	Desc     string
	Type     string
	CoverId  string
	RepoType string
	Protocol string
	Training string
	Tags     []string
	Version  int
}

type ProjectMapper interface {
	Insert(ProjectDO) (string, error)
	Update(string, ProjectDO) error
	Get(string) (ProjectDO, error)
}
