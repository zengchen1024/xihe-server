package app

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type ProjectUpdateCmd struct {
	Name     domain.ProjName
	Desc     domain.ProjDesc
	RepoType domain.RepoType
	CoverId  domain.CoverId
}

func (cmd *ProjectUpdateCmd) toProject(p *domain.Project) {
	if cmd.Name != nil {
		p.Name = cmd.Name
	}

	if cmd.Desc != nil {
		p.Desc = cmd.Desc
	}

	if cmd.RepoType != nil {
		p.RepoType = cmd.RepoType
	}

	if cmd.CoverId != nil {
		p.CoverId = cmd.CoverId
	}

	// tags
}

func (s projectService) Update(p *domain.Project, cmd *ProjectUpdateCmd) (dto ProjectDTO, err error) {
	cmd.toProject(p)

	v, err := s.repo.Save(p)
	if err != nil {
		return
	}

	s.toProjectDTO(&v, &dto)

	return
}

type ProjectListCmd struct {
	Name domain.ProjName
}

func (cmd *ProjectListCmd) toProjectListOption() (
	option repository.ProjectListOption,
) {
	option.Name = cmd.Name

	return
}

func (s projectService) List(owner string, cmd *ProjectListCmd) (
	dtos []ProjectDTO, err error,
) {
	v, err := s.repo.List(owner, cmd.toProjectListOption())
	if err != nil || len(v) == 0 {
		return
	}

	dtos = make([]ProjectDTO, len(v))
	for i := range v {
		s.toProjectDTO(&v[i], &dtos[i])
	}

	return
}
