package app

import (
	"github.com/opensourceways/xihe-server/domain"
)

type ProjectUpdateCmd struct {
	Name     domain.ProjName
	Desc     domain.ProjDesc
	RepoType domain.RepoType
	CoverId  domain.CoverId
}

func (cmd *ProjectUpdateCmd) toProject(p *domain.ProjectModifiableProperty) {
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
}

func (s projectService) Update(p *domain.Project, cmd *ProjectUpdateCmd) (dto ProjectDTO, err error) {
	cmd.toProject(&p.ProjectModifiableProperty)

	v, err := s.repo.Save(p)
	if err != nil {
		return
	}

	// TODO update repo visibility

	s.toProjectDTO(&v, &dto)

	return
}

func (s projectService) AddLike(owner domain.Account, rid string) error {
	return s.repo.AddLike(owner, rid)
}

func (s projectService) RemoveLike(owner domain.Account, rid string) error {
	return s.repo.RemoveLike(owner, rid)
}
