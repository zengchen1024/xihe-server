package platform

import (
	"github.com/opensourceways/xihe-server/domain"
)

type UserOption struct {
	Name     domain.Account
	Email    domain.Email
	Password domain.Password
}

type User interface {
	New(UserOption) (domain.PlatformUser, error)
	NewToken(domain.PlatformUser) (string, error)
}
