package repository

import "github.com/opensourceways/xihe-server/domain"

type Following interface {
	Save(*domain.Following) error
	Remove(owner domain.Account, folloing domain.Account) error
	Find(domain.Account) ([]domain.Following, error)
}

type Follower interface {
	Save(*domain.Follower) error
	Remove(owner domain.Account, follower domain.Account) error
	Find(me domain.Account) ([]domain.Follower, error)
}
