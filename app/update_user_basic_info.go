package app

import (
	"github.com/opensourceways/xihe-server/domain"
)

type UpdateUserBasicInfoCmd struct {
	AvatarId domain.AvatarId
	Bio      domain.Bio
}

func (cmd *UpdateUserBasicInfoCmd) toUser(u *domain.User) (changed bool) {
	if !domain.IsSameDomainValue(cmd.AvatarId, u.AvatarId) {
		u.AvatarId = cmd.AvatarId
		changed = true
	}

	if !domain.IsSameDomainValue(cmd.Bio, u.Bio) {
		u.Bio = cmd.Bio
		changed = true
	}

	return
}

func (s userService) UpdateBasicInfo(account domain.Account, cmd UpdateUserBasicInfoCmd) error {
	user, err := s.repo.GetByAccount(account)
	if err != nil {
		return err
	}

	if b := cmd.toUser(&user); !b {
		return nil
	}

	_, err = s.repo.Save(&user)

	return err
}
