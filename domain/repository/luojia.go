package repository

import (
	"github.com/opensourceways/xihe-server/domain"
)

type LuoJia interface {
	Save(*domain.UserLuoJiaRecord) (domain.LuoJiaRecord, error)
	List(domain.Account) ([]domain.LuoJiaRecord, error)
}
