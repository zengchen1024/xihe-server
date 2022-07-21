package repositories

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type UserMapper interface {
	Insert(UserDO) (string, error)
	Get(string) (UserDO, error)
}

// TODO: mapper can be mysql
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

	if r.Bio, err = domain.NewBio(do.Bio); err != nil {
		return
	}

	if r.Email, _ = domain.NewEmail(do.Email); err != nil {
		return
	}

	if r.Account, _ = domain.NewAccount(do.Account); err != nil {
		return
	}

	if r.Nickname, _ = domain.NewNickname(do.Nickname); err != nil {
		return
	}

	if r.AvatarId, _ = domain.NewAvatarId(do.AvatarId); err != nil {
		return
	}

	r.PhoneNumber, err = domain.NewPhoneNumber(do.PhoneNumber)

	return
}

func (impl user) Save(u *domain.User) (r domain.User, err error) {
	if u.Id != "" {
		return
	}

	do := UserDO{
		Bio:         u.Bio.Bio(),
		Email:       u.Email.Email(),
		Account:     u.Account.Account(),
		Password:    u.Password.Password(),
		Nickname:    u.Nickname.Nickname(),
		AvatarId:    u.AvatarId.AvatarId(),
		PhoneNumber: u.PhoneNumber.PhoneNumber(),
	}

	v, err := impl.mapper.Insert(do)
	if err != nil {
		err = convertError(err)
	} else {
		r = *u
		r.Id = v
	}

	return
}

type UserDO struct {
	Id          string `json:"id"`
	Bio         string `json:"bio"`
	Email       string `json:"email"`
	Account     string `json:"account"`
	Password    string `json:"password"`
	Nickname    string `json:"nickname"`
	AvatarId    string `json:"avatar_id"`
	PhoneNumber string `json:"phone_number"`
}
