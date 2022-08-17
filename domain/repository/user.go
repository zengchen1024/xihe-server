package repository

import (
	"github.com/opensourceways/xihe-server/domain"
)

type FollowFindOption struct {
}

type User interface {
	Save(*domain.User) (domain.User, error)
	GetByAccount(domain.Account) (domain.User, error)
	GetByFollower(owner, follower domain.Account) (domain.User, bool, error)
	FindUsers(Names []domain.Account) ([]domain.UserInfo, error)

	AddFollowing(*domain.Following) error
	RemoveFollowing(*domain.Following) error
	FindFollowing(domain.Account, FollowFindOption) ([]domain.FollowUserInfo, error)

	AddFollower(*domain.Follower) error
	RemoveFollower(*domain.Follower) error
	FindFollower(domain.Account, FollowFindOption) ([]domain.FollowUserInfo, error)
}
