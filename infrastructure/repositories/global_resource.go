package repositories

import (
	"github.com/opensourceways/xihe-server/domain/repository"
)

type GlobalResourceListDO = repository.GlobalResourceListOption

func (impl project) ListGlobalAndSortByUpdateTime(
	option *repository.GlobalResourceListOption,
) (repository.UserProjectsInfo, error) {
	return impl.doList(
		func() ([]ProjectSummaryDO, int, error) {
			return impl.mapper.ListGlobalAndSortByUpdateTime(option)
		},
	)
}

func (impl project) ListGlobalAndSortByFirstLetter(
	option *repository.GlobalResourceListOption,
) (repository.UserProjectsInfo, error) {
	return impl.doList(
		func() ([]ProjectSummaryDO, int, error) {
			return impl.mapper.ListGlobalAndSortByFirstLetter(option)
		},
	)
}

func (impl project) ListGlobalAndSortByDownloadCount(
	option *repository.GlobalResourceListOption,
) (repository.UserProjectsInfo, error) {
	return impl.doList(
		func() ([]ProjectSummaryDO, int, error) {
			return impl.mapper.ListGlobalAndSortByDownloadCount(option)
		},
	)
}
