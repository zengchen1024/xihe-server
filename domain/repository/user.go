package repository

import (
	"github.com/opensourceways/xihe-server/domain"
)

type FollowFindOption struct {
	Follower domain.Account

	CountPerPage int
	PageNum      int
}

type FollowerUsersInfo struct {
	Users []domain.FollowerUserInfo
	Total int
}

type User interface {
	Save(*domain.User) (domain.User, error)
	GetByAccount(domain.Account) (domain.User, error)
	GetByFollower(owner, follower domain.Account) (domain.User, bool, error)
	FindUsersInfo([]domain.Account) ([]domain.UserInfo, error)
	GetUserAvatarId(domain.Account) (domain.AvatarId, error)

	AddFollowing(*domain.FollowerInfo) error
	RemoveFollowing(*domain.FollowerInfo) error
	FindFollowing(domain.Account, *FollowFindOption) (FollowerUsersInfo, error)

	AddFollower(*domain.FollowerInfo) error
	RemoveFollower(*domain.FollowerInfo) error
	FindFollower(domain.Account, *FollowFindOption) (FollowerUsersInfo, error)
}
