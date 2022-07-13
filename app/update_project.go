package app

import (
	"github.com/opensourceways/xihe-server/domain"
)

type ProjectUpdateCmd struct {
	Name    domain.ProjName
	Desc    domain.ProjDesc
	Type    domain.RepoType
	CoverId domain.CoverId
}

func (cmd *ProjectUpdateCmd) toProject(p *domain.Project) {
	if cmd.Name != nil {
		p.Name = cmd.Name
	}

	if cmd.Desc != nil {
		p.Desc = cmd.Desc
	}

	if cmd.Type != nil {
		p.Type = cmd.Type
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
