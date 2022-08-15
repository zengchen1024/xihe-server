package repositories

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

func (impl user) AddFollowing(v *domain.Following) error {
	err := impl.mapper.AddFollowing(
		v.Owner.Account(),
		v.Account.Account(),
	)
	if err != nil {
		return convertError(err)
	}

	return nil
}

func (impl user) RemoveFollowing(v *domain.Following) error {
	err := impl.mapper.RemoveFollowing(
		v.Owner.Account(),
		v.Account.Account(),
	)
	if err != nil {
		return convertError(err)
	}

	return nil
}

func (impl user) FindFollowing(owner domain.Account, option repository.FollowFindOption) (
	[]domain.FollowUserInfo, error,
) {
	v, err := impl.mapper.ListFollowing(owner.Account())
	if err != nil {
		return nil, convertError(err)
	}

	if len(v) == 0 {
		return nil, nil
	}

	r := make([]domain.FollowUserInfo, len(v))
	for i := range v {
		if err := v[i].toFollowUserInfo(&r[i]); err != nil {
			return nil, err
		}
	}

	return r, nil
}

type UserInfoDO struct {
	Account  string
	AvatarId string
	Bio      string
}

func (do *UserInfoDO) toFollowUserInfo(r *domain.FollowUserInfo) (err error) {
	if r.Bio, err = domain.NewBio(do.Bio); err != nil {
		return
	}

	if r.Account, _ = domain.NewAccount(do.Account); err != nil {
		return
	}

	if r.AvatarId, _ = domain.NewAvatarId(do.AvatarId); err != nil {
		return
	}

	return
}
