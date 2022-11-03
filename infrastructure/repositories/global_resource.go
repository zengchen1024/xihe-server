package repositories

import (
	"github.com/opensourceways/xihe-server/domain/repository"
)

type GlobalResourceListDO = repository.GlobalResourceListOption

func (impl project) GlobalListAndSortByUpdateTime(
	option *repository.GlobalResourceListOption,
) (repository.UserProjectsInfo, error) {
	return impl.doList(func() ([]ProjectSummaryDO, int, error) {
		return impl.mapper.GlobalListAndSortByUpdateTime(option)
	})

}
