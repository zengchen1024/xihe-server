package repository

import (
	"github.com/opensourceways/xihe-server/domain"
)

type Login interface {
	Save(*domain.Login) error
	Get(domain.Account) (domain.Login, error)
}
