package repositories

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

func (impl dataset) AddLike(owner domain.Account, rid string) error {
	err := impl.mapper.AddLike(owner.Account(), rid)
	if err != nil {
		err = convertError(err)
	}

	return err
}

func (impl dataset) RemoveLike(owner domain.Account, rid string) error {
	err := impl.mapper.RemoveLike(owner.Account(), rid)
	if err != nil {
		err = convertError(err)
	}

	return err
}

func (impl dataset) UpdateProperty(info *repository.DatasetPropertyUpdateInfo) error {
	p := &info.Property

	do := DatasetPropertyDO{
		ResourceToUpdateDO: toResourceToUpdateDO(&info.ResourceToUpdate),

		FL:       p.Name.FirstLetterOfName(),
		Name:     p.Name.DatasetName(),
		Desc:     p.Desc.ResourceDesc(),
		RepoType: p.RepoType.RepoType(),
		Tags:     p.Tags,
	}

	if err := impl.mapper.UpdateProperty(&do); err != nil {
		return convertError(err)
	}

	return nil
}

type DatasetPropertyDO struct {
	ResourceToUpdateDO

	FL       byte
	Name     string
	Desc     string
	RepoType string
	Tags     []string
}

func (impl dataset) List(
	owner domain.Account, option *repository.ResourceListOption,
) (repository.UserDatasetsInfo, error) {
	return impl.list(
		owner, option, impl.mapper.List,
	)
}

func (impl dataset) ListAndSortByUpdateTime(
	owner domain.Account, option *repository.ResourceListOption,
) (repository.UserDatasetsInfo, error) {
	return impl.list(
		owner, option, impl.mapper.ListAndSortByUpdateTime,
	)
}

func (impl dataset) ListAndSortByFirstLetter(
	owner domain.Account, option *repository.ResourceListOption,
) (repository.UserDatasetsInfo, error) {
	return impl.list(
		owner, option, impl.mapper.ListAndSortByFirstLetter,
	)
}

func (impl dataset) ListAndSortByDownloadCount(
	owner domain.Account, option *repository.ResourceListOption,
) (repository.UserDatasetsInfo, error) {
	return impl.list(
		owner, option, impl.mapper.ListAndSortByDownloadCount,
	)
}

func (impl dataset) list(
	owner domain.Account,
	option *repository.ResourceListOption,
	f func(string, *ResourceListDO) ([]DatasetSummaryDO, int, error),
) (
	info repository.UserDatasetsInfo, err error,
) {
	do := toResourceListDO(option)

	v, total, err := f(owner.Account(), &do)
	if err != nil {
		err = convertError(err)

		return
	}

	if len(v) == 0 {
		return
	}

	r := make([]domain.DatasetSummary, len(v))
	for i := range v {
		if err = v[i].toDatasetSummary(&r[i]); err != nil {
			r = nil

			return
		}
	}

	info.Datasets = r
	info.Total = total

	return
}

type DatasetSummaryDO struct {
	Id            string
	Owner         string
	Name          string
	Desc          string
	Tags          []string
	UpdatedAt     int64
	LikeCount     int
	DownloadCount int
}

func (do *DatasetSummaryDO) toDatasetSummary(r *domain.DatasetSummary) (err error) {
	r.Id = do.Id

	if r.Owner, err = domain.NewAccount(do.Owner); err != nil {
		return
	}

	if r.Name, err = domain.NewDatasetName(do.Name); err != nil {
		return
	}

	if r.Desc, err = domain.NewResourceDesc(do.Desc); err != nil {
		return
	}

	r.Tags = do.Tags
	r.UpdatedAt = do.UpdatedAt
	r.LikeCount = do.LikeCount
	r.DownloadCount = do.DownloadCount

	return
}
