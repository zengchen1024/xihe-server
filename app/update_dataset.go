package app

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/platform"
)

type DatasetUpdateCmd struct {
	Name     domain.DatasetName
	Desc     domain.ResourceDesc
	RepoType domain.RepoType
}

func (cmd *DatasetUpdateCmd) toDataset(
	p *domain.DatasetModifiableProperty, repo *platform.RepoOption,
) (b bool) {
	f := func() {
		if !b {
			b = true
		}
	}

	if cmd.Name != nil && p.Name.DatasetName() != cmd.Name.DatasetName() {
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

	return
}

func (s datasetService) Update(
	p *domain.Dataset, cmd *DatasetUpdateCmd, pr platform.Repository,
) (dto DatasetDTO, err error) {
	opt := new(platform.RepoOption)
	if !cmd.toDataset(&p.DatasetModifiableProperty, opt) {
		s.toDatasetDTO(p, &dto)

		return

	}

	v, err := s.repo.Save(p)
	if err != nil {
		return
	}

	if opt.IsNotEmpty() {
		if err = pr.Update(p.RepoId, opt); err != nil {
			return
		}
	}

	s.toDatasetDTO(&v, &dto)

	return
}

func (s datasetService) AddLike(owner domain.Account, rid string) error {
	return s.repo.AddLike(owner, rid)
}

func (s datasetService) RemoveLike(owner domain.Account, rid string) error {
	return s.repo.RemoveLike(owner, rid)
}
