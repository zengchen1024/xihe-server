package repository

import (
	"github.com/opensourceways/xihe-server/domain"
)

type User interface {
	Get(string) (domain.User, error)
	GetByAccount(domain.Account) (domain.User, error)
	Save(*domain.User) (domain.User, error)
}
