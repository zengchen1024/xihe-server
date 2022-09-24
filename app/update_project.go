package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/platform"
	"github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
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

// the step1 must be done before step2.
// For example, it can't set the project's name to the one existing.
// gitlab will help to avoid this case.
func (s projectService) Update(
	p *domain.Project, cmd *ProjectUpdateCmd, pr platform.Repository,
) (dto ProjectDTO, err error) {
	opt := new(platform.RepoOption)
	if !cmd.toProject(&p.ProjectModifiableProperty, opt) {
		s.toProjectDTO(p, &dto)

		return
	}

	// step1
	if opt.IsNotEmpty() {
		if err = pr.Update(p.RepoId, opt); err != nil {
			return
		}
	}

	// step2
	info := repository.ProjectPropertyUpdateInfo{
		ResourceToUpdate: s.toResourceToUpdate(p),
		Property:         p.ProjectModifiableProperty,
	}
	if err = s.repo.UpdateProperty(&info); err != nil {
		return
	}

	s.toProjectDTO(p, &dto)

	return
}

func (s projectService) SetTags(p *domain.Project, cmd *ResourceTagsUpdateCmd) error {
	tags, b := cmd.toTags(p.ProjectModifiableProperty.Tags)
	if !b {
		return nil
	}

	p.ProjectModifiableProperty.Tags = tags

	info := repository.ProjectPropertyUpdateInfo{
		ResourceToUpdate: s.toResourceToUpdate(p),
		Property:         p.ProjectModifiableProperty,
	}

	return s.repo.UpdateProperty(&info)
}

func (s projectService) IncreaseFork(index domain.ResourceIndex) error {
	return s.repo.IncreaseFork(&index)
}

func (s projectService) AddLike(owner domain.Account, rid string) error {
	return s.repo.AddLike(owner, rid)
}

func (s projectService) RemoveLike(owner domain.Account, rid string) error {
	return s.repo.RemoveLike(owner, rid)
}

func (s projectService) AddRelatedModel(
	p *domain.Project, index *domain.ResourceIndex,
) error {
	return s.addRelatedResource(
		p, index, p.RelatedModels, s.repo.AddRelatedModel,
	)
}

func (s projectService) AddRelatedDataset(
	p *domain.Project, index *domain.ResourceIndex,
) error {
	return s.addRelatedResource(
		p, index, p.RelatedDatasets, s.repo.AddRelatedDataset,
	)
}

func (s projectService) addRelatedResource(
	p *domain.Project, index *domain.ResourceIndex,
	v domain.RelatedResources,
	f func(*repository.RelatedResourceInfo) error,

) error {
	if v.Has(index) {
		return nil
	}

	if v.Count()+1 > p.MaxRelatedResourceNum() {
		return ErrorExceedMaxRelatedResourceNum{
			errors.New("exceed max related reousrce num"),
		}
	}

	info := repository.RelatedResourceInfo{
		ResourceToUpdate: s.toResourceToUpdate(p),
		RelatedResource:  *index,
	}

	return f(&info)
}

func (s projectService) RemoveRelatedModel(
	p *domain.Project, index *domain.ResourceIndex,
) error {
	return s.removeRelatedResource(
		p, index, p.RelatedModels, s.repo.RemoveRelatedModel,
	)
}

func (s projectService) RemoveRelatedDataset(
	p *domain.Project, index *domain.ResourceIndex,
) error {
	return s.removeRelatedResource(
		p, index, p.RelatedDatasets, s.repo.RemoveRelatedDataset,
	)
}

func (s projectService) removeRelatedResource(
	p *domain.Project, index *domain.ResourceIndex,
	v domain.RelatedResources,
	f func(*repository.RelatedResourceInfo) error,
) error {
	if !v.Has(index) {
		return nil
	}

	info := repository.RelatedResourceInfo{
		ResourceToUpdate: s.toResourceToUpdate(p),
		RelatedResource:  *index,
	}

	return f(&info)
}

func (s projectService) toResourceToUpdate(p *domain.Project) repository.ResourceToUpdate {
	return repository.ResourceToUpdate{
		Owner:     p.Owner,
		Id:        p.Id,
		Version:   p.Version,
		UpdatedAt: utils.Now(),
	}
}
