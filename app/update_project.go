package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/platform"
	"github.com/opensourceways/xihe-server/domain/repository"
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

func (s projectService) Update(p *domain.Project, cmd *ProjectUpdateCmd, pr platform.Repository) (dto ProjectDTO, err error) {
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
		if err = pr.Update(p.RepoId, opt); err != nil {
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

type RelatedResourceModifyCmd domain.ResourceObj

func (cmd *RelatedResourceModifyCmd) ValidateForProject() error {
	if err := cmd.validate(); err != nil {
		return err
	}

	t := cmd.ResourceType.ResourceType()
	if t != domain.ResourceDataset && t != domain.ResourceModel {
		return errors.New("unspported resource type")
	}

	return nil
}

func (cmd *RelatedResourceModifyCmd) validate() error {
	b := cmd.ResourceOwner == nil ||
		cmd.ResourceId == "" ||
		cmd.ResourceType == nil
	if b {
		return errors.New("invalid related resource modify cmd")
	}

	return nil
}

func (s projectService) AddRelatedResource(
	p *domain.Project, cmd *RelatedResourceModifyCmd,
) error {
	// TODO limited num of related resources

	var v domain.RelatedResources
	switch cmd.ResourceType.ResourceType() {
	case domain.ResourceModel:
		v = p.RelatedModels
	case domain.ResourceDataset:
		v = p.RelatedDatasets
	}

	if v.Has(cmd.ResourceOwner, cmd.ResourceId) {
		return nil
	}

	info := repository.RelatedResourceInfo{
		Owner:       p.Owner,
		ResourceId:  p.Id,
		Version:     p.Version,
		ResourceObj: *(*domain.ResourceObj)(cmd),
	}

	return s.repo.AddRelatedResource(&info)
}

func (s projectService) RemoveRelatedResource(
	p *domain.Project, cmd *RelatedResourceModifyCmd,
) error {
	var v domain.RelatedResources
	switch cmd.ResourceType.ResourceType() {
	case domain.ResourceModel:
		v = p.RelatedModels
	case domain.ResourceDataset:
		v = p.RelatedDatasets
	}

	if !v.Has(cmd.ResourceOwner, cmd.ResourceId) {
		return nil
	}

	info := repository.RelatedResourceInfo{
		Owner:       p.Owner,
		ResourceId:  p.Id,
		Version:     p.Version,
		ResourceObj: *(*domain.ResourceObj)(cmd),
	}

	return s.repo.RemoveRelatedResource(&info)
}
