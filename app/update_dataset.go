package app

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/platform"
	"github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
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
	d *domain.Dataset, cmd *DatasetUpdateCmd, pr platform.Repository,
) (dto DatasetDTO, err error) {
	opt := new(platform.RepoOption)
	if !cmd.toDataset(&d.DatasetModifiableProperty, opt) {
		s.toDatasetDTO(d, &dto)

		return
	}

	if opt.IsNotEmpty() {
		if err = pr.Update(d.RepoId, opt); err != nil {
			return
		}
	}

	info := repository.DatasetPropertyUpdateInfo{
		ResourceToUpdate: s.toResourceToUpdate(d),
		Property:         d.DatasetModifiableProperty,
	}
	if err = s.repo.UpdateProperty(&info); err != nil {
		return
	}

	s.toDatasetDTO(d, &dto)

	return
}

func (s datasetService) SetTags(d *domain.Dataset, cmd *ResourceTagsUpdateCmd) error {
	tags, b := cmd.toTags(d.DatasetModifiableProperty.Tags)
	if !b {
		return nil
	}

	d.DatasetModifiableProperty.Tags = tags

	info := repository.DatasetPropertyUpdateInfo{
		ResourceToUpdate: s.toResourceToUpdate(d),
		Property:         d.DatasetModifiableProperty,
	}

	return s.repo.UpdateProperty(&info)
}

func (s datasetService) AddLike(owner domain.Account, rid string) error {
	return s.repo.AddLike(owner, rid)
}

func (s datasetService) RemoveLike(owner domain.Account, rid string) error {
	return s.repo.RemoveLike(owner, rid)
}

func (s datasetService) toResourceToUpdate(d *domain.Dataset) repository.ResourceToUpdate {
	return repository.ResourceToUpdate{
		Owner:     d.Owner,
		Id:        d.Id,
		Version:   d.Version,
		UpdatedAt: utils.Now(),
	}
}
