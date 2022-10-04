package app

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type FollowsListCmd struct {
	User domain.Account

	repository.FollowFindOption
}

type FollowsDTO struct {
	Total int         `json:"total"`
	Data  []FollowDTO `json:"data"`
}

type FollowDTO struct {
	Account    string `json:"account"`
	AvatarId   string `json:"avatar_id"`
	Bio        string `json:"bio"`
	IsFollower bool   `json:"is_follower"`
}

func (s userService) AddFollowing(f *domain.FollowerInfo) error {
	err := s.repo.AddFollowing(f)
	if err != nil {
		return err
	}

	// TODO: activity

	// send event
	_ = s.sender.AddFollowing(f)

	return nil
}

func (s userService) RemoveFollowing(f *domain.FollowerInfo) error {
	err := s.repo.RemoveFollowing(f)
	if err != nil {
		return err
	}

	// send event
	_ = s.sender.RemoveFollowing(f)

	return nil
}

func (s userService) ListFollowing(cmd *FollowsListCmd) (
	dto FollowsDTO, err error,
) {
	v, err := s.repo.FindFollowing(cmd.User, &cmd.FollowFindOption)
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
