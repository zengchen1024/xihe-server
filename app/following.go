package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type FollowingCreateCmd struct {
	Owner    domain.Account
	Account  domain.Account
	AvatarId domain.AvatarId
	Bio      domain.Bio
}

func (cmd *FollowingCreateCmd) Validate() error {
	b := cmd.Owner != nil &&
		cmd.Account != nil &&
		cmd.AvatarId != nil

	if !b {
		return errors.New("invalid cmd of creating following")
	}

	return nil
}

func (cmd *FollowingCreateCmd) toFollowing() domain.Following {
	return domain.Following{
		Owner:    cmd.Owner,
		Account:  cmd.Account,
		AvatarId: cmd.AvatarId,
		Bio:      cmd.Bio,
	}
}

type FollowingDTO struct {
	Account  string `json:"account"`
	AvatarId string `json:"avatar_id"`
	Bio      string `json:"bio"`
}

type FollowingService interface {
	Create(cmd *FollowingCreateCmd) (dto FollowingDTO, err error)
	Delete(owner domain.Account, following domain.Account) error
	List(owner domain.Account) (dtos []FollowingDTO, err error)
}

type followingService struct {
	repo repository.Following
}

func (s followingService) Create(cmd *FollowingCreateCmd) (dto FollowingDTO, err error) {
	f := cmd.toFollowing()

	if err = s.repo.Save(&f); err != nil {
		return
	}

	s.toFollowingDTO(&f, &dto)

	// TODO: activity

	// TODO: event

	return
}

func (s followingService) Delete(owner domain.Account, following domain.Account) error {
	return s.repo.Remove(owner, following)
}

func (s followingService) List(owner domain.Account) (
	dtos []FollowingDTO, err error,
) {
	v, err := s.repo.Find(owner)
	if err != nil || len(v) == 0 {
		return
	}

	dtos = make([]FollowingDTO, len(v))
	for i := range v {
		s.toFollowingDTO(&v[i], &dtos[i])
	}

	return
}

func (s followingService) toFollowingDTO(f *domain.Following, dto *FollowingDTO) {
	*dto = FollowingDTO{
		Account:  f.Account.Account(),
		AvatarId: f.AvatarId.AvatarId(),
	}

	if f.Bio != nil {
		dto.Bio = f.Bio.Bio()
	}
}
