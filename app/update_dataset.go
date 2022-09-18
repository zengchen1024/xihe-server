package app

import (
	"github.com/opensourceways/xihe-server/domain"
)

type DatasetUpdateCmd struct {
	Name     domain.DatasetName
	Desc     domain.ProjDesc
	RepoType domain.RepoType
}

func (cmd *DatasetUpdateCmd) toDataset(p *domain.DatasetModifiableProperty) {
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

func (s datasetService) Update(p *domain.Dataset, cmd *DatasetUpdateCmd) (dto DatasetDTO, err error) {
	cmd.toDataset(&p.DatasetModifiableProperty)

	v, err := s.repo.Save(p)
	if err != nil {
		return
	}

	// TODO update repo visibility

	s.toDatasetDTO(&v, &dto)

	return
}

func (s datasetService) AddLike(owner domain.Account, rid string) error {
	return s.repo.AddLike(owner, rid)
}

func (s datasetService) RemoveLike(owner domain.Account, rid string) error {
	return s.repo.RemoveLike(owner, rid)
}
