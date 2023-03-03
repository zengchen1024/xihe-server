package repository

import (
	types "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/user/domain"
)

type UserReg interface {
	AddUserRegInfo(*domain.UserRegInfo) error
	GetUserRegInfo(types.Account) (domain.UserRegInfo, error)
	UpdateUserRegInfo(u *domain.UserRegInfo, version int) error
}
