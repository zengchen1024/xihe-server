package repository

import (
	"github.com/opensourceways/xihe-server/domain"
)

type Dataset interface {
	Save(*domain.Dataset) (domain.Dataset, error)
	Get(domain.Account, string) (domain.Dataset, error)
}
