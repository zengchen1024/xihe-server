package app

import (
	"github.com/opensourceways/xihe-server/domain"
)

func (s modelService) AddLike(owner domain.Account, rid string) error {
	return s.repo.AddLike(owner, rid)
}

func (s modelService) RemoveLike(owner domain.Account, rid string) error {
	return s.repo.RemoveLike(owner, rid)
}
