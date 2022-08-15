package app

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

func (s userService) AddFollower(owner, follower domain.Account) error {
	return s.repo.AddFollower(&domain.Follower{
		Owner:   owner,
		Account: follower,
	})
}

func (s userService) RemoveFollower(owner, follower domain.Account) error {
	return s.repo.RemoveFollower(&domain.Follower{
		Owner:   owner,
		Account: follower,
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
