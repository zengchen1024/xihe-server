package repository

import (
	"github.com/opensourceways/xihe-server/domain"
)

type DatasetListOption struct {
	Name domain.ProjName
}

type Dataset interface {
	Save(*domain.Dataset) (domain.Dataset, error)
	Get(domain.Account, string) (domain.Dataset, error)
	List(domain.Account, DatasetListOption) ([]domain.Dataset, error)
}
