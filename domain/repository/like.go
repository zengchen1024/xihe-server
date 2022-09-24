package repository

import (
	"github.com/opensourceways/xihe-server/domain"
)

type LikeFindOption struct {
}

type Like interface {
	Save(*domain.UserLike) error
	Remove(*domain.UserLike) error
	Find(domain.Account, LikeFindOption) ([]domain.Like, error)
	HasLike(domain.Account, *domain.ResourceObj) (bool, error)
}
