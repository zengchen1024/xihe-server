package app

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type FollowDTO struct {
	Account    string `json:"account"`
	AvatarId   string `json:"avatar_id"`
	Bio        string `json:"bio"`
	IsFollower bool   `json:"is_follower"`
}

func (s userService) AddFollowing(user, follower domain.Account) error {
	f := domain.FollowerInfo{
		User:     user,
		Follower: follower,
	}
	err := s.repo.AddFollowing(&f)
	if err != nil {
		return err
	}

	// TODO: activity

	// send event
	_ = s.sender.AddFollowing(f)

	return nil
}

func (s userService) RemoveFollowing(user, follower domain.Account) error {
	f := domain.FollowerInfo{
		User:     user,
		Follower: follower,
	}
	err := s.repo.RemoveFollowing(&f)
	if err != nil {
		return err
	}

	// send event
	_ = s.sender.RemoveFollowing(f)

	return nil
}

func (s userService) ListFollowing(owner domain.Account) (
	dtos []FollowDTO, err error,
) {
	v, err := s.repo.FindFollowing(owner, &repository.FollowFindOption{})
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

func (s userService) toFollowDTO(f *domain.FollowerUserInfo, dto *FollowDTO) {
	*dto = FollowDTO{
		Account:    f.Account.Account(),
		AvatarId:   f.AvatarId.AvatarId(),
		IsFollower: f.IsFollower,
	}

	if f.Bio != nil {
		dto.Bio = f.Bio.Bio()
	}
}
