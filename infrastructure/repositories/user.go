package repositories

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type UserMapper interface {
	Insert(UserDO) (string, error)
	Update(UserDO) error
	Get(string) (UserDO, error)
	GetByAccount(string) (UserDO, error)
}

// TODO: mapper can be mysql
func NewUserRepository(mapper UserMapper) repository.User {
	return user{mapper}
}

type user struct {
	mapper UserMapper
}

func (impl user) Get(uid string) (r domain.User, err error) {
	do, err := impl.mapper.Get(uid)
	if err != nil {
		return
	}

	err = do.toUser(&r)

	return
}

func (impl user) GetByAccount(account domain.Account) (r domain.User, err error) {
	do, err := impl.mapper.GetByAccount(account.Account())
	if err != nil {
		return
	}

	err = do.toUser(&r)

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

	Version int
}

func (do *UserDO) toUser(r *domain.User) (err error) {
	r.Id = do.Id

	if r.Bio, err = domain.NewBio(do.Bio); err != nil {
		return
	}

	if r.Email, _ = domain.NewEmail(do.Email); err != nil {
		return
	}

	if r.Account, _ = domain.NewAccount(do.Account); err != nil {
		return
	}

	if r.AvatarId, _ = domain.NewAvatarId(do.AvatarId); err != nil {
		return
	}

	r.PlatformToken = do.Platform.Token
	r.PlatformUser.Id = do.Platform.UserId
	r.PlatformUser.NamespaceId = do.Platform.NamespaceId

	r.Version = do.Version

	return
}
