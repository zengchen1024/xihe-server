package repositories

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

func (impl user) AddFollowing(v *domain.FollowerInfo) error {
	err := impl.mapper.AddFollowing(v.User.Account(), v.Follower.Account())
	if err != nil {
		return convertError(err)
	}

	return nil
}

func (impl user) RemoveFollowing(v *domain.FollowerInfo) error {
	err := impl.mapper.RemoveFollowing(v.User.Account(), v.Follower.Account())
	if err != nil {
		return convertError(err)
	}

	return nil
}

func (impl user) FindFollowing(owner domain.Account, option *repository.FollowFindOption) (
	info repository.FollowerUsersInfo, err error,
) {
	opt := toFollowerUsersInfoListDO(owner, option)

	v, total, err := impl.mapper.ListFollowing(&opt)
	if err != nil {
		err = convertError(err)

		return
	}

	if len(v) == 0 {
		return
	}

	r := make([]domain.FollowerUserInfo, len(v))
	for i := range v {
		if err = v[i].toFollowUserInfo(&r[i]); err != nil {
			return
		}
	}

	info.Users = r
	info.Total = total

	return
}

type FollowerUsersInfoListDO struct {
	User         string
	Follower     string
	PageNum      int
	CountPerPage int
}

func toFollowerUsersInfoListDO(
	owner domain.Account, option *repository.FollowFindOption,
) FollowerUsersInfoListDO {
	return FollowerUsersInfoListDO{
		User:         owner.Account(),
		Follower:     option.Follower.Account(),
		PageNum:      option.PageNum,
		CountPerPage: option.CountPerPage,
	}
}

type FollowerUserInfoDO struct {
	Account    string
	AvatarId   string
	Bio        string
	IsFollower bool
}

func (do *FollowerUserInfoDO) toFollowUserInfo(r *domain.FollowerUserInfo) (err error) {
	if r.Bio, err = domain.NewBio(do.Bio); err != nil {
		return
	}

	if r.Account, err = domain.NewAccount(do.Account); err != nil {
		return
	}

	if r.AvatarId, err = domain.NewAvatarId(do.AvatarId); err != nil {
		return
	}

	r.IsFollower = do.IsFollower

	return
}
