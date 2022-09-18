package app

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/platform"
)

type ProjectUpdateCmd struct {
	Name     domain.ProjName
	Desc     domain.ResourceDesc
	RepoType domain.RepoType
	CoverId  domain.CoverId
}

func (cmd *ProjectUpdateCmd) toProject(
	p *domain.ProjectModifiableProperty, repo *platform.RepoOption,
) (b bool) {
	f := func() {
		if !b {
			b = true
		}
	}

	if cmd.Name != nil && p.Name.ProjName() != cmd.Name.ProjName() {
		p.Name = cmd.Name
		repo.Name = cmd.Name
		f()
	}

	if cmd.Desc != nil && p.Desc.ResourceDesc() != cmd.Desc.ResourceDesc() {
		p.Desc = cmd.Desc
		f()
	}

	if cmd.RepoType != nil && p.RepoType.RepoType() != cmd.RepoType.RepoType() {
		p.RepoType = cmd.RepoType
		repo.RepoType = cmd.RepoType
		f()
	}

	if cmd.CoverId != nil && p.CoverId.CoverId() == cmd.CoverId.CoverId() {
		p.CoverId = cmd.CoverId
		f()
	}

	return
}

func (s projectService) Update(p *domain.Project, cmd *ProjectUpdateCmd) (dto ProjectDTO, err error) {
	opt := new(platform.RepoOption)
	if !cmd.toProject(&p.ProjectModifiableProperty, opt) {
		s.toProjectDTO(p, &dto)

		return
	}

	v, err := s.repo.Save(p)
	if err != nil {
		return
	}

	if opt.IsNotEmpty() {
		if err = s.pr.Update(p.RepoId, opt); err != nil {
			return
		}
	}

	s.toProjectDTO(&v, &dto)

	return
}

func (s projectService) AddLike(owner domain.Account, rid string) error {
	return s.repo.AddLike(owner, rid)
}

func (s projectService) RemoveLike(owner domain.Account, rid string) error {
	return s.repo.RemoveLike(owner, rid)
}
