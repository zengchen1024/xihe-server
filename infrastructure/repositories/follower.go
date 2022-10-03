package repositories

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

func (impl user) GetByFollower(owner, follower domain.Account) (
	u domain.User, isFollower bool, err error,
) {
	account := ""
	if follower != nil {
		account = follower.Account()
	}

	do, isFollower, err := impl.mapper.GetByFollower(owner.Account(), account)
	if err != nil {
		err = convertError(err)
	} else {
		err = do.toUser(&u)
	}

	return
}

func (impl user) AddFollower(v *domain.FollowerInfo) error {
	err := impl.mapper.AddFollower(v.User.Account(), v.Follower.Account())
	if err != nil {
		return convertError(err)
	}

	return nil
}

func (impl user) RemoveFollower(v *domain.FollowerInfo) error {
	err := impl.mapper.RemoveFollower(v.User.Account(), v.Follower.Account())
	if err != nil {
		return convertError(err)
	}

	return nil
}

func (impl user) FindFollower(owner domain.Account, option *repository.FollowFindOption) (
	info repository.FollowerUsersInfo, err error,
) {
	opt := toFollowerUsersInfoListDO(owner, option)

	v, total, err := impl.mapper.ListFollower(&opt)
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
