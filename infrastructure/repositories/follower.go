package repositories

import (
	"github.com/opensourceways/xihe-server/domain"
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

		return
	}

	err = do.toUser(&u)

	return
}
