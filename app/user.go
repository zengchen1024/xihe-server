package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type UserCreateCmd struct {
	Bio         domain.Bio
	Email       domain.Email
	Account     domain.Account
	Password    domain.Password
	Nickname    domain.Nickname
	AvatarId    domain.AvatarId
	PhoneNumber domain.PhoneNumber
}

func (cmd *UserCreateCmd) Validate() error {
	b := cmd.Email != nil &&
		cmd.Account != nil &&
		cmd.Password != nil &&
		cmd.AvatarId != nil

	if !b {
		return errors.New("invalid cmd of creating user")
	}

	return nil
}

func (cmd *UserCreateCmd) toUser() domain.User {
	return domain.User{
		Bio:         cmd.Bio,
		Email:       cmd.Email,
		Account:     cmd.Account,
		Password:    cmd.Password,
		Nickname:    cmd.Nickname,
		AvatarId:    cmd.AvatarId,
		PhoneNumber: cmd.PhoneNumber,
	}
}

type UserDTO struct {
	Id          string `json:"id"`
	Bio         string `json:"bio"`
	Email       string `json:"email"`
	Account     string `json:"account"`
	Password    string `json:"-"`
	Nickname    string `json:"nickname"`
	AvatarId    string `json:"avatar_id"`
	PhoneNumber string `json:"phone_number"`
}

type UserService interface {
	Create(*UserCreateCmd) (UserDTO, error)
	UpdateBasicInfo(userId string, cmd UpdateUserBasicInfoCmd) error
}

func NewUserService(repo repository.User) UserService {
	return userService{repo}
}

type userService struct {
	repo repository.User
}

func (s userService) Create(cmd *UserCreateCmd) (dto UserDTO, err error) {
	m := cmd.toUser()

	// TODO encrypt password
	v, err := s.repo.Save(&m)
	if err != nil {
		return
	}

	s.toUserDTO(&v, &dto)

	// TODO send event

	return
}

func (s userService) toUserDTO(u *domain.User, dto *UserDTO) {
	// TODO cecrypt password

	*dto = UserDTO{
		Id:          u.Id,
		Bio:         u.Bio.Bio(),
		Email:       u.Email.Email(),
		Account:     u.Account.Account(),
		Password:    u.Password.Password(),
		Nickname:    u.Nickname.Nickname(),
		AvatarId:    u.AvatarId.AvatarId(),
		PhoneNumber: u.PhoneNumber.PhoneNumber(),
	}
}
