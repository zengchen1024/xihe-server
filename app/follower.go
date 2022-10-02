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
	v, err := s.repo.FindFollower(owner, repository.FollowFindOption{})
	if err != nil || len(v) == 0 {
		return
	}

	dtos = make([]FollowDTO, len(v))
	for i := range v {
		s.toFollowDTO(&v[i], &dtos[i])
	}

	return
}
