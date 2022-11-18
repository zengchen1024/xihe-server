package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/platform"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type UserCreateCmd struct {
	Email    domain.Email
	Account  domain.Account
	Password domain.Password

	Bio      domain.Bio
	AvatarId domain.AvatarId
}

func (cmd *UserCreateCmd) Validate() error {
	b := cmd.Email != nil &&
		cmd.Account != nil &&
		cmd.Password != nil

	if !b {
		return errors.New("invalid cmd of creating user")
	}

	return nil
}

func (cmd *UserCreateCmd) toUser() domain.User {
	return domain.User{
		Email:   cmd.Email,
		Account: cmd.Account,

		Bio:      cmd.Bio,
		AvatarId: cmd.AvatarId,
	}
}

type UserDTO struct {
	Id      string `json:"id"`
	Email   string `json:"email"`
	Account string `json:"account"`

	Bio      string `json:"bio"`
	AvatarId string `json:"avatar_id"`

	FollowerCount  int `json:"follower_count"`
	FollowingCount int `json:"following_count"`

	Platform struct {
		UserId      string
		Token       string
		NamespaceId string
	} `json:"-"`
}

type UserService interface {
	Create(*UserCreateCmd) (UserDTO, error)
	UpdateBasicInfo(domain.Account, UpdateUserBasicInfoCmd) error

	GetByAccount(domain.Account) (UserDTO, error)
	GetByFollower(owner, follower domain.Account) (UserDTO, bool, error)

	AddFollowing(*domain.FollowerInfo) error
	RemoveFollowing(*domain.FollowerInfo) error
	ListFollowing(*FollowsListCmd) (FollowsDTO, error)

	AddFollower(*domain.FollowerInfo) error
	RemoveFollower(*domain.FollowerInfo) error
	ListFollower(*FollowsListCmd) (FollowsDTO, error)
}

// ps: platform user service
func NewUserService(
	repo repository.User,
	ps platform.User, sender message.Sender,
) UserService {
	return userService{
		ps:     ps,
		repo:   repo,
		sender: sender,
	}
}

type userService struct {
	ps     platform.User
	repo   repository.User
	sender message.Sender
}

func (s userService) Create(cmd *UserCreateCmd) (dto UserDTO, err error) {
	// TODO keep transaction

	v := cmd.toUser()

	// new code platform user
	pu, err := s.ps.New(platform.UserOption{
		Email:    v.Email,
		Name:     v.Account,
		Password: cmd.Password,
	})
	if err != nil {
		return
	}

	v.PlatformUser = pu

	// apply token
	token, err := s.ps.NewToken(pu)
	if err != nil {
		return
	}

	v.PlatformToken = token

	// update user
	u, err := s.repo.Save(&v)
	if err != nil {
		return
	}

	s.toUserDTO(&u, &dto)

	return
}

func (s userService) GetByAccount(account domain.Account) (dto UserDTO, err error) {
	v, err := s.repo.GetByAccount(account)
	if err != nil {
		return
	}

	s.toUserDTO(&v, &dto)

	return
}

func (s userService) GetByFollower(owner, follower domain.Account) (
	dto UserDTO, isFollower bool, err error,
) {
	v, isFollower, err := s.repo.GetByFollower(owner, follower)
	if err != nil {
		return
	}

	s.toUserDTO(&v, &dto)

	return
}

func (s userService) toUserDTO(u *domain.User, dto *UserDTO) {
	*dto = UserDTO{
		Id:      u.Id,
		Email:   u.Email.Email(),
		Account: u.Account.Account(),
	}

	if u.Bio != nil {
		dto.Bio = u.Bio.Bio()
	}

	if u.AvatarId != nil {
		dto.AvatarId = u.AvatarId.AvatarId()
	}

	dto.FollowerCount = u.FollowerCount
	dto.FollowingCount = u.FollowingCount

	dto.Platform.Token = u.PlatformToken
	dto.Platform.UserId = u.PlatformUser.Id
	dto.Platform.NamespaceId = u.PlatformUser.NamespaceId
}
