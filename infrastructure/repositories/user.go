package repositories

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

func NewUserRepository(mapper UserMapper) repository.User {
	return user{mapper}
}

type user struct {
	mapper UserMapper
}

func (impl user) Get(index string) (r domain.User, err error) {
	do, err := impl.mapper.Get(index)
	if err != nil {
		return
	}

	r.Id = do.Id
	r.Bio, _ = domain.NewBio(do.Bio)
	r.Email, _ = domain.NewEmail(do.Email)
	r.Account, _ = domain.NewAccount(do.Account)
	r.Nickname, _ = domain.NewNickname(do.Nickname)
	r.AvatarId, _ = domain.NewAvatarId(do.AvatarId)
	r.PhoneNumber, _ = domain.NewPhoneNumber(do.PhoneNumber)

	return
}

func (impl user) Save(u domain.User) error {
	do := UserDO{
		Id:          u.Id,
		Bio:         u.Bio.Bio(),
		Email:       u.Email.Email(),
		Account:     u.Account.Account(),
		Nickname:    u.Nickname.Nickname(),
		AvatarId:    u.AvatarId.AvatarId(),
		PhoneNumber: u.PhoneNumber.PhoneNumber(),
	}

	return impl.mapper.Update(do)
}

type UserDO struct {
	Id          string
	Bio         string
	Email       string
	Account     string
	Nickname    string
	AvatarId    string
	PhoneNumber string
}

type UserMapper interface {
	Get(string) (UserDO, error)
	Update(UserDO) error
}
