package repository

import (
	"github.com/opensourceways/xihe-server/domain"
)

type User interface {
	Get(string) (domain.User, error)
	GetByAccount(owner, follower domain.Account) (domain.User, bool, error)
	Save(*domain.User) (domain.User, error)

	AddFollowing(*domain.Following) (domain.User, error)
	RemoveFollowing(*domain.Following) (domain.User, error)
}
