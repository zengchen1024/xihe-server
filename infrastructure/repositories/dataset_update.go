package repositories

import (
	"github.com/opensourceways/xihe-server/domain"
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
