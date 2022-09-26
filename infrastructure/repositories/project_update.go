package repositories

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

func (impl project) IncreaseFork(index *domain.ResourceIndex) error {
	err := impl.mapper.IncreaseFork(
		index.Owner.Account(),
		index.Id,
	)
	if err != nil {
		err = convertError(err)
	}

	return err
}

func (impl project) AddLike(owner domain.Account, pid string) error {
	err := impl.mapper.AddLike(owner.Account(), pid)
	if err != nil {
		err = convertError(err)
	}

	return err
}

func (impl project) RemoveLike(owner domain.Account, pid string) error {
	err := impl.mapper.RemoveLike(owner.Account(), pid)
	if err != nil {
		err = convertError(err)
	}

	return err
}

func (impl project) AddRelatedModel(info *repository.RelatedResourceInfo) error {
	do := toRelatedResourceDO(info)

	if err := impl.mapper.AddRelatedModel(&do); err != nil {
		return convertError(err)
	}

	return nil
}

func (impl project) RemoveRelatedModel(info *repository.RelatedResourceInfo) error {
	do := toRelatedResourceDO(info)

	if err := impl.mapper.RemoveRelatedModel(&do); err != nil {
		return convertError(err)
	}

	return nil
}

func (impl project) AddRelatedDataset(info *repository.RelatedResourceInfo) error {
	do := toRelatedResourceDO(info)

	if err := impl.mapper.AddRelatedDataset(&do); err != nil {
		return convertError(err)
	}

	return nil
}

func (impl project) RemoveRelatedDataset(info *repository.RelatedResourceInfo) error {
	do := toRelatedResourceDO(info)

	if err := impl.mapper.RemoveRelatedDataset(&do); err != nil {
		return convertError(err)
	}

	return nil
}

func (impl project) UpdateProperty(info *repository.ProjectPropertyUpdateInfo) error {
	p := &info.Property

	do := ProjectPropertyDO{
		ResourceToUpdateDO: toResourceToUpdateDO(&info.ResourceToUpdate),

		Name:     p.Name.ProjName(),
		Desc:     p.Desc.ResourceDesc(),
		CoverId:  p.CoverId.CoverId(),
		RepoType: p.RepoType.RepoType(),
		Tags:     p.Tags,
	}

	if err := impl.mapper.UpdateProperty(&do); err != nil {
		return convertError(err)
	}

	return nil
}

type ProjectPropertyDO struct {
	ResourceToUpdateDO

	Name     string
	Desc     string
	CoverId  string
	RepoType string
	Tags     []string
}

func toRelatedResourceDO(info *repository.RelatedResourceInfo) RelatedResourceDO {
	return RelatedResourceDO{
		ResourceToUpdateDO: toResourceToUpdateDO(&info.ResourceToUpdate),
		ResourceOwner:      info.RelatedResource.Owner.Account(),
		ResourceId:         info.RelatedResource.Id,
	}
}

type RelatedResourceDO struct {
	ResourceToUpdateDO

	ResourceOwner string
	ResourceId    string
}

type ResourceToUpdateDO struct {
	Id        string
	Owner     string
	Version   int
	UpdatedAt int64
}

func toResourceToUpdateDO(info *repository.ResourceToUpdate) ResourceToUpdateDO {
	return ResourceToUpdateDO{
		Id:        info.Id,
		Owner:     info.Owner.Account(),
		Version:   info.Version,
		UpdatedAt: info.UpdatedAt,
	}
}

func (impl project) List(
	owner domain.Account, option *repository.ResourceListOption,
) (
	r []domain.Project, err error,
) {
	return impl.list(
		owner, option, impl.mapper.List,
	)
}

func (impl project) ListAndSortByUpdateTime(
	owner domain.Account, option *repository.ResourceListOption,
) (
	r []domain.Project, err error,
) {
	return impl.list(
		owner, option, impl.mapper.ListAndSortByUpdateTime,
	)
}

func (impl project) ListAndSortByFirtLetter(
	owner domain.Account, option *repository.ResourceListOption,
) (
	r []domain.Project, err error,
) {
	return impl.list(
		owner, option, impl.mapper.ListAndSortByFirtLetter,
	)
}

func (impl project) ListAndSortByDownloadCount(
	owner domain.Account, option *repository.ResourceListOption,
) (
	r []domain.Project, err error,
) {
	return impl.list(
		owner, option, impl.mapper.ListAndSortByDownloadCount,
	)
}

func (impl project) list(
	owner domain.Account,
	option *repository.ResourceListOption,
	f func(string, *ResourceListDO) ([]ProjectDO, error),
) (
	r []domain.Project, err error,
) {
	do := toResourceListDO(option)

	v, err := f(owner.Account(), &do)
	if err != nil {
		err = convertError(err)

		return
	}

	r = make([]domain.Project, len(v))
	for i := range v {
		//TODO no need to return detail
		if err = v[i].toProject(&r[i]); err != nil {
			return
		}
	}

	return
}
