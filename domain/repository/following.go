package repository

import "github.com/opensourceways/xihe-server/domain"

type FollowFindOption struct {
}

type Following interface {
	Save(*domain.Following) error
	Remove(*domain.Following) error
	Find(domain.Account, FollowFindOption) ([]domain.FollowUserInfo, error)
}

type Follower interface {
	Save(*domain.Follower) error
	Remove(*domain.Follower) error
	Find(domain.Account, FollowFindOption) ([]domain.FollowUserInfo, error)
}
