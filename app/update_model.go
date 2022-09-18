package app

import (
	"github.com/opensourceways/xihe-server/domain"
)

type ModelUpdateCmd struct {
	Name     domain.ModelName
	Desc     domain.ResourceDesc
	RepoType domain.RepoType
}

func (cmd *ModelUpdateCmd) toModel(p *domain.ModelModifiableProperty) {
	if cmd.Name != nil {
		p.Name = cmd.Name
	}

	if cmd.Desc != nil {
		p.Desc = cmd.Desc
	}

	if cmd.RepoType != nil {
		p.RepoType = cmd.RepoType
	}
}

func (s modelService) Update(p *domain.Model, cmd *ModelUpdateCmd) (dto ModelDTO, err error) {
	cmd.toModel(&p.ModelModifiableProperty)

	v, err := s.repo.Save(p)
	if err != nil {
		return
	}

	// TODO update repo visibility

	s.toModelDTO(&v, &dto)

	return
}

func (s modelService) AddLike(owner domain.Account, rid string) error {
	return s.repo.AddLike(owner, rid)
}

func (s modelService) RemoveLike(owner domain.Account, rid string) error {
	return s.repo.RemoveLike(owner, rid)
}
