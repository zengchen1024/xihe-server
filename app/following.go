package app

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type FollowDTO struct {
	Account  string `json:"account"`
	AvatarId string `json:"avatar_id"`
	Bio      string `json:"bio"`
}

func (s userService) AddFollowing(owner, following domain.Account) error {
	f := domain.Following{
		Owner:   owner,
		Account: following,
	}
	err := s.repo.AddFollowing(&f)
	if err != nil {
		return err
	}

	// TODO: activity

	// send event
	return s.sender.AddFollowing(f)
}

func (s userService) RemoveFollowing(owner, following domain.Account) error {
	f := domain.Following{
		Owner:   owner,
		Account: following,
	}
	err := s.repo.RemoveFollowing(&f)
	if err != nil {
		return err
	}

	// send event
	return s.sender.RemoveFollowing(f)
}

func (s userService) ListFollowing(owner domain.Account) (
	dtos []FollowDTO, err error,
) {
	v, err := s.repo.FindFollowing(owner, repository.FollowFindOption{})
	if err != nil || len(v) == 0 {
		return
	}

	dtos = make([]FollowDTO, len(v))
	for i := range v {
		s.toFollowDTO(&v[i], &dtos[i])
	}

	return
}

func (s userService) toFollowDTO(f *domain.FollowUserInfo, dto *FollowDTO) {
	*dto = FollowDTO{
		Account:  f.Account.Account(),
		AvatarId: f.AvatarId.AvatarId(),
	}

	if f.Bio != nil {
		dto.Bio = f.Bio.Bio()
	}
}
