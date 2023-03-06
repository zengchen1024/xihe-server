package user

import (
	"github.com/opensourceways/xihe-server/course/domain"
	types "github.com/opensourceways/xihe-server/domain"
)

type User interface {
	AddUserRegInfo(*domain.Student) error
	GetUserRegInfo(types.Account) (domain.Student, error)
}
