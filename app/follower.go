package app

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

func (s userService) AddFollower(user, follower domain.Account) error {
	return s.repo.AddFollower(&domain.FollowerInfo{
		User:     user,
		Follower: follower,
	})
}

func (s userService) RemoveFollower(user, follower domain.Account) error {
	return s.repo.RemoveFollower(&domain.FollowerInfo{
		User:     user,
		Follower: follower,
	})
}

func (s userService) ListFollower(owner domain.Account) (
	dtos []FollowDTO, err error,
) {
	opt := repository.FollowFindOption{Follower: owner}

	v, err := s.repo.FindFollower(owner, &opt)
	items := v.Users
	if err != nil || len(items) == 0 {
		return
	}

	dtos = make([]FollowDTO, len(items))
	for i := range items {
		s.toFollowDTO(&items[i], &dtos[i])
	}

	return
}
