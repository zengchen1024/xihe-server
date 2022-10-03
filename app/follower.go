package app

import (
	"github.com/opensourceways/xihe-server/domain"
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

func (s userService) ListFollower(cmd *FollowsListCmd) (
	dto FollowsDTO, err error,
) {
	v, err := s.repo.FindFollower(cmd.User, &cmd.FollowFindOption)
	items := v.Users
	if err != nil || len(items) == 0 {
		return
	}

	dtos := make([]FollowDTO, len(items))
	for i := range items {
		s.toFollowDTO(&items[i], &dtos[i])
	}

	dto.Total = v.Total
	dto.Data = dtos

	return
}
