package repositories

import (
	"github.com/opensourceways/xihe-server/domain"
)

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
