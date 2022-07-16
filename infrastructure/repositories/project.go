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

func (impl project) Get(owner, identity string) (r domain.Project, err error) {
	v, err := impl.mapper.Get(owner, identity)
	if err != nil {
		err = convertError(err)

		return
	}

	err = v.toProject(&r)

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

func (do *ProjectDO) toProject(r *domain.Project) (err error) {
	r.Id = do.Id
	r.Owner = do.Owner

	if r.Name, err = domain.NewProjName(do.Name); err != nil {
		return
	}

	if r.Desc, err = domain.NewProjDesc(do.Desc); err != nil {
		return
	}

	if r.Type, err = domain.NewProjType(do.Type); err != nil {
		return
	}

	if r.CoverId, err = domain.NewConverId(do.CoverId); err != nil {
		return
	}

	if r.RepoType, err = domain.NewRepoType(do.RepoType); err != nil {
		return
	}

	if r.Protocol, err = domain.NewProtocolName(do.Protocol); err != nil {
		return
	}

	if r.Training, err = domain.NewTrainingPlatform(do.Training); err != nil {
		return
	}

	r.Tags = do.Tags

	r.Version = do.Version

	return
}

type ProjectMapper interface {
	Insert(ProjectDO) (string, error)
	Update(string, ProjectDO) error
	Get(string, string) (ProjectDO, error)
}
