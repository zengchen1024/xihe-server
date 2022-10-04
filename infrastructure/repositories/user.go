package repositories

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type UserMapper interface {
	Insert(UserDO) (string, error)
	Update(UserDO) error
	GetByFollower(FollowerInfoDO) (do UserDO, isFollower bool, err error)
	ListUsersInfo([]string) ([]UserInfoDO, error)
	GetUserAvatarId(string) (string, error)

	AddFollowing(FollowerInfoDO) error
	RemoveFollowing(FollowerInfoDO) error
	ListFollowing(*FollowerUserInfoListDO) ([]FollowerUserInfoDO, int, error)

	AddFollower(FollowerInfoDO) error
	RemoveFollower(FollowerInfoDO) error
	ListFollower(*FollowerUserInfoListDO) ([]FollowerUserInfoDO, int, error)
}

// TODO: mapper can be mysql
func NewUserRepository(mapper UserMapper) repository.User {
	return user{mapper}
}

type user struct {
	mapper UserMapper
}

func (impl user) GetByAccount(account domain.Account) (r domain.User, err error) {
	do, _, err := impl.mapper.GetByFollower(
		FollowerInfoDO{User: account.Account()},
	)
	if err != nil {
		err = convertError(err)
	} else {
		err = do.toUser(&r)
	}

	return
}

func (impl user) Save(u *domain.User) (r domain.User, err error) {
	if u.Id != "" {
		if err = impl.mapper.Update(impl.toUserDO(u)); err != nil {
			err = convertError(err)
		} else {
			r = *u
			r.Version += 1
		}

		return
	}

	v, err := impl.mapper.Insert(impl.toUserDO(u))
	if err != nil {
		err = convertError(err)
	} else {
		r = *u
		r.Id = v
	}

	return
}

func (impl user) GetUserAvatarId(account domain.Account) (domain.AvatarId, error) {
	d, err := impl.mapper.GetUserAvatarId(account.Account())
	if err != nil {
		return nil, convertError(err)
	}

	return domain.NewAvatarId(d)
}

func (impl user) FindUsersInfo(accounts []domain.Account) (r []domain.UserInfo, err error) {
	v := make([]string, len(accounts))
	for i := range accounts {
		v[i] = accounts[i].Account()
	}

	d, err := impl.mapper.ListUsersInfo(v)
	if err != nil {
		return nil, convertError(err)
	}

	r = make([]domain.UserInfo, len(d))
	for i := range d {
		if r[i].Account, err = domain.NewAccount(d[i].Account); err != nil {
			return nil, err
		}

		if r[i].AvatarId, err = domain.NewAvatarId(d[i].AvatarId); err != nil {
			return nil, err
		}
	}

	return
}

func (impl user) toUserDO(u *domain.User) UserDO {
	do := UserDO{
		Id:      u.Id,
		Email:   u.Email.Email(),
		Account: u.Account.Account(),
	}

	if u.Bio != nil {
		do.Bio = u.Bio.Bio()
	}

	if u.AvatarId != nil {
		do.AvatarId = u.AvatarId.AvatarId()
	}

	do.Platform.Token = u.PlatformToken
	do.Platform.UserId = u.PlatformUser.Id
	do.Platform.NamespaceId = u.PlatformUser.NamespaceId

	do.Version = u.Version

	return do
}

type UserInfoDO struct {
	Account  string
	AvatarId string
}

type UserDO struct {
	Id      string
	Email   string
	Account string

	Bio      string
	AvatarId string

	Platform struct {
		UserId      string
		Token       string
		NamespaceId string
	}

	FollowerCount  int
	FollowingCount int

	Version int
}

func (do *UserDO) toUser(r *domain.User) (err error) {
	r.Id = do.Id

	if r.Bio, err = domain.NewBio(do.Bio); err != nil {
		return
	}

	if r.Email, err = domain.NewEmail(do.Email); err != nil {
		return
	}

	if r.Account, err = domain.NewAccount(do.Account); err != nil {
		return
	}

	if r.AvatarId, err = domain.NewAvatarId(do.AvatarId); err != nil {
		return
	}

	r.FollowerCount = do.FollowerCount
	r.FollowingCount = do.FollowingCount

	r.PlatformToken = do.Platform.Token
	r.PlatformUser.Id = do.Platform.UserId
	r.PlatformUser.NamespaceId = do.Platform.NamespaceId

	r.Version = do.Version

	return
}
